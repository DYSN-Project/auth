package grpc

import (
	"github.com/DYSN-Project/auth/internal/service"
	"github.com/DYSN-Project/auth/internal/transport/grpc/pb"
	"github.com/DYSN-Project/auth/pkg/log"
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

func NewGrpcServer(port string,
	authSrv service.AuthorizationInterface,
	registerSrv service.RegistrationInterface,
	logger *log.Logger) *GrpcServer {
	srv := grpc.NewServer()
	authServer := NewAuthServer(authSrv, registerSrv, logger)
	pb.RegisterAuthServer(srv, authServer)

	return &GrpcServer{server: srv, port: port, logger: logger}
}

func (g *GrpcServer) StartServer() {
	g.logger.InfoLog.Println("Auth transport starting...")
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
	g.logger.InfoLog.Println("Auth transport stopping...")
	g.server.Stop()
}
