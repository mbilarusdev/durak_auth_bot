package endpoint

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/mbilarusdev/durak_auth_bot/internal/models"
	"github.com/mbilarusdev/durak_auth_bot/internal/service"
)

type SendCodeEndpoint struct {
	codeService    service.CodeManager
	messageService service.MessageManager
	playerService  service.PlayerManager
}

func NewSendCodeEndpoint(
	codeService *service.CodeService,
	messageService *service.MessageService,
	playerService *service.PlayerService,
) *SendCodeEndpoint {
	endpoint := &SendCodeEndpoint{}
	endpoint.codeService = codeService
	endpoint.messageService = messageService
	endpoint.playerService = playerService

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
	player, err := endpoint.playerService.FindByPhone(getCodeRequest.PhoneNumber)
	if err != nil ||
		player == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Контакт с данным номером телефона не найден"))
		return
	}
	code, err := endpoint.codeService.CreateCode(getCodeRequest.PhoneNumber)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка сохранения кода"))
		return
	}
	if err = endpoint.messageService.Send(fmt.Sprintf("Ваш код подтверждения: %v", code), player.ChatID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка отправки кода"))
		return
	}
	w.WriteHeader(http.StatusOK)
	log.Printf("Код успешно отправлен для номера телефона %v", getCodeRequest.PhoneNumber)
	w.Write([]byte("Код успешно отправлен"))
}
