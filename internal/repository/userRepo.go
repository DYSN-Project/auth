package repository

import (
	"github.com/DYSN-Project/auth/internal/models"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type UserRepoInterface interface {
	GetUserByEmail(email string) *models.User
	GetUserById(id uuid.UUID) *models.User
	CreateUser(user *models.User) (*models.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) GetUserByEmail(email string) *models.User {
	user := models.NewUserIngot()
	u.db.Where("email = ?", email).First(user)

	return user
}

func (u *UserRepository) GetUserById(id uint) *models.User {
	user := models.NewUserIngot()
	u.db.Where("id = ?", id).First(user)

	return user
}

func (u *UserRepository) CreateUser(user *models.User) (*models.User, error) {
	if err := u.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
