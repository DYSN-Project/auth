package dto

import (
	"dysn/auth/internal/helper"
	"dysn/auth/internal/model/consts"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"strings"
)

const minPasswordEntropy = 60

type Register struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
	Lang     string `form:"lang" json:"lang"`
}

func NewRegister(email, password, lang string) *Register {
	if lang == "" {
		lang = helper.GetDefaultLang()
	}
	return &Register{
		Email:    strings.TrimSpace(strings.ToLower(email)),
		Password: strings.TrimSpace(password),
		Lang:     strings.TrimSpace(strings.ToLower(lang)),
	}
}

func (r *Register) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Email, validation.Required.Error(consts.ErrFieldRequired),
			is.Email.Error(consts.ErrFieldIncorrectFormat)),
		validation.Field(&r.Password,
			validation.Required.Error(consts.ErrFieldRequired),
			validation.By(checkPassword)),
		validation.Field(&r.Lang,
			validation.In(helper.LangList...).Error(consts.ErrFieldIncorrectFormat)),
	)
}

func checkPassword(value interface{}) error {
	err := passwordvalidator.Validate(fmt.Sprintf("%v", value), minPasswordEntropy)
	if err != nil {
		return errors.New(consts.ErrFieldIncorrectFormat)
	}

	return nil
}
