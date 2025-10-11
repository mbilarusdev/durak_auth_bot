package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/mbilarusdev/durak_auth_bot/internal/interfaces"
	"github.com/mbilarusdev/durak_auth_bot/internal/models"
)

type TokenProvider interface {
	Insert(token *models.Token) (uint64, error)
	FindOne(options *models.TokenFindOptions) (*models.Token, error)
	UpdateStatus(ID uint64, status string) error
}

type TokenRepository struct {
	pool interfaces.DBPool
}

func NewTokenRepository(pool interfaces.DBPool) *TokenRepository {
	repository := new(TokenRepository)
	repository.pool = pool

	return repository
}

func (repository *TokenRepository) Insert(token *models.Token) (uint64, error) {
	ctx := context.Background()
	conn, err := repository.pool.Acquire(ctx)
	if err != nil {
		log.Println("Ошибка при открытии соединения pgx")
		return 0, err
	}
	defer conn.Release()
	var tokenID uint64
	if err := conn.QueryRow(
		ctx,
		"INSERT INTO tokens (jwt, player_id, status) VALUES ($1, $2, $3) RETURNING id;",
		token.Jwt,
		token.PlayerID,
		token.Status,
	).Scan(&tokenID); err != nil {
		log.Println("Ошибка при вставке нового токена")
		return 0, err
	}
	return tokenID, nil
}

func (repository *TokenRepository) FindOne(
	options *models.TokenFindOptions,
) (*models.Token, error) {
	ctx := context.Background()
	conn, err := repository.pool.Acquire(ctx)
	if err != nil {
		log.Println("Ошибка при открытии соединения pgx")
		return nil, err
	}
	defer conn.Release()
	query := "SELECT * FROM tokens WHERE "
	args := []any{}
	argNum := 0

	if options.ID != 0 {
		argNum += 1
		query += fmt.Sprintf("id = $%v AND ", argNum)
		args = append(args, options.ID)
	}

	if options.PlayerID != 0 {
		argNum += 1
		query += fmt.Sprintf("player_id = $%v AND ", argNum)
		args = append(args, options.PlayerID)
	}

	query = strings.TrimSuffix(query, "AND ") + "LIMIT 1;"
	findedToken := new(models.Token)
	if err := conn.QueryRow(ctx, query, args).Scan(findedToken); err != nil {
		if err == sql.ErrNoRows {
			log.Println("Не найдено токена по данному поисковому запросу")
			return nil, err
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
	if err := conn.QueryRow(ctx, "UPDATE tokens SET status = $1 WHERE id = $2;", status, ID).Scan(); err != nil {
		if err == sql.ErrNoRows {
			log.Println("Не найдено токена с данным id: ", ID)
			return err
		}
		log.Println("Ошибка при обновлении токена")
		return err
	}
	return nil
}
