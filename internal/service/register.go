package service

import (
	"github.com/DYSN-Project/auth/config"
	"github.com/DYSN-Project/auth/internal/helper"
	"github.com/DYSN-Project/auth/internal/model/entity"
	"github.com/DYSN-Project/auth/internal/repository"
	"github.com/DYSN-Project/auth/pkg/jwt"
	"github.com/DYSN-Project/auth/pkg/log"
	"github.com/google/uuid"
)

type RegistrationInterface interface {
	Register(email string, password string) (string, error)
	ConfirmRegister(token string) (*entity.User, error)
}

type Registration struct {
	cfg        *config.Config
	jwtService jwt.JwtInterface
	userRepo   repository.UserRepoInterface
	logger     *log.Logger
}

func NewRegistration(cfg *config.Config,
	jwtSrv jwt.JwtInterface,
	usrRepo repository.UserRepoInterface,
	logger *log.Logger) *Registration {
	return &Registration{
		cfg:        cfg,
		jwtService: jwtSrv,
		userRepo:   usrRepo,
		logger:     logger,
	}
}

func (r *Registration) Register(email string, password string) (string, error) {
	existUser := r.userRepo.GetUserByEmail(email)
	if !existUser.IsEmpty() {
		return "", errUserAlreadyExist
	}

	passwordHash, err := helper.GetPwdHash(password, r.cfg.GetPwdSalt())
	if err != nil {
		r.logger.ErrorLog.Println("hash password error: ", err)

		return "", errInternalServer
	}
	regClaims := map[string]interface{}{
		"email": email,
		"pwd":   passwordHash,
	}
	token, err := r.jwtService.GenerateToken(regClaims,
		r.cfg.GetJwtRegSecretKey(),
		r.cfg.GetRegisterDuration(),
	)
	if err != nil {
		r.logger.ErrorLog.Println("generate token error: ", err)

		return "", errInternalServer
	}

	return token, nil
}

func (r *Registration) ConfirmRegister(token string) (*entity.User, error) {
	tokenClaims, err := r.jwtService.ParseToken(token, r.cfg.GetJwtRegSecretKey())
	if err != nil {
		r.logger.ErrorLog.Println("parse token error: ", err)

		return nil, errTokenInvalid
	}

	if _, ok := tokenClaims["email"]; !ok {
		return nil, errTokenInvalid
	}
	email := tokenClaims["email"].(string)

	if _, ok := tokenClaims["pwd"]; !ok {
		return nil, errTokenInvalid
	}
	pwd := tokenClaims["pwd"].(string)

	existUser := r.userRepo.GetUserByEmail(email)
	if !existUser.IsEmpty() {
		return nil, errUserAlreadyExist
	}

	user := entity.NewUser(email, pwd)
	user, err = r.userRepo.CreateUser(user)
	if err != nil {
		r.logger.ErrorLog.Println("create user error: ", err)

		return nil, errInternalServer
	}

	return user, nil
}

func (r *Registration) getTokens(userId uuid.UUID) (*entity.Tokens, error) {
	data := map[string]interface{}{
		"uid": userId.String(),
	}
	accessToken, err := r.jwtService.GenerateToken(data, r.cfg.GetJwtAccessSecretKey(), r.cfg.GetAccessDuration())
	if err != nil {
		return nil, err
	}

	refreshToken, err := r.jwtService.GenerateToken(data, r.cfg.GetJwtRefreshSecretKey(), r.cfg.GetRefreshDuration())
	if err != nil {
		return nil, err
	}

	return entity.NewTokens(accessToken, refreshToken), nil
}
