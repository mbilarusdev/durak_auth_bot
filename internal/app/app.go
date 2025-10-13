package app

import (
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mbilarusdev/durak_auth_bot/internal/adapter"
	"github.com/mbilarusdev/durak_auth_bot/internal/bot"
	"github.com/mbilarusdev/durak_auth_bot/internal/client"
	"github.com/mbilarusdev/durak_auth_bot/internal/common"
	"github.com/mbilarusdev/durak_auth_bot/internal/endpoint"
	"github.com/mbilarusdev/durak_auth_bot/internal/repository"
	"github.com/mbilarusdev/durak_auth_bot/internal/service"
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
	playerRepository := repository.NewPlayerRepository(adapter.NewAdapterPool(dbPool))
	tokenRepository := repository.NewTokenRepository(adapter.NewAdapterPool(dbPool))
	codeRepository := repository.NewCodeRepository(adapter.NewCacheManagerAdapter(cacheManager))

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
	r := mux.NewRouter()
	swaggerDir := filepath.Join("docs", "swagger-ui")
	r.PathPrefix("/swagger/").
		Handler(http.StripPrefix("/swagger/", http.FileServer(http.Dir(swaggerDir))))
	r.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})
	r.HandleFunc("/code/send", sendCodeEndpoint.Call).Methods(http.MethodPost)
	r.HandleFunc("/code/confirm", confirmCodeEndpoint.Call).Methods(http.MethodPost)
	r.HandleFunc("/login/check", checkAuthEndpoint.Call).Methods(http.MethodPost)
	r.HandleFunc("/logout", logoutEndpoint.Call).Methods(http.MethodPost)

	// Serving
	go bot.StartPolling()
	server := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(server.ListenAndServe())
}
