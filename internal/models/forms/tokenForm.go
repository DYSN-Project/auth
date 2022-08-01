package forms

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"strings"
)

type TokenFrom struct {
	Token string `form:"token" json:"token"`
}

func NewTokenForm(token string) *TokenFrom {
	token = strings.TrimSpace(token)
	return &TokenFrom{token}
}

func NewTokenFormIngot() *TokenFrom {
	return &TokenFrom{}
}

func (t *TokenFrom) Validate() error {
	return validation.ValidateStruct(t, validation.Field(&t.Token, validation.Required))
}
