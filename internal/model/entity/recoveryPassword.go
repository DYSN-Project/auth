package entity

import (
	"dysn/auth/internal/model/consts"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

type RecoveryPassword struct {
	Id          uuid.UUID `gorm:"primary_key"`
	Email       string
	ConfirmCode string
	Status      int
	Date
}

func NewRecoveryIngot() *RecoveryPassword {
	return &RecoveryPassword{}
}

func NewRecovery(email, code string) *RecoveryPassword {
	return &RecoveryPassword{
		Email:       email,
		ConfirmCode: code,
	}
}

func (r *RecoveryPassword) BeforeCreate(tx *gorm.DB) (err error) {
	r.Id = uuid.New()
	r.CreatedAt = time.Now()
	r.Status = consts.StatusActive

	return
}

func (r *RecoveryPassword) TableName() string {
	return "recovery_password"
}

func (r *RecoveryPassword) IsEmpty() bool {
	return r.Id == uuid.Nil
}

func (r *RecoveryPassword) IsConfirm() bool {
	return r.Status == consts.StatusConfirmed
}
