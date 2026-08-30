// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package models

import "time"

type Client struct {
	ID        int64
	Name      string
	Email     string
	Company   string
	Status    string
	Notes     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
