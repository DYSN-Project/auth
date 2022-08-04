package grpc

import (
	"context"
	"github.com/DYSN-Project/auth/internal/models/forms"
	"github.com/DYSN-Project/auth/internal/service"
	"github.com/DYSN-Project/auth/internal/transport/grpc/pb"
	"github.com/DYSN-Project/auth/pkg/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthServer struct {
	useCaseManager service.UseCaseInterface
	logger         *log.Logger
	pb.UnimplementedAuthServer
}

func NewAuthServer(useCaseManager service.UseCaseInterface, logger *log.Logger) *AuthServer {
	return &AuthServer{
		useCaseManager: useCaseManager,
		logger:         logger,
	}
}

func (a *AuthServer) Register(_ context.Context, request *pb.RegisterRequest) (*pb.Token, error) {
	registerForm := forms.NewRegisterForm(request.GetEmail(), request.GetPassword())
	if err := registerForm.Validate(); err != nil {
		return nil, err
	}

	token, err := a.useCaseManager.RegisterUser(registerForm.Email, registerForm.Password)
	if err != nil {
		return nil, err
	}

	return &pb.Token{
		Token: token,
	}, nil
}

func (a *AuthServer) ConfirmRegister(_ context.Context, request *pb.Token) (*pb.User, error) {
	tokenForm := forms.NewTokenForm(request.GetToken())
	if err := tokenForm.Validate(); err != nil {
		return nil, err
	}

	user, err := a.useCaseManager.ConfirmRegister(tokenForm.Token)
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
	loginForm := forms.NewLoginForm(request.GetEmail(), request.GetPassword())
	if err := loginForm.Validate(); err != nil {
		return nil, err
	}

	tokens, err := a.useCaseManager.Login(loginForm.Email, loginForm.Password)
	if err != nil {
		return nil, err
	}

	return &pb.Tokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (a *AuthServer) UpdateTokens(_ context.Context, request *pb.Token) (*pb.Tokens, error) {
	tokenForm := forms.NewTokenForm(request.GetToken())
	if err := tokenForm.Validate(); err != nil {
		return nil, err
	}

	tokens, err := a.useCaseManager.GetTokensByRefresh(tokenForm.Token)
	if err != nil {
		return nil, err
	}

	return &pb.Tokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (a *AuthServer) Verify(_ context.Context, request *pb.Token) (*emptypb.Empty, error) {
	tokenForm := forms.NewTokenForm(request.GetToken())
	if err := tokenForm.Validate(); err != nil {
		return nil, err
	}

	if err := a.useCaseManager.Verify(tokenForm.Token); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
