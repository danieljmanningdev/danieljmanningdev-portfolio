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
	"fmt"
	"testing"
	"time"
)

func newLoginLimiterTest(
	now time.Time,
) *LoginLimiter {
	limiter := NewLoginLimiter()

	limiter.now = func() time.Time {
		return now
	}

	return limiter
}

func TestLoginLimiterInitiallyAllows(
	t *testing.T,
) {
	now := time.Date(
		2026,
		time.August,
		13,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	limiter := newLoginLimiterTest(now)

	allowed, retryAfter :=
		limiter.Check(
			"192.0.2.1",
			"admin@example.com",
		)

	if !allowed {
		t.Fatal(
			"expected initial login attempt to be allowed",
		)
	}

	if retryAfter != 0 {
		t.Fatalf(
			"expected zero retry delay, got %s",
			retryAfter,
		)
	}
}

func TestLoginLimiterBlocksCredentialAfterFailures(
	t *testing.T,
) {
	now := time.Date(
		2026,
		time.August,
		13,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	limiter := newLoginLimiterTest(now)

	for range maxLoginFailuresPerCredential {
		limiter.RecordFailure(
			"192.0.2.1",
			"admin@example.com",
		)
	}

	allowed, retryAfter :=
		limiter.Check(
			"192.0.2.1",
			"admin@example.com",
		)

	if allowed {
		t.Fatal(
			"expected credential to be throttled",
		)
	}

	if retryAfter != loginFailureWindow {
		t.Fatalf(
			"expected retry after %s, got %s",
			loginFailureWindow,
			retryAfter,
		)
	}
}

func TestLoginLimiterAllowsCredentialAfterWindow(
	t *testing.T,
) {
	start := time.Date(
		2026,
		time.August,
		13,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	currentTime := start

	limiter := NewLoginLimiter()

	limiter.now = func() time.Time {
		return currentTime
	}

	for range maxLoginFailuresPerCredential {
		limiter.RecordFailure(
			"192.0.2.1",
			"admin@example.com",
		)
	}

	currentTime = start.Add(
		loginFailureWindow,
	)

	allowed, _ := limiter.Check(
		"192.0.2.1",
		"admin@example.com",
	)

	if !allowed {
		t.Fatal(
			"expected credential to be allowed after window",
		)
	}
}

func TestLoginLimiterBlocksIPAcrossEmails(
	t *testing.T,
) {
	now := time.Date(
		2026,
		time.August,
		13,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	limiter := newLoginLimiterTest(now)

	for i := range maxLoginFailuresPerIP {
		limiter.RecordFailure(
			"192.0.2.1",
			fmt.Sprintf(
				"user-%d@example.com",
				i,
			),
		)
	}

	allowed, retryAfter :=
		limiter.Check(
			"192.0.2.1",
			"another@example.com",
		)

	if allowed {
		t.Fatal(
			"expected IP to be throttled",
		)
	}

	if retryAfter <= 0 {
		t.Fatal(
			"expected positive retry delay",
		)
	}
}

func TestLoginLimiterDoesNotBlockDifferentIP(
	t *testing.T,
) {
	now := time.Date(
		2026,
		time.August,
		13,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	limiter := newLoginLimiterTest(now)

	for range maxLoginFailuresPerCredential {
		limiter.RecordFailure(
			"192.0.2.1",
			"admin@example.com",
		)
	}

	allowed, _ := limiter.Check(
		"192.0.2.2",
		"admin@example.com",
	)

	if !allowed {
		t.Fatal(
			"expected different IP to remain allowed",
		)
	}
}

func TestLoginLimiterResetCredential(
	t *testing.T,
) {
	now := time.Date(
		2026,
		time.August,
		13,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	limiter := newLoginLimiterTest(now)

	for range maxLoginFailuresPerCredential {
		limiter.RecordFailure(
			"192.0.2.1",
			"admin@example.com",
		)
	}

	limiter.ResetCredential(
		"192.0.2.1",
		"admin@example.com",
	)

	allowed, _ := limiter.Check(
		"192.0.2.1",
		"admin@example.com",
	)

	if !allowed {
		t.Fatal(
			"expected credential throttle to reset",
		)
	}
}
