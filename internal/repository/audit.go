package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/models"
)

type AuditRepository struct {
	db *sql.DB
}

func NewAuditRepository(db *sql.DB) *AuditRepository {
	return &AuditRepository{
		db: db,
	}
}

func (r *AuditRepository) ListRecent(
	ctx context.Context,
	limit int,
	entityType string,
) ([]models.AuditEvent, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	entityType = strings.TrimSpace(entityType)

	query := `
		SELECT
			id,
			action,
			entity_type,
			entity_id,
			summary,
			details,
			created_at
		FROM audit_events
	`

	arguments := make([]any, 0, 2)
	if entityType != "" {
		query += "WHERE entity_type = ?\n"
		arguments = append(arguments, entityType)
	}

	query += "ORDER BY created_at DESC, id DESC\nLIMIT ?"
	arguments = append(arguments, limit)

	rows, err := r.db.QueryContext(
		ctx,
		query,
		arguments...,
	)
	if err != nil {
		return nil, fmt.Errorf("list audit events: %w", err)
	}
	defer rows.Close()

	events := make([]models.AuditEvent, 0, limit)

	for rows.Next() {
		var event models.AuditEvent
		var entityID sql.NullInt64

		if err := rows.Scan(
			&event.ID,
			&event.Action,
			&event.EntityType,
			&entityID,
			&event.Summary,
			&event.Details,
			&event.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan audit event: %w", err)
		}

		if entityID.Valid {
			id := entityID.Int64
			event.EntityID = &id
		}

		hydrateAuditEvent(&event)
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate audit events: %w", err)
	}

	return events, nil
}

func hydrateAuditEvent(event *models.AuditEvent) {
	if event == nil {
		return
	}

	switch event.EntityType {
	case "client":
		event.EntityLabel = "Client"
	case "project":
		event.EntityLabel = "Project"
	case "contract":
		event.EntityLabel = "Contract"
	case "blog_post":
		event.EntityLabel = "Journal"
	default:
		event.EntityLabel = "Workspace"
	}

	if event.EntityID == nil || strings.HasSuffix(event.Action, ".deleted") {
		return
	}

	id := strconv.FormatInt(*event.EntityID, 10)

	switch event.EntityType {
	case "client":
		event.URL = "/dashboard/clients/" + id
	case "project":
		event.URL = "/dashboard/projects/" + id
	case "contract":
		event.URL = "/dashboard/contracts/" + id
	case "blog_post":
		event.URL = "/dashboard/blog/" + id + "/edit"
	}
}
