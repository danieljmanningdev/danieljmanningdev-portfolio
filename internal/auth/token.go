package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

const sessionTokenBytes = 32

func GenerateSessionToken() (string, error) {
	tokenBytes := make(
		[]byte,
		sessionTokenBytes,
	)

	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(
		tokenBytes,
	), nil
}

func HashSessionToken(
	token string,
) string {
	sum := sha256.Sum256(
		[]byte(token),
	)

	return hex.EncodeToString(
		sum[:],
	)
}
