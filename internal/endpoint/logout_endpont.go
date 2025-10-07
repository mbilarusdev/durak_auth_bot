package endpoint

import (
	"net/http"
	"strings"

	"github.com/mbilarusdev/durak_auth_bot/internal/models"
	"github.com/mbilarusdev/durak_auth_bot/internal/service"
	"github.com/mbilarusdev/durak_network/network"
)

type LogoutEndpoint struct {
	tokenService service.TokenManager
}

func NewLogoutEndpoint(tokenService *service.TokenService) *LogoutEndpoint {
	endpoint := new(LogoutEndpoint)
	endpoint.tokenService = tokenService
	return endpoint
}

func (endpoint *LogoutEndpoint) Call(
	w http.ResponseWriter,
	r *http.Request,
) *network.DurakHandlerResult {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return network.TokenNotProvided(w)
	}
	segments := strings.Split(authHeader, " ")
	if len(segments) != 2 {
		return network.TokenIncorrect(w)
	}
	token := segments[1]
	actualToken, err := endpoint.tokenService.FindActualByToken(token)
	if err != nil {
		return network.ServerError(w, "Ошибка проверки информации о токене")
	}
	if actualToken == nil || actualToken.Status != models.TokenAvailable {
		return network.TokenIncorrect(w)
	}

	if err := endpoint.tokenService.BlockToken(actualToken.ID); err != nil {
		return network.ServerError(w, "Ошибка при блокировке токена")
	}
	return network.SuccessString(w, "Вы успешно удалили токен!")
}
