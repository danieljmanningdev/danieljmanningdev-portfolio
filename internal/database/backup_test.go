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

func TestBackupAndRestore(
	t *testing.T,
) {
	ctx := context.Background()
	temporaryDirectory := t.TempDir()

	databasePath := filepath.Join(
		temporaryDirectory,
		"source.db",
	)
	backupPath := filepath.Join(
		temporaryDirectory,
		"backups",
		"source-backup.db",
	)
	restoredPath := filepath.Join(
		temporaryDirectory,
		"restored.db",
	)

	database, err := Open(ctx, databasePath)
	if err != nil {
		t.Fatalf("open source database: %v", err)
	}

	if _, err := database.SQL.Exec(`
		CREATE TABLE notes (
			id INTEGER PRIMARY KEY,
			body TEXT NOT NULL
		);

		INSERT INTO notes (body)
		VALUES ('backup test');
	`); err != nil {
		_ = database.Close()
		t.Fatalf("seed source database: %v", err)
	}

	if err := database.Close(); err != nil {
		t.Fatalf("close source database: %v", err)
	}

	if err := Backup(
		ctx,
		databasePath,
		backupPath,
	); err != nil {
		t.Fatalf("create backup: %v", err)
	}

	if err := VerifyBackup(ctx, backupPath); err != nil {
		t.Fatalf("verify backup: %v", err)
	}

	backupInfo, err := os.Stat(backupPath)
	if err != nil {
		t.Fatalf("stat backup: %v", err)
	}

	if backupInfo.Mode().Perm() != 0o600 {
		t.Fatalf(
			"expected backup permissions 0600, got %o",
			backupInfo.Mode().Perm(),
		)
	}

	if err := Restore(
		ctx,
		backupPath,
		restoredPath,
		false,
	); err != nil {
		t.Fatalf("restore backup: %v", err)
	}

	restored, err := Open(ctx, restoredPath)
	if err != nil {
		t.Fatalf("open restored database: %v", err)
	}
	defer func() {
		_ = restored.Close()
	}()

	var body string
	if err := restored.SQL.QueryRow(
		"SELECT body FROM notes WHERE id = 1",
	).Scan(&body); err != nil {
		t.Fatalf("read restored data: %v", err)
	}

	if body != "backup test" {
		t.Fatalf(
			"expected restored body %q, got %q",
			"backup test",
			body,
		)
	}
}

func TestBackupRefusesToOverwriteExistingFile(
	t *testing.T,
) {
	ctx := context.Background()
	temporaryDirectory := t.TempDir()
	databasePath := filepath.Join(temporaryDirectory, "source.db")
	backupPath := filepath.Join(temporaryDirectory, "backup.db")

	database, err := Open(ctx, databasePath)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := database.Close(); err != nil {
		t.Fatalf("close database: %v", err)
	}

	if err := os.WriteFile(
		backupPath,
		[]byte("existing"),
		0o600,
	); err != nil {
		t.Fatalf("create existing backup: %v", err)
	}

	if err := Backup(
		ctx,
		databasePath,
		backupPath,
	); err == nil {
		t.Fatal("expected overwrite protection error")
	}
}

func TestRestoreRequiresForceForExistingDestination(
	t *testing.T,
) {
	ctx := context.Background()
	temporaryDirectory := t.TempDir()
	sourcePath := filepath.Join(temporaryDirectory, "source.db")
	backupPath := filepath.Join(temporaryDirectory, "backup.db")
	destinationPath := filepath.Join(temporaryDirectory, "destination.db")

	database, err := Open(ctx, sourcePath)
	if err != nil {
		t.Fatalf("open source database: %v", err)
	}
	if err := database.Close(); err != nil {
		t.Fatalf("close source database: %v", err)
	}

	if err := Backup(ctx, sourcePath, backupPath); err != nil {
		t.Fatalf("create backup: %v", err)
	}

	if err := os.WriteFile(
		destinationPath,
		[]byte("existing"),
		0o600,
	); err != nil {
		t.Fatalf("create destination: %v", err)
	}

	if err := Restore(
		ctx,
		backupPath,
		destinationPath,
		false,
	); err == nil {
		t.Fatal("expected restore overwrite protection error")
	}
}

func TestVerifyBackupRejectsInvalidSQLiteFile(
	t *testing.T,
) {
	path := filepath.Join(t.TempDir(), "invalid.db")

	if err := os.WriteFile(
		path,
		[]byte("not sqlite"),
		0o600,
	); err != nil {
		t.Fatalf("write invalid backup: %v", err)
	}

	if err := VerifyBackup(
		context.Background(),
		path,
	); err == nil {
		t.Fatal("expected invalid SQLite backup error")
	}
}
