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
