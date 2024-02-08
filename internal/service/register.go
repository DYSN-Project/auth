package service

import (
	"context"
	"dysn/auth/internal/helper"
	"dysn/auth/internal/model/dto"
	"dysn/auth/internal/model/entity"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"time"
)

type Register struct {
	userRepo  UserRepoInterface
	userKafka *kafka.Writer
	*baseService
}

func NewRegister(userRepo UserRepoInterface,
	service *baseService,
	writer *kafka.Writer) *Register {
	return &Register{userRepo, writer, service}
}

func (r *Register) RegisterUser(ctx context.Context, registerDto *dto.Register) (*entity.User, error) {
	err := r.sendRegisterToKafka(registerDto.Email, "123456", registerDto.Lang)
	if err != nil {
		fmt.Println("Kafka err :", err)
		return nil, err
	}
	fmt.Println("Kafka was sender")
	return nil, nil

	if r.userRepo.ExistUserByEmail(ctx, registerDto.Email) {
		return nil, helper.MakeGrpcBadRequestError(errUserAlreadyExist)
	}

	passwordHash, err := helper.GetHash(registerDto.Password, r.cfg.GetPwdSalt())
	if err != nil {
		r.logger.ErrorLog.Println("err hash password: ", err)

		return nil, err
	}

	code := helper.RandStringInt(r.cfg.GetCodeLength())
	codeHash, err := helper.GetHash(code, r.cfg.GetCodeSalt())
	if err != nil {
		r.logger.ErrorLog.Println("err hash code: ", err)

		return nil, err
	}

	user := entity.NewUser(registerDto.Email, passwordHash, codeHash, registerDto.Lang)
	if err = r.userRepo.CreateUser(ctx, user); err != nil {
		r.logger.ErrorLog.Println("err create user: ", err)

		return nil, err
	}

	r.sendRegisterToKafka(user.Email, code, registerDto.Lang)

	return user, nil
}

func (r *Register) ConfirmRegister(ctx context.Context, confirmRegisterDto *dto.Confirm) (*entity.Tokens, error) {
	user, err := r.userRepo.GetUserByEmail(ctx, confirmRegisterDto.Email)
	if user == nil || user.IsEmpty() {
		return nil, helper.MakeGrpcBadRequestError(errUserNotFound)
	}
	if user.IsConfirmed {
		return nil, helper.MakeGrpcBadRequestError(errAlreadyConfirmed)
	}

	if err = helper.CompareHash(confirmRegisterDto.Password,
		user.Password,
		r.cfg.GetPwdSalt()); err != nil {
		r.logger.ErrorLog.Println("err compare hash password for confirm: ", err)

		return nil, helper.MakeGrpcBadRequestError(errInvalidUserData)
	}
	if !r.compareCode(confirmRegisterDto.Code, user) {
		return nil, helper.MakeGrpcBadRequestError(errInvalidUserCode)
	}

	if err = r.userRepo.ConfirmUser(ctx, user.Id); err != nil {
		r.logger.ErrorLog.Println("err confirm user: ", err)

		return nil, errInternalServer
	}

	tokens, err := r.getTokens(user.Id, user.Lang)
	if err != nil {
		r.logger.ErrorLog.Println("err generate tokens: ", err)

		return nil, errInternalServer
	}

	r.logger.InfoLog.Println("New user was registered: ", user.Email)
	//TODO::statistic

	return tokens, nil
}

/*
func (r *Register) notifyConfirmRegister(email, code, lang string) {
	go func() {
		err := r.notify.ConfirmRegister(email, code, lang)
		if err != nil {
			r.logger.ErrorLog.Printf("Sending confirm register for $s error: $s ", email, err)

			return
		}
		r.logger.InfoLog.Printf("notify for %s was sent", email)

		return
	}()
}*/

func (r *Register) compareCode(code string, user *entity.User) bool {
	return helper.CompareHash(code,
		user.ConfirmCode,
		r.cfg.GetCodeSalt()) == nil
}

func (r *Register) sendRegisterToKafka(email, code, lang string) error {
	const retries = 3
	var err error
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		val := struct {
			Email string `json:"email"`
			Lang  string `json:"lang"`
			Code  string `json:"code"`
		}{
			Email: email,
			Lang:  lang,
			Code:  code,
		}
		resVal, _ := json.Marshal(val)

		fmt.Println(string(resVal))
		err := r.userKafka.WriteMessages(ctx, kafka.Message{
			Topic: r.cfg.GetTopicUserRegister(),
			Value: resVal,
		})
		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(time.Millisecond * 250)
			continue
		}

		if err != nil {
			r.logger.ErrorLog.Printf("unexpected error %v", err)
			return err
		}
		break
	}

	r.logger.InfoLog.Println("send user to kafka was send to topic:", r.cfg.GetTopicUserRegister())

	return err
}
