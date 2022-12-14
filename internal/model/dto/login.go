package dto

import (
	"dysn/auth/internal/model/consts"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"strings"
)

type Login struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

func NewLogin(email, password string) *Login {
	return &Login{
		Email:    strings.TrimSpace(strings.ToLower(email)),
		Password: strings.TrimSpace(password),
	}
}

func (l *Login) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.Email,
			validation.Required.Error(consts.ErrFieldRequired),
			is.Email.Error(consts.ErrFieldIncorrectFormat)),
		validation.Field(&l.Password,
			validation.Required.Error(consts.ErrFieldRequired)))
}
