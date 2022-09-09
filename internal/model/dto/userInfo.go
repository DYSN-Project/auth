package dto

import (
	"github.com/google/uuid"
	"time"
)

type UserInfo struct {
	Id             uuid.UUID
	QueryTimestamp time.Time
}

func NewUserInfo(id uuid.UUID, queryTimestamp time.Time) *UserInfo {
	return &UserInfo{Id: id, QueryTimestamp: queryTimestamp}
}
