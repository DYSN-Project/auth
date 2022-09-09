package service

import (
	"dysn/auth/config"
	"dysn/auth/internal/helper"
	"dysn/auth/internal/model/consts"
	"dysn/auth/internal/model/entity"
	"dysn/auth/internal/repository"
	"dysn/auth/internal/transport/grpc/client"
	"dysn/auth/pkg/log"
)

type RecoveryInterface interface {
	CreateRecovery(email string) (*entity.RecoveryPassword, error)
	ConfirmRecovery(email, code string) error
	ChangePassword(email, password string) error
}

type Recovery struct {
	cfg          *config.Config
	userRepo     repository.UserRepoInterface
	recoveryRepo repository.RecoveryRepoInterface
	notify       *client.Notify
	logger       *log.Logger
	BaseService
}

func NewRecovery(cfg *config.Config,
	userRepo repository.UserRepoInterface,
	recoveryRepo repository.RecoveryRepoInterface,
	notify *client.Notify,
	logger *log.Logger) *Recovery {
	return &Recovery{
		cfg:          cfg,
		userRepo:     userRepo,
		recoveryRepo: recoveryRepo,
		notify:       notify,
		logger:       logger,
	}
}

func (r *Recovery) CreateRecovery(email string) (*entity.RecoveryPassword, error) {
	user := r.userRepo.GetUserByEmail(email)
	if err := r.CheckUser(user); err != nil {
		return nil, err
	}

	code := helper.RandStringInt(r.cfg.GetCodeLength())
	codeHash, _ := helper.GetHash(code, r.cfg.GetCodeSalt())

	recovery := r.recoveryRepo.GetRecoveryByEmail(email)
	if !recovery.IsEmpty() {
		err := r.recoveryRepo.UpdateRecovery(recovery.Id,
			map[string]interface{}{
				"status":       consts.StatusActive,
				"confirm_code": codeHash,
			})
		if err != nil {
			r.logger.ErrorLog.Println("err confirm recovery: ", err)

			return nil, err
		}
	} else {
		recovery = entity.NewRecovery(email, codeHash)
		var err error
		recovery, err = r.recoveryRepo.CreateRecovery(recovery)
		if err != nil {
			r.logger.ErrorLog.Println("err create recovery: ", err)

			return nil, err
		}
	}

	r.notifyRecoveryPassword(user.Email, code, user.Lang)

	return recovery, nil
}

func (r *Recovery) ConfirmRecovery(email, code string) error {
	recovery := r.recoveryRepo.GetRecovery(email, consts.StatusActive)
	if recovery.IsEmpty() {
		return helper.MakeGrpcBadRequestError(errRecoveryRequestNotFound)
	}

	if !r.compareCode(code, recovery) {
		return helper.MakeGrpcBadRequestError(errInvalidUserCode)
	}

	err := r.recoveryRepo.UpdateRecovery(recovery.Id,
		map[string]interface{}{
			"status":       consts.StatusConfirmed,
			"confirm_code": nil,
		})
	if err != nil {
		r.logger.ErrorLog.Println("err confirm recovery: ", err)

		return err
	}

	r.logger.InfoLog.Println("user %s confirm recovery pass", email)

	return nil
}

func (r *Recovery) ChangePassword(email, password string) error {
	recovery := r.recoveryRepo.GetRecovery(email, consts.StatusConfirmed)
	if recovery.IsEmpty() {
		return helper.MakeGrpcBadRequestError(errRecoveryRequestNotFound)
	}

	passwordHash, err := helper.GetHash(password, r.cfg.GetPwdSalt())
	if err != nil {
		r.logger.ErrorLog.Println("hash password error: ", err)

		return errInternalServer
	}

	if err = r.userRepo.ChangePassword(email, passwordHash); err != nil {
		r.logger.ErrorLog.Println("err change password: ", err)

		return err
	}

	if err = r.recoveryRepo.UpdateRecovery(recovery.Id,
		map[string]interface{}{
			"status": consts.StatusCompleted,
		}); err != nil {
		r.logger.ErrorLog.Println("err confirm recovery: ", err)

		return err
	}

	r.logger.InfoLog.Println("user %s change pass", email)

	return nil
}

func (r *Recovery) compareCode(code string, recovery *entity.RecoveryPassword) bool {
	err := helper.CompareHash(code, recovery.ConfirmCode, r.cfg.GetCodeSalt())

	return err == nil
}

func (r *Recovery) notifyRecoveryPassword(email, code, lang string) {
	go func() {
		err := r.notify.RecoveryPassword(email, code, lang)
		if err != nil {
			r.logger.ErrorLog.Printf("Sending recovery password for $s error: $s ", email, err)

			return
		}
		r.logger.InfoLog.Printf("notify for %s was sent", email)

		return
	}()
}
