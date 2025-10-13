package endpoint

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mbilarusdev/durak_auth_bot/internal/service"
	app_model "github.com/mbilarusdev/durak_auth_bot/internal/structs/app/model"
	app_request "github.com/mbilarusdev/durak_auth_bot/internal/structs/app/request"
	net_util "github.com/mbilarusdev/durak_network/net_util"
	net_res "github.com/mbilarusdev/durak_network/response"
	net_err "github.com/mbilarusdev/durak_network/response/err"
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

// @Summary Confirm code
// @Description Подтверждает код, который был отправлен в Telegram
// @Tags Codes
// @Param confirm_code_request body app_request.ConfirmCodeRequest true "Запрос на подтверждение кода"
// @Success 201 {object} net_res.Body{code=int,err=nil,data=app_model.Token} "Код успешно подтвержден, токен выпущен"
// @Header 201 {string} Authorization "JWT"
// @Success 200 {object} net_res.Body{code=int,err=nil,data=app_model.Token} "Код уже был подтвержден, токен повторно отправлен"
// @Header 200 {string} Authorization "JWT"
// @Failure 400 {object} net_res.Body{code=int,err=net_err.ReadBodyErr,data=nil} "Ошибка чтения body"
// @Failure 400 {object} net_res.Body{code=int,err=net_err.ReadBodyErr,data=nil} "Ошибка декодирования body"
// @Failure 404 {object} net_res.Body{code=int,err=net_err.NotFoundErr,data=nil} "Игрок с данным номером телефона не найден"
// @Failure 500 {object} net_res.Body{code=int,err=net_err.ServerErr,data=nil} "Ошибка при проверке/выпуске токена"
// @Router /code/confirm [post]
func (endpoint *ConfirmCodeEndpoint) Call(
	w http.ResponseWriter,
	r *http.Request,
) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		code := http.StatusBadRequest
		net_util.SendResponse(
			w,
			code,
			net_res.NewBody(code, net_err.NewReadBodyErr(), nil),
			map[string]string{},
		)
		return
	}
	confirmCodeRequest := &app_request.ConfirmCodeRequest{}
	err = json.Unmarshal(bodyBytes, confirmCodeRequest)
	if err != nil {
		code := http.StatusBadRequest
		net_util.SendResponse(
			w,
			code,
			net_res.NewBody(code, net_err.NewUnmarshalingErr(), nil),
			map[string]string{},
		)
		return
	}
	player, err := endpoint.playerService.FindByPhone(confirmCodeRequest.PhoneNumber)
	if err != nil ||
		player == nil {
		code := http.StatusNotFound
		net_util.SendResponse(
			w,
			code,
			net_res.NewBody(code, net_err.NewNotFoundErr(), nil),
			map[string]string{},
		)
	}
	isRightCode, err := endpoint.codeService.ConsumeIsRightCode(
		player.PhoneNumber,
		confirmCodeRequest.Code,
	)
	if err != nil {
		code := http.StatusInternalServerError
		net_util.SendResponse(
			w,
			code,
			net_res.NewBody(code, net_err.NewServerErr(), nil),
			map[string]string{},
		)
		return

	}
	if !isRightCode {
		code := http.StatusBadRequest
		net_util.SendResponse(
			w,
			code,
			net_res.NewBody(code, net_err.NewWrongInfoErr(), nil),
			map[string]string{},
		)
		return
	}
	existedToken, err := endpoint.tokenService.FindActualByPlayerID(player.ID)
	if err != nil {
		code := http.StatusInternalServerError
		net_util.SendResponse(
			w,
			code,
			net_res.NewBody(code, net_err.NewServerErr(), nil),
			map[string]string{},
		)
		return
	}
	if existedToken != nil && existedToken.Status == app_model.TokenAvailable {
		code := http.StatusOK
		net_util.SendResponse(
			w,
			code,
			net_res.NewBody(code, nil, existedToken),
			map[string]string{"Authorization": fmt.Sprintf("Bearer %v", existedToken.Jwt)},
		)
		return
	}
	newToken, err := endpoint.tokenService.IssueToken(player.ID)
	if err != nil {
		code := http.StatusInternalServerError
		net_util.SendResponse(
			w,
			code,
			net_res.NewBody(code, net_err.NewServerErr(), nil),
			map[string]string{},
		)
		return
	}
	code := http.StatusOK
	net_util.SendResponse(
		w,
		code,
		net_res.NewBody(code, nil, newToken),
		map[string]string{"Authorization": fmt.Sprintf("Bearer %v", newToken.Jwt)},
	)
}
