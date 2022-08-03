package app

import (
	"context"
	"github.com/DYSN-Project/auth/config"
	"github.com/DYSN-Project/auth/internal/packages/db"
	"github.com/DYSN-Project/auth/internal/packages/jwt"
	"github.com/DYSN-Project/auth/internal/packages/log"
	"github.com/DYSN-Project/auth/internal/repository"
	"github.com/DYSN-Project/auth/internal/server/grpc"
	"github.com/DYSN-Project/auth/internal/usecases"
	"os"
	"os/signal"
	"syscall"
)

func Run(ctx context.Context) {
	cfg := config.NewConfig()

	logger := log.NewLogger()
	database := db.StartDB(cfg, logger)
	defer db.CloseDB(database, logger)

	jwtService := jwt.NewJwtService()

	userRepo := repository.NewUserRepository(database)
	useCaseManager := usecases.NewUseCase(cfg, jwtService, userRepo, logger)

	srv := grpc.NewGrpcServer(cfg.GetGrpcPort(), useCaseManager, logger)
	go srv.StartServer()
	defer srv.StopServer()

	sgn := make(chan os.Signal, 1)
	signal.Notify(sgn, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
	case <-sgn:
	}
}
