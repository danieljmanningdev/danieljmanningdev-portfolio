package models

import "time"

type AdminSession struct {
	ID         int64
	AdminID    int64
	TokenHash  string
	ExpiresAt  time.Time
	LastSeenAt time.Time
	CreatedAt  time.Time
}
