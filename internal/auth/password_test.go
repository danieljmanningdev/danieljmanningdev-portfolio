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
	"errors"
	"strings"
	"testing"
)

func TestHashPasswordAndVerify(
	t *testing.T,
) {
	password := "a-secure-test-password"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	if hash == password {
		t.Fatal("expected password hash to differ from plaintext")
	}

	if !VerifyPassword(hash, password) {
		t.Fatal("expected correct password to verify")
	}

	if VerifyPassword(
		hash,
		"wrong-password",
	) {
		t.Fatal("expected incorrect password to fail")
	}
}

func TestHashPasswordUsesUniqueSalt(
	t *testing.T,
) {
	password := "same-password-that-is-long"

	firstHash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("hash first password: %v", err)
	}

	secondHash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("hash second password: %v", err)
	}

	if firstHash == secondHash {
		t.Fatal(
			"expected bcrypt hashes to differ because of unique salts",
		)
	}

	if !VerifyPassword(
		firstHash,
		password,
	) {
		t.Fatal("expected first hash to verify")
	}

	if !VerifyPassword(
		secondHash,
		password,
	) {
		t.Fatal("expected second hash to verify")
	}
}

func TestHashPasswordRejectsEmptyPassword(
	t *testing.T,
) {
	_, err := HashPassword("")

	if !errors.Is(
		err,
		ErrPasswordEmpty,
	) {
		t.Fatalf(
			"expected ErrPasswordEmpty, got %v",
			err,
		)
	}
}

func TestHashPasswordRejectsFewerThan15Characters(
	t *testing.T,
) {
	_, err := HashPassword(
		"short-password",
	)

	if !errors.Is(
		err,
		ErrPasswordTooShort,
	) {
		t.Fatalf(
			"expected ErrPasswordTooShort, got %v",
			err,
		)
	}
}

func TestHashPasswordAccepts15Characters(
	t *testing.T,
) {
	password := "123456789012345"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf(
			"expected 15-character password to be accepted: %v",
			err,
		)
	}

	if !VerifyPassword(
		hash,
		password,
	) {
		t.Fatal(
			"expected accepted password to verify",
		)
	}
}

func TestHashPasswordRejectsMoreThan72Bytes(
	t *testing.T,
) {
	password := strings.Repeat(
		"a",
		73,
	)

	_, err := HashPassword(password)

	if !errors.Is(
		err,
		ErrPasswordTooLong,
	) {
		t.Fatalf(
			"expected ErrPasswordTooLong, got %v",
			err,
		)
	}
}

func TestHashPasswordCountsBytesNotCharacters(
	t *testing.T,
) {
	password := strings.Repeat(
		"£",
		37,
	)

	_, err := HashPassword(password)

	if !errors.Is(
		err,
		ErrPasswordTooLong,
	) {
		t.Fatalf(
			"expected multibyte password to exceed bcrypt byte limit, got %v",
			err,
		)
	}
}

func TestVerifyPasswordRejectsEmptyValues(
	t *testing.T,
) {
	if VerifyPassword(
		"",
		"password",
	) {
		t.Fatal("expected empty hash to fail")
	}

	hash, err := HashPassword("password-long-enough")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	if VerifyPassword(
		hash,
		"",
	) {
		t.Fatal("expected empty password to fail")
	}
}
