package client

import (
	"context"
	pb "dysn/auth/internal/transport/grpc/pb/notify"
	"dysn/auth/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Notify struct {
	client pb.NotifyClient
	logger *log.Logger
}
type NotifyInterface interface {
	ConfirmRegister(email, code, lang string) error
	DisableGa(email, code, lang string) error
	RecoveryPassword(email, code, lang string) error
}

func NewNotify(address string, logger *log.Logger) *Notify {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.ErrorLog.Panic(err)
	}

	return &Notify{
		client: pb.NewNotifyClient(conn),
		logger: logger,
	}
}

func (n *Notify) ConfirmRegister(email, code, lang string) error {
	request := &pb.EmailWithCode{
		Email: email,
		Code:  code,
		Lang:  lang,
	}
	if _, err := n.client.ConfirmRegister(context.Background(), request); err != nil {
		return err
	}

	return nil
}

func (n *Notify) RecoveryPassword(email, code, lang string) error {
	request := &pb.EmailWithCode{
		Email: email,
		Code:  code,
		Lang:  lang,
	}
	if _, err := n.client.RecoveryPassword(context.Background(), request); err != nil {
		return err
	}

	return nil
}
