package endpoint

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mbilarusdev/durak_auth_bot/internal/models"
	"github.com/mbilarusdev/durak_auth_bot/internal/service"
)

type ConfirmCodeEndpoint struct {
	codeService    service.CodeManager
	messageService service.MessageManager
	playerService  service.PlayerManager
	tokenService   service.TokenManager
}

func NewConfirmCodeEndpoint(
	codeService *service.CodeService,
	messageService *service.MessageService,
	playerService *service.PlayerService,
	tokenService *service.TokenService,
) *ConfirmCodeEndpoint {
	endpoint := &ConfirmCodeEndpoint{}
	endpoint.codeService = codeService
	endpoint.messageService = messageService
	endpoint.playerService = playerService
	endpoint.tokenService = tokenService

	return endpoint
}

func (endpoint *ConfirmCodeEndpoint) Call(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Ошибка при чтении Body"))
		return
	}
	confirmCodeRequest := &models.ConfirmCodeRequest{}
	err = json.Unmarshal(bodyBytes, confirmCodeRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Ошибка при десериализации JSON"))
		return
	}
	player, err := endpoint.playerService.FindByPhone(confirmCodeRequest.PhoneNumber)
	if err != nil ||
		player == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Контакт с данным номером телефона не найден"))
		return
	}
	isRightCode, err := endpoint.codeService.ConsumeIsRightCode(
		player.PhoneNumber,
		confirmCodeRequest.Code,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка получения кода"))
		return
	}
	if !isRightCode {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Неправильный код!"))
		return
	}
	existedToken, err := endpoint.tokenService.FindActualByPlayerID(player.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка создания токена"))
		return
	}
	if existedToken != nil && existedToken.Status == models.TokenAvailable {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Authorization", fmt.Sprintf("Bearer %v", existedToken.Jwt))
		w.Write([]byte("У вас уже есть токен!"))
		return
	}
	newToken, err := endpoint.tokenService.IssueToken(player.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка создания токена"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %v", newToken.Jwt))
	w.Write([]byte("Вы успешно авторизованы!"))
}
