package entity

import "time"

type Date struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
