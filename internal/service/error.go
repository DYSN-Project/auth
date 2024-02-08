package service

import (
	"dysn/auth/internal/model/consts"
	"errors"
	"fmt"
)

var (
	errTokenInvalid            = errors.New(consts.ErrInvalidToken)
	errInvalidUserData         = errors.New(consts.ErrInvalidEmailOrPassword)
	errInvalidUserCode         = errors.New(consts.ErrInvalidUserCode)
	errUserAlreadyExist        = errors.New(consts.ErrUserAlreadyExist)
	errUserNotFound            = errors.New(consts.ErrUserNotFound)
	errInternalServer          = errors.New(consts.ErrInternalServer)
	errAlreadyConfirmed        = errors.New(consts.ErrAlreadyConfirmed)
	errRecoveryRequestNotFound = errors.New(consts.ErrRecoveryRequestNotFound)
)

type BadRequestError struct {
	Code int
	Err  error
}

func (r *BadRequestError) Error() string {
	return fmt.Sprintf("status %d: err %v", r.Code, r.Err)
}
