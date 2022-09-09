package dto

import (
	"dysn/auth/internal/model/consts"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"regexp"
	"strings"
)

var codeRegexp = regexp.MustCompile("^[0-9]{6,6}$")

type Confirm struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
	Code     string `form:"codeWord" json:"code"`
}

func NewConfirm(email, password, code string) *Confirm {
	return &Confirm{
		Email:    strings.TrimSpace(strings.ToLower(email)),
		Password: strings.TrimSpace(password),
		Code:     strings.TrimSpace(code),
	}
}

func (r *Confirm) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Email, validation.Required.Error(consts.ErrFieldRequired),
			is.Email.Error(consts.ErrFieldIncorrectFormat)),
		validation.Field(&r.Password, validation.Required.Error(consts.ErrFieldRequired)),
		validation.Field(&r.Code, validation.Required.Error(consts.ErrFieldRequired),
			validation.Match(codeRegexp).Error(consts.ErrFieldIncorrectFormat)))
}
