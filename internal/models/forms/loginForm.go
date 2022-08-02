package forms

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"strings"
)

type LoginForm struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

func NewLoginForm(email string, password string) *LoginForm {
	email = strings.TrimSpace(strings.ToLower(email))
	password = strings.TrimSpace(password)

	return &LoginForm{email, password}
}

func NewLoginFormIngot() *LoginForm {
	return &LoginForm{}
}

func (l *LoginForm) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.Email, validation.Required, is.Email),
		validation.Field(&l.Password, validation.Required))
}
