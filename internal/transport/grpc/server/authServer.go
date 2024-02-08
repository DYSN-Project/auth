package server

import (
	"context"
	"dysn/auth/internal/helper"
	"dysn/auth/internal/model/dto"
	"dysn/auth/internal/model/entity"
	pb "dysn/auth/internal/transport/grpc/pb/auth"
	"dysn/auth/pkg/log"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RegisterInterface interface {
	RegisterUser(ctx context.Context, registerDto *dto.Register) (*entity.User, error)
	ConfirmRegister(ctx context.Context, confirmRegisterDto *dto.Confirm) (*entity.Tokens, error)
}

type AuthInterface interface {
	Login(ctx context.Context, loginDto *dto.Login) (*entity.Tokens, error)
	GetTokensByRefresh(ctx context.Context, refreshToken string) (*entity.Tokens, error)
	VerifyAndGetId(ctx context.Context, accessToken string) (userId uuid.UUID, err error)
	SetLanguage(ctx context.Context, langDto *dto.ChangeLang) error
}

type RecoveryInterface interface {
	CreateRecovery(ctx context.Context, email string) (*entity.RecoveryPassword, error)
	ConfirmRecovery(ctx context.Context, confirmDto *dto.ConfirmRemovePass) error
	ChangePassword(ctx context.Context, passDto *dto.ChangePass) error
}

type AuthServer struct {
	authSrv     AuthInterface
	regSrv      RegisterInterface
	recoverySrv RecoveryInterface
	logger      *log.Logger
	pb.UnimplementedAuthServer
}

func NewAuthServer(authService AuthInterface,
	registerService RegisterInterface,
	recoverySrv RecoveryInterface,
	logger *log.Logger) *AuthServer {
	return &AuthServer{
		authSrv:     authService,
		regSrv:      registerService,
		recoverySrv: recoverySrv,
		logger:      logger,
	}
}

func (a *AuthServer) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.User, error) {
	registerDto := dto.NewRegister(request.GetEmail(),
		request.GetPassword(),
		request.GetLang(),
	)
	if err := registerDto.Validate(); err != nil {
		errVld := helper.MakeGrpcValidationError(err)

		return nil, errVld
	}

	user, err := a.regSrv.RegisterUser(ctx, registerDto)
	if err != nil {

		return nil, err
	}

	return &pb.User{
		Id:        user.Id.String(),
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}, nil
}

func (a *AuthServer) ConfirmRegister(ctx context.Context, request *pb.ConfirmRequest) (*pb.Tokens, error) {
	confirmDto := dto.NewConfirm(
		request.GetEmail(),
		request.GetPassword(),
		request.GetCode())
	if err := confirmDto.Validate(); err != nil {
		errVld := helper.MakeGrpcValidationError(err)

		return nil, errVld
	}

	tokens, err := a.regSrv.ConfirmRegister(ctx, confirmDto)
	if err != nil {
		return nil, err
	}

	return &pb.Tokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (a *AuthServer) Login(ctx context.Context, request *pb.LoginRequest) (*pb.Tokens, error) {
	loginDto := dto.NewLogin(request.GetEmail(), request.GetPassword())
	if err := loginDto.Validate(); err != nil {
		errVld := helper.MakeGrpcValidationError(err)

		return nil, errVld
	}

	tokens, err := a.authSrv.Login(ctx, loginDto)
	if err != nil {
		return nil, err
	}

	return &pb.Tokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (a *AuthServer) UpdateTokens(ctx context.Context, request *pb.Token) (*pb.Tokens, error) {
	updateTokensDto := dto.NewToken(request.GetToken())
	if err := updateTokensDto.Validate(); err != nil {
		errVld := helper.MakeGrpcValidationError(err)

		return nil, errVld
	}

	tokens, err := a.authSrv.GetTokensByRefresh(ctx, updateTokensDto.Token)
	if err != nil {
		return nil, err
	}

	return &pb.Tokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (a *AuthServer) GetUserByToken(ctx context.Context, request *pb.Token) (*pb.User, error) {
	tokenDto := dto.NewToken(request.GetToken())

	if err := tokenDto.Validate(); err != nil {
		errVld := helper.MakeGrpcValidationError(err)

		return nil, errVld
	}

	id, err := a.authSrv.VerifyAndGetId(ctx, tokenDto.Token)
	if err != nil || id == uuid.Nil {
		return nil, err
	}

	return &pb.User{
		Id: id.String(),
	}, nil
}

func (a *AuthServer) RemovePassword(ctx context.Context, pbRemovePass *pb.RemovePasswordRequest) (*emptypb.Empty, error) {
	removePassDto := dto.NewRemovePass(pbRemovePass.GetEmail())
	if err := removePassDto.Validate(); err != nil {
		errVld := helper.MakeGrpcValidationError(err)

		return nil, errVld
	}

	_, err := a.recoverySrv.CreateRecovery(ctx, removePassDto.Email)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (a *AuthServer) RemovePasswordConfirm(ctx context.Context, pbConfirmPass *pb.ConfirmRemovePasswordRequest) (*emptypb.Empty, error) {
	confirmPassDto := dto.NewConfirmRemovePass(pbConfirmPass.GetEmail(), pbConfirmPass.GetCode())
	if err := confirmPassDto.Validate(); err != nil {
		errVld := helper.MakeGrpcValidationError(err)

		return nil, errVld
	}

	err := a.recoverySrv.ConfirmRecovery(ctx, confirmPassDto)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (a *AuthServer) ChangePassword(ctx context.Context, pbChangePass *pb.ChangePasswordRequest) (*emptypb.Empty, error) {
	changePassDto := dto.NewChangePass(pbChangePass.GetEmail(), pbChangePass.GetPassword())
	if err := changePassDto.Validate(); err != nil {
		errVld := helper.MakeGrpcValidationError(err)

		return nil, errVld
	}

	err := a.recoverySrv.ChangePassword(ctx, changePassDto)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (a *AuthServer) SetLanguage(ctx context.Context, request *pb.LanguageRequest) (*emptypb.Empty, error) {
	userId := helper.GetUserIdFromGrpcContext(ctx)
	if userId == uuid.Nil {
		return nil, helper.GetGrpcUnauthenticatedError()
	}

	changeLangDto := dto.NewChangeLang(request.GetLang(), userId)
	if err := changeLangDto.Validate(); err != nil {
		errVld := helper.MakeGrpcValidationError(err)

		return nil, errVld
	}

	err := a.authSrv.SetLanguage(ctx, changeLangDto)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
