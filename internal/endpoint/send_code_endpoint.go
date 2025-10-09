package endpoint

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mbilarusdev/durak_auth_bot/internal/models"
	"github.com/mbilarusdev/durak_auth_bot/internal/service"
	"github.com/mbilarusdev/durak_network/network"
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

func (endpoint *SendCodeEndpoint) Call(
	w http.ResponseWriter,
	r *http.Request,
) *network.Result {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return network.ReadBodyError(w)
	}
	getCodeRequest := &models.GetCodeRequest{}
	err = json.Unmarshal(bodyBytes, getCodeRequest)
	if err != nil {
		return network.UnmarshalingError(w)
	}
	player, err := endpoint.playerService.FindByPhone(getCodeRequest.PhoneNumber)
	if err != nil ||
		player == nil {
		return network.NotFound(w, "Игрок с данным номером телефона не найден")
	}
	code, err := endpoint.codeService.CreateCode(getCodeRequest.PhoneNumber)
	if err != nil {
		return network.ServerError(w, "Ошибка создания кода")
	}
	if err = endpoint.messageService.Send(fmt.Sprintf("Ваш код подтверждения: %v", code), player.ChatID); err != nil {
		return network.ServerError(w, "Ошибка отправки кода")

	}
	return network.SuccessString(w, "Вы успешно отправили код!")
}
