package endpoint

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/mbilarusdev/durak_auth_bot/internal/models"
	"github.com/mbilarusdev/durak_auth_bot/internal/repository"
	"github.com/mbilarusdev/durak_auth_bot/internal/service"
	"github.com/mbilarusdev/durak_auth_bot/internal/utils"
)

type SendCodeEndpoint struct {
	codeRepository   repository.CodeProvider
	messageService   service.MessageManager
	playerRepository repository.PlayerProvider
}

func NewSendCodeEndpoint(
	codeRepository *repository.CodeRepository,
	messageService *service.MessageService,
	playerRepository *repository.PlayerRepository,
) *SendCodeEndpoint {
	endpoint := &SendCodeEndpoint{}
	endpoint.codeRepository = codeRepository
	endpoint.messageService = messageService
	endpoint.playerRepository = playerRepository

	return endpoint
}

func (endpoint *SendCodeEndpoint) Call(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Ошибка при чтении Body"))
		return
	}
	getCodeRequest := &models.GetCodeRequest{}
	err = json.Unmarshal(bodyBytes, getCodeRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Ошибка при десериализации JSON"))
		return
	}
	player, err := endpoint.playerRepository.FindOne(
		&models.FindOptions{PhoneNumber: getCodeRequest.PhoneNumber},
	)
	if err != nil ||
		player == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Контакт с данным номером телефона не найден"))
		return
	}
	smsCode := utils.GenerateRandomCode()
	if err = endpoint.codeRepository.SaveCode(getCodeRequest.PhoneNumber, smsCode); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка сохранения кода"))
		return
	}
	if err = endpoint.messageService.Send(fmt.Sprintf("Ваш код подтверждения: %v", smsCode), player.ChatID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка отправки кода"))
		return
	}
	w.WriteHeader(http.StatusOK)
	log.Printf("Код успешно отправлен для номера телефона %v", getCodeRequest.PhoneNumber)
	w.Write([]byte("Код успешно отправлен"))
}
