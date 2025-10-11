package repository_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/mbilarusdev/durak_auth_bot/internal/interfaces"
	"github.com/mbilarusdev/durak_auth_bot/internal/models"
	"github.com/mbilarusdev/durak_auth_bot/internal/repository"
	"github.com/stretchr/testify/mock"
)

func TestPlayerRepository(t *testing.T) {
	playerID := uint64(123)
	player := &models.Player{
		Username:    "Vasiliy",
		PhoneNumber: "+79680719568",
		ChatID:      50,
		CreatedAt:   time.Now().UTC().UnixMilli(),
	}
	t.Run("Player Insert with success", func(t *testing.T) {
		pool := interfaces.NewMockDBPool(t)
		conn := interfaces.NewMockDBConn(t)
		row := interfaces.NewMockDBRow(t)
		pool.On("Acquire", mock.Anything).Once().Return(conn, nil)
		conn.On("Release").Once()
		conn.On(
			"QueryRow",
			mock.Anything,
			"INSERT INTO players (username, phone_number, chat_id, created_at) VALUES ($1, $2, $3, $4) RETURNING id;",
			[]any{player.Username, player.PhoneNumber, player.ChatID, player.CreatedAt},
		).Once().Return(row)

		var playerID uint64
		row.On("Scan", []any{&playerID}).
			Once().
			Return(nil)

		repository := repository.NewPlayerRepository(pool)

		res, err := repository.Insert(player)

		if res != playerID || err != nil {
			t.Errorf("Insert() failed: %v", err)
		}
	})
	t.Run("Player FindOne by ID with success", func(t *testing.T) {
		pool := interfaces.NewMockDBPool(t)
		conn := interfaces.NewMockDBConn(t)
		row := interfaces.NewMockDBRow(t)
		pool.On("Acquire", mock.Anything).Once().Return(conn, nil)
		conn.On("Release").Once()
		conn.On(
			"QueryRow",
			mock.Anything,
			"SELECT * FROM players WHERE id = $1 LIMIT 1;",
			[]any{[]any{playerID}},
		).Once().Return(row)

		findedPlayer := new(models.Player)
		row.On("Scan", []any{findedPlayer}).
			Once().
			Return(nil)

		repository := repository.NewPlayerRepository(pool)

		_, err := repository.FindOne(&models.PlayerFindOptions{ID: playerID})

		if err != nil {
			t.Errorf("FindOne() by ID failed: %v", err)
		}
	})

	t.Run("Player FindOne by ChatID with success", func(t *testing.T) {
		pool := interfaces.NewMockDBPool(t)
		conn := interfaces.NewMockDBConn(t)
		row := interfaces.NewMockDBRow(t)
		pool.On("Acquire", mock.Anything).Once().Return(conn, nil)
		conn.On("Release").Once()
		conn.On(
			"QueryRow",
			mock.Anything,
			"SELECT * FROM players WHERE chat_id = $1 LIMIT 1;",
			[]any{[]any{player.ChatID}},
		).Once().Return(row)

		findedPlayer := new(models.Player)
		row.On("Scan", []any{findedPlayer}).
			Once().
			Return(nil)

		repository := repository.NewPlayerRepository(pool)

		_, err := repository.FindOne(&models.PlayerFindOptions{ChatID: player.ChatID})

		if err != nil {
			t.Errorf("FindOne() by ChatID failed: %v", err)
		}
	})

	t.Run("Player FindOne by phone number with success", func(t *testing.T) {
		pool := interfaces.NewMockDBPool(t)
		conn := interfaces.NewMockDBConn(t)
		row := interfaces.NewMockDBRow(t)
		pool.On("Acquire", mock.Anything).Once().Return(conn, nil)
		conn.On("Release").Once()
		conn.On(
			"QueryRow",
			mock.Anything,
			"SELECT * FROM players WHERE phone_number = $1 LIMIT 1;",
			[]any{[]any{player.PhoneNumber}},
		).Once().Return(row)

		findedPlayer := new(models.Player)
		row.On("Scan", []any{findedPlayer}).
			Once().
			Return(nil)

		repository := repository.NewPlayerRepository(pool)

		_, err := repository.FindOne(&models.PlayerFindOptions{PhoneNumber: player.PhoneNumber})

		if err != nil {
			t.Errorf("FindOne() by PhoneNumber failed: %v", err)
		}
	})

	t.Run("Player FindOne by all options with success", func(t *testing.T) {
		pool := interfaces.NewMockDBPool(t)
		conn := interfaces.NewMockDBConn(t)
		row := interfaces.NewMockDBRow(t)
		pool.On("Acquire", mock.Anything).Once().Return(conn, nil)
		conn.On("Release").Once()
		conn.On(
			"QueryRow",
			mock.Anything,
			"SELECT * FROM players WHERE id = $1 AND phone_number = $2 AND chat_id = $3 LIMIT 1;",
			[]any{[]any{playerID, player.PhoneNumber, player.ChatID}},
		).Once().Return(row)

		findedPlayer := new(models.Player)
		row.On("Scan", []any{findedPlayer}).
			Once().
			Return(nil)

		repository := repository.NewPlayerRepository(pool)

		_, err := repository.FindOne(
			&models.PlayerFindOptions{
				ID:          playerID,
				PhoneNumber: player.PhoneNumber,
				ChatID:      player.ChatID,
			},
		)

		if err != nil {
			t.Errorf("FindOne() by all options failed: %v", err)
		}
	})

	t.Run("Player FindOne by all options when not finded", func(t *testing.T) {
		pool := interfaces.NewMockDBPool(t)
		conn := interfaces.NewMockDBConn(t)
		row := interfaces.NewMockDBRow(t)
		pool.On("Acquire", mock.Anything).Once().Return(conn, nil)
		conn.On("Release").Once()
		conn.On(
			"QueryRow",
			mock.Anything,
			"SELECT * FROM players WHERE id = $1 AND phone_number = $2 AND chat_id = $3 LIMIT 1;",
			[]any{[]any{playerID, player.PhoneNumber, player.ChatID}},
		).Once().Return(row)

		findedPlayer := new(models.Player)
		row.On("Scan", []any{findedPlayer}).
			Once().
			Return(sql.ErrNoRows)

		repository := repository.NewPlayerRepository(pool)

		_, err := repository.FindOne(
			&models.PlayerFindOptions{
				ID:          playerID,
				PhoneNumber: player.PhoneNumber,
				ChatID:      player.ChatID,
			},
		)
		fmt.Println(err)

		if err == nil || err != sql.ErrNoRows {
			t.Errorf(
				"FindOne() expected to return sql.ErrNoRows error but returned nil or other error",
			)
		}
	})
}
