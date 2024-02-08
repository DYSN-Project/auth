package service

import (
	"dysn/auth/config"
	"dysn/auth/internal/helper"
	"dysn/auth/internal/model/entity"
	"dysn/auth/pkg/jwt"
	"dysn/auth/pkg/log"
	"github.com/google/uuid"
)

var jwtUserIdKey = "uid"
var langUserKey = "lg"

type baseService struct {
	cfg        *config.Config
	jwtService jwt.JwtInterface
	logger     *log.Logger
}

func NewService(cfg *config.Config, jwtService jwt.JwtInterface, logger *log.Logger) *baseService {
	return &baseService{cfg, jwtService, logger}
}

func (b *baseService) checkUser(user *entity.User) error {
	if user == nil {
		return helper.MakeGrpcBadRequestError(errUserNotFound)
	}
	if user.IsEmpty() {
		return helper.MakeGrpcBadRequestError(errInvalidUserData)
	}
	if !user.IsUserConfirmed() {
		return helper.MakeGrpcBadRequestError(errInvalidUserData)
	}

	return nil
}

func (b *baseService) getTokens(userId uuid.UUID, lang string) (*entity.Tokens, error) {
	accessToken, err := b.jwtService.GenerateToken(map[string]interface{}{
		jwtUserIdKey: userId.String(),
		langUserKey:  lang,
	}, b.cfg.GetJwtAccessSecretKey(),
		b.cfg.GetAccessDuration())
	if err != nil {
		return nil, err
	}

	refreshToken, err := b.jwtService.GenerateToken(map[string]interface{}{
		jwtUserIdKey: userId.String(),
	}, b.cfg.GetJwtRefreshSecretKey(),
		b.cfg.GetRefreshDuration())
	if err != nil {
		return nil, err
	}

	return entity.NewTokens(accessToken, refreshToken), nil
}
