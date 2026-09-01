// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package database

import (
	"context"
	"testing"
)

func TestOpenEnablesForeignKeyEnforcement(t *testing.T) {
	ctx := context.Background()

	db, err := Open(ctx, ":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer db.Close()

	var enabled int

	if err := db.SQL.QueryRowContext(
		ctx,
		"PRAGMA foreign_keys",
	).Scan(&enabled); err != nil {
		t.Fatalf(
			"query foreign key setting: %v",
			err,
		)
	}

	if enabled != 1 {
		t.Fatalf(
			"expected foreign keys enabled, got %d",
			enabled,
		)
	}

	_, err = db.SQL.ExecContext(ctx, `
		CREATE TABLE parents (
			id INTEGER PRIMARY KEY
		);

		CREATE TABLE children (
			id INTEGER PRIMARY KEY,
			parent_id INTEGER NOT NULL,
			FOREIGN KEY (parent_id)
				REFERENCES parents(id)
		);
	`)
	if err != nil {
		t.Fatalf(
			"create foreign key test tables: %v",
			err,
		)
	}

	_, err = db.SQL.ExecContext(
		ctx,
		`
			INSERT INTO children (
				id,
				parent_id
			)
			VALUES (?, ?)
		`,
		1,
		999,
	)

	if err == nil {
		t.Fatal(
			"expected invalid foreign key insert to fail",
		)
	}
}
