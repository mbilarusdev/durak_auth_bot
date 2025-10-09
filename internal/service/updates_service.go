package service

import (
	"encoding/json"
	"fmt"

	"github.com/mbilarusdev/durak_auth_bot/internal/client"
	"github.com/mbilarusdev/durak_auth_bot/internal/models"
)

type UpdatesManager interface {
	Process(updates []models.Update)
	Get(offset int) ([]models.Update, error)
}

type UpdatesService struct {
	tgClient         client.ApiClient
	updHandleService UpdatesHandleManager
}

func NewUpdatesService(
	tgClient *client.TelegramClient,
	updHandleService *UpdatesHandleService,
) *UpdatesService {
	service := &UpdatesService{}
	service.tgClient = tgClient
	service.updHandleService = updHandleService

	return service
}

func (service *UpdatesService) Process(updates []models.Update) {
	for _, upd := range updates {
		if upd.Message != nil && upd.Message.Text == "/start" {
			service.updHandleService.HandleStartMsg(upd)
		}
		if upd.Message != nil && upd.Message.Contact != nil {
			service.updHandleService.HandleMsgWithContact(upd)
		}
	}
}

func (service *UpdatesService) Get(offset int) ([]models.Update, error) {
	client := service.tgClient
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
