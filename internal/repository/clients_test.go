// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package repository

import (
	"context"
	"testing"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/database"
)

func TestClientRepositoryUpdateRefreshesUpdatedAt(
	t *testing.T,
) {
	ctx := context.Background()

	db, err := database.Open(ctx, ":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer db.Close()

	_, err = db.SQL.ExecContext(ctx, `
		CREATE TABLE clients (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			company TEXT,
			status TEXT NOT NULL DEFAULT 'active',
			notes TEXT,
			created_at DATETIME NOT NULL
				DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL
				DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("create clients table: %v", err)
	}

	result, err := db.SQL.ExecContext(ctx, `
		INSERT INTO clients (
			name,
			email,
			company,
			status,
			notes,
			created_at,
			updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`,
		"Original Client",
		"original@example.com",
		"Original Company",
		"active",
		"Original notes",
		"2000-01-01 00:00:00",
		"2000-01-01 00:00:00",
	)
	if err != nil {
		t.Fatalf("insert client: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("get inserted client id: %v", err)
	}

	clientRepository := NewClientRepository(db.SQL)

	if err := clientRepository.Update(
		ctx,
		id,
		"Updated Client",
		"updated@example.com",
		"Updated Company",
		"inactive",
		"Updated notes",
	); err != nil {
		t.Fatalf("update client: %v", err)
	}

	var timestampChanged int

	if err := db.SQL.QueryRowContext(ctx, `
		SELECT CASE
			WHEN updated_at > ? THEN 1
			ELSE 0
		END
		FROM clients
		WHERE id = ?
	`,
		"2000-01-01 00:00:00",
		id,
	).Scan(&timestampChanged); err != nil {
		t.Fatalf("check updated_at: %v", err)
	}

	if timestampChanged != 1 {
		t.Fatal(
			"expected updated_at to change after client update",
		)
	}
}
