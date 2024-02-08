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
	"github.com/segmentio/kafka-go"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(ctx context.Context) {
	cfg := config.NewConfig()

	logger := log.NewLogger()
	jwtService := jwt.NewJwtService()

	userKafkaProducer := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.GetKafkaBroker1()),
		Balancer:               &kafka.LeastBytes{},
		WriteTimeout:           30 * time.Second,
		ReadTimeout:            30 * time.Second,
		Async:                  false,
		AllowAutoTopicCreation: true,
		/*	Transport: &kafka.Transport{
			SASL: plain.Mechanism{
				Username: cfg.GetAnalyticsKafkaUsername(),
				Password: cfg.GetAnalyticsKafkaPassword(),
			},
		},*/
	}
	defer userKafkaProducer.Close()

	database := db.StartDB(cfg, logger)
	defer db.CloseDB(database, logger)

	userRepo := repository.NewUserRepository(database)
	recoveryRepo := repository.NewRecoveryRepository(database)

	notifyCli := client.NewNotify(cfg.GetNotifyAddress(), logger)

	baseService := service.NewService(cfg, jwtService, logger)
	authSrv := service.NewAuth(userRepo, baseService)
	registerSrv := service.NewRegister(userRepo, baseService, userKafkaProducer)
	recoverySrv := service.NewRecovery(userRepo, recoveryRepo, notifyCli, baseService)

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
