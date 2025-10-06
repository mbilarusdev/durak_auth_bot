package endpoint

import (
	"net/http"
	"strings"

	"github.com/mbilarusdev/durak_auth_bot/internal/models"
	"github.com/mbilarusdev/durak_auth_bot/internal/service"
)

type CheckAuthEndpoint struct {
	tokenService service.TokenManager
}

func NewCheckAuthEndpoint(tokenService *service.TokenService) *CheckAuthEndpoint {
	endpoint := new(CheckAuthEndpoint)
	endpoint.tokenService = tokenService
	return endpoint
}

func (endpoint *CheckAuthEndpoint) Call(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Токен не передан"))
		return
	}
	segments := strings.Split(authHeader, " ")
	if len(segments) != 2 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Неверный формат передачи токена"))
		return
	}
	token := segments[1]
	actualToken, err := endpoint.tokenService.FindActualByToken(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка при проверке токена"))
		return
	}
	if actualToken == nil || actualToken.Status != models.TokenAvailable {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Неверный токен или срок его действия истек"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Токен действителен!"))
}
