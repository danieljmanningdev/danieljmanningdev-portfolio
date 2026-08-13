package auth

import (
	"encoding/base64"
	"testing"
)

func TestGenerateCSRFToken(t *testing.T) {
	token, err := GenerateCSRFToken()
	if err != nil {
		t.Fatalf(
			"generate CSRF token: %v",
			err,
		)
	}

	decoded, err :=
		base64.RawURLEncoding.DecodeString(
			token,
		)
	if err != nil {
		t.Fatalf(
			"decode CSRF token: %v",
			err,
		)
	}

	if len(decoded) != csrfTokenBytes {
		t.Fatalf(
			"expected %d bytes, got %d",
			csrfTokenBytes,
			len(decoded),
		)
	}
}

func TestGenerateCSRFTokenProducesUniqueValues(
	t *testing.T,
) {
	first, err := GenerateCSRFToken()
	if err != nil {
		t.Fatalf(
			"generate first token: %v",
			err,
		)
	}

	second, err := GenerateCSRFToken()
	if err != nil {
		t.Fatalf(
			"generate second token: %v",
			err,
		)
	}

	if first == second {
		t.Fatal(
			"expected independently generated CSRF tokens to differ",
		)
	}
}

func TestVerifyCSRFToken(t *testing.T) {
	token, err := GenerateCSRFToken()
	if err != nil {
		t.Fatalf(
			"generate token: %v",
			err,
		)
	}

	if !VerifyCSRFToken(
		token,
		token,
	) {
		t.Fatal(
			"expected matching CSRF token to verify",
		)
	}

	if VerifyCSRFToken(
		token,
		"different-token",
	) {
		t.Fatal(
			"expected different token to fail",
		)
	}
}

func TestVerifyCSRFTokenRejectsEmptyValues(
	t *testing.T,
) {
	if VerifyCSRFToken("", "") {
		t.Fatal(
			"expected empty tokens to fail",
		)
	}

	if VerifyCSRFToken(
		"token",
		"",
	) {
		t.Fatal(
			"expected missing submitted token to fail",
		)
	}
}

func TestLogoutCSRFTokenIsSessionBound(
	t *testing.T,
) {
	first := LogoutCSRFToken(
		"first-session-token",
	)

	second := LogoutCSRFToken(
		"second-session-token",
	)

	if first == "" {
		t.Fatal(
			"expected non-empty logout CSRF token",
		)
	}

	if first == second {
		t.Fatal(
			"expected different sessions to produce different logout CSRF tokens",
		)
	}

	if first != LogoutCSRFToken(
		"first-session-token",
	) {
		t.Fatal(
			"expected logout CSRF token to be deterministic for a session",
		)
	}
}

func TestLogoutCSRFTokenRejectsEmptySession(
	t *testing.T,
) {
	if token := LogoutCSRFToken(""); token != "" {
		t.Fatalf(
			"expected empty token, got %q",
			token,
		)
	}
}
