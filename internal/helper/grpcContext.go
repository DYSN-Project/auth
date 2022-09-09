package helper

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const userIdKey = "uid"

func SetUserToGrpcContext(ctx context.Context, userId uuid.UUID) context.Context {
	return metadata.AppendToOutgoingContext(ctx,
		userIdKey, userId.String(),
		"time", timestamppb.Now().String())
}

func GetUserIdFromGrpcContext(ctx context.Context) uuid.UUID {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return uuid.Nil
	}
	userMd := md.Get(userIdKey)
	if len(userMd) == 0 {
		return uuid.Nil
	}

	userId, err := uuid.Parse(userMd[0])
	if err != nil {
		return uuid.Nil
	}

	return userId
}
