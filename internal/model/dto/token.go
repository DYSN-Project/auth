package dto

import (
	"dysn/auth/internal/model/consts"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"strings"
)

type Token struct {
	Token string `form:"token" json:"token"`
}

func NewToken(token string) *Token {
	return &Token{strings.TrimSpace(token)}
}

func (t *Token) Validate() error {
	return validation.ValidateStruct(t, validation.Field(&t.Token,
		validation.Required.Error(consts.ErrFieldRequired)))
}
