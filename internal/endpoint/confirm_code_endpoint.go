package endpoint

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mbilarusdev/durak_auth_bot/internal/common"
	"github.com/mbilarusdev/durak_auth_bot/internal/locator"
	"github.com/mbilarusdev/durak_auth_bot/internal/models"
	"github.com/mbilarusdev/durak_auth_bot/internal/repository"
	"github.com/mbilarusdev/durak_auth_bot/internal/service"
	"github.com/mbilarusdev/jwt/jwt"
	jwtmodels "github.com/mbilarusdev/jwt/models"
)

type ConfirmCodeEndpoint struct {
	codeRepository   repository.CodeProvider
	messageService   service.MessageManager
	playerRepository repository.PlayerProvider
	tokenRepository  repository.TokenProvider
}

func NewConfirmCodeEndpoint(
	codeRepository *repository.CodeRepository,
	messageService *service.MessageService,
	playerRepository *repository.PlayerRepository,
	tokenRepository *repository.TokenRepository,
) *ConfirmCodeEndpoint {
	endpoint := &ConfirmCodeEndpoint{}
	endpoint.codeRepository = codeRepository
	endpoint.messageService = messageService
	endpoint.playerRepository = playerRepository
	endpoint.tokenRepository = tokenRepository

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
	player, err := endpoint.playerRepository.FindOne(
		&models.FindOptions{PhoneNumber: confirmCodeRequest.PhoneNumber},
	)
	if err != nil ||
		player == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Контакт с данным номером телефона не найден"))
		return
	}
	code, err := endpoint.codeRepository.GetCode(player.PhoneNumber)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка получения кода"))
		return
	}
	if code != confirmCodeRequest.Code {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Неправильный код!"))
		return
	}
	endpoint.codeRepository.DelCode(confirmCodeRequest.PhoneNumber)
	config := locator.Instance.Get("bot_config").(*common.AuthBotConfig)
	existedToken, err := endpoint.tokenRepository.FindActual(player.ID)
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
	newJwt := jwt.IssueShort(
		&jwtmodels.JwtShortPayload{Iss: "durak", Sub: fmt.Sprint(player.ID)},
		config.SecretKey,
	)
	newToken, err := endpoint.tokenRepository.Insert(
		&models.Token{PlayerID: player.ID, Jwt: newJwt, Status: models.TokenAvailable},
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка создания токена"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %v", newToken.Jwt))
	w.Write([]byte("Вы успешно авторизованы!"))
}
