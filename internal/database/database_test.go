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

func TestOpenAndPing(t *testing.T) {
	ctx := context.Background()

	db, err := Open(ctx, ":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}

	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		t.Fatalf("ping database: %v", err)
	}
}

func TestClose(t *testing.T) {
	ctx := context.Background()

	db, err := Open(ctx, ":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}

	if err := db.Close(); err != nil {
		t.Fatalf("close database: %v", err)
	}
}

func TestOpenCreatesDatabaseDirectory(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/nested/app.db"

	db, err := Open(context.Background(), path)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}

	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		t.Fatalf("ping database: %v", err)
	}
}
