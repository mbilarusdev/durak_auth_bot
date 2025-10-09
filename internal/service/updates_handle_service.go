package service

import (
	"log"

	"github.com/mbilarusdev/durak_auth_bot/internal/models"
)

type UpdatesHandleManager interface {
	HandleMsgWithContact(upd models.Update) error
	HandleStartMsg(upd models.Update) error
}

type UpdatesHandleService struct {
	msgService    MessageManager
	playerService PlayerManager
}

func NewUpdatesHandleService(
	msgService *MessageService,
	playerService *PlayerService,
) *UpdatesHandleService {
	updHandleService := new(UpdatesHandleService)
	updHandleService.msgService = msgService
	updHandleService.playerService = playerService

	return updHandleService
}

func (service *UpdatesHandleService) HandleMsgWithContact(upd models.Update) error {
	log.Printf(
		"Контакт получен: Имя=%s, Телефон=%s\n",
		upd.Message.Contact.FirstName,
		upd.Message.Contact.PhoneNumber,
	)
	player, err := service.playerService.FindByPhone(upd.Message.Contact.PhoneNumber)
	if err != nil {
		service.msgService.Send("Произошла ошибка, попробуйте позже", upd.Message.Chat.ID)
		return err
	}
	if player != nil {
		service.msgService.Send("Ваш контакт уже был сохранен!", upd.Message.Chat.ID)
		return err
	}
	if _, err = service.playerService.CreatePlayer(
		upd.Message.Contact.PhoneNumber,
		upd.Message.Chat.ID,
	); err != nil {
		service.msgService.Send("Произошла ошибка, попробуйте позже", upd.Message.Chat.ID)
		return err
	}
	service.msgService.Send("Ваш контакт успешно сохранен!", upd.Message.Chat.ID)
	return nil
}

func (service *UpdatesHandleService) HandleStartMsg(upd models.Update) error {
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
