package service

import (
	"errors"
	"github.com/DYSN-Project/auth/config"
	"github.com/DYSN-Project/auth/internal/helper"
	"github.com/DYSN-Project/auth/internal/model/consts"
	"github.com/DYSN-Project/auth/internal/model/entity"
	"github.com/DYSN-Project/auth/internal/repository"
	"github.com/DYSN-Project/auth/pkg/jwt"
	"github.com/DYSN-Project/auth/pkg/log"
	"github.com/google/uuid"
)

type AuthorizationInterface interface {
	Login(email string, password string) (*entity.Tokens, error)
	GetTokensByRefresh(refreshToken string) (*entity.Tokens, error)
	VerifyAndGetId(token string) (uuid.UUID, error)
}

var (
	errInvalidUserData  = errors.New(consts.ErrInvalidEmailOrPassword)
	errUserAlreadyExist = errors.New(consts.ErrUserAlreadyExist)
	errUserNotFound     = errors.New(consts.ErrUserNotFound)
	errInternalServer   = errors.New(consts.ErrInternalServer)
	errUserStatus       = errors.New(consts.ErrUserNotActive)
	errTokenInvalid     = errors.New(consts.ErrInvalidToken)
)

type Authorization struct {
	cfg        *config.Config
	jwtService jwt.JwtInterface
	userRepo   repository.UserRepoInterface
	logger     *log.Logger
}

func NewAuthorization(cfg *config.Config,
	jwtSrv jwt.JwtInterface,
	usrRepo repository.UserRepoInterface,
	logger *log.Logger) *Authorization {
	return &Authorization{
		cfg:        cfg,
		jwtService: jwtSrv,
		userRepo:   usrRepo,
		logger:     logger,
	}
}

func (a *Authorization) Login(email string, password string) (*entity.Tokens, error) {
	user := a.userRepo.GetUserByEmail(email)
	if user.IsEmpty() {
		return nil, errInvalidUserData
	}
	if !user.IsActive() {
		return nil, errUserStatus
	}

	comparedPassword := helper.ComparePassword(password, user.Password, a.cfg.GetPwdSalt())
	if comparedPassword != nil {
		return nil, errInvalidUserData
	}

	return a.getTokens(user.ID)
}

func (a *Authorization) GetTokensByRefresh(token string) (*entity.Tokens, error) {
	tokenClaims, err := a.jwtService.ParseToken(token, a.cfg.GetJwtRefreshSecretKey())
	if err != nil {
		a.logger.ErrorLog.Println("parse token error: ", err)

		return nil, errTokenInvalid
	}

	userId, err := a.getUserIdFromToken(tokenClaims)
	if err != nil {
		a.logger.ErrorLog.Println("parse token error: ", err)

		return nil, errTokenInvalid
	}

	user := a.userRepo.GetUserById(userId)
	if user.IsEmpty() {
		return nil, errUserNotFound
	}
	if !user.IsActive() {
		return nil, errUserStatus
	}

	return a.getTokens(user.ID)
}

func (a *Authorization) VerifyAndGetId(token string) (uuid.UUID, error) {
	tokenClaims, err := a.jwtService.ParseToken(token, a.cfg.GetJwtAccessSecretKey())
	if err != nil {
		a.logger.ErrorLog.Println("parse token error: ", err)

		return uuid.Nil, errTokenInvalid
	}

	userId, err := a.getUserIdFromToken(tokenClaims)
	if err != nil {
		a.logger.ErrorLog.Println("parse token error: ", err)

		return uuid.Nil, errTokenInvalid
	}

	return userId, nil
}

func (a *Authorization) getUserIdFromToken(tokenClaims map[string]interface{}) (uuid.UUID, error) {
	if _, ok := tokenClaims["uid"]; !ok {
		return uuid.Nil, errTokenInvalid
	}
	userId, err := uuid.Parse(tokenClaims["uid"].(string))
	if err != nil {
		return uuid.Nil, err
	}

	return userId, nil
}

func (a *Authorization) getTokens(userId uuid.UUID) (*entity.Tokens, error) {
	data := map[string]interface{}{
		"uid": userId.String(),
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
