package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"
)

func TestAdminSessionRepositoryCreateAndGet(
	t *testing.T,
) {
	db := newAuthRepositoryTestDB(t)

	adminRepo := NewAdminRepository(db)
	sessionRepo := NewAdminSessionRepository(db)

	adminID, err := adminRepo.Create(
		context.Background(),
		"admin@example.com",
		"hash",
		"Admin",
	)
	if err != nil {
		t.Fatalf("create admin: %v", err)
	}

	expiresAt := time.Now().
		UTC().
		Add(24 * time.Hour).
		Truncate(time.Second)

	sessionID, err := sessionRepo.Create(
		context.Background(),
		adminID,
		"token-hash",
		expiresAt,

		time.Now().UTC(),	)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	session, err := sessionRepo.GetByTokenHash(
		context.Background(),
		"token-hash",
	)
	if err != nil {
		t.Fatalf("get session: %v", err)
	}

	if session.ID != sessionID {
		t.Fatalf(
			"expected session ID %d, got %d",
			sessionID,
			session.ID,
		)
	}

	if session.AdminID != adminID {
		t.Fatalf(
			"expected admin ID %d, got %d",
			adminID,
			session.AdminID,
		)
	}

	if session.TokenHash != "token-hash" {
		t.Fatalf(
			"unexpected token hash %q",
			session.TokenHash,
		)
	}
}

func TestAdminSessionRepositoryTouch(
	t *testing.T,
) {
	db := newAuthRepositoryTestDB(t)

	adminRepo := NewAdminRepository(db)
	sessionRepo := NewAdminSessionRepository(db)

	adminID, err := adminRepo.Create(
		context.Background(),
		"admin@example.com",
		"hash",
		"Admin",
	)
	if err != nil {
		t.Fatalf("create admin: %v", err)
	}

	sessionID, err := sessionRepo.Create(
		context.Background(),
		adminID,
		"token-hash",
		time.Now().UTC().Add(time.Hour),

		time.Now().UTC(),	)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	touchedAt := time.Now().
		UTC().
		Add(5 * time.Minute).
		Truncate(time.Second)

	if err := sessionRepo.Touch(
		context.Background(),
		sessionID,
		touchedAt,
	); err != nil {
		t.Fatalf("touch session: %v", err)
	}

	session, err := sessionRepo.GetByTokenHash(
		context.Background(),
		"token-hash",
	)
	if err != nil {
		t.Fatalf("get session: %v", err)
	}

	if session.LastSeenAt.Before(
		touchedAt.Add(-time.Second),
	) {
		t.Fatalf(
			"expected last seen near %s, got %s",
			touchedAt,
			session.LastSeenAt,
		)
	}
}

func TestAdminSessionRepositoryDeleteExpired(
	t *testing.T,
) {
	db := newAuthRepositoryTestDB(t)

	adminRepo := NewAdminRepository(db)
	sessionRepo := NewAdminSessionRepository(db)

	adminID, err := adminRepo.Create(
		context.Background(),
		"admin@example.com",
		"hash",
		"Admin",
	)
	if err != nil {
		t.Fatalf("create admin: %v", err)
	}

	now := time.Now().UTC()

	if _, err := sessionRepo.Create(
		context.Background(),
		adminID,
		"expired-token",
		now.Add(-time.Hour),

		now,	); err != nil {
		t.Fatalf("create expired session: %v", err)
	}

	if _, err := sessionRepo.Create(
		context.Background(),
		adminID,
		"active-token",
		now.Add(time.Hour),

		now,	); err != nil {
		t.Fatalf("create active session: %v", err)
	}

	if err := sessionRepo.DeleteExpired(
		context.Background(),
		now,
	); err != nil {
		t.Fatalf("delete expired sessions: %v", err)
	}

	_, err = sessionRepo.GetByTokenHash(
		context.Background(),
		"expired-token",
	)

	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf(
			"expected expired session to be deleted, got %v",
			err,
		)
	}

	if _, err := sessionRepo.GetByTokenHash(
		context.Background(),
		"active-token",
	); err != nil {
		t.Fatalf(
			"expected active session to remain: %v",
			err,
		)
	}
}

func TestAdminSessionRepositoryCascadeDelete(
	t *testing.T,
) {
	db := newAuthRepositoryTestDB(t)

	adminRepo := NewAdminRepository(db)
	sessionRepo := NewAdminSessionRepository(db)

	adminID, err := adminRepo.Create(
		context.Background(),
		"admin@example.com",
		"hash",
		"Admin",
	)
	if err != nil {
		t.Fatalf("create admin: %v", err)
	}

	if _, err := sessionRepo.Create(
		context.Background(),
		adminID,
		"token-hash",
		time.Now().UTC().Add(time.Hour),

		time.Now().UTC(),	); err != nil {
		t.Fatalf("create session: %v", err)
	}

	if _, err := db.Exec(
		"DELETE FROM admins WHERE id = ?",
		adminID,
	); err != nil {
		t.Fatalf("delete admin: %v", err)
	}

	_, err = sessionRepo.GetByTokenHash(
		context.Background(),
		"token-hash",
	)

	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf(
			"expected session cascade delete, got %v",
			err,
		)
	}
}

func TestAdminSessionRepositoryMissingTouch(
	t *testing.T,
) {
	db := newAuthRepositoryTestDB(t)

	sessionRepo := NewAdminSessionRepository(db)

	err := sessionRepo.Touch(
		context.Background(),
		999,
		time.Now().UTC(),
	)

	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf(
			"expected sql.ErrNoRows, got %v",
			err,
		)
	}
}
