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
	"encoding/base64"
	"testing"
)

func TestGenerateSessionToken(
	t *testing.T,
) {
	token, err := GenerateSessionToken()
	if err != nil {
		t.Fatalf(
			"generate session token: %v",
			err,
		)
	}

	if token == "" {
		t.Fatal("expected non-empty token")
	}

	decoded, err := base64.RawURLEncoding.DecodeString(
		token,
	)
	if err != nil {
		t.Fatalf(
			"decode session token: %v",
			err,
		)
	}

	if len(decoded) != sessionTokenBytes {
		t.Fatalf(
			"expected %d random bytes, got %d",
			sessionTokenBytes,
			len(decoded),
		)
	}
}

func TestGenerateSessionTokenProducesUniqueTokens(
	t *testing.T,
) {
	first, err := GenerateSessionToken()
	if err != nil {
		t.Fatalf(
			"generate first token: %v",
			err,
		)
	}

	second, err := GenerateSessionToken()
	if err != nil {
		t.Fatalf(
			"generate second token: %v",
			err,
		)
	}

	if first == second {
		t.Fatal(
			"expected independently generated tokens to differ",
		)
	}
}

func TestHashSessionTokenDeterministic(
	t *testing.T,
) {
	const token = "test-session-token"

	first := HashSessionToken(token)
	second := HashSessionToken(token)

	if first != second {
		t.Fatal(
			"expected identical token hashes",
		)
	}

	if first == token {
		t.Fatal(
			"expected stored hash to differ from raw token",
		)
	}

	if len(first) != 64 {
		t.Fatalf(
			"expected SHA-256 hex hash length 64, got %d",
			len(first),
		)
	}
}

func TestHashSessionTokenDifferentTokensDiffer(
	t *testing.T,
) {
	first := HashSessionToken(
		"first-token",
	)

	second := HashSessionToken(
		"second-token",
	)

	if first == second {
		t.Fatal(
			"expected different tokens to have different hashes",
		)
	}
}
