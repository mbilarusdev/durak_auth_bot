package endpoint

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/mbilarusdev/durak_auth_bot/internal/models"
	"github.com/mbilarusdev/durak_auth_bot/internal/service"
	"github.com/mbilarusdev/durak_network/network"
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

func (endpoint *ConfirmCodeEndpoint) Call(
	w http.ResponseWriter,
	r *http.Request,
) *network.Result {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return network.ReadBodyError(w)
	}
	confirmCodeRequest := &models.ConfirmCodeRequest{}
	err = json.Unmarshal(bodyBytes, confirmCodeRequest)
	if err != nil {
		return network.UnmarshalingError(w)
	}
	player, err := endpoint.playerService.FindByPhone(confirmCodeRequest.PhoneNumber)
	if err != nil ||
		player == nil {
		return network.NotFound(w, "Игрок с данным номером телефона не найден")
	}
	isRightCode, err := endpoint.codeService.ConsumeIsRightCode(
		player.PhoneNumber,
		confirmCodeRequest.Code,
	)
	if err != nil {
		return network.ServerError(w, "Ошибка при проверке кода")

	}
	if !isRightCode {
		return network.WrongInfo(w, "Неправильный код")
	}
	existedToken, err := endpoint.tokenService.FindActualByPlayerID(player.ID)
	if err != nil {
		return network.ServerError(w, "Ошибка выпуска токена")
	}
	if existedToken != nil && existedToken.Status == models.TokenAvailable {
		return network.AlreadyExistJWT(w, existedToken.Jwt)
	}
	newToken, err := endpoint.tokenService.IssueToken(player.ID)
	if err != nil {
		return network.ServerError(w, "Ошибка выпуска токена")
	}
	return network.SuccessJWT(w, newToken.Jwt)
}
