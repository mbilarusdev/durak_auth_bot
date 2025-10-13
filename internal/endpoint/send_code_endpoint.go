package endpoint

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mbilarusdev/durak_auth_bot/internal/service"
	app_request "github.com/mbilarusdev/durak_auth_bot/internal/structs/app/request"
	net_util "github.com/mbilarusdev/durak_network/net_util"
	net_res "github.com/mbilarusdev/durak_network/response"
	net_err "github.com/mbilarusdev/durak_network/response/err"
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

// @Summary Send code
// @Description Отправляет код в Telegram чат, по номеру телефона
// @Tags Codes
// @Param send_code_request body app_request.SendCodeRequest true "Запрос на отправку кода"
// @Success 201 {object} net_res.Body{code=int,err=nil,data=nil} "Код успешно отправлен"
// @Failure 400 {object} net_res.Body{code=int,err=net_err.ReadBodyErr,data=nil} "Ошибка чтения body"
// @Failure 400 {object} net_res.Body{code=int,err=net_err.ReadBodyErr,data=nil} "Ошибка декодирования body"
// @Failure 404 {object} net_res.Body{code=int,err=net_err.NotFoundErr,data=nil} "Игрок с данным номером телефона не найден"
// @Failure 500 {object} net_res.Body{code=int,err=net_err.ServerErr,data=nil} "Ошибка при создании/отправке токена"
// @Router /code/send [post]
func (endpoint *SendCodeEndpoint) Call(
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
	sendCodeRequest := &app_request.SendCodeRequest{}
	err = json.Unmarshal(bodyBytes, sendCodeRequest)
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
	player, err := endpoint.playerService.FindByPhone(sendCodeRequest.PhoneNumber)
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
	tgCode, err := endpoint.codeService.CreateCode(sendCodeRequest.PhoneNumber)
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
	if err := endpoint.messageService.Send(fmt.Sprintf("Ваш код подтверждения: %v", tgCode), player.ChatID); err != nil {
		code := http.StatusInternalServerError
		net_util.SendResponse(
			w,
			code,
			net_res.NewBody(code, net_err.NewServerErr(), nil),
			map[string]string{},
		)
		return
	}

	code := http.StatusCreated
	net_util.SendResponse(
		w,
		code,
		net_res.NewBody(code, nil, nil),
		map[string]string{},
	)
}
