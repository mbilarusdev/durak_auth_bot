package endpoint

import (
	"net/http"
	"strings"

	"github.com/mbilarusdev/durak_auth_bot/internal/models"
	"github.com/mbilarusdev/durak_auth_bot/internal/service"
)

type LogoutEndpoint struct {
	tokenService service.TokenManager
}

func NewLogoutEndpoint(tokenService *service.TokenService) *LogoutEndpoint {
	endpoint := new(LogoutEndpoint)
	endpoint.tokenService = tokenService
	return endpoint
}

func (endpoint *LogoutEndpoint) Call(w http.ResponseWriter, r *http.Request) {
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

	if err := endpoint.tokenService.BlockToken(actualToken.ID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка блокировки токена"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Вы успешно удалили токен"))
}
