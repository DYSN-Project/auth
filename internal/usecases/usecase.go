package usecases

import (
	"errors"
	"github.com/DYSN-Project/auth/config"
	"github.com/DYSN-Project/auth/internal/helpers"
	"github.com/DYSN-Project/auth/internal/models"
	"github.com/DYSN-Project/auth/internal/models/consts"
	"github.com/DYSN-Project/auth/internal/packages/jwt"
	"github.com/DYSN-Project/auth/internal/packages/log"
	"github.com/DYSN-Project/auth/internal/repository"
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
	token, err := u.jwtService.GenerateRegisterToken(email, passwordHash, u.cfg.GetRegisterDuration())
	if err != nil {
		u.logger.ErrorLog.Println("generate token error: ", err)

		return "", errInternalServer
	}

	return token, nil
}

func (u *Usecase) ConfirmRegister(token string) (*models.User, error) {
	tokenClaims, err := u.jwtService.ParseRegisterToken(token)
	if err != nil {
		u.logger.ErrorLog.Println("parse token error: ", err)

		return nil, errTokenInvalid
	}

	existUser := u.userRepo.GetUserByEmail(tokenClaims.Email)
	if !existUser.IsEmpty() {
		return nil, errUserAlreadyExist
	}

	user := models.NewUser(tokenClaims.Email, tokenClaims.Password)
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
	tokenClaims, err := u.jwtService.ParseAuthToken(refreshToken)
	if err != nil {
		u.logger.ErrorLog.Println("parse token error: ", err)

		return nil, errTokenInvalid
	}

	user := u.userRepo.GetUserById(tokenClaims.Uid)
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
	accessToken, err := u.jwtService.GenerateAuthToken(userId, u.cfg.GetAccessDuration())
	if err != nil {
		return nil, err
	}
	refreshToken, err := u.jwtService.GenerateAuthToken(userId, u.cfg.GetRefreshDuration())
	if err != nil {
		return nil, err
	}

	return models.NewTokens(accessToken, refreshToken), nil
}
