package server

import (
	"log"

	"github.com/mbilarusdev/durak_auth_bot/internal/service"
)

type Bot interface {
	StartPolling()
}

type AuthBot struct {
	updatesService service.UpdatesManager
}

func NewAuthBot(updatesService *service.UpdatesService) *AuthBot {
	bot := &AuthBot{}
	bot.updatesService = updatesService

	return bot
}

func (bot *AuthBot) StartPolling() {
	updatesService := bot.updatesService

	offset := 0
	log.Printf("Starting TG Bot pooling")
	for {
		updates, err := updatesService.Get(offset)
		if err != nil {
			log.Println("Ошибка получения обновлений:", err)
			continue
		}
		updatesService.Process(updates)
		if len(updates) > 0 {
			offset = updates[len(updates)-1].UpdateID
		}
	}
}
