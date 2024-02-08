package entity

import (
	"dysn/auth/internal/model/consts"
	"github.com/google/uuid"
)

type RecoveryPassword struct {
	Id          uuid.UUID
	Email       string
	ConfirmCode string
	Status      int
	Date
}

func NewRecovery(email, code string) *RecoveryPassword {
	return &RecoveryPassword{
		Email:       email,
		ConfirmCode: code,
	}
}

func (r *RecoveryPassword) IsEmpty() bool {
	return r.Id == uuid.Nil
}

func (r *RecoveryPassword) IsConfirm() bool {
	return r.Status == consts.StatusConfirmed
}
