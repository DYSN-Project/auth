package server

import (
	"context"
	"dysn/auth/internal/helper"
	"dysn/auth/internal/model/dto"
	"dysn/auth/internal/service"
	pb "dysn/auth/internal/transport/grpc/pb/auth"
	"dysn/auth/pkg/log"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthServer struct {
	authSrv     service.AuthInterface
	regSrv      service.RegisterInterface
	recoverySrv service.RecoveryInterface
	logger      *log.Logger
	pb.UnimplementedAuthServer
}

func NewAuthServer(authService service.AuthInterface,
	registerService service.RegisterInterface,
	recoverySrv service.RecoveryInterface,
	logger *log.Logger) *AuthServer {
	return &AuthServer{
		authSrv:     authService,
		regSrv:      registerService,
		recoverySrv: recoverySrv,
		logger:      logger,
	}
}

func (a *AuthServer) Register(_ context.Context, request *pb.RegisterRequest) (*pb.User, error) {
	registerDto := dto.NewRegister(request.GetEmail(),
		request.GetPassword(),
		request.GetLang(),
	)
	if err := registerDto.Validate(); err != nil {
		errVld := helper.MakeGrpcValidationError(err)

		return nil, errVld
	}

	user, err := a.regSrv.RegisterUser(registerDto.Email,
		registerDto.Password,
		registerDto.Lang,
	)
	if err != nil {

		return nil, err
	}

	return &pb.User{
		Id:        user.Id.String(),
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}, nil
}

func (a *AuthServer) ConfirmRegister(_ context.Context, request *pb.ConfirmRequest) (*pb.Tokens, error) {
	confirmDto := dto.NewConfirm(
		request.GetEmail(),
		request.GetPassword(),
		request.GetCode())
	if err := confirmDto.Validate(); err != nil {
		errVld := helper.MakeGrpcValidationError(err)

		return nil, errVld
	}

	tokens, err := a.regSrv.ConfirmRegister(confirmDto.Email,
		confirmDto.Password,
		confirmDto.Code)
	if err != nil {
		return nil, err
	}

	return &pb.Tokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (a *AuthServer) Login(_ context.Context, request *pb.LoginRequest) (*pb.Tokens, error) {
	login := dto.NewLogin(request.GetEmail(), request.GetPassword())
	if err := login.Validate(); err != nil {
		errVld := helper.MakeGrpcValidationError(err)

		return nil, errVld
	}

	tokens, err := a.authSrv.Login(login.Email, login.Password)
	if err != nil {
		return nil, err
	}

	return &pb.Tokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (a *AuthServer) UpdateTokens(_ context.Context, request *pb.Token) (*pb.Tokens, error) {
	tokenForm := dto.NewToken(request.GetToken())

	if err := tokenForm.Validate(); err != nil {
		errVld := helper.MakeGrpcValidationError(err)

		return nil, errVld
	}

	tokens, err := a.authSrv.GetTokensByRefresh(tokenForm.Token)
	if err != nil {
		return nil, err
	}

	return &pb.Tokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (a *AuthServer) GetUserByToken(_ context.Context, request *pb.Token) (*pb.User, error) {
	tokenForm := dto.NewToken(request.GetToken())

	if err := tokenForm.Validate(); err != nil {
		errVld := helper.MakeGrpcValidationError(err)

		return nil, errVld
	}

	id, err := a.authSrv.VerifyAndGetId(tokenForm.Token)
	if err != nil || id == nil {
		return nil, err
	}

	return &pb.User{
		Id: *id,
	}, nil
}

func (a *AuthServer) RemovePassword(_ context.Context, pbRemovePass *pb.RemovePasswordRequest) (*emptypb.Empty, error) {
	removePassDto := dto.NewRemovePass(pbRemovePass.GetEmail())
	if err := removePassDto.Validate(); err != nil {
		errVld := helper.MakeGrpcValidationError(err)

		return nil, errVld
	}

	_, err := a.recoverySrv.CreateRecovery(removePassDto.Email)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (a *AuthServer) RemovePasswordConfirm(_ context.Context, pbConfirmPass *pb.ConfirmRemovePasswordRequest) (*emptypb.Empty, error) {
	confirmPassDto := dto.NewConfirmRemovePass(pbConfirmPass.GetEmail(),
		pbConfirmPass.GetCode())
	if err := confirmPassDto.Validate(); err != nil {
		errVld := helper.MakeGrpcValidationError(err)

		return nil, errVld
	}

	err := a.recoverySrv.ConfirmRecovery(confirmPassDto.Email,
		confirmPassDto.Code)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (a *AuthServer) ChangePassword(_ context.Context, pbChangePass *pb.ChangePasswordRequest) (*emptypb.Empty, error) {
	changePassDto := dto.NewChangePass(pbChangePass.GetEmail(),
		pbChangePass.GetPassword())
	if err := changePassDto.Validate(); err != nil {
		errVld := helper.MakeGrpcValidationError(err)

		return nil, errVld
	}

	err := a.recoverySrv.ChangePassword(changePassDto.Email,
		changePassDto.Password)
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

	changeLangDto := dto.NewChangeLang(request.GetLang())
	if err := changeLangDto.Validate(); err != nil {
		errVld := helper.MakeGrpcValidationError(err)

		return nil, errVld
	}

	err := a.authSrv.SetLanguage(userId, changeLangDto.Lang)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
