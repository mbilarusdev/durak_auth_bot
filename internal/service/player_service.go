package service

import (
	"time"

	"github.com/mbilarusdev/durak_auth_bot/internal/models"
	"github.com/mbilarusdev/durak_auth_bot/internal/repository"
)

type PlayerManager interface {
	FindByPhone(phone string) (*models.Player, error)
	FindByChatID(chatID int) (*models.Player, error)
	CreatePlayer(phone string, chatID int) (*models.Player, error)
}

type PlayerService struct {
	playerRepository repository.PlayerProvider
}

func NewPlayerService(playerRepository repository.PlayerProvider) *PlayerService {
	service := new(PlayerService)
	service.playerRepository = playerRepository
	return service
}

func (service *PlayerService) FindByPhone(phone string) (*models.Player, error) {
	player, err := service.playerRepository.FindOne(
		&models.FindOptions{PhoneNumber: phone},
	)

	return player, err
}

func (service *PlayerService) FindByChatID(chatID int) (*models.Player, error) {
	player, err := service.playerRepository.FindOne(
		&models.FindOptions{ChatID: chatID},
	)

	return player, err
}

func (service *PlayerService) CreatePlayer(phone string, chatID int) (*models.Player, error) {
	player, err := service.playerRepository.Insert(&models.Player{
		PhoneNumber: phone,
		ChatID:      chatID,
		CreatedAt:   time.Now().UTC().UnixMilli(),
	})

	return player, err
}
