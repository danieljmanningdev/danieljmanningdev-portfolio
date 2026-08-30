// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package models

type DashboardSummary struct {
	ActiveClients     int
	ActiveProjects    int
	OverdueProjects   int
	DueSoonProjects   int
	DraftContracts    int
	SentContracts     int
	DraftPosts        int
	PublishedPosts    int
	AttentionTotal    int
	RecentAuditEvents []AuditEvent
}
