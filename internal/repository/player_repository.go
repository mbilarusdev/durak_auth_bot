package repository

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/mbilarusdev/durak_auth_bot/internal/interfaces"
	"github.com/mbilarusdev/durak_auth_bot/internal/models"
)

type PlayerProvider interface {
	Insert(player *models.Player) (uint64, error)
	FindOne(options *models.PlayerFindOptions) (*models.Player, error)
}

type PlayerRepository struct {
	pool interfaces.DBPool
}

func NewPlayerRepository(pool interfaces.DBPool) *PlayerRepository {
	repository := new(PlayerRepository)
	repository.pool = pool
	return repository
}

func (repository *PlayerRepository) Insert(player *models.Player) (uint64, error) {
	ctx := context.Background()
	conn, err := repository.pool.Acquire(ctx)
	if err != nil {
		log.Println("Ошибка при открытии соединения pgx")
		return 0, err
	}
	defer conn.Release()
	var playerID uint64
	if err := conn.QueryRow(
		ctx,
		"INSERT INTO players (username, phone_number, chat_id, created_at) VALUES ($1, $2, $3, $4) RETURNING id;",
		player.Username,
		player.PhoneNumber,
		player.ChatID,
		player.CreatedAt,
	).Scan(&playerID); err != nil {
		log.Println("Ошибка при вставке нового игрока")
		return 0, err
	}

	return playerID, nil
}

func (repository *PlayerRepository) FindOne(
	options *models.PlayerFindOptions,
) (*models.Player, error) {
	ctx := context.Background()
	conn, err := repository.pool.Acquire(ctx)
	if err != nil {
		log.Println("Ошибка при открытии соединения pgx")
		return nil, err
	}
	defer conn.Release()
	query := "SELECT * FROM players WHERE "
	args := []any{}
	argNum := 0

	if options.ID != 0 {
		argNum += 1
		query += fmt.Sprintf("id = $%v AND ", argNum)
		args = append(args, options.ID)
	}

	if options.PhoneNumber != "" {
		argNum += 1
		query += fmt.Sprintf("phone_number = $%v AND ", argNum)
		args = append(args, options.PhoneNumber)
	}

	if options.ChatID != 0 {
		argNum += 1
		query += fmt.Sprintf("chat_id = $%v AND ", argNum)
		args = append(args, options.ChatID)
	}

	query = strings.TrimSuffix(query, "AND ") + "LIMIT 1;"
	findedPlayer := new(models.Player)
	if err := conn.QueryRow(
		ctx,
		query,
		args...,
	).Scan(&findedPlayer.ID, &findedPlayer.Username, &findedPlayer.PhoneNumber, &findedPlayer.ChatID, &findedPlayer.CreatedAt); err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("Не найдено игрока по данному поисковому запросу")
			return nil, err
		}
		log.Println("Ошибка при поиске игрока")
		return nil, err
	}
	return findedPlayer, nil
}
