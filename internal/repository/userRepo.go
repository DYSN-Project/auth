package repository

import (
	"dysn/auth/internal/model/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type UserRepoInterface interface {
	GetUserByEmail(email string) *entity.User
	GetUserById(id uuid.UUID) *entity.User
	CreateUser(user *entity.User) (*entity.User, error)
	ConfirmUser(id uuid.UUID) error
	Add2FaCode(id uuid.UUID, code string) error
	Remove2FaCode(id uuid.UUID) error
	SetConfirmCode(id uuid.UUID, code string) error
	Confirm2FaCode(id uuid.UUID) error
	ChangePassword(email, password string) error
	UpdateLang(userId uuid.UUID, lang string) error
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

func (u *UserRepository) ConfirmUser(id uuid.UUID) error {
	return u.db.Model(entity.NewUserIngot()).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_confirmed": true,
			"confirm_code": nil,
			"updated_at":   time.Now(),
		}).Error
}

func (u *UserRepository) Add2FaCode(id uuid.UUID, code string) error {
	return u.db.Model(entity.NewUserIngot()).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"two_factor_code": code,
			"updated_at":      time.Now(),
		}).Error
}

func (u *UserRepository) Confirm2FaCode(id uuid.UUID) error {
	return u.db.Model(entity.NewUserIngot()).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_included_2fa": true,
			"updated_at":      time.Now(),
		}).Error
}

func (u *UserRepository) Remove2FaCode(id uuid.UUID) error {
	return u.db.Model(entity.NewUserIngot()).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_included_2fa": false,
			"confirm_code":    nil,
			"updated_at":      time.Now(),
		}).Error
}

func (u *UserRepository) SetConfirmCode(id uuid.UUID, code string) error {
	return u.db.Model(entity.NewUserIngot()).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"confirm_code": code,
			"updated_at":   time.Now(),
		}).Error
}

func (u *UserRepository) ChangePassword(email, password string) error {
	return u.db.Model(entity.NewUserIngot()).
		Where("email = ?", email).
		Updates(map[string]interface{}{
			"password":   password,
			"updated_at": time.Now(),
		}).Error
}

func (u *UserRepository) UpdateLang(userId uuid.UUID, lang string) error {
	return u.db.Model(entity.NewUserIngot()).
		Where("id = ?", userId).
		Updates(map[string]interface{}{
			"lang":       lang,
			"updated_at": time.Now(),
		}).Error
}
