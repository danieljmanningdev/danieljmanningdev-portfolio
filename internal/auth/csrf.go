package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
)

const csrfTokenBytes = 32

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
