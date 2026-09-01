// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package repository

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/database"
)

func newProjectRepositoryTestDatabase(
	t *testing.T,
) (*ProjectRepository, *sql.DB) {
	t.Helper()

	ctx := context.Background()

	db, err := database.Open(ctx, ":memory:")
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate project repository test file")
	}

	migrationsDir := filepath.Join(
		filepath.Dir(filename),
		"..",
		"..",
		"migrations",
	)

	if err := database.RunMigrations(
		db.SQL,
		migrationsDir,
	); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	return NewProjectRepository(db.SQL), db.SQL
}

func createProjectTestClient(
	t *testing.T,
	db *sql.DB,
	name string,
	email string,
) int64 {
	t.Helper()

	result, err := db.Exec(`
		INSERT INTO clients (
			name,
			email
		)
		VALUES (?, ?)
	`, name, email)
	if err != nil {
		t.Fatalf("insert test client: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("get test client id: %v", err)
	}

	return id
}

func TestProjectRepositoryCreateAndGetByID(
	t *testing.T,
) {
	repository, db := newProjectRepositoryTestDatabase(t)
	ctx := context.Background()

	clientID := createProjectTestClient(
		t,
		db,
		"Acme Studio",
		"hello@acme.test",
	)

	startDate := time.Date(
		2026,
		time.August,
		10,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	dueDate := time.Date(
		2026,
		time.September,
		30,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	id, err := repository.Create(
		ctx,
		clientID,
		"Website Platform",
		"Design and build a new client platform.",
		"active",
		&startDate,
		&dueDate,
	)
	if err != nil {
		t.Fatalf("create project: %v", err)
	}

	project, err := repository.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("get project: %v", err)
	}

	if project.ID != id {
		t.Fatalf(
			"expected project ID %d, got %d",
			id,
			project.ID,
		)
	}

	if project.ClientID != clientID {
		t.Fatalf(
			"expected client ID %d, got %d",
			clientID,
			project.ClientID,
		)
	}

	if project.ClientName != "Acme Studio" {
		t.Fatalf(
			"expected client name %q, got %q",
			"Acme Studio",
			project.ClientName,
		)
	}

	if project.Name != "Website Platform" {
		t.Fatalf(
			"expected project name %q, got %q",
			"Website Platform",
			project.Name,
		)
	}

	if project.Description !=
		"Design and build a new client platform." {
		t.Fatalf(
			"unexpected description %q",
			project.Description,
		)
	}

	if project.Status != "active" {
		t.Fatalf(
			"expected status %q, got %q",
			"active",
			project.Status,
		)
	}

	if project.StartDate == nil ||
		!project.StartDate.Equal(startDate) {
		t.Fatalf(
			"expected start date %v, got %v",
			startDate,
			project.StartDate,
		)
	}

	if project.DueDate == nil ||
		!project.DueDate.Equal(dueDate) {
		t.Fatalf(
			"expected due date %v, got %v",
			dueDate,
			project.DueDate,
		)
	}

	if project.CreatedAt.IsZero() {
		t.Fatal("expected created_at to be populated")
	}

	if project.UpdatedAt.IsZero() {
		t.Fatal("expected updated_at to be populated")
	}
}

func TestProjectRepositoryListAndListByClientID(
	t *testing.T,
) {
	repository, db := newProjectRepositoryTestDatabase(t)
	ctx := context.Background()

	acmeID := createProjectTestClient(
		t,
		db,
		"Acme Studio",
		"hello@acme.test",
	)

	orbitID := createProjectTestClient(
		t,
		db,
		"Orbit Labs",
		"hello@orbit.test",
	)

	_, err := repository.Create(
		ctx,
		acmeID,
		"Website Platform",
		"",
		"active",
		nil,
		nil,
	)
	if err != nil {
		t.Fatalf("create first Acme project: %v", err)
	}

	_, err = repository.Create(
		ctx,
		acmeID,
		"Client Portal",
		"",
		"planned",
		nil,
		nil,
	)
	if err != nil {
		t.Fatalf("create second Acme project: %v", err)
	}

	_, err = repository.Create(
		ctx,
		orbitID,
		"API Integration",
		"",
		"completed",
		nil,
		nil,
	)
	if err != nil {
		t.Fatalf("create Orbit project: %v", err)
	}

	projects, err := repository.List(ctx)
	if err != nil {
		t.Fatalf("list projects: %v", err)
	}

	if len(projects) != 3 {
		t.Fatalf(
			"expected 3 projects, got %d",
			len(projects),
		)
	}

	byName := make(
		map[string]string,
		len(projects),
	)

	for _, project := range projects {
		byName[project.Name] = project.ClientName
	}

	if byName["Website Platform"] != "Acme Studio" {
		t.Fatalf(
			"expected Website Platform to belong "+
				"to Acme Studio, got %q",
			byName["Website Platform"],
		)
	}

	if byName["API Integration"] != "Orbit Labs" {
		t.Fatalf(
			"expected API Integration to belong "+
				"to Orbit Labs, got %q",
			byName["API Integration"],
		)
	}

	acmeProjects, err := repository.ListByClientID(
		ctx,
		acmeID,
	)
	if err != nil {
		t.Fatalf("list Acme projects: %v", err)
	}

	if len(acmeProjects) != 2 {
		t.Fatalf(
			"expected 2 Acme projects, got %d",
			len(acmeProjects),
		)
	}

	for _, project := range acmeProjects {
		if project.ClientID != acmeID {
			t.Fatalf(
				"expected client ID %d, got %d",
				acmeID,
				project.ClientID,
			)
		}
	}
}

func TestProjectRepositoryGetByIDReturnsNotFound(
	t *testing.T,
) {
	repository, _ := newProjectRepositoryTestDatabase(t)

	_, err := repository.GetByID(
		context.Background(),
		999,
	)

	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf(
			"expected sql.ErrNoRows, got %v",
			err,
		)
	}
}

func TestProjectRepositoryCreateRejectsMissingClient(
	t *testing.T,
) {
	repository, _ := newProjectRepositoryTestDatabase(t)

	_, err := repository.Create(
		context.Background(),
		999,
		"Invalid Project",
		"",
		"planned",
		nil,
		nil,
	)

	if err == nil {
		t.Fatal(
			"expected project creation to reject a missing client",
		)
	}
}

func TestProjectRepositoryUpdate(t *testing.T) {
	repository, db := newProjectRepositoryTestDatabase(t)
	ctx := context.Background()

	firstClientID := createProjectTestClient(
		t,
		db,
		"First Client",
		"first@example.test",
	)

	secondClientID := createProjectTestClient(
		t,
		db,
		"Second Client",
		"second@example.test",
	)

	id, err := repository.Create(
		ctx,
		firstClientID,
		"Original Project",
		"Original description",
		"planned",
		nil,
		nil,
	)
	if err != nil {
		t.Fatalf("create project: %v", err)
	}

	const oldTimestamp = "2000-01-01 00:00:00"

	_, err = db.Exec(
		`
			UPDATE projects
			SET updated_at = ?
			WHERE id = ?
		`,
		oldTimestamp,
		id,
	)
	if err != nil {
		t.Fatalf(
			"set old updated_at value: %v",
			err,
		)
	}

	startDate := time.Date(
		2026,
		time.October,
		1,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	dueDate := time.Date(
		2026,
		time.December,
		15,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	if err := repository.Update(
		ctx,
		id,
		secondClientID,
		"Updated Project",
		"Updated description",
		"active",
		&startDate,
		&dueDate,
	); err != nil {
		t.Fatalf("update project: %v", err)
	}

	project, err := repository.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("get updated project: %v", err)
	}

	if project.ClientID != secondClientID {
		t.Fatalf(
			"expected client ID %d, got %d",
			secondClientID,
			project.ClientID,
		)
	}

	if project.ClientName != "Second Client" {
		t.Fatalf(
			"expected client name %q, got %q",
			"Second Client",
			project.ClientName,
		)
	}

	if project.Name != "Updated Project" {
		t.Fatalf(
			"unexpected project name %q",
			project.Name,
		)
	}

	if project.Description != "Updated description" {
		t.Fatalf(
			"unexpected description %q",
			project.Description,
		)
	}

	if project.Status != "active" {
		t.Fatalf(
			"expected status %q, got %q",
			"active",
			project.Status,
		)
	}

	if project.StartDate == nil ||
		!project.StartDate.Equal(startDate) {
		t.Fatalf(
			"expected start date %v, got %v",
			startDate,
			project.StartDate,
		)
	}

	if project.DueDate == nil ||
		!project.DueDate.Equal(dueDate) {
		t.Fatalf(
			"expected due date %v, got %v",
			dueDate,
			project.DueDate,
		)
	}

	oldTime, err := time.Parse(
		"2006-01-02 15:04:05",
		oldTimestamp,
	)
	if err != nil {
		t.Fatalf("parse old timestamp: %v", err)
	}

	if !project.UpdatedAt.After(oldTime) {
		t.Fatalf(
			"expected updated_at after %v, got %v",
			oldTime,
			project.UpdatedAt,
		)
	}
}

func TestProjectRepositoryUpdateReturnsNotFound(
	t *testing.T,
) {
	repository, db := newProjectRepositoryTestDatabase(t)

	clientID := createProjectTestClient(
		t,
		db,
		"Test Client",
		"test@example.test",
	)

	err := repository.Update(
		context.Background(),
		999,
		clientID,
		"Missing Project",
		"",
		"active",
		nil,
		nil,
	)

	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf(
			"expected sql.ErrNoRows, got %v",
			err,
		)
	}
}

func TestProjectRepositoryArchive(t *testing.T) {
	repository, db := newProjectRepositoryTestDatabase(t)
	ctx := context.Background()

	clientID := createProjectTestClient(
		t,
		db,
		"Test Client",
		"test@example.test",
	)

	id, err := repository.Create(
		ctx,
		clientID,
		"Project to Archive",
		"",
		"active",
		nil,
		nil,
	)
	if err != nil {
		t.Fatalf("create project: %v", err)
	}

	if err := repository.Archive(ctx, id); err != nil {
		t.Fatalf("archive project: %v", err)
	}

	project, err := repository.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("get archived project: %v", err)
	}

	if project.Status != "archived" {
		t.Fatalf(
			"expected status %q, got %q",
			"archived",
			project.Status,
		)
	}

	var count int

	if err := db.QueryRow(
		`
			SELECT COUNT(*)
			FROM projects
			WHERE id = ?
		`,
		id,
	).Scan(&count); err != nil {
		t.Fatalf(
			"count archived project: %v",
			err,
		)
	}

	if count != 1 {
		t.Fatalf(
			"expected archived project to remain, "+
				"found %d rows",
			count,
		)
	}
}

func TestProjectRepositoryArchiveReturnsNotFound(
	t *testing.T,
) {
	repository, _ := newProjectRepositoryTestDatabase(t)

	err := repository.Archive(
		context.Background(),
		999,
	)

	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf(
			"expected sql.ErrNoRows, got %v",
			err,
		)
	}
}
