package service

import (
	"log"
	"strings"

	tg_model "github.com/mbilarusdev/durak_auth_bot/internal/structs/tg/model"
)

type UpdatesHandleManager interface {
	HandleMsgWithContact(upd tg_model.Update) error
	HandleStartMsg(upd tg_model.Update) error
}

type UpdatesHandleService struct {
	msgService    MessageManager
	playerService PlayerManager
}

func NewUpdatesHandleService(
	msgService MessageManager,
	playerService PlayerManager,
) *UpdatesHandleService {
	updHandleService := new(UpdatesHandleService)
	updHandleService.msgService = msgService
	updHandleService.playerService = playerService

	return updHandleService
}

func (service *UpdatesHandleService) HandleMsgWithContact(upd tg_model.Update) error {
	purePhone := strings.ReplaceAll(upd.Message.Contact.PhoneNumber, " ", "")
	log.Printf(
		"Контакт получен: Имя=%s, Телефон=%s\n",
		upd.Message.Contact.FirstName,
		purePhone,
	)
	player, err := service.playerService.FindByPhone(purePhone)
	if err != nil {
		service.msgService.Send("Произошла ошибка, попробуйте позже", upd.Message.Chat.ID)
		return err
	}
	if player != nil {
		service.msgService.Send("Ваш контакт уже был сохранен!", upd.Message.Chat.ID)
		return err
	}
	if _, err = service.playerService.CreatePlayer(
		purePhone,
		upd.Message.Chat.ID,
	); err != nil {
		service.msgService.Send("Произошла ошибка, попробуйте позже", upd.Message.Chat.ID)
		return err
	}
	service.msgService.Send("Ваш контакт успешно сохранен!", upd.Message.Chat.ID)
	return nil
}

func (service *UpdatesHandleService) HandleStartMsg(upd tg_model.Update) error {
	player, err := service.playerService.FindByChatID(upd.Message.Chat.ID)
	if err != nil {
		service.msgService.Send("Произошла ошибка, попробуйте позже", upd.Message.Chat.ID)
		return err
	}
	if player == nil {
		service.msgService.Send("Ваш контакт уже был сохранен!", upd.Message.Chat.ID)
		return err
	}
	service.msgService.SendWithContactButton(upd.Message.Chat.ID)
	return nil
}
