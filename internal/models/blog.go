// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package models

import "time"

type BlogPost struct {
	ID          int64
	Title       string
	Slug        string
	Excerpt     string
	Content     string
	Status      string
	PublishedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
