package service

import (
	"errors"
	"github.com/DYSN-Project/auth/config"
	"github.com/DYSN-Project/auth/internal/helpers"
	"github.com/DYSN-Project/auth/internal/models"
	"github.com/DYSN-Project/auth/internal/models/consts"
	"github.com/DYSN-Project/auth/internal/repository"
	"github.com/DYSN-Project/auth/pkg/jwt"
	"github.com/DYSN-Project/auth/pkg/log"
	"github.com/google/uuid"
)

type UseCaseInterface interface {
	RegisterUser(email string, password string) (string, error)
	ConfirmRegister(token string) (*models.User, error)
	Login(email string, password string) (*models.Tokens, error)
	GetTokensByRefresh(refreshToken string) (*models.Tokens, error)
	Verify(accessToken string) error
}

var (
	errInvalidUserData  = errors.New(consts.ErrInvalidEmailOrPassword)
	errUserAlreadyExist = errors.New(consts.ErrUserAlreadyExist)
	errUserNotFound     = errors.New(consts.ErrUserNotFound)
	errInternalServer   = errors.New(consts.ErrInternalServer)
	errUserStatus       = errors.New(consts.ErrUserNotActive)
	errTokenInvalid     = errors.New(consts.ErrInvalidToken)
)

type Usecase struct {
	cfg        *config.Config
	jwtService jwt.JwtInterface
	userRepo   repository.UserRepoInterface
	logger     *log.Logger
}

func NewUseCase(cfg *config.Config,
	jwtSrv jwt.JwtInterface,
	usrRepo repository.UserRepoInterface,
	logger *log.Logger) *Usecase {
	return &Usecase{
		cfg:        cfg,
		jwtService: jwtSrv,
		userRepo:   usrRepo,
		logger:     logger,
	}
}

func (u *Usecase) RegisterUser(email string, password string) (string, error) {
	existUser := u.userRepo.GetUserByEmail(email)
	if !existUser.IsEmpty() {
		return "", errUserAlreadyExist
	}

	passwordHash, err := helpers.GetPwdHash(password, u.cfg.GetPwdSalt())
	if err != nil {
		u.logger.ErrorLog.Println("hash password error: ", err)

		return "", errInternalServer
	}
	regClaims := map[string]interface{}{
		"email": email,
		"pwd":   passwordHash,
	}
	token, err := u.jwtService.GenerateToken(regClaims, u.cfg.GetJwtRegSecretKey(), u.cfg.GetRegisterDuration())
	if err != nil {
		u.logger.ErrorLog.Println("generate token error: ", err)

		return "", errInternalServer
	}

	return token, nil
}

func (u *Usecase) ConfirmRegister(token string) (*models.User, error) {
	tokenClaims, err := u.jwtService.ParseToken(token, u.cfg.GetJwtRegSecretKey())
	if err != nil {
		u.logger.ErrorLog.Println("parse token error: ", err)

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

	existUser := u.userRepo.GetUserByEmail(email)
	if !existUser.IsEmpty() {
		return nil, errUserAlreadyExist
	}

	user := models.NewUser(email, pwd)
	user, err = u.userRepo.CreateUser(user)
	if err != nil {
		u.logger.ErrorLog.Println("create user error: ", err)

		return nil, errInternalServer
	}

	return user, nil
}

func (u *Usecase) Login(email string, password string) (*models.Tokens, error) {
	user := u.userRepo.GetUserByEmail(email)
	if user.IsEmpty() {
		return nil, errInvalidUserData
	}
	if !user.IsActive() {
		return nil, errUserStatus
	}

	comparedPassword := helpers.ComparePassword(password, user.Password, u.cfg.GetPwdSalt())
	if comparedPassword != nil {
		return nil, errInvalidUserData
	}

	return u.getTokens(user.ID)
}

func (u *Usecase) GetTokensByRefresh(refreshToken string) (*models.Tokens, error) {
	tokenClaims, err := u.jwtService.ParseToken(refreshToken, u.cfg.GetJwtRefreshSecretKey())
	if err != nil {
		u.logger.ErrorLog.Println("parse token error: ", err)

		return nil, errTokenInvalid
	}

	if _, ok := tokenClaims["uid"]; !ok {
		return nil, errTokenInvalid
	}
	userId, _ := uuid.Parse(tokenClaims["uid"].(string))

	user := u.userRepo.GetUserById(userId)
	if user.IsEmpty() {
		return nil, errUserNotFound
	}

	return u.getTokens(user.ID)
}

func (u *Usecase) Verify(accessToken string) error {
	verify, err := u.jwtService.Verify(accessToken, u.cfg.GetJwtAccessSecretKey())
	if err != nil || !verify {
		u.logger.ErrorLog.Println("verify access token error: ", err)

		return errTokenInvalid
	}
	return nil
}

func (u *Usecase) getTokens(userId uuid.UUID) (*models.Tokens, error) {
	data := map[string]interface{}{
		"uid": userId.String(),
	}
	accessToken, err := u.jwtService.GenerateToken(data, u.cfg.GetJwtAccessSecretKey(), u.cfg.GetAccessDuration())
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.jwtService.GenerateToken(data, u.cfg.GetJwtRefreshSecretKey(), u.cfg.GetRefreshDuration())
	if err != nil {
		return nil, err
	}

	return models.NewTokens(accessToken, refreshToken), nil
}
