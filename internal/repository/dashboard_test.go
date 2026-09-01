// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package repository

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/database"
)

func TestDashboardRepositorySummary(
	t *testing.T,
) {
	ctx := context.Background()

	databaseHandle, err := database.Open(ctx, ":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer func() {
		_ = databaseHandle.Close()
	}()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve test path")
	}

	migrationsDirectory := filepath.Join(
		filepath.Dir(filename),
		"..",
		"..",
		"migrations",
	)

	if err := database.RunMigrations(
		databaseHandle.SQL,
		migrationsDirectory,
	); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	clientResult, err := databaseHandle.SQL.ExecContext(ctx, `
		INSERT INTO clients (name, email, status)
		VALUES ('Active Client', 'active@example.com', 'active')
	`)
	if err != nil {
		t.Fatalf("insert active client: %v", err)
	}

	clientID, err := clientResult.LastInsertId()
	if err != nil {
		t.Fatalf("read client ID: %v", err)
	}

	if _, err := databaseHandle.SQL.ExecContext(ctx, `
		INSERT INTO clients (name, email, status)
		VALUES ('Inactive Client', 'inactive@example.com', 'inactive')
	`); err != nil {
		t.Fatalf("insert inactive client: %v", err)
	}

	if _, err := databaseHandle.SQL.ExecContext(ctx, `
		INSERT INTO projects (
			client_id,
			name,
			status,
			start_date,
			due_date
		)
		VALUES
			(?, 'Overdue Project', 'active', date('now', '-30 days'), date('now', '-1 day')),
			(?, 'Due Soon Project', 'active', date('now', '-7 days'), date('now', '+7 days')),
			(?, 'Completed Project', 'completed', date('now', '-30 days'), date('now', '-2 days'))
	`, clientID, clientID, clientID); err != nil {
		t.Fatalf("insert projects: %v", err)
	}

	if _, err := databaseHandle.SQL.ExecContext(ctx, `
		INSERT INTO contracts (
			client_id,
			title,
			status,
			value_cents
		)
		VALUES
			(?, 'Draft Contract', 'draft', 10000),
			(?, 'Sent Contract', 'sent', 20000),
			(?, 'Accepted Contract', 'accepted', 30000)
	`, clientID, clientID, clientID); err != nil {
		t.Fatalf("insert contracts: %v", err)
	}

	if _, err := databaseHandle.SQL.ExecContext(ctx, `
		INSERT INTO blog_posts (
			title,
			slug,
			excerpt,
			content,
			status,
			published_at
		)
		VALUES
			('Draft Post', 'draft-post', 'Draft excerpt', 'Draft content', 'draft', NULL),
			('Published Post', 'published-post', 'Published excerpt', 'Published content', 'published', CURRENT_TIMESTAMP)
	`); err != nil {
		t.Fatalf("insert blog posts: %v", err)
	}

	repository := NewDashboardRepository(databaseHandle.SQL)

	summary, err := repository.Summary(ctx)
	if err != nil {
		t.Fatalf("load dashboard summary: %v", err)
	}

	checks := map[string]struct {
		got  int
		want int
	}{
		"active clients": {
			got:  summary.ActiveClients,
			want: 1,
		},
		"active projects": {
			got:  summary.ActiveProjects,
			want: 2,
		},
		"overdue projects": {
			got:  summary.OverdueProjects,
			want: 1,
		},
		"projects due soon": {
			got:  summary.DueSoonProjects,
			want: 1,
		},
		"draft contracts": {
			got:  summary.DraftContracts,
			want: 1,
		},
		"sent contracts": {
			got:  summary.SentContracts,
			want: 1,
		},
		"draft posts": {
			got:  summary.DraftPosts,
			want: 1,
		},
		"published posts": {
			got:  summary.PublishedPosts,
			want: 1,
		},
		"attention total": {
			got:  summary.AttentionTotal,
			want: 5,
		},
	}

	for name, check := range checks {
		t.Run(name, func(t *testing.T) {
			if check.got != check.want {
				t.Fatalf(
					"expected %d, got %d",
					check.want,
					check.got,
				)
			}
		})
	}

	if len(summary.RecentAuditEvents) == 0 {
		t.Fatal("expected recent audit events")
	}
}

func TestAuditRepositoryFiltersAndBuildsEntityURLs(
	t *testing.T,
) {
	ctx := context.Background()

	databaseHandle, err := database.Open(ctx, ":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer func() {
		_ = databaseHandle.Close()
	}()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve test path")
	}

	migrationsDirectory := filepath.Join(
		filepath.Dir(filename),
		"..",
		"..",
		"migrations",
	)

	if err := database.RunMigrations(
		databaseHandle.SQL,
		migrationsDirectory,
	); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	result, err := databaseHandle.SQL.ExecContext(ctx, `
		INSERT INTO clients (name, email, status)
		VALUES ('Audit Client', 'audit@example.com', 'active')
	`)
	if err != nil {
		t.Fatalf("insert client: %v", err)
	}

	clientID, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("read client ID: %v", err)
	}

	repository := NewAuditRepository(databaseHandle.SQL)

	events, err := repository.ListRecent(ctx, 20, "client")
	if err != nil {
		t.Fatalf("list client events: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("expected 1 client event, got %d", len(events))
	}

	if events[0].EntityLabel != "Client" {
		t.Fatalf(
			"expected Client label, got %q",
			events[0].EntityLabel,
		)
	}

	wantURL := "/dashboard/clients/" + stringID(clientID)
	if events[0].URL != wantURL {
		t.Fatalf(
			"expected URL %q, got %q",
			wantURL,
			events[0].URL,
		)
	}
}

func stringID(id int64) string {
	const digits = "0123456789"

	if id == 0 {
		return "0"
	}

	var reversed [20]byte
	index := len(reversed)

	for id > 0 {
		index--
		reversed[index] = digits[id%10]
		id /= 10
	}

	return string(reversed[index:])
}
