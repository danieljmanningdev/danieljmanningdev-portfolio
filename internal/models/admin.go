package models

import "time"

type Admin struct {
	ID           int64
	Email        string
	PasswordHash string
	DisplayName  string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
