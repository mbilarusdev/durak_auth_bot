package service

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/mbilarusdev/durak_auth_bot/internal/client"
	"github.com/mbilarusdev/durak_auth_bot/internal/locator"
	"github.com/mbilarusdev/durak_auth_bot/internal/models"
	"github.com/mbilarusdev/durak_auth_bot/internal/repository"
)

type UpdatesManager interface {
	Process(updates []models.Update)
	Get(offset int) ([]models.Update, error)
}

type UpdatesService struct {
	messageService   MessageManager
	playerRepository repository.PlayerProvider
}

func NewUpdatesService(
	messageService *MessageService,
	playerRepository *repository.PlayerRepository,
) *UpdatesService {
	service := &UpdatesService{}
	service.messageService = messageService
	service.playerRepository = playerRepository

	return service
}

func (service *UpdatesService) Process(updates []models.Update) {
	errMsg := "Произошла техническая ошибка, попробуйте позже..."
	for _, upd := range updates {
		if upd.Message != nil && upd.Message.Contact != nil {
			log.Printf(
				"Контакт получен: Имя=%s, Телефон=%s\n",
				upd.Message.Contact.FirstName,
				upd.Message.Contact.PhoneNumber,
			)
			player, err := service.playerRepository.FindOne(
				&models.FindOptions{PhoneNumber: upd.Message.Contact.PhoneNumber},
			)
			if err != nil {
				service.messageService.Send(errMsg, upd.Message.Chat.ID)
				continue
			}
			if player != nil {
				service.messageService.Send("Ваш контакт уже был сохранен!", upd.Message.Chat.ID)
				continue
			}
			_, err = service.playerRepository.Insert(&models.Player{
				PhoneNumber: upd.Message.Contact.PhoneNumber,
				ChatID:      upd.Message.Chat.ID,
				CreatedAt:   time.Now().UTC().UnixMilli(),
			})
			if err != nil {
				service.messageService.Send(errMsg, upd.Message.Chat.ID)
			}
			service.messageService.Send("Ваш контакт успешно сохранен!", upd.Message.Chat.ID)
		} else if upd.Message != nil && upd.Message.Text == "/start" {
			player, err := service.playerRepository.FindOne(&models.FindOptions{ChatID: upd.Message.Chat.ID})
			if err != nil {
				service.messageService.Send(errMsg, upd.Message.Chat.ID)
				continue
			}
			if player == nil {
				service.messageService.Send("Ваш контакт уже был сохранен!", upd.Message.Chat.ID)
				continue
			}
			service.messageService.SendWithContactButton(upd.Message.Chat.ID)
		}
	}
}

func (service *UpdatesService) Get(offset int) ([]models.Update, error) {
	client := locator.Instance.Get("tg_client").(*client.TelegramClient)
	getURL := fmt.Sprintf("getUpdates?offset=%d&timeout=10", offset+1)
	rawData, err := client.Get(getURL)
	if err != nil {
		return nil, err
	}
	var result struct {
		Ok     bool            `json:"ok"`
		Result []models.Update `json:"result"`
	}
	err = json.Unmarshal(rawData, &result)
	if err != nil {
		return nil, err
	}
	return result.Result, nil
}
