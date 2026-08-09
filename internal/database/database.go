package database

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Database struct {
	SQL *sql.DB
}

func Open(ctx context.Context, path string) (*Database, error) {
	if path != ":memory:" {
		dir := filepath.Dir(path)

		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("create database directory: %w", err)
		}
	}

	dsn, err := sqliteDSN(path)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if path == ":memory:" {
		// SQLite in-memory databases exist per physical connection.
		// Keeping one connection ensures every query sees the same database.
		db.SetMaxOpenConns(1)
		db.SetMaxIdleConns(1)
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	var foreignKeysEnabled int

	if err := db.QueryRowContext(
		ctx,
		"PRAGMA foreign_keys",
	).Scan(&foreignKeysEnabled); err != nil {
		_ = db.Close()

		return nil, fmt.Errorf(
			"check foreign key enforcement: %w",
			err,
		)
	}

	if foreignKeysEnabled != 1 {
		_ = db.Close()
		return nil, fmt.Errorf("foreign key enforcement is disabled")
	}

	return &Database{
		SQL: db,
	}, nil
}

func sqliteDSN(path string) (string, error) {
	if path == ":memory:" {
		return "file::memory:?_foreign_keys=on", nil
	}

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolve database path: %w", err)
	}

	dsn := &url.URL{
		Scheme: "file",
		Path:   filepath.ToSlash(absolutePath),
	}

	query := dsn.Query()
	query.Set("_foreign_keys", "on")
	dsn.RawQuery = query.Encode()

	return dsn.String(), nil
}

func (db *Database) Ping(ctx context.Context) error {
	if db == nil || db.SQL == nil {
		return fmt.Errorf("database is not initialized")
	}

	return db.SQL.PingContext(ctx)
}

func (db *Database) Close() error {
	if db == nil || db.SQL == nil {
		return nil
	}

	return db.SQL.Close()
}
