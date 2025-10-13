package repository

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/mbilarusdev/durak_auth_bot/internal/interfaces"
	app_model "github.com/mbilarusdev/durak_auth_bot/internal/structs/app/model"
	app_option "github.com/mbilarusdev/durak_auth_bot/internal/structs/app/option"
)

type TokenProvider interface {
	Insert(token *app_model.Token) (uint64, error)
	FindOne(options *app_option.TokenFindOptions) (*app_model.Token, error)
	UpdateStatus(ID uint64, status app_model.TokenStatus) error
}

type TokenRepository struct {
	pool interfaces.DBPool
}

func NewTokenRepository(pool interfaces.DBPool) *TokenRepository {
	repository := new(TokenRepository)
	repository.pool = pool

	return repository
}

func (repository *TokenRepository) Insert(token *app_model.Token) (uint64, error) {
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
	options *app_option.TokenFindOptions,
) (*app_model.Token, error) {
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
	findedToken := new(app_model.Token)
	if err := conn.QueryRow(ctx, query, args...).Scan(&findedToken.ID, &findedToken.PlayerID, &findedToken.Jwt, &findedToken.Status); err != nil {
		if err == pgx.ErrNoRows {
			log.Println("Не найдено токена по данному поисковому запросу")
			return nil, err
		}
		log.Println("Ошибка при поиске токена")
		return nil, err
	}
	return findedToken, nil
}

func (repository *TokenRepository) UpdateStatus(ID uint64, status app_model.TokenStatus) error {
	ctx := context.Background()
	conn, err := repository.pool.Acquire(ctx)
	if err != nil {
		log.Println("Ошибка при открытии соединения pgx")
		return err
	}
	defer conn.Release()
	if err := conn.QueryRow(ctx, "UPDATE tokens SET status = $1 WHERE id = $2;", status, ID).Scan(); err != nil {
		if err == pgx.ErrNoRows {
			return nil
		}
		log.Println("Ошибка при обновлении токена")
		return err
	}
	return nil
}
