package service

import (
	"dysn/auth/config"
	"dysn/auth/internal/helper"
	"dysn/auth/internal/model/entity"
	"dysn/auth/internal/repository"
	"dysn/auth/pkg/jwt"
	"dysn/auth/pkg/log"
	"github.com/google/uuid"
)

type AuthInterface interface {
	Login(email, password string) (*entity.Tokens, error)
	GetTokensByRefresh(refreshToken string) (*entity.Tokens, error)
	VerifyAndGetId(accessToken string) (*string, error)
	SetLanguage(userId uuid.UUID, lang string) error
}

type Auth struct {
	cfg        *config.Config
	jwtService jwt.JwtInterface
	userRepo   repository.UserRepoInterface
	logger     *log.Logger
	BaseService
}

func NewAuth(cfg *config.Config,
	jwtSrv jwt.JwtInterface,
	usrRepo repository.UserRepoInterface,
	logger *log.Logger) *Auth {
	return &Auth{
		cfg:        cfg,
		jwtService: jwtSrv,
		userRepo:   usrRepo,
		logger:     logger,
	}
}

func (a *Auth) Login(email, password string) (*entity.Tokens, error) {
	user := a.userRepo.GetUserByEmail(email)
	if err := a.CheckUser(user); err != nil {
		return nil, err
	}

	comparedPassword := helper.CompareHash(password, user.Password, a.cfg.GetPwdSalt())
	if comparedPassword != nil {
		return nil, helper.MakeGrpcBadRequestError(errInvalidUserData)
	}

	return a.getTokens(user.Id, user.Lang)
}

func (a *Auth) GetTokensByRefresh(refreshToken string) (*entity.Tokens, error) {
	tokenClaims, err := a.jwtService.ParseToken(refreshToken, a.cfg.GetJwtRefreshSecretKey())
	if err != nil {
		a.logger.ErrorLog.Println("parse token error: ", err)

		return nil, helper.MakeGrpcBadRequestError(errTokenInvalid)
	}

	if _, ok := tokenClaims["uid"]; !ok {
		return nil, helper.MakeGrpcBadRequestError(errTokenInvalid)
	}
	userId, _ := uuid.Parse(tokenClaims["uid"].(string))

	user := a.userRepo.GetUserById(userId)
	if err = a.CheckUser(user); err != nil {
		return nil, err
	}

	return a.getTokens(user.Id, user.Lang)
}

func (a *Auth) VerifyAndGetId(accessToken string) (*string, error) {
	tokenClaims, err := a.jwtService.ParseToken(accessToken, a.cfg.GetJwtAccessSecretKey())
	if err != nil || len(tokenClaims) == 0 {
		a.logger.ErrorLog.Println("verify access token error: ", err)

		return nil, helper.MakeGrpcBadRequestError(errTokenInvalid)
	}

	id, ok := tokenClaims["uid"]
	if !ok || id == "" {
		return nil, helper.MakeGrpcBadRequestError(errTokenInvalid)
	}
	result := id.(string)

	return &result, nil
}

func (a *Auth) SetLanguage(userId uuid.UUID, lang string) error {
	user := a.userRepo.GetUserById(userId)
	if err := a.CheckUser(user); err != nil {
		return err
	}
	if err := a.userRepo.UpdateLang(userId, lang); err != nil {
		a.logger.ErrorLog.Println("Set lang err: ", err)

		return err
	}

	return nil
}
func (a *Auth) getTokens(userId uuid.UUID, lang string) (*entity.Tokens, error) {
	data := map[string]interface{}{
		"uid":  userId.String(),
		"lang": lang,
	}
	accessToken, err := a.jwtService.GenerateToken(data, a.cfg.GetJwtAccessSecretKey(), a.cfg.GetAccessDuration())
	if err != nil {
		return nil, err
	}

	refreshToken, err := a.jwtService.GenerateToken(data, a.cfg.GetJwtRefreshSecretKey(), a.cfg.GetRefreshDuration())
	if err != nil {
		return nil, err
	}

	return entity.NewTokens(accessToken, refreshToken), nil
}
