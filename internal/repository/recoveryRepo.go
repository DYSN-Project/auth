package repository

import (
	"dysn/auth/internal/model/entity"
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type RecoveryRepoInterface interface {
	CreateRecovery(recovery *entity.RecoveryPassword) (*entity.RecoveryPassword, error)
	GetRecoveryByEmail(email string) *entity.RecoveryPassword
	GetRecovery(email string, status int) *entity.RecoveryPassword
	UpdateRecovery(id uuid.UUID, data map[string]interface{}) error
	ExistEmail(email string) bool
}

type RecoveryRepository struct {
	db *gorm.DB
}

func NewRecoveryRepository(db *gorm.DB) *RecoveryRepository {
	return &RecoveryRepository{
		db: db,
	}
}

func (r *RecoveryRepository) GetRecoveryByEmail(email string) *entity.RecoveryPassword {
	recovery := entity.NewRecoveryIngot()
	r.db.Where("email = ? ", email).
		First(recovery)

	return recovery
}

func (r *RecoveryRepository) GetRecovery(email string,
	status int) *entity.RecoveryPassword {
	recovery := entity.NewRecoveryIngot()
	r.db.Where("email = ? AND status = ? ", email, status).
		First(recovery)

	return recovery
}

func (r *RecoveryRepository) CreateRecovery(recovery *entity.RecoveryPassword) (*entity.RecoveryPassword, error) {
	if err := r.db.Create(recovery).Error; err != nil {
		return nil, err
	}

	return recovery, nil
}

func (r *RecoveryRepository) UpdateRecovery(id uuid.UUID,
	data map[string]interface{}) error {
	return r.db.Model(entity.NewRecoveryIngot()).
		Where("id = ?", id).
		Updates(data).Error
}

func (r *RecoveryRepository) ExistEmail(email string) bool {
	var result struct {
		Found bool
	}
	err := r.db.Raw("SELECT EXISTS(SELECT 1 "+
		"FROM users "+
		"WHERE email = ? ) AS found",
		email).Scan(&result).Error

	if err != nil {
		fmt.Println("exist err: ", err)

		return false
	}

	return result.Found
}

func (r *RecoveryRepository) ExistRecovery(email string) bool {
	var result struct {
		Found bool
	}
	err := r.db.Raw("SELECT EXISTS(SELECT 1 "+
		"FROM recovery_password "+
		"WHERE email = ? ) AS found",
		email).Scan(&result).Error

	if err != nil {
		fmt.Println("exist err: ", err)

		return false
	}

	return result.Found
}
