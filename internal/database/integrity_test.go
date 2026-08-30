// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package database

import (
	"context"
	"database/sql"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func openIntegrityTestDatabase(
	t *testing.T,
) *sql.DB {
	t.Helper()

	database, err := Open(
		context.Background(),
		":memory:",
	)
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}

	t.Cleanup(func() {
		_ = database.Close()
	})

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

	if err := RunMigrations(
		database.SQL,
		migrationsDirectory,
	); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	return database.SQL
}

func TestDatabaseIntegrityTriggersRejectInvalidRecords(
	t *testing.T,
) {
	database := openIntegrityTestDatabase(t)

	clientResult, err := database.Exec(`
		INSERT INTO clients (
			name,
			email,
			status
		)
		VALUES ('Client One', 'one@example.com', 'active')
	`)
	if err != nil {
		t.Fatalf("insert valid client: %v", err)
	}

	clientID, err := clientResult.LastInsertId()
	if err != nil {
		t.Fatalf("read client ID: %v", err)
	}

	secondClientResult, err := database.Exec(`
		INSERT INTO clients (
			name,
			email,
			status
		)
		VALUES ('Client Two', 'two@example.com', 'active')
	`)
	if err != nil {
		t.Fatalf("insert second client: %v", err)
	}

	secondClientID, err := secondClientResult.LastInsertId()
	if err != nil {
		t.Fatalf("read second client ID: %v", err)
	}

	projectResult, err := database.Exec(`
		INSERT INTO projects (
			client_id,
			name,
			status,
			start_date,
			due_date
		)
		VALUES (?, 'Valid Project', 'active', '2026-08-01', '2026-08-31')
	`, clientID)
	if err != nil {
		t.Fatalf("insert valid project: %v", err)
	}

	projectID, err := projectResult.LastInsertId()
	if err != nil {
		t.Fatalf("read project ID: %v", err)
	}

	tests := []struct {
		name      string
		statement string
		arguments []any
		contains  string
	}{
		{
			name: "invalid client status",
			statement: `
				INSERT INTO clients (name, email, status)
				VALUES ('Invalid', 'invalid@example.com', 'unknown')
			`,
			contains: "invalid client record",
		},
		{
			name: "invalid project status",
			statement: `
				INSERT INTO projects (client_id, name, status)
				VALUES (?, 'Invalid Project', 'unknown')
			`,
			arguments: []any{clientID},
			contains:  "invalid project record",
		},
		{
			name: "project due date before start date",
			statement: `
				INSERT INTO projects (
					client_id,
					name,
					status,
					start_date,
					due_date
				)
				VALUES (?, 'Invalid Dates', 'planned', '2026-09-30', '2026-09-01')
			`,
			arguments: []any{clientID},
			contains:  "invalid project record",
		},
		{
			name: "negative contract value",
			statement: `
				INSERT INTO contracts (
					client_id,
					title,
					status,
					value_cents
				)
				VALUES (?, 'Invalid Value', 'draft', -1)
			`,
			arguments: []any{clientID},
			contains:  "invalid contract record",
		},
		{
			name: "contract project belongs to another client",
			statement: `
				INSERT INTO contracts (
					client_id,
					project_id,
					title,
					status
				)
				VALUES (?, ?, 'Mismatched Project', 'draft')
			`,
			arguments: []any{secondClientID, projectID},
			contains:  "invalid contract record",
		},
		{
			name: "published blog post without publication date",
			statement: `
				INSERT INTO blog_posts (
					title,
					slug,
					excerpt,
					content,
					status
				)
				VALUES (
					'Invalid Publication',
					'invalid-publication',
					'Excerpt',
					'Content',
					'published'
				)
			`,
			contains: "invalid blog post record",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := database.Exec(
				test.statement,
				test.arguments...,
			)
			if err == nil {
				t.Fatal("expected database constraint error")
			}

			if !strings.Contains(err.Error(), test.contains) {
				t.Fatalf(
					"expected error containing %q, got %q",
					test.contains,
					err.Error(),
				)
			}
		})
	}
}

func TestAuditTriggersCaptureWorkspaceHistory(
	t *testing.T,
) {
	database := openIntegrityTestDatabase(t)

	result, err := database.Exec(`
		INSERT INTO clients (
			name,
			email,
			status
		)
		VALUES ('Audit Client', 'audit@example.com', 'active')
	`)
	if err != nil {
		t.Fatalf("insert client: %v", err)
	}

	clientID, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("read client ID: %v", err)
	}

	if _, err := database.Exec(`
		UPDATE clients
		SET status = 'inactive'
		WHERE id = ?
	`, clientID); err != nil {
		t.Fatalf("update client status: %v", err)
	}

	if _, err := database.Exec(`
		INSERT INTO blog_posts (
			title,
			slug,
			excerpt,
			content,
			status,
			published_at
		)
		VALUES (
			'Audited Post',
			'audited-post',
			'A useful excerpt',
			'Article content',
			'draft',
			NULL
		)
	`); err != nil {
		t.Fatalf("insert blog post: %v", err)
	}

	if _, err := database.Exec(`
		UPDATE blog_posts
		SET
			status = 'published',
			published_at = CURRENT_TIMESTAMP
		WHERE slug = 'audited-post'
	`); err != nil {
		t.Fatalf("publish blog post: %v", err)
	}

	rows, err := database.Query(`
		SELECT action, summary, details
		FROM audit_events
		ORDER BY id
	`)
	if err != nil {
		t.Fatalf("query audit events: %v", err)
	}
	defer rows.Close()

	var actions []string
	var summaries []string
	var details []string

	for rows.Next() {
		var action string
		var summary string
		var detail string

		if err := rows.Scan(
			&action,
			&summary,
			&detail,
		); err != nil {
			t.Fatalf("scan audit event: %v", err)
		}

		actions = append(actions, action)
		summaries = append(summaries, summary)
		details = append(details, detail)
	}

	if err := rows.Err(); err != nil {
		t.Fatalf("iterate audit events: %v", err)
	}

	wantActions := []string{
		"client.created",
		"client.status_changed",
		"blog_post.created",
		"blog_post.published",
	}

	if len(actions) != len(wantActions) {
		t.Fatalf(
			"expected %d audit events, got %d (%v)",
			len(wantActions),
			len(actions),
			actions,
		)
	}

	for index, want := range wantActions {
		if actions[index] != want {
			t.Fatalf(
				"expected action %q at index %d, got %q",
				want,
				index,
				actions[index],
			)
		}
	}

	if !strings.Contains(summaries[1], "Audit Client") {
		t.Fatalf(
			"expected client name in status event, got %q",
			summaries[1],
		)
	}

	if details[1] != "active → inactive" {
		t.Fatalf(
			"expected status transition details, got %q",
			details[1],
		)
	}
}
