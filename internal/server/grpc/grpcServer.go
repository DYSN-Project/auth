package grpc

import (
	"github.com/DYSN-Project/auth/internal/packages/log"
	"github.com/DYSN-Project/auth/internal/server/grpc/pb"
	"github.com/DYSN-Project/auth/internal/usecases"
	"google.golang.org/grpc"
	"net"
)

type GrpcServerInterface interface {
	StartServer()
	StopServer()
}

type GrpcServer struct {
	server *grpc.Server
	port   string
	logger *log.Logger
}

func NewGrpcServer(port string, useCaseManager usecases.UseCaseInterface, logger *log.Logger) *GrpcServer {
	srv := grpc.NewServer()
	authSrv := NewAuthServer(useCaseManager, logger)
	pb.RegisterAuthServer(srv, authSrv)

	return &GrpcServer{server: srv, port: port, logger: logger}
}

func (g *GrpcServer) StartServer() {
	g.logger.InfoLog.Println("Auth server starting...")
	l, err := net.Listen("tcp", g.port)
	if err != nil {
		g.logger.ErrorLog.Panic(err)
	}
	err = g.server.Serve(l)
	if err != nil {
		g.logger.ErrorLog.Panic(err)
	}
}

func (g *GrpcServer) StopServer() {
	g.logger.InfoLog.Println("Auth server stopping...")
	g.server.Stop()
}
