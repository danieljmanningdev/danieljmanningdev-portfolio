// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/auth"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/config"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/database"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/repository"
)

func main() {
	email := flag.String(
		"email",
		"",
		"admin email address",
	)

	name := flag.String(
		"name",
		"",
		"admin display name",
	)

	flag.Parse()

	normalizedEmail := strings.ToLower(
		strings.TrimSpace(*email),
	)

	displayName := strings.TrimSpace(*name)

	if normalizedEmail == "" {
		fmt.Fprintln(
			os.Stderr,
			"error: -email is required",
		)
		os.Exit(1)
	}

	if displayName == "" {
		fmt.Fprintln(
			os.Stderr,
			"error: -name is required",
		)
		os.Exit(1)
	}

	fmt.Print("Password: ")

	passwordBytes, err := term.ReadPassword(
		int(syscall.Stdin),
	)
	fmt.Println()

	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"error reading password: %v\n",
			err,
		)
		os.Exit(1)
	}

	fmt.Print("Confirm password: ")

	confirmationBytes, err := term.ReadPassword(
		int(syscall.Stdin),
	)
	fmt.Println()

	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"error reading password confirmation: %v\n",
			err,
		)
		os.Exit(1)
	}

	password := string(passwordBytes)
	confirmation := string(confirmationBytes)

	if password != confirmation {
		fmt.Fprintln(
			os.Stderr,
			"error: passwords do not match",
		)
		os.Exit(1)
	}

	passwordHash, err := auth.HashPassword(
		password,
	)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"error hashing password: %v\n",
			err,
		)
		os.Exit(1)
	}

	cfg := config.Load()

	ctx := context.Background()

	db, err := database.Open(
		ctx,
		cfg.DatabasePath,
	)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"error opening database: %v\n",
			err,
		)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Fprintf(
				os.Stderr,
				"error closing database: %v\n",
				err,
			)
		}
	}()

	if err := database.RunMigrations(
		db.SQL,
		"migrations",
	); err != nil {
		fmt.Fprintf(
			os.Stderr,
			"error running migrations: %v\n",
			err,
		)
		os.Exit(1)
	}

	adminRepository :=
		repository.NewAdminRepository(
			db.SQL,
		)

	adminID, err := adminRepository.Create(
		ctx,
		normalizedEmail,
		passwordHash,
		displayName,
	)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"error creating admin: %v\n",
			err,
		)
		os.Exit(1)
	}

	fmt.Printf(
		"Admin created successfully (id=%d, email=%s)\n",
		adminID,
		normalizedEmail,
	)
}
