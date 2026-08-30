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
	"strings"
	"sync"
	"time"
)

const (
	loginFailureWindow            = 10 * time.Minute
	maxLoginFailuresPerCredential = 5
	maxLoginFailuresPerIP         = 20
)

type loginFailureEntry struct {
	count       int
	windowStart time.Time
}

type LoginLimiter struct {
	mu sync.Mutex

	ipFailures         map[string]loginFailureEntry
	credentialFailures map[string]loginFailureEntry

	window           time.Duration
	maxPerCredential int
	maxPerIP         int

	lastCleanup time.Time
	now         func() time.Time
}

func NewLoginLimiter() *LoginLimiter {
	return &LoginLimiter{
		ipFailures: make(
			map[string]loginFailureEntry,
		),
		credentialFailures: make(
			map[string]loginFailureEntry,
		),
		window:           loginFailureWindow,
		maxPerCredential: maxLoginFailuresPerCredential,
		maxPerIP:         maxLoginFailuresPerIP,
		now: func() time.Time {
			return time.Now().UTC()
		},
	}
}

func (l *LoginLimiter) Check(
	ip string,
	email string,
) (bool, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now().UTC()

	l.cleanupLocked(now)

	ip = normalizeLoginLimiterValue(ip)
	email = normalizeLoginLimiterValue(email)

	ipEntry := l.activeEntry(
		l.ipFailures,
		ip,
		now,
	)

	credentialKey := loginCredentialKey(
		ip,
		email,
	)

	credentialEntry := l.activeEntry(
		l.credentialFailures,
		credentialKey,
		now,
	)

	var retryAfter time.Duration

	if ipEntry.count >= l.maxPerIP {
		retryAfter = remainingLoginWindow(
			now,
			ipEntry.windowStart,
			l.window,
		)
	}

	if credentialEntry.count >=
		l.maxPerCredential {
		credentialRetry := remainingLoginWindow(
			now,
			credentialEntry.windowStart,
			l.window,
		)

		if credentialRetry > retryAfter {
			retryAfter = credentialRetry
		}
	}

	if retryAfter > 0 {
		return false, retryAfter
	}

	return true, 0
}

func (l *LoginLimiter) RecordFailure(
	ip string,
	email string,
) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now().UTC()

	l.cleanupLocked(now)

	ip = normalizeLoginLimiterValue(ip)
	email = normalizeLoginLimiterValue(email)

	incrementLoginFailure(
		l.ipFailures,
		ip,
		now,
		l.window,
	)

	incrementLoginFailure(
		l.credentialFailures,
		loginCredentialKey(
			ip,
			email,
		),
		now,
		l.window,
	)
}

func (l *LoginLimiter) ResetCredential(
	ip string,
	email string,
) {
	l.mu.Lock()
	defer l.mu.Unlock()

	ip = normalizeLoginLimiterValue(ip)
	email = normalizeLoginLimiterValue(email)

	delete(
		l.credentialFailures,
		loginCredentialKey(
			ip,
			email,
		),
	)
}

func (l *LoginLimiter) activeEntry(
	entries map[string]loginFailureEntry,
	key string,
	now time.Time,
) loginFailureEntry {
	entry, ok := entries[key]
	if !ok {
		return loginFailureEntry{}
	}

	if now.Sub(entry.windowStart) >=
		l.window {
		delete(entries, key)

		return loginFailureEntry{}
	}

	return entry
}

func (l *LoginLimiter) cleanupLocked(
	now time.Time,
) {
	if !l.lastCleanup.IsZero() &&
		now.Sub(l.lastCleanup) < l.window {
		return
	}

	for key, entry := range l.ipFailures {
		if now.Sub(entry.windowStart) >=
			l.window {
			delete(
				l.ipFailures,
				key,
			)
		}
	}

	for key, entry := range l.credentialFailures {
		if now.Sub(entry.windowStart) >=
			l.window {
			delete(
				l.credentialFailures,
				key,
			)
		}
	}

	l.lastCleanup = now
}

func incrementLoginFailure(
	entries map[string]loginFailureEntry,
	key string,
	now time.Time,
	window time.Duration,
) {
	entry, ok := entries[key]

	if !ok ||
		now.Sub(entry.windowStart) >= window {
		entries[key] = loginFailureEntry{
			count:       1,
			windowStart: now,
		}

		return
	}

	entry.count++

	entries[key] = entry
}

func loginCredentialKey(
	ip string,
	email string,
) string {
	return ip + "\x00" + email
}

func normalizeLoginLimiterValue(
	value string,
) string {
	value = strings.TrimSpace(
		strings.ToLower(value),
	)

	if value == "" {
		return "unknown"
	}

	return value
}

func remainingLoginWindow(
	now time.Time,
	windowStart time.Time,
	window time.Duration,
) time.Duration {
	remaining := window -
		now.Sub(windowStart)

	if remaining < 0 {
		return 0
	}

	return remaining
}
