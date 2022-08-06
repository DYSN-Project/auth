package repository

import (
	"github.com/DYSN-Project/auth/internal/model/entity"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type UserRepoInterface interface {
	GetUserByEmail(email string) *entity.User
	GetUserById(id uuid.UUID) *entity.User
	CreateUser(user *entity.User) (*entity.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) GetUserByEmail(email string) *entity.User {
	user := entity.NewUserIngot()
	u.db.Where("email = ?", email).First(user)

	return user
}

func (u *UserRepository) GetUserById(id uuid.UUID) *entity.User {
	user := entity.NewUserIngot()
	u.db.Where("id = ?", id).First(user)

	return user
}

func (u *UserRepository) CreateUser(user *entity.User) (*entity.User, error) {
	if err := u.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
