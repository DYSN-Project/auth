package dto

import (
	"dysn/auth/internal/model/consts"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"strings"
)

type RemovePass struct {
	Email string `form:"email" json:"email"`
}

func NewRemovePass(email string) *RemovePass {
	return &RemovePass{
		Email: strings.TrimSpace(strings.ToLower(email)),
	}
}

func (r *RemovePass) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Email, validation.Required.Error(consts.ErrFieldRequired),
			is.Email.Error(consts.ErrFieldIncorrectFormat)),
	)
}
