package dto

import (
	"dysn/auth/internal/helper"
	"dysn/auth/internal/model/consts"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"strings"
)

type ChangeLang struct {
	Lang   string    `json:"ru"`
	UserId uuid.UUID `json:"userId"`
}

func NewChangeLang(lang string, userId uuid.UUID) *ChangeLang {
	return &ChangeLang{
		Lang:   strings.TrimSpace(strings.ToLower(lang)),
		UserId: userId,
	}
}

func (l *ChangeLang) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.Lang,
			validation.Required.Error(consts.ErrFieldRequired),
			validation.In(helper.LangList...).Error(consts.ErrFieldIncorrectFormat)),
		validation.Field(&l.UserId,
			validation.Required.Error(consts.ErrFieldRequired)))
}
