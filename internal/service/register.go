package service

import (
	"dysn/auth/config"
	"dysn/auth/internal/helper"
	"dysn/auth/internal/model/entity"
	"dysn/auth/internal/repository"
	"dysn/auth/internal/transport/grpc/client"
	"dysn/auth/pkg/jwt"
	"dysn/auth/pkg/log"
	"github.com/google/uuid"
)

type RegisterInterface interface {
	RegisterUser(email, password, lang string) (*entity.User, error)
	ConfirmRegister(email, password, code string) (*entity.Tokens, error)
}

type Register struct {
	cfg        *config.Config
	jwtService jwt.JwtInterface
	userRepo   repository.UserRepoInterface
	logger     *log.Logger
	notify     *client.Notify
	BaseService
}

func NewRegister(cfg *config.Config,
	jwtSrv jwt.JwtInterface,
	usrRepo repository.UserRepoInterface,
	logger *log.Logger,
	notify *client.Notify) *Register {
	return &Register{
		cfg:        cfg,
		jwtService: jwtSrv,
		userRepo:   usrRepo,
		logger:     logger,
		notify:     notify,
	}
}

func (r *Register) RegisterUser(email, password, lang string) (*entity.User, error) {
	existUser := r.userRepo.GetUserByEmail(email)
	if !existUser.IsEmpty() {
		return nil, helper.MakeGrpcBadRequestError(errAlreadyConfirmed)
	}

	passwordHash, err := helper.GetHash(password, r.cfg.GetPwdSalt())
	if err != nil {
		r.logger.ErrorLog.Println("hash password error: ", err)

		return nil, err
	}

	code := helper.RandStringInt(r.cfg.GetCodeLength())
	codeHash, _ := helper.GetHash(code, r.cfg.GetCodeSalt())

	user := entity.NewUser(email, passwordHash, codeHash, lang)
	user, err = r.userRepo.CreateUser(user)
	if err != nil {
		r.logger.ErrorLog.Println("Save user error: ", err)

		return nil, err
	}

	r.notifyConfirmRegister(user.Email, code, lang)

	return user, nil
}

func (r *Register) ConfirmRegister(email, password, code string) (*entity.Tokens, error) {
	user := r.userRepo.GetUserByEmail(email)
	if user.IsEmpty() {
		return nil, helper.MakeGrpcBadRequestError(errInvalidUserData)
	}
	if user.IsConfirmed {
		return nil, helper.MakeGrpcBadRequestError(errAlreadyConfirmed)
	}

	err := helper.CompareHash(password, user.Password, r.cfg.GetPwdSalt())
	if err != nil {
		return nil, helper.MakeGrpcBadRequestError(errInvalidUserData)
	}
	if !r.compareCode(code, user) {
		return nil, helper.MakeGrpcBadRequestError(errInvalidUserCode)
	}

	if err = r.userRepo.ConfirmUser(user.Id); err != nil {
		r.logger.ErrorLog.Println("confirm user error: ", err)

		return nil, errInternalServer
	}

	tokens, err := r.getTokens(user.Id, user.Lang)
	if err != nil {
		r.logger.ErrorLog.Println("generate tokens error: ", err)

		return nil, errInternalServer
	}

	r.logger.InfoLog.Println("New user was registered: ", user.Email)

	return tokens, nil
}

func (r *Register) notifyConfirmRegister(email, code, lang string) {
	go func() {
		err := r.notify.ConfirmRegister(email, code, lang)
		if err != nil {
			r.logger.ErrorLog.Printf("Sending confirm register for $s error: $s ", email, err)

			return
		}
		r.logger.InfoLog.Printf("notify for %s was sent", email)

		return
	}()
}

func (r *Register) getTokens(userId uuid.UUID, lang string) (*entity.Tokens, error) {
	data := map[string]interface{}{
		"uid":  userId.String(),
		"lang": lang,
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

func (r *Register) compareCode(code string, user *entity.User) bool {
	return helper.CompareHash(code,
		user.ConfirmCode,
		r.cfg.GetCodeSalt()) == nil
}
