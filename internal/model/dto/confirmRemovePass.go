package dto

import (
	"dysn/auth/internal/model/consts"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"strings"
)

type ConfirmRemovePass struct {
	Email string `form:"email" json:"email"`
	Code  string `form:"code" json:"code"`
}

func NewConfirmRemovePass(email, code string) *ConfirmRemovePass {
	return &ConfirmRemovePass{
		Email: strings.TrimSpace(strings.ToLower(email)),
		Code:  strings.TrimSpace(code),
	}
}

func (l *ConfirmRemovePass) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.Email,
			validation.Required.Error(consts.ErrFieldRequired),
			is.Email.Error(consts.ErrFieldIncorrectFormat)),
		validation.Field(&l.Code,
			validation.Required.Error(consts.ErrFieldRequired),
			validation.Match(codeRegexp).Error(consts.ErrFieldIncorrectFormat)))
}
