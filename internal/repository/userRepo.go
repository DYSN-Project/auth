package repository

import (
	"context"
	"database/sql"
	"dysn/auth/internal/model/entity"
	"errors"
	"github.com/google/uuid"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT id,
       email,
       password,
       confirm_code,
       lang,
       is_confirmed,
       created_at,
       updated_at FROM users
             WHERE email = $1`
	rows, err := u.db.QueryContext(ctx, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()
	var user *entity.User

	err = rows.Scan(&user.Id,
		&user.Email,
		user.Password,
		&user.ConfirmCode,
		&user.Lang,
		&user.IsConfirmed,
		&user.CreatedAt,
		&user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepository) GetUserById(ctx context.Context, userId uuid.UUID) (*entity.User, error) {
	query := `SELECT id,
       email,
       password,
       confirm_code,
       lang,
       is_confirmed,
       created_at,
       updated_at FROM users
             WHERE id = $1`
	rows, err := u.db.QueryContext(ctx, query, userId.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()
	var user *entity.User

	err = rows.Scan(&user.Id,
		&user.Email,
		user.Password,
		&user.ConfirmCode,
		&user.Lang,
		&user.IsConfirmed,
		&user.CreatedAt,
		&user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users (id,
                  email,
                  password,
                  confirm_code,
                  lang,
                  is_confirmed) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := u.db.ExecContext(ctx, query,
		user.Id,
		user.Email,
		user.Password,
		user.ConfirmCode,
		user.Lang,
		false)

	return err
}

func (u *UserRepository) ConfirmUser(ctx context.Context, userId uuid.UUID) error {
	query := `UPDATE users SET 
                  is_confirmed = $1,
                  confirm_code = $2,
                  updated_at = $3
             WHERE id = $4`
	_, err := u.db.ExecContext(ctx, query, true,
		"",
		time.Now(),
		userId)

	return err
}

func (u *UserRepository) Add2FaCode(ctx context.Context, userId uuid.UUID, code string) error {
	query := `UPDATE users SET 
                  two_factor_code = $1,
                  updated_at = $2
             WHERE id = $3`
	_, err := u.db.ExecContext(ctx, query, code,
		time.Now(),
		userId)

	return err
}

func (u *UserRepository) Confirm2FaCode(ctx context.Context, userId uuid.UUID) error {
	query := `UPDATE users SET 
                  is_included_2fa = $1,
                  updated_at = $2
             WHERE id = $3`
	_, err := u.db.ExecContext(ctx, query, false,
		time.Now(),
		userId)

	return err
}

func (u *UserRepository) Remove2FaCode(ctx context.Context, userId uuid.UUID) error {
	query := `UPDATE users SET 
                  is_included_2fa = $1,
                  confirm_code = $2,
                  updated_at = $3
             WHERE id = $4`
	_, err := u.db.ExecContext(ctx, query, false,
		nil,
		time.Now(),
		userId)

	return err
}

func (u *UserRepository) SetConfirmCode(ctx context.Context, userId uuid.UUID, code string) error {
	query := `UPDATE users SET 
                  confirm_code = $1,
                  updated_at = $2
             WHERE id = $3`
	_, err := u.db.ExecContext(ctx, query, code,
		time.Now(),
		userId)

	return err
}

func (u *UserRepository) UpdateLang(ctx context.Context, userId uuid.UUID, lang string) error {
	query := `UPDATE users SET 
                  lang = $1,
                  updated_at = $2
             WHERE id = $3`
	_, err := u.db.ExecContext(ctx, query, lang,
		time.Now(),
		userId)

	return err
}

func (u *UserRepository) ChangePasswordByEmail(ctx context.Context, email, password string) error {
	query := `UPDATE users SET 
                  password = $1,
                  updated_at = $2
             WHERE email = $3`
	_, err := u.db.ExecContext(ctx, query, password,
		time.Now(),
		email)

	return err
}

func (u *UserRepository) ExistUserByEmail(ctx context.Context, email string) bool {
	var exist bool
	query := `SELECT EXISTS( SELECT 1 FROM users where email = $1)`
	if err := u.db.QueryRowContext(ctx, query, email).Scan(&exist); err != nil {
		return false
	}

	return exist
}
