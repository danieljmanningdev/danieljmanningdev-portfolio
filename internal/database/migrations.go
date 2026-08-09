package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Migration struct {
	Version int
	Name    string
	Path    string
}

func RunMigrations(db *sql.DB, directory string) error {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`); err != nil {
		return fmt.Errorf("create schema_migrations table: %w", err)
	}

	migrations, err := discoverMigrations(directory)
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		applied, err := migrationApplied(db, migration.Version)
		if err != nil {
			return err
		}

		if applied {
			continue
		}

		if err := applyMigration(db, migration); err != nil {
			return err
		}
	}

	return nil
}

func discoverMigrations(directory string) ([]Migration, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("read migrations directory: %w", err)
	}

	var migrations []Migration

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		parts := strings.SplitN(entry.Name(), "_", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid migration filename %q", entry.Name())
		}

		version, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf(
				"invalid migration version in %q: %w",
				entry.Name(),
				err,
			)
		}

		migrations = append(migrations, Migration{
			Version: version,
			Name:    strings.TrimSuffix(parts[1], ".sql"),
			Path:    filepath.Join(directory, entry.Name()),
		})
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

func migrationApplied(db *sql.DB, version int) (bool, error) {
	var exists bool

	err := db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM schema_migrations
			WHERE version = ?
		)`,
		version,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf(
			"check migration %d: %w",
			version,
			err,
		)
	}

	return exists, nil
}

func applyMigration(db *sql.DB, migration Migration) error {
	content, err := os.ReadFile(migration.Path)
	if err != nil {
		return fmt.Errorf(
			"read migration %s: %w",
			migration.Path,
			err,
		)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf(
			"begin migration %d: %w",
			migration.Version,
			err,
		)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	if _, err := tx.Exec(string(content)); err != nil {
		return fmt.Errorf(
			"execute migration %d: %w",
			migration.Version,
			err,
		)
	}

	if _, err := tx.Exec(
		`INSERT INTO schema_migrations (version, name)
		 VALUES (?, ?)`,
		migration.Version,
		migration.Name,
	); err != nil {
		return fmt.Errorf(
			"record migration %d: %w",
			migration.Version,
			err,
		)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf(
			"commit migration %d: %w",
			migration.Version,
			err,
		)
	}

	return nil
}
