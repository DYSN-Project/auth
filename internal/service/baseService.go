package service

import (
	"dysn/auth/internal/helper"
	"dysn/auth/internal/model/entity"
)

type BaseService struct{}

func (b *BaseService) CheckUser(user *entity.User) error {
	if user.IsEmpty() {
		return helper.MakeGrpcBadRequestError(errInvalidUserData)
	}
	if !user.IsUserConfirmed() {
		return helper.MakeGrpcBadRequestError(errInvalidUserData)
	}

	return nil
}
