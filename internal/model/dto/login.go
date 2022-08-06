package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"strings"
)

type Login struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

func NewLogin(email string, password string) *Login {
	email = strings.TrimSpace(strings.ToLower(email))
	password = strings.TrimSpace(password)

	return &Login{email, password}
}

func (l *Login) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.Email, validation.Required, is.Email),
		validation.Field(&l.Password, validation.Required))
}
