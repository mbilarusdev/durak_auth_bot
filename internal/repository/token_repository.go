package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mbilarusdev/durak_auth_bot/internal/interfaces"
	"github.com/mbilarusdev/durak_auth_bot/internal/models"
)

type TokenProvider interface {
	Insert(token *models.Token) (*models.Token, error)
	FindOne(playerID uint64) (*models.Token, error)
	UpdateStatus(ID uint64, status string) error
}

type TokenRepository struct {
	pool interfaces.DBPool
}

func NewTokenRepository(pool *pgxpool.Pool) *TokenRepository {
	repository := new(TokenRepository)
	repository.pool = pool

	return repository
}

func (repository *TokenRepository) Insert(token *models.Token) (*models.Token, error) {
	ctx := context.Background()
	conn, err := repository.pool.Acquire(ctx)
	if err != nil {
		log.Println("Ошибка при открытии соединения pgx")
		return nil, err
	}
	defer conn.Release()
	var tokenID uint64
	if err := conn.QueryRow(
		ctx,
		"INSERT INTO tokens (jwt, player_id) VALUES ($1, $2) RETURNING id;",
		token.Jwt,
		token.PlayerID,
	).Scan(&tokenID); err != nil {
		log.Println("Ошибка при вставке нового токена")
		return nil, err
	}
	newToken := new(models.Token)
	if err := conn.QueryRow(ctx, "SELECT * FROM tokens WHERE id = $1 LIMIT 1;", tokenID).Scan(newToken); err != nil {
		if err == sql.ErrNoRows {
			log.Println("Не найдено токена с данным идентификатором: ", tokenID)
			return nil, nil
		}
		log.Println("Ошибка при выборке нового токена")
		return nil, err
	}
	return newToken, nil
}

func (repository *TokenRepository) FindOne(playerID uint64) (*models.Token, error) {
	ctx := context.Background()
	conn, err := repository.pool.Acquire(ctx)
	if err != nil {
		log.Println("Ошибка при открытии соединения pgx")
		return nil, err
	}
	defer conn.Release()
	findedToken := new(models.Token)
	if err := conn.QueryRow(ctx, "SELECT * FROM tokens WHERE player_id = $1 LIMIT 1;", playerID).Scan(findedToken); err != nil {
		if err == sql.ErrNoRows {
			log.Println("Не найдено токена для игрока с данным ID: ", playerID)
			return nil, nil
		}
		log.Println("Ошибка при поиске токена")
		return nil, err
	}
	return findedToken, nil
}

func (repository *TokenRepository) UpdateStatus(ID uint64, status string) error {
	ctx := context.Background()
	conn, err := repository.pool.Acquire(ctx)
	if err != nil {
		log.Println("Ошибка при открытии соединения pgx")
		return err
	}
	defer conn.Release()
	updatedToken := new(models.Token)
	if err := conn.QueryRow(ctx, "UPDATE tokens SET status = $1 WHERE id = $2;", status, ID).Scan(updatedToken); err != nil {
		if err == sql.ErrNoRows {
			log.Println("Не найдено токена с данным id: ", ID)
			return nil
		}
		log.Println("Ошибка при обновлении токена")
		return err
	}
	return nil
}
