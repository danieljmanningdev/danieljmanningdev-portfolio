package models

import "time"

type User struct {
	ID           int64     `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	StripeID     string    `db:"stripe_id"`
	Role         string    `db:"role"` // admin/client
	CreatedAt    time.Time `db:"created_at"`
}
