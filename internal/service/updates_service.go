package service

import (
	"encoding/json"
	"fmt"

	"github.com/mbilarusdev/durak_auth_bot/internal/client"
	tg_model "github.com/mbilarusdev/durak_auth_bot/internal/structs/tg/model"
)

type UpdatesManager interface {
	Process(updates []tg_model.Update)
	Get(offset int) ([]tg_model.Update, error)
}

type UpdatesService struct {
	tgClient         client.TgNetworkManager
	updHandleService UpdatesHandleManager
}

func NewUpdatesService(
	tgClient client.TgNetworkManager,
	updHandleService UpdatesHandleManager,
) *UpdatesService {
	service := &UpdatesService{}
	service.tgClient = tgClient
	service.updHandleService = updHandleService

	return service
}

func (service *UpdatesService) Process(updates []tg_model.Update) {
	for _, upd := range updates {
		if upd.Message != nil && upd.Message.Text == "/start" {
			service.updHandleService.HandleStartMsg(upd)
		}
		if upd.Message != nil && upd.Message.Contact != nil {
			service.updHandleService.HandleMsgWithContact(upd)
		}
	}
}

func (service *UpdatesService) Get(offset int) ([]tg_model.Update, error) {
	client := service.tgClient
	getURL := fmt.Sprintf("getUpdates?offset=%d&timeout=10", offset+1)
	rawData, err := client.Get(getURL)
	if err != nil {
		return nil, err
	}
	var result struct {
		Ok     bool              `json:"ok"`
		Result []tg_model.Update `json:"result"`
	}
	err = json.Unmarshal(rawData, &result)
	if err != nil {
		return nil, err
	}
	return result.Result, nil
}
