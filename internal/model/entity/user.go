package entity

import (
	"github.com/google/uuid"
)

type User struct {
	Id          uuid.UUID
	Email       string
	Password    string
	ConfirmCode string
	Lang        string
	IsConfirmed bool
	Date
}

func NewUser(email,
	password,
	confirmCode,
	lang string) *User {
	return &User{
		Id:          uuid.New(),
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
