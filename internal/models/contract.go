package models

import "time"

type Contract struct {
	ID        int64
	ClientID  int64
	ProjectID *int64

	ClientName  string
	ProjectName string

	Title      string
	Status     string
	StartDate  *time.Time
	EndDate    *time.Time
	ValueCents *int64
	Notes      string

	CreatedAt time.Time
	UpdatedAt time.Time
}
