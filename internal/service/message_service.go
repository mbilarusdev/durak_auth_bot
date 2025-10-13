package service

import (
	"encoding/json"
	"log"

	"github.com/mbilarusdev/durak_auth_bot/internal/client"
	tg_model "github.com/mbilarusdev/durak_auth_bot/internal/structs/tg/model"
	tg_request "github.com/mbilarusdev/durak_auth_bot/internal/structs/tg/request"
)

type MessageManager interface {
	Send(message string, chatID int64) error
	SendWithContactButton(chatID int64) error
}

type MessageService struct {
	tgClient client.TgNetworkManager
}

func NewMessageService(tgClient client.TgNetworkManager) *MessageService {
	service := &MessageService{}
	service.tgClient = tgClient

	return service
}

func (service *MessageService) Send(message string, chatID int64) error {
	client := service.tgClient
	msgReq := tg_request.SendMessageRequest{
		ChatID: chatID,
		Text:   message,
	}
	data, err := json.Marshal(msgReq)
	if err != nil {
		log.Println("Ошибка маршалинга при отправке сообщения:", err)
		return err
	}
	sendResp, err := client.Post("sendMessage", data)
	if err != nil || !sendResp.Ok {
		log.Println("Ошибка отправки сообщения:", err, sendResp.ErrorDescription)
		return err
	}

	return nil
}

func (service *MessageService) SendWithContactButton(chatID int64) error {
	client := service.tgClient
	msgReq := tg_request.SendMessageRequest{
		ChatID: chatID,
		Text:   "Чтобы получить временный код, отправьте свой контакт",
		ReplyMarkup: tg_model.ReplyKeyboardMarkup{
			Keyboard: [][]tg_model.KeyboardButton{
				{
					{Text: "Поделиться контактом", RequestContact: true},
				},
			},
			ResizeKeyboard:  true,
			OneTimeKeyboard: true,
		},
	}
	data, err := json.Marshal(msgReq)
	if err != nil {
		log.Println("Ошибка маршалинга при отправке сообщения:", err)
		return err
	}
	sendResp, err := client.Post("sendMessage", data)
	if err != nil || !sendResp.Ok {
		log.Println("Ошибка отправки сообщения:", err, sendResp.ErrorDescription)
		return err
	}
	return nil
}
