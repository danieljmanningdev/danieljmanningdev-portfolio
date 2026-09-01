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

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/config"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/database"
)

func main() {
	action := flag.String(
		"action",
		"",
		"database action: backup, verify, or restore",
	)

	databasePath := flag.String(
		"database",
		"",
		"application SQLite database path; defaults to DATABASE_PATH",
	)

	filePath := flag.String(
		"file",
		"",
		"backup file path",
	)

	force := flag.Bool(
		"force",
		false,
		"allow restore to replace an existing database",
	)

	flag.Parse()

	selectedAction := strings.ToLower(
		strings.TrimSpace(*action),
	)
	selectedFile := strings.TrimSpace(*filePath)

	cfg := config.Load()
	selectedDatabase := strings.TrimSpace(*databasePath)
	if selectedDatabase == "" {
		selectedDatabase = cfg.DatabasePath
	}

	ctx := context.Background()

	var err error

	switch selectedAction {
	case "backup":
		err = database.Backup(
			ctx,
			selectedDatabase,
			selectedFile,
		)
		if err == nil {
			fmt.Printf(
				"SQLite backup created: %s\n",
				selectedFile,
			)
		}

	case "verify":
		err = database.VerifyBackup(
			ctx,
			selectedFile,
		)
		if err == nil {
			fmt.Printf(
				"SQLite backup verified: %s\n",
				selectedFile,
			)
		}

	case "restore":
		if !*force {
			err = fmt.Errorf(
				"restore requires -force after the application has been stopped",
			)
			break
		}

		err = database.Restore(
			ctx,
			selectedFile,
			selectedDatabase,
			true,
		)
		if err == nil {
			fmt.Printf(
				"SQLite database restored: %s\n",
				selectedDatabase,
			)
		}

	default:
		err = fmt.Errorf(
			"-action must be one of: backup, verify, restore",
		)
	}

	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"database operation failed: %v\n",
			err,
		)
		os.Exit(1)
	}
}
