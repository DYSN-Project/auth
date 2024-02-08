package repository

import (
	"context"
	"database/sql"
	"dysn/auth/internal/model/consts"
	"dysn/auth/internal/model/entity"
	"errors"
	"github.com/google/uuid"
	"time"
)

type RecoveryRepository struct {
	db *sql.DB
}

func NewRecoveryRepository(db *sql.DB) *RecoveryRepository {
	return &RecoveryRepository{
		db: db,
	}
}

func (r *RecoveryRepository) GetRecoveryByEmail(ctx context.Context, email string) (*entity.RecoveryPassword, error) {
	query := `SELECT id,
       email,
       confirm_code,
       status,
       created_at,
       updated_at FROM recovery_password 
                  WHERE email = $1`
	rows, err := r.db.QueryContext(ctx, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	var recovery *entity.RecoveryPassword
	err = rows.Scan(&recovery.Id,
		&recovery.Email,
		&recovery.ConfirmCode,
		&recovery.Status,
		&recovery.CreatedAt,
		&recovery.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return recovery, nil
}

func (r *RecoveryRepository) GetRecoveryByStatus(ctx context.Context, email string, status int) (*entity.RecoveryPassword, error) {
	query := `SELECT id,
       email,
       confirm_code,
       status,
       created_at,
       updated_at FROM recovery_password 
                  WHERE email = $1 AND status = $2`
	rows, err := r.db.QueryContext(ctx, query, email, status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	var recovery *entity.RecoveryPassword
	err = rows.Scan(&recovery.Id,
		&recovery.Email,
		&recovery.ConfirmCode,
		&recovery.Status,
		&recovery.CreatedAt,
		&recovery.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return recovery, nil
}

func (r *RecoveryRepository) CreateRecovery(ctx context.Context, recovery *entity.RecoveryPassword) error {
	query := `INSERT INTO recovery_password (id,
                      email,
                      confirm_code,
                      status) VALUES ($1, $2, $3, $4)`

	_, err := r.db.ExecContext(ctx, query,
		uuid.New(),
		recovery.Email,
		recovery.ConfirmCode,
		consts.StatusActive)

	return err
}

func (r *RecoveryRepository) UpdateRecovery(ctx context.Context, id uuid.UUID, status int, code string) error {
	query := `UPDATE recovery_password SET 
                  status = $1,
                  code = $2,
                  updated_at = $3
             WHERE id = $4`
	_, err := r.db.ExecContext(ctx, query, status,
		code,
		time.Now(),
		id)

	return err
}

func (r *RecoveryRepository) ExistRecovery(ctx context.Context, email string) bool {
	var exist bool
	query := `SELECT EXISTS( SELECT 1 FROM recovery_password where email = $1)`
	if err := r.db.QueryRowContext(ctx, query, email).Scan(&exist); err != nil {
		return false
	}

	return exist
}
