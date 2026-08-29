package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/models"
)

type DashboardRepository struct {
	db              *sql.DB
	auditRepository *AuditRepository
}

func NewDashboardRepository(db *sql.DB) *DashboardRepository {
	return &DashboardRepository{
		db:              db,
		auditRepository: NewAuditRepository(db),
	}
}

func (r *DashboardRepository) Summary(
	ctx context.Context,
) (models.DashboardSummary, error) {
	var summary models.DashboardSummary

	err := r.db.QueryRowContext(ctx, `
		SELECT
			(
				SELECT COUNT(*)
				FROM clients
				WHERE status = 'active'
			),
			(
				SELECT COUNT(*)
				FROM projects
				WHERE status = 'active'
			),
			(
				SELECT COUNT(*)
				FROM projects
				WHERE status IN ('planned', 'active')
				  AND due_date IS NOT NULL
				  AND date(due_date) < date('now')
			),
			(
				SELECT COUNT(*)
				FROM projects
				WHERE status IN ('planned', 'active')
				  AND due_date IS NOT NULL
				  AND date(due_date) >= date('now')
				  AND date(due_date) <= date('now', '+14 days')
			),
			(
				SELECT COUNT(*)
				FROM contracts
				WHERE status = 'draft'
			),
			(
				SELECT COUNT(*)
				FROM contracts
				WHERE status = 'sent'
			),
			(
				SELECT COUNT(*)
				FROM blog_posts
				WHERE status = 'draft'
			),
			(
				SELECT COUNT(*)
				FROM blog_posts
				WHERE status = 'published'
			)
	`).Scan(
		&summary.ActiveClients,
		&summary.ActiveProjects,
		&summary.OverdueProjects,
		&summary.DueSoonProjects,
		&summary.DraftContracts,
		&summary.SentContracts,
		&summary.DraftPosts,
		&summary.PublishedPosts,
	)
	if err != nil {
		return models.DashboardSummary{}, fmt.Errorf(
			"load dashboard summary: %w",
			err,
		)
	}

	summary.AttentionTotal =
		summary.OverdueProjects +
			summary.DueSoonProjects +
			summary.DraftContracts +
			summary.SentContracts +
			summary.DraftPosts

	recentEvents, err := r.auditRepository.ListRecent(
		ctx,
		8,
		"",
	)
	if err != nil {
		return models.DashboardSummary{}, err
	}

	summary.RecentAuditEvents = recentEvents

	return summary, nil
}
