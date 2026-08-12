package repository

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/database"
)

func newAuthRepositoryTestDB(
	t *testing.T,
) *sql.DB {
	t.Helper()

	db, err := database.Open(
		context.Background(),
		":memory:",
	)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve test file path")
	}

	migrationsDir := filepath.Join(
		filepath.Dir(filename),
		"..",
		"..",
		"migrations",
	)

	if err := database.RunMigrations(
		db.SQL,
		migrationsDir,
	); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	return db.SQL
}

func TestAdminRepositoryCreateAndGetByEmail(
	t *testing.T,
) {
	db := newAuthRepositoryTestDB(t)

	repo := NewAdminRepository(db)

	id, err := repo.Create(
		context.Background(),
		"  ADMIN@Example.com  ",
		"hashed-password",
		" Daniel Manning ",
	)
	if err != nil {
		t.Fatalf("create admin: %v", err)
	}

	admin, err := repo.GetByEmail(
		context.Background(),
		"admin@example.com",
	)
	if err != nil {
		t.Fatalf("get admin by email: %v", err)
	}

	if admin.ID != id {
		t.Fatalf(
			"expected ID %d, got %d",
			id,
			admin.ID,
		)
	}

	if admin.Email != "admin@example.com" {
		t.Fatalf(
			"unexpected email %q",
			admin.Email,
		)
	}

	if admin.PasswordHash != "hashed-password" {
		t.Fatalf(
			"unexpected password hash %q",
			admin.PasswordHash,
		)
	}

	if admin.DisplayName != "Daniel Manning" {
		t.Fatalf(
			"unexpected display name %q",
			admin.DisplayName,
		)
	}

	if !admin.IsActive {
		t.Fatal("expected admin to be active")
	}
}

func TestAdminRepositoryEmailIsUnique(
	t *testing.T,
) {
	db := newAuthRepositoryTestDB(t)

	repo := NewAdminRepository(db)

	_, err := repo.Create(
		context.Background(),
		"admin@example.com",
		"hash-one",
		"Admin",
	)
	if err != nil {
		t.Fatalf("create first admin: %v", err)
	}

	_, err = repo.Create(
		context.Background(),
		"ADMIN@example.com",
		"hash-two",
		"Second Admin",
	)
	if err == nil {
		t.Fatal("expected duplicate email error")
	}
}

func TestAdminRepositorySetActive(
	t *testing.T,
) {
	db := newAuthRepositoryTestDB(t)

	repo := NewAdminRepository(db)

	id, err := repo.Create(
		context.Background(),
		"admin@example.com",
		"hash",
		"Admin",
	)
	if err != nil {
		t.Fatalf("create admin: %v", err)
	}

	if err := repo.SetActive(
		context.Background(),
		id,
		false,
	); err != nil {
		t.Fatalf("deactivate admin: %v", err)
	}

	admin, err := repo.GetByID(
		context.Background(),
		id,
	)
	if err != nil {
		t.Fatalf("get admin: %v", err)
	}

	if admin.IsActive {
		t.Fatal("expected inactive admin")
	}
}

func TestAdminRepositoryMissingAdmin(
	t *testing.T,
) {
	db := newAuthRepositoryTestDB(t)

	repo := NewAdminRepository(db)

	_, err := repo.GetByID(
		context.Background(),
		999,
	)

	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf(
			"expected sql.ErrNoRows, got %v",
			err,
		)
	}
}
