// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package database

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestRunMigrations(t *testing.T) {
	dir := t.TempDir()

	writeMigration(t, dir, "002_second.sql", `
		CREATE TABLE second_table (
			id INTEGER PRIMARY KEY
		);
	`)

	writeMigration(t, dir, "001_first.sql", `
		CREATE TABLE first_table (
			id INTEGER PRIMARY KEY
		);
	`)

	db, err := Open(context.Background(), ":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}

	defer db.Close()

	if err := RunMigrations(db.SQL, dir); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	for _, table := range []string{
		"first_table",
		"second_table",
	} {
		var exists bool

		err := db.SQL.QueryRow(
			`SELECT EXISTS(
				SELECT 1
				FROM sqlite_master
				WHERE type = 'table'
				AND name = ?
			)`,
			table,
		).Scan(&exists)

		if err != nil {
			t.Fatalf("check table %s: %v", table, err)
		}

		if !exists {
			t.Errorf("expected table %s", table)
		}
	}

	var count int

	if err := db.SQL.QueryRow(
		`SELECT COUNT(*) FROM schema_migrations`,
	).Scan(&count); err != nil {
		t.Fatalf("count migrations: %v", err)
	}

	if count != 2 {
		t.Fatalf("expected 2 migrations, got %d", count)
	}

	if err := RunMigrations(db.SQL, dir); err != nil {
		t.Fatalf("second migration run: %v", err)
	}

	if err := db.SQL.QueryRow(
		`SELECT COUNT(*) FROM schema_migrations`,
	).Scan(&count); err != nil {
		t.Fatalf("count migrations after second run: %v", err)
	}

	if count != 2 {
		t.Fatalf(
			"expected 2 migrations after second run, got %d",
			count,
		)
	}
}

func TestFailedMigrationIsNotRecorded(t *testing.T) {
	dir := t.TempDir()

	writeMigration(t, dir, "001_broken.sql", `
		CREATE TABLE broken_table (
	`)

	db, err := Open(context.Background(), ":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}

	defer db.Close()

	if err := RunMigrations(db.SQL, dir); err == nil {
		t.Fatal("expected migration to fail")
	}

	var count int

	if err := db.SQL.QueryRow(
		`SELECT COUNT(*) FROM schema_migrations`,
	).Scan(&count); err != nil {
		t.Fatalf("count migrations: %v", err)
	}

	if count != 0 {
		t.Fatalf("expected 0 recorded migrations, got %d", count)
	}
}

func writeMigration(t *testing.T, dir, name, content string) {
	t.Helper()

	path := filepath.Join(dir, name)

	if err := os.WriteFile(
		path,
		[]byte(content),
		0o644,
	); err != nil {
		t.Fatalf("write migration: %v", err)
	}
}
