package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
)

const csrfTokenBytes = 32

type adminCSRFContextKey struct{}

func WithAdminCSRFToken(
	ctx context.Context,
	token string,
) context.Context {
	return context.WithValue(
		ctx,
		adminCSRFContextKey{},
		token,
	)
}

func AdminCSRFTokenFromContext(
	ctx context.Context,
) (string, bool) {
	token, ok := ctx.Value(
		adminCSRFContextKey{},
	).(string)

	return token, ok && token != ""
}

func GenerateCSRFToken() (string, error) {
	tokenBytes := make(
		[]byte,
		csrfTokenBytes,
	)

	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(
		tokenBytes,
	), nil
}

func VerifyCSRFToken(
	expected string,
	provided string,
) bool {
	if expected == "" || provided == "" {
		return false
	}

	expectedBytes := []byte(expected)
	providedBytes := []byte(provided)

	if len(expectedBytes) != len(providedBytes) {
		return false
	}

	return subtle.ConstantTimeCompare(
		expectedBytes,
		providedBytes,
	) == 1
}

func LogoutCSRFToken(
	sessionToken string,
) string {
	if sessionToken == "" {
		return ""
	}

	sum := sha256.Sum256(
		[]byte(
			"admin-logout-csrf:" +
				sessionToken,
		),
	)

	return base64.RawURLEncoding.EncodeToString(
		sum[:],
	)
}

func AdminCSRFToken(
	sessionToken string,
) string {
	if sessionToken == "" {
		return ""
	}

	sum := sha256.Sum256(
		[]byte(
			"admin-csrf:" +
				sessionToken,
		),
	)

	return base64.RawURLEncoding.EncodeToString(
		sum[:],
	)
}
