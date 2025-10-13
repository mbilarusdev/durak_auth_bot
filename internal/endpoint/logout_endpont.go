package endpoint

import (
	"net/http"
	"strings"

	"github.com/mbilarusdev/durak_auth_bot/internal/service"
	app_model "github.com/mbilarusdev/durak_auth_bot/internal/structs/app/model"
	net_util "github.com/mbilarusdev/durak_network/net_util"
	net_res "github.com/mbilarusdev/durak_network/response"
	net_err "github.com/mbilarusdev/durak_network/response/err"
)

type LogoutEndpoint struct {
	tokenService service.TokenManager
}

func NewLogoutEndpoint(tokenService *service.TokenService) *LogoutEndpoint {
	endpoint := new(LogoutEndpoint)
	endpoint.tokenService = tokenService
	return endpoint
}

// @Summary Logout
// @Description Делает логаут, блокирует токен
// @Tags Auth
// @Param Authorization header string true "Bearer JWT" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9)
// @Success 200 {object} net_res.Body{code=int,err=nil,data=nil} "Вы успешно удалили токен!"
// @Failure 400 {object} net_res.Body{code=int,err=net_err.TokenNotProvidedErr,data=nil} "Токен не передан"
// @Failure 404 {object} net_res.Body{code=int,err=net_err.TokenIncorrectErr,data=nil} "Токен инкорректен"
// @Failure 500 {object} net_res.Body{code=int,err=net_err.ServerErr,data=nil} "Ошибка при проверке/блокировке токена"
// @Router /logout [post]
func (endpoint *LogoutEndpoint) Call(
	w http.ResponseWriter,
	r *http.Request,
) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		code := http.StatusBadRequest
		net_util.SendResponse(
			w,
			code,
			net_res.NewBody(code, net_err.NewTokenNotProvidedErr(), nil),
			map[string]string{},
		)
		return
	}
	segments := strings.Split(authHeader, " ")
	if len(segments) != 2 {
		code := http.StatusUnauthorized
		net_util.SendResponse(
			w,
			code,
			net_res.NewBody(code, net_err.NewTokenIncorrectErr(), nil),
			map[string]string{},
		)
		return
	}
	token := segments[1]
	actualToken, err := endpoint.tokenService.FindActualByToken(token)
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
	if actualToken == nil || actualToken.Status != app_model.TokenAvailable {
		code := http.StatusUnauthorized
		net_util.SendResponse(
			w,
			code,
			net_res.NewBody(code, net_err.NewTokenIncorrectErr(), nil),
			map[string]string{},
		)
		return
	}

	if err := endpoint.tokenService.BlockToken(actualToken.ID); err != nil {
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
		net_res.NewBody(code, nil, nil),
		map[string]string{},
	)
}
