package service

import (
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/mbilarusdev/durak_auth_bot/internal/repository"
	app_model "github.com/mbilarusdev/durak_auth_bot/internal/structs/app/model"
	app_option "github.com/mbilarusdev/durak_auth_bot/internal/structs/app/option"
)

type PlayerManager interface {
	FindByPhone(phone string) (*app_model.Player, error)
	FindByChatID(chatID int64) (*app_model.Player, error)
	CreatePlayer(phone string, chatID int64) (*app_model.Player, error)
}

type PlayerService struct {
	playerRepository repository.PlayerProvider
}

func NewPlayerService(playerRepository repository.PlayerProvider) *PlayerService {
	service := new(PlayerService)
	service.playerRepository = playerRepository
	return service
}

func (service *PlayerService) FindByPhone(phone string) (*app_model.Player, error) {
	player, err := service.playerRepository.FindOne(
		&app_option.PlayerFindOptions{PhoneNumber: phone},
	)

	return player, err
}

func (service *PlayerService) FindByChatID(chatID int64) (*app_model.Player, error) {
	player, err := service.playerRepository.FindOne(
		&app_option.PlayerFindOptions{ChatID: chatID},
	)

	return player, err
}

func (service *PlayerService) CreatePlayer(phone string, chatID int64) (*app_model.Player, error) {
	newPlayerID, err := service.playerRepository.Insert(&app_model.Player{
		PhoneNumber: phone,
		ChatID:      chatID,
		CreatedAt:   time.Now().UTC().UnixMilli(),
	})
	if err != nil {
		return nil, err
	}

	player, err := service.playerRepository.FindOne(&app_option.PlayerFindOptions{ID: newPlayerID})
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	return player, err
}
