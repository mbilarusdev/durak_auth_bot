package service

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mbilarusdev/durak_auth_bot/internal/client"
	"github.com/mbilarusdev/durak_auth_bot/internal/locator"
	"github.com/mbilarusdev/durak_auth_bot/internal/models"
)

type UpdatesManager interface {
	Process(updates []models.Update)
	Get(offset int) ([]models.Update, error)
}

type UpdatesService struct {
	messageService MessageManager
	playerService  PlayerManager
}

func NewUpdatesService(
	messageService *MessageService,
	playerService *PlayerService,
) *UpdatesService {
	service := &UpdatesService{}
	service.messageService = messageService
	service.playerService = playerService
	return service
}

func (service *UpdatesService) Process(updates []models.Update) {
	errMsg := "Произошла ошибка, попробуйте позже"
	for _, upd := range updates {
		isMsgWithContact := upd.Message != nil && upd.Message.Contact != nil
		isStartMsg := upd.Message != nil && upd.Message.Text == "/start"
		if isMsgWithContact {
			log.Printf(
				"Контакт получен: Имя=%s, Телефон=%s\n",
				upd.Message.Contact.FirstName,
				upd.Message.Contact.PhoneNumber,
			)
			player, err := service.playerService.FindByPhone(upd.Message.Contact.PhoneNumber)
			if err != nil {
				service.messageService.Send(errMsg, upd.Message.Chat.ID)
				continue
			}
			if player != nil {
				service.messageService.Send("Ваш контакт уже был сохранен!", upd.Message.Chat.ID)
				continue
			}
			if _, err = service.playerService.CreatePlayer(
				upd.Message.Contact.PhoneNumber,
				upd.Message.Chat.ID,
			); err != nil {
				service.messageService.Send(errMsg, upd.Message.Chat.ID)
				continue
			}
			service.messageService.Send("Ваш контакт успешно сохранен!", upd.Message.Chat.ID)
		}
		if isStartMsg {
			player, err := service.playerService.FindByChatID(upd.Message.Chat.ID)
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
