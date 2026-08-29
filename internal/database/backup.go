package database

import (
	"context"
	"database/sql"
	"fmt"
	"io"
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

	temporaryPath := destinationAbsolute + ".tmp"
	_ = os.Remove(temporaryPath)
	defer func() {
		_ = os.Remove(temporaryPath)
	}()

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

	if err := os.Rename(
		temporaryPath,
		destinationAbsolute,
	); err != nil {
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

	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("inspect backup file: %w", err)
	}

	dsn, err := sqliteDSN(path)
	if err != nil {
		return err
	}

	database, err := sql.Open("sqlite", dsn)
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

	temporaryPath := databaseAbsolute + ".restore.tmp"
	_ = os.Remove(temporaryPath)
	defer func() {
		_ = os.Remove(temporaryPath)
	}()

	if err := copyFile(
		backupAbsolute,
		temporaryPath,
	); err != nil {
		return err
	}

	if err := VerifyBackup(ctx, temporaryPath); err != nil {
		return err
	}

	if err := os.Rename(
		temporaryPath,
		databaseAbsolute,
	); err != nil {
		return fmt.Errorf("replace database atomically: %w", err)
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
