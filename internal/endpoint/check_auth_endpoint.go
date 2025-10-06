package endpoint

import (
	"net/http"
	"strings"

	"github.com/mbilarusdev/durak_auth_bot/internal/common"
	"github.com/mbilarusdev/durak_auth_bot/internal/locator"
	"github.com/mbilarusdev/durak_auth_bot/internal/models"
	"github.com/mbilarusdev/durak_auth_bot/internal/repository"
	"github.com/mbilarusdev/jwt/jwt"
)

type CheckAuthEndpoint struct {
	tokenRepository repository.TokenProvider
}

func NewCheckAuthEndpoint(tokenRepository *repository.TokenRepository) *CheckAuthEndpoint {
	endpoint := new(CheckAuthEndpoint)
	endpoint.tokenRepository = tokenRepository
	return endpoint
}

func (endpoint *CheckAuthEndpoint) Call(w http.ResponseWriter, r *http.Request) {
	config := locator.Instance.Get("bot_config").(*common.AuthBotConfig)
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
	actualToken, err := endpoint.tokenRepository.FindActual(jwt.GetSubID(token, config.SecretKey))
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
