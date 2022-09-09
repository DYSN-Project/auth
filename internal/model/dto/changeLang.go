package dto

import (
	"dysn/auth/internal/helper"
	"dysn/auth/internal/model/consts"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"strings"
)

type ChangeLang struct {
	Lang string `form:"lang" json:"ru"`
}

func NewChangeLang(lang string) *ChangeLang {
	return &ChangeLang{
		Lang: strings.TrimSpace(strings.ToLower(lang)),
	}
}

func (l *ChangeLang) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.Lang,
			validation.Required.Error(consts.ErrFieldRequired),
			validation.In(helper.LangList...).Error(consts.ErrFieldIncorrectFormat)))
}
