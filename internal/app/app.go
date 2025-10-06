package app

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mbilarusdev/durak_auth_bot/internal/bot"
	"github.com/mbilarusdev/durak_auth_bot/internal/client"
	"github.com/mbilarusdev/durak_auth_bot/internal/common"
	"github.com/mbilarusdev/durak_auth_bot/internal/endpoint"
	"github.com/mbilarusdev/durak_auth_bot/internal/locator"
	"github.com/mbilarusdev/durak_auth_bot/internal/repository"
	"github.com/mbilarusdev/durak_auth_bot/internal/service"
)

func Run() {
	// Locator
	locator.Setup()
	locator.Instance.Register("bot_config", common.NewAuthBotConfig())
	locator.Instance.Register("tg_client", client.NewTelegramClient())

	// Repositories
	playerRepository := repository.NewPlayerRepository()
	codeRepository := repository.NewCodeRepository()
	tokenRepository := repository.NewTokenRepository()

	// Services
	codeService := service.NewCodeService(codeRepository)
	messageService := service.NewMessageService()
	playerService := service.NewPlayerService(playerRepository)
	tokenService := service.NewTokenService(tokenRepository)
	updatesService := service.NewUpdatesService(messageService, playerService)

	// Bot
	bot := bot.NewAuthBot(updatesService)

	// Endpoints
	confirmCodeEndpoint := endpoint.NewConfirmCodeEndpoint(
		codeService,
		messageService,
		playerService,
		tokenService,
	)
	sendCodeEndpoint := endpoint.NewSendCodeEndpoint(
		codeService,
		messageService,
		playerService,
	)
	checkAuthEndpoint := endpoint.NewCheckAuthEndpoint(tokenService)
	logoutEndpoint := endpoint.NewLogoutEndpoint(tokenService)

	// Router
	router := mux.NewRouter()
	router.HandleFunc("/code/send", sendCodeEndpoint.Call).Methods(http.MethodPost)
	router.HandleFunc("/code/confirm", confirmCodeEndpoint.Call).Methods(http.MethodPost)
	router.HandleFunc("/auth/check", checkAuthEndpoint.Call).Methods(http.MethodGet)
	router.HandleFunc("/logout", logoutEndpoint.Call).Methods(http.MethodPost)

	// Serving
	go bot.StartPolling()
	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(server.ListenAndServe())
}
