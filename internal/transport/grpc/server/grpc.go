package server

import (
	"dysn/auth/internal/service"
	pb "dysn/auth/internal/transport/grpc/pb/auth"
	"dysn/auth/pkg/log"
	"google.golang.org/grpc"
	"net"
)

type TransportInterface interface {
	StartServer()
	StopServer()
}

type Grpc struct {
	server *grpc.Server
	port   string
	logger *log.Logger
}

func NewGrpc(port string,
	authSrv service.AuthInterface,
	regSrv service.RegisterInterface,
	recoverySrv service.RecoveryInterface,
	logger *log.Logger) *Grpc {
	srv := grpc.NewServer()
	auth := NewAuthServer(authSrv, regSrv, recoverySrv, logger)
	pb.RegisterAuthServer(srv, auth)

	return &Grpc{
		server: srv,
		port:   port,
		logger: logger,
	}
}

func (g *Grpc) StartServer() {
	g.logger.InfoLog.Println("Server transport starting...")

	connection, err := net.Listen("tcp", g.port)
	if err != nil {
		g.logger.ErrorLog.Panic(err)
	}

	err = g.server.Serve(connection)
	if err != nil {
		g.logger.ErrorLog.Panic(err)
	}
}

func (g *Grpc) StopServer() {
	g.logger.InfoLog.Println("Server transport stopping...")
	g.server.Stop()
}
