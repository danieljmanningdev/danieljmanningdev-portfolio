package models

import "time"

type AuditEvent struct {
	ID          int64
	Action      string
	EntityType  string
	EntityID    *int64
	EntityLabel string
	Summary     string
	Details     string
	URL         string
	CreatedAt   time.Time
}
