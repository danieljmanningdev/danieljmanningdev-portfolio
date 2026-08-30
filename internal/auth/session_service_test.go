// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
//
// Description
// -----------
// Authentication and session handling for the application.
//
// Security
// --------
// Changes to this package may affect authentication, session integrity,
// credential handling, and access control.
// -----------------------------------------------------------------------------
package auth

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/database"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/repository"
)

func newSessionServiceTest(
	t *testing.T,
) (
	*SessionService,
	*repository.AdminRepository,
	*repository.AdminSessionRepository,
	*sql.DB,
) {
	t.Helper()

	db, err := database.Open(
		context.Background(),
		":memory:",
	)
	if err != nil {
		t.Fatalf("open test database: %v", err)
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

	adminRepository := repository.NewAdminRepository(
		db.SQL,
	)

	sessionRepository := repository.NewAdminSessionRepository(
		db.SQL,
	)

	service := NewSessionService(
		adminRepository,
		sessionRepository,
	)

	return service,
		adminRepository,
		sessionRepository,
		db.SQL
}

func createSessionTestAdmin(
	t *testing.T,
	repository *repository.AdminRepository,
) int64 {
	t.Helper()

	id, err := repository.Create(
		context.Background(),
		"admin@example.com",
		"password-hash",
		"Daniel Manning",
	)
	if err != nil {
		t.Fatalf("create admin: %v", err)
	}

	return id
}

func TestSessionServiceCreateAndAuthenticate(
	t *testing.T,
) {
	service, adminRepository, sessionRepository, _ :=
		newSessionServiceTest(t)

	adminID := createSessionTestAdmin(
		t,
		adminRepository,
	)

	now := time.Date(
		2026,
		time.August,
		12,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	service.now = func() time.Time {
		return now
	}

	rawToken, expiresAt, err := service.CreateSession(
		context.Background(),
		adminID,
	)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	if rawToken == "" {
		t.Fatal("expected raw session token")
	}

	expectedExpiry := now.Add(
		adminSessionLifetime,
	)

	if !expiresAt.Equal(expectedExpiry) {
		t.Fatalf(
			"expected expiry %s, got %s",
			expectedExpiry,
			expiresAt,
		)
	}

	session, err := sessionRepository.GetByTokenHash(
		context.Background(),
		HashSessionToken(rawToken),
	)
	if err != nil {
		t.Fatalf("get stored session: %v", err)
	}

	if session.TokenHash == rawToken {
		t.Fatal(
			"raw session token must not be stored in the database",
		)
	}

	admin, err := service.Authenticate(
		context.Background(),
		rawToken,
	)
	if err != nil {
		t.Fatalf("authenticate session: %v", err)
	}

	if admin.ID != adminID {
		t.Fatalf(
			"expected admin ID %d, got %d",
			adminID,
			admin.ID,
		)
	}

	if admin.Email != "admin@example.com" {
		t.Fatalf(
			"unexpected admin email %q",
			admin.Email,
		)
	}
}

func TestSessionServiceRejectsUnknownToken(
	t *testing.T,
) {
	service, _, _, _ := newSessionServiceTest(t)

	_, err := service.Authenticate(
		context.Background(),
		"unknown-token",
	)

	if !errors.Is(err, ErrInvalidSession) {
		t.Fatalf(
			"expected ErrInvalidSession, got %v",
			err,
		)
	}
}

func TestSessionServiceRejectsEmptyToken(
	t *testing.T,
) {
	service, _, _, _ := newSessionServiceTest(t)

	_, err := service.Authenticate(
		context.Background(),
		"",
	)

	if !errors.Is(err, ErrInvalidSession) {
		t.Fatalf(
			"expected ErrInvalidSession, got %v",
			err,
		)
	}
}

func TestSessionServiceRejectsAbsoluteExpiry(
	t *testing.T,
) {
	service, adminRepository, sessionRepository, _ :=
		newSessionServiceTest(t)

	adminID := createSessionTestAdmin(
		t,
		adminRepository,
	)

	start := time.Date(
		2026,
		time.August,
		12,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	service.now = func() time.Time {
		return start
	}

	rawToken, _, err := service.CreateSession(
		context.Background(),
		adminID,
	)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	/*
		Prevent the idle timeout from becoming the reason this
		session is rejected. We want this test to specifically
		exercise the absolute expiry.
	*/
	session, err := sessionRepository.GetByTokenHash(
		context.Background(),
		HashSessionToken(rawToken),
	)
	if err != nil {
		t.Fatalf("get session: %v", err)
	}

	nearExpiry := start.
		Add(adminSessionLifetime).
		Add(-time.Minute)

	if err := sessionRepository.Touch(
		context.Background(),
		session.ID,
		nearExpiry,
	); err != nil {
		t.Fatalf("touch session: %v", err)
	}

	service.now = func() time.Time {
		return start.Add(adminSessionLifetime)
	}

	_, err = service.Authenticate(
		context.Background(),
		rawToken,
	)

	if !errors.Is(err, ErrExpiredSession) {
		t.Fatalf(
			"expected ErrExpiredSession, got %v",
			err,
		)
	}

	_, err = sessionRepository.GetByTokenHash(
		context.Background(),
		HashSessionToken(rawToken),
	)

	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf(
			"expected expired session to be deleted, got %v",
			err,
		)
	}
}

func TestSessionServiceRejectsIdleExpiry(
	t *testing.T,
) {
	service, adminRepository, sessionRepository, _ :=
		newSessionServiceTest(t)

	adminID := createSessionTestAdmin(
		t,
		adminRepository,
	)

	start := time.Date(
		2026,
		time.August,
		12,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	service.now = func() time.Time {
		return start
	}

	rawToken, _, err := service.CreateSession(
		context.Background(),
		adminID,
	)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	service.now = func() time.Time {
		return start.Add(
			adminSessionIdleTimeout,
		)
	}

	_, err = service.Authenticate(
		context.Background(),
		rawToken,
	)

	if !errors.Is(err, ErrExpiredSession) {
		t.Fatalf(
			"expected ErrExpiredSession, got %v",
			err,
		)
	}

	_, err = sessionRepository.GetByTokenHash(
		context.Background(),
		HashSessionToken(rawToken),
	)

	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf(
			"expected idle session to be deleted, got %v",
			err,
		)
	}
}

func TestSessionServiceTouchesActiveSession(
	t *testing.T,
) {
	service, adminRepository, sessionRepository, _ :=
		newSessionServiceTest(t)

	adminID := createSessionTestAdmin(
		t,
		adminRepository,
	)

	start := time.Date(
		2026,
		time.August,
		12,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	service.now = func() time.Time {
		return start
	}

	rawToken, _, err := service.CreateSession(
		context.Background(),
		adminID,
	)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	sessionBefore, err := sessionRepository.GetByTokenHash(
		context.Background(),
		HashSessionToken(rawToken),
	)
	if err != nil {
		t.Fatalf("get session before touch: %v", err)
	}

	later := start.Add(
		adminSessionTouchInterval +
			time.Minute,
	)

	service.now = func() time.Time {
		return later
	}

	if _, err := service.Authenticate(
		context.Background(),
		rawToken,
	); err != nil {
		t.Fatalf("authenticate session: %v", err)
	}

	sessionAfter, err := sessionRepository.GetByTokenHash(
		context.Background(),
		HashSessionToken(rawToken),
	)
	if err != nil {
		t.Fatalf("get session after touch: %v", err)
	}

	if !sessionAfter.LastSeenAt.After(
		sessionBefore.LastSeenAt,
	) {
		t.Fatalf(
			"expected last seen to advance from %s to %s",
			sessionBefore.LastSeenAt,
			sessionAfter.LastSeenAt,
		)
	}
}

func TestSessionServiceDoesNotTouchTooFrequently(
	t *testing.T,
) {
	service, adminRepository, sessionRepository, _ :=
		newSessionServiceTest(t)

	adminID := createSessionTestAdmin(
		t,
		adminRepository,
	)

	start := time.Date(
		2026,
		time.August,
		12,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	service.now = func() time.Time {
		return start
	}

	rawToken, _, err := service.CreateSession(
		context.Background(),
		adminID,
	)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	before, err := sessionRepository.GetByTokenHash(
		context.Background(),
		HashSessionToken(rawToken),
	)
	if err != nil {
		t.Fatalf("get session before: %v", err)
	}

	service.now = func() time.Time {
		return start.Add(time.Minute)
	}

	if _, err := service.Authenticate(
		context.Background(),
		rawToken,
	); err != nil {
		t.Fatalf("authenticate session: %v", err)
	}

	after, err := sessionRepository.GetByTokenHash(
		context.Background(),
		HashSessionToken(rawToken),
	)
	if err != nil {
		t.Fatalf("get session after: %v", err)
	}

	if !after.LastSeenAt.Equal(
		before.LastSeenAt,
	) {
		t.Fatalf(
			"expected last seen to remain %s, got %s",
			before.LastSeenAt,
			after.LastSeenAt,
		)
	}
}

func TestSessionServiceRejectsInactiveAdmin(
	t *testing.T,
) {
	service, adminRepository, sessionRepository, _ :=
		newSessionServiceTest(t)

	adminID := createSessionTestAdmin(
		t,
		adminRepository,
	)

	rawToken, _, err := service.CreateSession(
		context.Background(),
		adminID,
	)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	if err := adminRepository.SetActive(
		context.Background(),
		adminID,
		false,
	); err != nil {
		t.Fatalf("deactivate admin: %v", err)
	}

	_, err = service.Authenticate(
		context.Background(),
		rawToken,
	)

	if !errors.Is(err, ErrInactiveAdmin) {
		t.Fatalf(
			"expected ErrInactiveAdmin, got %v",
			err,
		)
	}

	_, err = sessionRepository.GetByTokenHash(
		context.Background(),
		HashSessionToken(rawToken),
	)

	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf(
			"expected inactive admin session to be revoked, got %v",
			err,
		)
	}
}

func TestSessionServiceCannotCreateForInactiveAdmin(
	t *testing.T,
) {
	service, adminRepository, _, _ :=
		newSessionServiceTest(t)

	adminID := createSessionTestAdmin(
		t,
		adminRepository,
	)

	if err := adminRepository.SetActive(
		context.Background(),
		adminID,
		false,
	); err != nil {
		t.Fatalf("deactivate admin: %v", err)
	}

	_, _, err := service.CreateSession(
		context.Background(),
		adminID,
	)

	if !errors.Is(err, ErrInactiveAdmin) {
		t.Fatalf(
			"expected ErrInactiveAdmin, got %v",
			err,
		)
	}
}

func TestSessionServiceRevokeSession(
	t *testing.T,
) {
	service, adminRepository, sessionRepository, _ :=
		newSessionServiceTest(t)

	adminID := createSessionTestAdmin(
		t,
		adminRepository,
	)

	rawToken, _, err := service.CreateSession(
		context.Background(),
		adminID,
	)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	if err := service.RevokeSession(
		context.Background(),
		rawToken,
	); err != nil {
		t.Fatalf("revoke session: %v", err)
	}

	_, err = sessionRepository.GetByTokenHash(
		context.Background(),
		HashSessionToken(rawToken),
	)

	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf(
			"expected revoked session to be deleted, got %v",
			err,
		)
	}
}

func TestSessionServiceRevokeAllAdminSessions(
	t *testing.T,
) {
	service, adminRepository, sessionRepository, _ :=
		newSessionServiceTest(t)

	adminID := createSessionTestAdmin(
		t,
		adminRepository,
	)

	firstToken, _, err := service.CreateSession(
		context.Background(),
		adminID,
	)
	if err != nil {
		t.Fatalf("create first session: %v", err)
	}

	secondToken, _, err := service.CreateSession(
		context.Background(),
		adminID,
	)
	if err != nil {
		t.Fatalf("create second session: %v", err)
	}

	if err := service.RevokeAllAdminSessions(
		context.Background(),
		adminID,
	); err != nil {
		t.Fatalf(
			"revoke all admin sessions: %v",
			err,
		)
	}

	for _, rawToken := range []string{
		firstToken,
		secondToken,
	} {
		_, err := sessionRepository.GetByTokenHash(
			context.Background(),
			HashSessionToken(rawToken),
		)

		if !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf(
				"expected session to be deleted, got %v",
				err,
			)
		}
	}
}
