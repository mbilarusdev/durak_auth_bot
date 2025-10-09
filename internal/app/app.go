package app

import (
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mbilarusdev/durak_auth_bot/internal/bot"
	"github.com/mbilarusdev/durak_auth_bot/internal/client"
	"github.com/mbilarusdev/durak_auth_bot/internal/common"
	"github.com/mbilarusdev/durak_auth_bot/internal/endpoint"
	"github.com/mbilarusdev/durak_auth_bot/internal/repository"
	"github.com/mbilarusdev/durak_auth_bot/internal/service"
	"github.com/mbilarusdev/durak_network/network"
)

func Run() {
	// Config
	common.NewAuthBotConfig()

	// Clients
	tgClient := client.NewTelegramClient()

	// TODO: change connections
	dbPool := &pgxpool.Pool{}
	cacheManager := &redis.Client{}

	// Repositories
	playerRepository := repository.NewPlayerRepository(dbPool)
	tokenRepository := repository.NewTokenRepository(dbPool)
	codeRepository := repository.NewCodeRepository(cacheManager)

	// Services
	codeService := service.NewCodeService(codeRepository)
	messageService := service.NewMessageService(tgClient)
	playerService := service.NewPlayerService(playerRepository)
	tokenService := service.NewTokenService(tokenRepository)
	updatesHandleService := service.NewUpdatesHandleService(messageService, playerService)
	updatesService := service.NewUpdatesService(tgClient, updatesHandleService)

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
	router.HandleFunc("/code/send", network.Handler(sendCodeEndpoint.Call)).Methods(http.MethodPost)
	router.HandleFunc("/code/confirm", network.Handler(confirmCodeEndpoint.Call)).
		Methods(http.MethodPost)
	router.HandleFunc("/login/check", network.Handler(checkAuthEndpoint.Call)).
		Methods(http.MethodGet)
	router.HandleFunc("/logout", network.Handler(logoutEndpoint.Call)).Methods(http.MethodPost)

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
