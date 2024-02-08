package service

import (
	"context"
	"dysn/auth/internal/helper"
	"dysn/auth/internal/model/dto"
	"dysn/auth/internal/model/entity"
	"github.com/google/uuid"
)

type UserRepoInterface interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserById(ctx context.Context, userId uuid.UUID) (*entity.User, error)
	CreateUser(ctx context.Context, userId *entity.User) error
	ConfirmUser(ctx context.Context, userId uuid.UUID) error
	Add2FaCode(ctx context.Context, userId uuid.UUID, code string) error
	Remove2FaCode(ctx context.Context, userId uuid.UUID) error
	SetConfirmCode(ctx context.Context, userId uuid.UUID, code string) error
	Confirm2FaCode(ctx context.Context, userId uuid.UUID) error
	ChangePasswordByEmail(ctx context.Context, email, password string) error
	UpdateLang(ctx context.Context, userId uuid.UUID, lang string) error
	ExistUserByEmail(ctx context.Context, email string) bool
}

type Auth struct {
	userRepo UserRepoInterface
	*baseService
}

func NewAuth(userRepo UserRepoInterface, service *baseService) *Auth {
	return &Auth{userRepo, service}
}

func (a *Auth) Login(ctx context.Context, loginDto *dto.Login) (*entity.Tokens, error) {
	user, err := a.userRepo.GetUserByEmail(ctx, loginDto.Email)
	if err != nil {
		a.logger.ErrorLog.Println("err get user by email: ", err)

		return nil, err
	}

	if err := a.checkUser(user); err != nil {
		return nil, err
	}

	comparedPassword := helper.CompareHash(loginDto.Password, user.Password, a.cfg.GetPwdSalt())
	if comparedPassword != nil {
		return nil, helper.MakeGrpcBadRequestError(errInvalidUserData)
	}

	//TODO:: authHistory
	return a.getTokens(user.Id, user.Lang)
}

func (a *Auth) GetTokensByRefresh(ctx context.Context, refreshToken string) (*entity.Tokens, error) {
	tokenClaims, err := a.jwtService.ParseToken(refreshToken, a.cfg.GetJwtRefreshSecretKey())
	if err != nil {
		a.logger.ErrorLog.Println("err parse token: ", err)

		return nil, helper.MakeGrpcBadRequestError(errTokenInvalid)
	}

	if _, ok := tokenClaims[jwtUserIdKey]; !ok {
		return nil, helper.MakeGrpcBadRequestError(errTokenInvalid)
	}

	userId, err := uuid.Parse(tokenClaims[jwtUserIdKey].(string))
	if err != nil {
		return nil, helper.MakeGrpcBadRequestError(errTokenInvalid)
	}

	user, err := a.userRepo.GetUserById(ctx, userId)
	if err != nil {
		a.logger.ErrorLog.Println("err get user by id: ", err)

		return nil, err
	}

	if err = a.checkUser(user); err != nil {
		return nil, err
	}

	return a.getTokens(user.Id, user.Lang)
}

func (a *Auth) VerifyAndGetId(ctx context.Context, accessToken string) (userId uuid.UUID, err error) {
	tokenClaims, err := a.jwtService.ParseToken(accessToken, a.cfg.GetJwtAccessSecretKey())
	if err != nil || len(tokenClaims) == 0 {
		a.logger.ErrorLog.Println("err verify access token : ", err)
		err = helper.MakeGrpcBadRequestError(errTokenInvalid)

		return
	}

	id, ok := tokenClaims[jwtUserIdKey]
	if !ok || id == "" {
		err = helper.MakeGrpcBadRequestError(errTokenInvalid)

		return
	}
	userId, err = uuid.Parse(id.(string))
	if err != nil {
		err = helper.MakeGrpcBadRequestError(errTokenInvalid)

		return
	}

	return
}

func (a *Auth) SetLanguage(ctx context.Context, langDto *dto.ChangeLang) error {
	user, err := a.userRepo.GetUserById(ctx, langDto.UserId)
	if err != nil {
		a.logger.ErrorLog.Println("err get user by id: ", err)

		return err
	}
	if err := a.checkUser(user); err != nil {
		return err
	}
	if err := a.userRepo.UpdateLang(ctx, langDto.UserId, langDto.Lang); err != nil {
		a.logger.ErrorLog.Println("err set lang: ", err)

		return err
	}

	return nil
}
