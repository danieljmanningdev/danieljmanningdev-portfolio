package auth

import (
	"errors"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost       = 12
	minPasswordChars = 15
	maxPasswordBytes = 72
)

var (
	ErrPasswordEmpty    = errors.New("password is required")
	ErrPasswordTooShort = errors.New("password must be at least 15 characters")
	ErrPasswordTooLong  = errors.New("password exceeds bcrypt 72-byte limit")
)

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", ErrPasswordEmpty
	}

	if !utf8.ValidString(password) {
		return "", errors.New("password is not valid UTF-8")
	}

	if utf8.RuneCountInString(password) < minPasswordChars {
		return "", ErrPasswordTooShort
	}

	if len([]byte(password)) > maxPasswordBytes {
		return "", ErrPasswordTooLong
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcryptCost,
	)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func VerifyPassword(
	passwordHash string,
	password string,
) bool {
	if passwordHash == "" || password == "" {
		return false
	}

	if len([]byte(password)) > maxPasswordBytes {
		return false
	}

	return bcrypt.CompareHashAndPassword(
		[]byte(passwordHash),
		[]byte(password),
	) == nil
}
