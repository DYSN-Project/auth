package entity

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	Id          uuid.UUID `gorm:"primary_key"`
	Email       string
	Password    string
	ConfirmCode string
	Lang        string
	IsConfirmed bool
	Date
}

func NewUserIngot() *User {
	return &User{}
}

func NewUser(email,
	password,
	confirmCode,
	lang string) *User {
	return &User{
		Email:       email,
		Password:    password,
		ConfirmCode: confirmCode,
		Lang:        lang,
	}
}

func (u *User) IsEmpty() bool {
	return u.Id == uuid.Nil
}

func (u *User) IsUserConfirmed() bool {
	return u.IsConfirmed
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.Id = uuid.New()
	u.CreatedAt = time.Now()
	u.IsConfirmed = false

	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()

	return
}

func (u *User) TableName() string {
	return "users"
}
