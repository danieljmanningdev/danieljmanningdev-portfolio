package database

import (
	"context"
	"database/sql"
	"fmt"
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

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &Database{
		SQL: db,
	}, nil
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
