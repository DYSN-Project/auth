package service

import (
	"context"
	"dysn/auth/internal/helper"
	"dysn/auth/internal/model/consts"
	"dysn/auth/internal/model/dto"
	"dysn/auth/internal/model/entity"
	"dysn/auth/internal/transport/grpc/client"
	"github.com/google/uuid"
)

type RecoveryRepoInterface interface {
	CreateRecovery(ctx context.Context, recovery *entity.RecoveryPassword) error
	GetRecoveryByStatus(ctx context.Context, email string, status int) (*entity.RecoveryPassword, error)
	GetRecoveryByEmail(ctx context.Context, email string) (*entity.RecoveryPassword, error)
	UpdateRecovery(ctx context.Context, id uuid.UUID, status int, code string) error
	ExistRecovery(ctx context.Context, email string) bool
}

type Recovery struct {
	userRepo     UserRepoInterface
	recoveryRepo RecoveryRepoInterface
	notify       *client.Notify
	*baseService
}

func NewRecovery(userRepo UserRepoInterface,
	recoveryRepo RecoveryRepoInterface,
	notify *client.Notify, service *baseService) *Recovery {
	return &Recovery{
		userRepo:     userRepo,
		recoveryRepo: recoveryRepo,
		notify:       notify,
		baseService:  service,
	}
}

func (r *Recovery) CreateRecovery(ctx context.Context, email string) (*entity.RecoveryPassword, error) {
	user, err := r.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		r.logger.ErrorLog.Println("err get user by email for recovery:", err)

		return nil, err
	}
	if err := r.checkUser(user); err != nil {
		return nil, err
	}

	code := helper.RandStringInt(r.cfg.GetCodeLength())
	codeHash, _ := helper.GetHash(code, r.cfg.GetCodeSalt())

	recovery, err := r.recoveryRepo.GetRecoveryByEmail(ctx, email)
	if err != nil {
		r.logger.ErrorLog.Println("err get recovery:", err)

		return nil, err
	}

	if recovery != nil && !recovery.IsEmpty() {
		err := r.recoveryRepo.UpdateRecovery(ctx, recovery.Id, consts.StatusActive, codeHash)
		if err != nil {
			r.logger.ErrorLog.Println("err update recovery: ", err)

			return nil, err
		}
	} else {
		recovery = entity.NewRecovery(email, codeHash)
		var err error
		err = r.recoveryRepo.CreateRecovery(ctx, recovery)
		if err != nil {
			r.logger.ErrorLog.Println("err create recovery: ", err)

			return nil, err
		}
	}

	//TODO::kafka
	r.notifyRecoveryPassword(user.Email, code, user.Lang)

	return recovery, nil
}

func (r *Recovery) ConfirmRecovery(ctx context.Context, confirmDto *dto.ConfirmRemovePass) error {
	recovery, err := r.recoveryRepo.GetRecoveryByStatus(ctx, confirmDto.Email, consts.StatusActive)
	if err != nil {
		r.logger.ErrorLog.Println("err get active recovery: ", err)

		return err
	}
	if recovery == nil || recovery.IsEmpty() {
		return helper.MakeGrpcBadRequestError(errRecoveryRequestNotFound)
	}

	err = helper.CompareHash(confirmDto.Code, recovery.ConfirmCode, r.cfg.GetCodeSalt())
	if err != nil {
		r.logger.ErrorLog.Println("err compare hash: ", err)

		return helper.MakeGrpcBadRequestError(errInvalidUserCode)
	}

	err = r.recoveryRepo.UpdateRecovery(ctx, recovery.Id, consts.StatusConfirmed, "")
	if err != nil {
		r.logger.ErrorLog.Println("err confirm recovery: ", err)

		return err
	}

	r.logger.InfoLog.Println("user %s confirm recovery pass", confirmDto.Email)

	return nil
}

func (r *Recovery) ChangePassword(ctx context.Context, passDto *dto.ChangePass) error {
	recovery, err := r.recoveryRepo.GetRecoveryByStatus(ctx, passDto.Email, consts.StatusConfirmed)
	if err != nil {
		r.logger.ErrorLog.Println("err get active recovery: ", err)

		return err
	}
	if recovery == nil || recovery.IsEmpty() {
		return helper.MakeGrpcBadRequestError(errRecoveryRequestNotFound)
	}

	passwordHash, err := helper.GetHash(passDto.Password, r.cfg.GetPwdSalt())
	if err != nil {
		r.logger.ErrorLog.Println("err hash password: ", err)

		return errInternalServer
	}

	if err = r.userRepo.ChangePasswordByEmail(ctx, passDto.Email, passwordHash); err != nil {
		r.logger.ErrorLog.Println("err change password: ", err)

		return err
	}

	err = r.recoveryRepo.UpdateRecovery(ctx, recovery.Id, consts.StatusCompleted, "")
	if err != nil {
		r.logger.ErrorLog.Println("err confirm recovery: ", err)

		return err
	}

	r.logger.InfoLog.Println("user %s change pass", passDto.Email)

	return nil
}

func (r *Recovery) notifyRecoveryPassword(email, code, lang string) {
	go func() {
		err := r.notify.RecoveryPassword(email, code, lang)
		if err != nil {
			r.logger.ErrorLog.Printf("err sending recovery password for %s : %s ", email, err)

			return
		}
		r.logger.InfoLog.Printf("notify for %s was sent", email)

		return
	}()
}
