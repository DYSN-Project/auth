package models

import (
	"github.com/DYSN-Project/auth/internal/models/consts"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	ID        uuid.UUID `gorm:"primary_key"`
	Email     string
	Password  string
	Status    int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func NewUserIngot() *User {
	return &User{}
}

func NewUser(email, password string) *User {
	return &User{
		Email:    email,
		Password: password,
	}
}

func (u *User) IsEmpty() bool {
	return u.ID == uuid.Nil
}

func (u *User) IsActive() bool {
	return u.Status == consts.UserStatusActive
}

func (u *User) IsNotActive() bool {
	return u.Status == consts.UserStatusNotActive
}

func (u *User) IsBanned() bool {
	return u.Status == consts.UserStatusBanned
}

func (u *User) IsDeleted() bool {
	return u.Status == consts.UserStatusDeleted
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	u.CreatedAt = time.Now()
	u.Status = consts.UserStatusActive

	return
}
