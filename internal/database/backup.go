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
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
)

func Backup(
	ctx context.Context,
	databasePath string,
	destinationPath string,
) error {
	if databasePath == "" || databasePath == ":memory:" {
		return fmt.Errorf("backup requires a file-backed database")
	}

	info, err := os.Stat(databasePath)
	if err != nil {
		return fmt.Errorf("inspect source database: %w", err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("source database must be a regular file")
	}

	if destinationPath == "" {
		return fmt.Errorf("backup destination is required")
	}

	sourceAbsolute, err := filepath.Abs(databasePath)
	if err != nil {
		return fmt.Errorf("resolve source database path: %w", err)
	}

	destinationAbsolute, err := filepath.Abs(destinationPath)
	if err != nil {
		return fmt.Errorf("resolve backup destination path: %w", err)
	}

	if sourceAbsolute == destinationAbsolute {
		return fmt.Errorf("backup destination must differ from database path")
	}

	if err := os.MkdirAll(
		filepath.Dir(destinationAbsolute),
		0o750,
	); err != nil {
		return fmt.Errorf("create backup directory: %w", err)
	}

	if _, err := os.Stat(destinationAbsolute); err == nil {
		return fmt.Errorf("backup destination already exists")
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("inspect backup destination: %w", err)
	}

	temporaryDirectory, err := os.MkdirTemp(filepath.Dir(destinationAbsolute), ".backup-*")
	if err != nil {
		return fmt.Errorf("create private backup workspace: %w", err)
	}
	defer func() { _ = os.RemoveAll(temporaryDirectory) }()
	temporaryPath := filepath.Join(temporaryDirectory, "snapshot.db")

	database, err := Open(ctx, sourceAbsolute)
	if err != nil {
		return err
	}
	defer func() {
		_ = database.Close()
	}()

	if _, err := database.SQL.ExecContext(
		ctx,
		"VACUUM INTO ?",
		temporaryPath,
	); err != nil {
		return fmt.Errorf("create consistent SQLite backup: %w", err)
	}

	if err := VerifyBackup(ctx, temporaryPath); err != nil {
		return err
	}

	if err := os.Chmod(temporaryPath, 0o600); err != nil {
		return fmt.Errorf("secure backup permissions: %w", err)
	}

	// A hard link publishes the completed file atomically without replacing a
	// destination created concurrently after the initial existence check.
	if err := os.Link(temporaryPath, destinationAbsolute); err != nil {
		return fmt.Errorf("publish backup atomically: %w", err)
	}

	return nil
}

func VerifyBackup(
	ctx context.Context,
	path string,
) error {
	if path == "" || path == ":memory:" {
		return fmt.Errorf("backup path is required")
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("inspect backup file: %w", err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("backup must be a regular file")
	}

	dsn, err := sqliteDSN(path)
	if err != nil {
		return err
	}

	uri, err := url.Parse(dsn)
	if err != nil {
		return fmt.Errorf("parse backup URI: %w", err)
	}
	query := uri.Query()
	query.Set("mode", "ro")
	uri.RawQuery = query.Encode()

	database, err := sql.Open("sqlite", uri.String())
	if err != nil {
		return fmt.Errorf("open backup for verification: %w", err)
	}
	defer func() {
		_ = database.Close()
	}()

	var result string
	if err := database.QueryRowContext(
		ctx,
		"PRAGMA quick_check",
	).Scan(&result); err != nil {
		return fmt.Errorf("run SQLite backup integrity check: %w", err)
	}

	if result != "ok" {
		return fmt.Errorf("SQLite backup integrity check failed: %s", result)
	}

	return nil
}

func Restore(
	ctx context.Context,
	backupPath string,
	databasePath string,
	force bool,
) error {
	if backupPath == "" {
		return fmt.Errorf("backup path is required")
	}

	if databasePath == "" || databasePath == ":memory:" {
		return fmt.Errorf("restore requires a file-backed destination")
	}

	if err := VerifyBackup(ctx, backupPath); err != nil {
		return err
	}

	backupAbsolute, err := filepath.Abs(backupPath)
	if err != nil {
		return fmt.Errorf("resolve backup path: %w", err)
	}

	databaseAbsolute, err := filepath.Abs(databasePath)
	if err != nil {
		return fmt.Errorf("resolve database path: %w", err)
	}

	if backupAbsolute == databaseAbsolute {
		return fmt.Errorf("backup and destination paths must differ")
	}

	if _, err := os.Stat(databaseAbsolute); err == nil && !force {
		return fmt.Errorf("destination database already exists; use force to replace it")
	} else if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("inspect destination database: %w", err)
	}

	if err := os.MkdirAll(
		filepath.Dir(databaseAbsolute),
		0o750,
	); err != nil {
		return fmt.Errorf("create database directory: %w", err)
	}

	// Do not replace a database while SQLite may have live recovery state.
	// The operator must stop the application and checkpoint/recover it first.
	for _, suffix := range []string{"-wal", "-shm", "-journal"} {
		if _, err := os.Stat(databaseAbsolute + suffix); err == nil {
			return fmt.Errorf("destination has SQLite sidecar %s; stop and checkpoint the application before restoring", suffix)
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("inspect SQLite sidecar: %w", err)
		}
	}
	temporaryDirectory, err := os.MkdirTemp(filepath.Dir(databaseAbsolute), ".restore-*")
	if err != nil {
		return fmt.Errorf("create private restore workspace: %w", err)
	}
	defer func() { _ = os.RemoveAll(temporaryDirectory) }()
	temporaryPath := filepath.Join(temporaryDirectory, "snapshot.db")

	if err := copyFile(
		backupAbsolute,
		temporaryPath,
	); err != nil {
		return err
	}

	if err := VerifyBackup(ctx, temporaryPath); err != nil {
		return err
	}

	if force {
		err = os.Rename(temporaryPath, databaseAbsolute)
	} else {
		err = os.Link(temporaryPath, databaseAbsolute)
	}
	if err != nil {
		return fmt.Errorf("publish restored database atomically: %w", err)
	}

	return nil
}

func copyFile(sourcePath string, destinationPath string) error {
	source, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("open backup source: %w", err)
	}
	defer func() {
		_ = source.Close()
	}()

	destination, err := os.OpenFile(
		destinationPath,
		os.O_CREATE|os.O_EXCL|os.O_WRONLY,
		0o600,
	)
	if err != nil {
		return fmt.Errorf("create restore file: %w", err)
	}

	copySucceeded := false
	defer func() {
		_ = destination.Close()
		if !copySucceeded {
			_ = os.Remove(destinationPath)
		}
	}()

	if _, err := io.Copy(destination, source); err != nil {
		return fmt.Errorf("copy backup: %w", err)
	}

	if err := destination.Sync(); err != nil {
		return fmt.Errorf("sync restored database: %w", err)
	}

	if err := destination.Close(); err != nil {
		return fmt.Errorf("close restored database: %w", err)
	}

	copySucceeded = true
	return nil
}
