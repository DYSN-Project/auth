package app

import (
	"context"
	"dysn/auth/config"
	"dysn/auth/internal/repository"
	"dysn/auth/internal/service"
	"dysn/auth/internal/transport/grpc/client"
	"dysn/auth/internal/transport/grpc/server"
	"dysn/auth/pkg/db"
	"dysn/auth/pkg/jwt"
	"dysn/auth/pkg/log"
	"os"
	"os/signal"
	"syscall"
)

func Run(ctx context.Context) {
	cfg := config.NewConfig()

	logger := log.NewLogger()
	jwtService := jwt.NewJwtService()

	database := db.StartDB(cfg, logger)
	defer db.CloseDB(database, logger)

	userRepo := repository.NewUserRepository(database)
	recoveryRepo := repository.NewRecoveryRepository(database)

	notifyCli := client.NewNotify(cfg.GetNotifyAddress(), logger)

	authSrv := service.NewAuth(cfg, jwtService, userRepo, logger)
	registerSrv := service.NewRegister(cfg, jwtService, userRepo, logger, notifyCli)
	recoverySrv := service.NewRecovery(cfg, userRepo, recoveryRepo, notifyCli, logger)

	srv := server.NewGrpc(cfg.GetGrpcPort(),
		authSrv,
		registerSrv,
		recoverySrv,
		logger)
	go srv.StartServer()
	defer srv.StopServer()

	sgn := make(chan os.Signal, 1)
	signal.Notify(sgn, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
	case <-sgn:
	}
}
