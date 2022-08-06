package grpc

import (
	"context"
	"github.com/DYSN-Project/auth/internal/model/dto"
	"github.com/DYSN-Project/auth/internal/service"
	"github.com/DYSN-Project/auth/internal/transport/grpc/pb"
	"github.com/DYSN-Project/auth/pkg/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthServer struct {
	registerSrv service.RegistrationInterface
	authSrv     service.AuthorizationInterface
	logger      *log.Logger
	pb.UnimplementedAuthServer
}

func NewAuthServer(authSrv service.AuthorizationInterface,
	registerSrv service.RegistrationInterface,
	logger *log.Logger) *AuthServer {
	return &AuthServer{
		authSrv:     authSrv,
		registerSrv: registerSrv,
		logger:      logger,
	}
}

func (a *AuthServer) Register(_ context.Context, request *pb.RegisterRequest) (*pb.Token, error) {
	registerDto := dto.NewRegister(request.GetEmail(), request.GetPassword())
	if err := registerDto.Validate(); err != nil {
		return nil, err
	}

	token, err := a.registerSrv.Register(registerDto.Email, registerDto.Password)
	if err != nil {
		return nil, err
	}

	return &pb.Token{
		Token: token,
	}, nil
}

func (a *AuthServer) ConfirmRegister(_ context.Context, request *pb.Token) (*pb.User, error) {
	tokenDto := dto.NewToken(request.GetToken())
	if err := tokenDto.Validate(); err != nil {
		return nil, err
	}

	user, err := a.registerSrv.ConfirmRegister(tokenDto.Token)
	if err != nil {
		return nil, err
	}

	return &pb.User{
		Email:     user.Email,
		Id:        user.ID.String(),
		CreatedAt: timestamppb.New(user.CreatedAt),
		Status:    int32(user.Status),
	}, nil
}

func (a *AuthServer) Login(_ context.Context, request *pb.LoginRequest) (*pb.Tokens, error) {
	loginDto := dto.NewLogin(request.GetEmail(), request.GetPassword())
	if err := loginDto.Validate(); err != nil {
		return nil, err
	}

	tokens, err := a.authSrv.Login(loginDto.Email, loginDto.Password)
	if err != nil {
		return nil, err
	}

	return &pb.Tokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (a *AuthServer) RefreshTokens(_ context.Context, request *pb.Token) (*pb.Tokens, error) {
	tokenDto := dto.NewToken(request.GetToken())
	if err := tokenDto.Validate(); err != nil {
		return nil, err
	}

	tokens, err := a.authSrv.GetTokensByRefresh(tokenDto.Token)
	if err != nil {
		return nil, err
	}

	return &pb.Tokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (a *AuthServer) VerifyTokenAndGetId(_ context.Context, request *pb.Token) (*pb.Identity, error) {
	tokenDto := dto.NewToken(request.GetToken())
	if err := tokenDto.Validate(); err != nil {
		return nil, err
	}

	userId, err := a.authSrv.VerifyAndGetId(tokenDto.Token)
	if err != nil {
		return nil, err
	}

	return &pb.Identity{
		UserId: userId.String(),
	}, nil
}
