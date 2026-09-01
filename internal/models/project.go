// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package models

import "time"

type Project struct {
	ID          int64
	ClientID    int64
	ClientName  string
	Name        string
	Description string
	Status      string
	StartDate   *time.Time
	DueDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
