package app

import (
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mbilarusdev/durak_auth_bot/internal/adapter"
	"github.com/mbilarusdev/durak_auth_bot/internal/client"
	"github.com/mbilarusdev/durak_auth_bot/internal/common"
	"github.com/mbilarusdev/durak_auth_bot/internal/endpoint"
	grpcendpoint "github.com/mbilarusdev/durak_auth_bot/internal/grpc_endpoint"
	"github.com/mbilarusdev/durak_auth_bot/internal/repository"
	"github.com/mbilarusdev/durak_auth_bot/internal/server"
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

	// Bot
	bot := server.NewAuthBot(updatesService)
	go bot.StartPolling()

	// Grpc Endpoints
	grpcCheckAuthEndpoint := grpcendpoint.NewGrpcCheckAuthEndpoint(tokenService)

	// Grpc-server
	grpcServer := server.NewGrpcServer(grpcCheckAuthEndpoint)
	go grpcServer.ListenAndServe()

	// HTTP-server
	httpServer := server.NewHttpServer(
		sendCodeEndpoint,
		confirmCodeEndpoint,
		checkAuthEndpoint,
		logoutEndpoint,
	)
	httpServer.ListenAndServe()
}
