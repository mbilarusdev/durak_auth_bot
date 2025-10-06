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

var Config common.Config

func Run() {
	locator.Setup()
	locator.Instance.Register("bot_config", common.NewAuthBotConfig())
	locator.Instance.Register("tg_client", client.NewTelegramClient())

	playerRepository := repository.NewPlayerRepository()
	codeRepository := repository.NewCodeRepository()
	tokenRepository := repository.NewTokenRepository()

	messageService := service.NewMessageService()
	updatesService := service.NewUpdatesService(messageService, playerRepository)

	bot := bot.NewAuthBot(updatesService)

	confirmCodeEndpoint := endpoint.NewConfirmCodeEndpoint(
		codeRepository,
		messageService,
		playerRepository,
		tokenRepository,
	)
	sendCodeEndpoint := endpoint.NewSendCodeEndpoint(
		codeRepository,
		messageService,
		playerRepository,
	)
	checkAuthEndpoint := endpoint.NewCheckAuthEndpoint(tokenRepository)
	logoutEndpoint := endpoint.NewLogoutEndpoint(tokenRepository)

	go bot.StartPolling()

	router := mux.NewRouter()
	router.HandleFunc("/code/send", sendCodeEndpoint.Call).Methods(http.MethodPost)
	router.HandleFunc("/code/confirm", confirmCodeEndpoint.Call).Methods(http.MethodPost)
	router.HandleFunc("/auth/check", checkAuthEndpoint.Call).Methods(http.MethodGet)
	router.HandleFunc("/logout", logoutEndpoint.Call).Methods(http.MethodPost)

	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(server.ListenAndServe())
}
