package repository_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/mbilarusdev/durak_auth_bot/internal/interfaces"
	"github.com/mbilarusdev/durak_auth_bot/internal/models"
	"github.com/mbilarusdev/durak_auth_bot/internal/repository"
	"github.com/stretchr/testify/mock"
)

func TestTokenRepository(t *testing.T) {
	tokenID := uint64(123)
	token := &models.Token{
		PlayerID: 10,
		Jwt:      "abvgd",
		Status:   models.TokenAvailable,
	}
	t.Run("Token Insert with success", func(t *testing.T) {
		pool := interfaces.NewMockDBPool(t)
		conn := interfaces.NewMockDBConn(t)
		row := interfaces.NewMockDBRow(t)
		pool.On("Acquire", mock.Anything).Once().Return(conn, nil)
		conn.On("Release").Once()
		conn.On(
			"QueryRow",
			mock.Anything,
			"INSERT INTO tokens (jwt, player_id, status) VALUES ($1, $2, $3) RETURNING id;",
			[]any{token.Jwt, token.PlayerID, token.Status},
		).Once().Return(row)

		var tokenID uint64
		row.On("Scan", []any{&tokenID}).
			Once().
			Return(nil)

		repository := repository.NewTokenRepository(pool)

		_, err := repository.Insert(token)

		if err != nil {
			t.Errorf("Insert() failed: %v", err)
		}
	})

	t.Run("Token FindOne by ID with success", func(t *testing.T) {
		pool := interfaces.NewMockDBPool(t)
		conn := interfaces.NewMockDBConn(t)
		row := interfaces.NewMockDBRow(t)
		pool.On("Acquire", mock.Anything).Once().Return(conn, nil)
		conn.On("Release").Once()
		conn.On(
			"QueryRow",
			mock.Anything,
			"SELECT * FROM tokens WHERE id = $1 LIMIT 1;",
			[]any{[]any{tokenID}},
		).Once().Return(row)

		findedToken := new(models.Token)
		row.On("Scan", []any{findedToken}).
			Once().
			Return(nil)

		repository := repository.NewTokenRepository(pool)

		_, err := repository.FindOne(&models.TokenFindOptions{ID: tokenID})

		if err != nil {
			t.Errorf("FindOne() by ID failed: %v", err)
		}
	})

	t.Run("Token FindOne by Player ID with success", func(t *testing.T) {
		pool := interfaces.NewMockDBPool(t)
		conn := interfaces.NewMockDBConn(t)
		row := interfaces.NewMockDBRow(t)
		pool.On("Acquire", mock.Anything).Once().Return(conn, nil)
		conn.On("Release").Once()
		conn.On(
			"QueryRow",
			mock.Anything,
			"SELECT * FROM tokens WHERE player_id = $1 LIMIT 1;",
			[]any{[]any{token.PlayerID}},
		).Once().Return(row)

		findedToken := new(models.Token)
		row.On("Scan", []any{findedToken}).
			Once().
			Return(nil)

		repository := repository.NewTokenRepository(pool)

		_, err := repository.FindOne(&models.TokenFindOptions{PlayerID: token.PlayerID})

		if err != nil {
			t.Errorf("FindOne() by Player ID failed: %v", err)
		}
	})

	t.Run("Token FindOne by all options with success", func(t *testing.T) {
		pool := interfaces.NewMockDBPool(t)
		conn := interfaces.NewMockDBConn(t)
		row := interfaces.NewMockDBRow(t)
		pool.On("Acquire", mock.Anything).Once().Return(conn, nil)
		conn.On("Release").Once()
		conn.On(
			"QueryRow",
			mock.Anything,
			"SELECT * FROM tokens WHERE id = $1 AND player_id = $2 LIMIT 1;",
			[]any{[]any{tokenID, token.PlayerID}},
		).Once().Return(row)

		findedToken := new(models.Token)
		row.On("Scan", []any{findedToken}).
			Once().
			Return(nil)

		repository := repository.NewTokenRepository(pool)

		_, err := repository.FindOne(
			&models.TokenFindOptions{ID: tokenID, PlayerID: token.PlayerID},
		)

		if err != nil {
			t.Errorf("FindOne() by all options failed: %v", err)
		}
	})

	t.Run("Token FindOne by all options when not finded", func(t *testing.T) {
		pool := interfaces.NewMockDBPool(t)
		conn := interfaces.NewMockDBConn(t)
		row := interfaces.NewMockDBRow(t)
		pool.On("Acquire", mock.Anything).Once().Return(conn, nil)
		conn.On("Release").Once()
		conn.On(
			"QueryRow",
			mock.Anything,
			"SELECT * FROM tokens WHERE id = $1 AND player_id = $2 LIMIT 1;",
			[]any{[]any{tokenID, token.PlayerID}},
		).Once().Return(row)

		findedToken := new(models.Token)
		row.On("Scan", []any{findedToken}).
			Once().
			Return(sql.ErrNoRows)

		repository := repository.NewTokenRepository(pool)

		_, err := repository.FindOne(
			&models.TokenFindOptions{
				ID:       tokenID,
				PlayerID: token.PlayerID,
			},
		)
		fmt.Println(err)

		if err == nil || err != sql.ErrNoRows {
			t.Errorf(
				"FindOne() expected to return sql.ErrNoRows error but returned nil or other error",
			)
		}
	})

	t.Run("Token Update Status by ID with success", func(t *testing.T) {
		pool := interfaces.NewMockDBPool(t)
		conn := interfaces.NewMockDBConn(t)
		row := interfaces.NewMockDBRow(t)
		pool.On("Acquire", mock.Anything).Once().Return(conn, nil)
		conn.On("Release").Once()
		conn.On(
			"QueryRow",
			mock.Anything,
			"UPDATE tokens SET status = $1 WHERE id = $2;",
			[]any{models.TokenBlocked, tokenID},
		).Once().Return(row)

		row.On("Scan").
			Once().
			Return(nil)

		repository := repository.NewTokenRepository(pool)

		err := repository.UpdateStatus(tokenID, models.TokenBlocked)

		if err != nil {
			t.Errorf("UpdateStatus() by ID failed: %v", err)
		}
	})

	t.Run("Token Update Status by ID when not finded", func(t *testing.T) {
		pool := interfaces.NewMockDBPool(t)
		conn := interfaces.NewMockDBConn(t)
		row := interfaces.NewMockDBRow(t)
		pool.On("Acquire", mock.Anything).Once().Return(conn, nil)
		conn.On("Release").Once()
		conn.On(
			"QueryRow",
			mock.Anything,
			"UPDATE tokens SET status = $1 WHERE id = $2;",
			[]any{models.TokenBlocked, tokenID},
		).Once().Return(row)

		row.On("Scan").
			Once().
			Return(sql.ErrNoRows)

		repository := repository.NewTokenRepository(pool)

		err := repository.UpdateStatus(tokenID, models.TokenBlocked)

		if err == nil || err != sql.ErrNoRows {
			t.Errorf(
				"UpdateStatus() expected to return sql.ErrNoRows error but returned nil or other error",
			)
		}
	})
}
