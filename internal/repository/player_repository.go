package repository

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mbilarusdev/durak_auth_bot/internal/locator"
	"github.com/mbilarusdev/durak_auth_bot/internal/models"
)

type PlayerProvider interface {
	Insert(player *models.Player) (*models.Player, error)
	FindOne(options *models.FindOptions) (*models.Player, error)
}

type PlayerRepository struct{}

func NewPlayerRepository() *PlayerRepository {
	return new(PlayerRepository)
}

func (repository *PlayerRepository) Insert(player *models.Player) (*models.Player, error) {
	pool := locator.Instance.Get("pgx_pool").(*pgxpool.Pool)
	ctx := context.Background()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Println("Ошибка при открытии соединения pgx")
		return nil, err
	}
	defer conn.Release()
	var playerID uint64
	if err := conn.QueryRow(
		ctx,
		"INSERT INTO players (username, phone_number, chat_id, created_at) VALUES ($1, $2, $3) RETURNING id;",
		player.Username,
		player.PhoneNumber,
		player.ChatID,
		player.CreatedAt,
	).Scan(&playerID); err != nil {
		log.Println("Ошибка при вставке нового игрока")
		return nil, err
	}
	newPlayer := new(models.Player)
	if err := conn.QueryRow(ctx, "SELECT * FROM players WHERE id = $1 LIMIT 1;", playerID).Scan(newPlayer); err != nil {
		if err == sql.ErrNoRows {
			log.Println("Не найдено игрока с данным идентификатором: ", playerID)
			return nil, nil
		}
		log.Println("Ошибка при выборке нового игрока")
		return nil, err
	}
	return newPlayer, nil
}

func (repository *PlayerRepository) FindOne(options *models.FindOptions) (*models.Player, error) {
	pool := locator.Instance.Get("pgx_pool").(*pgxpool.Pool)
	ctx := context.Background()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Println("Ошибка при открытии соединения pgx")
		return nil, err
	}
	defer conn.Release()
	query := "SELECT * FROM players WHERE "
	args := []any{}

	if options.PhoneNumber != "" {
		query += "phone_number = $1 AND "
		args = append(args, options.PhoneNumber)
	}

	if options.ChatID != 0 {
		query += "chat_id = $2 AND "
		args = append(args, options.ChatID)
	}

	query = strings.TrimSuffix(query, "AND ") + "LIMIT 1;"
	findedPlayer := new(models.Player)
	if err := conn.QueryRow(
		ctx,
		query,
		args,
	).Scan(findedPlayer); err != nil {
		if err == sql.ErrNoRows {
			log.Printf(
				"Не найдено игрока с номером телефона = %v, и/или айди чата = %v",
				options.PhoneNumber,
				options.ChatID,
			)
			return nil, nil
		}
		log.Println("Ошибка при поиске игрока по номеру телефона")
		return nil, err
	}
	return findedPlayer, nil
}
