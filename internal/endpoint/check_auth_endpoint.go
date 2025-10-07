package endpoint

import (
	"net/http"
	"strings"

	"github.com/mbilarusdev/durak_auth_bot/internal/models"
	"github.com/mbilarusdev/durak_auth_bot/internal/service"
	"github.com/mbilarusdev/durak_network/network"
)

type CheckAuthEndpoint struct {
	tokenService service.TokenManager
}

func NewCheckAuthEndpoint(tokenService *service.TokenService) *CheckAuthEndpoint {
	endpoint := new(CheckAuthEndpoint)
	endpoint.tokenService = tokenService
	return endpoint
}

func (endpoint *CheckAuthEndpoint) Call(
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
	return network.SuccessString(w, "Токен верен!")
}
