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

func newContractRepositoryTestDatabase(
	t *testing.T,
) (*ContractRepository, *sql.DB) {
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
		t.Fatal("locate contract repository test file")
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

	return NewContractRepository(db.SQL), db.SQL
}

func createContractTestClient(
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

func createContractTestProject(
	t *testing.T,
	db *sql.DB,
	clientID int64,
	name string,
) int64 {
	t.Helper()

	result, err := db.Exec(`
		INSERT INTO projects (
			client_id,
			name
		)
		VALUES (?, ?)
	`, clientID, name)
	if err != nil {
		t.Fatalf("insert test project: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("get test project id: %v", err)
	}

	return id
}

func TestContractRepositoryCreateAndGetByID(
	t *testing.T,
) {
	repository, db := newContractRepositoryTestDatabase(t)
	ctx := context.Background()

	clientID := createContractTestClient(
		t,
		db,
		"Acme Studio",
		"hello@acme.test",
	)

	projectID := createContractTestProject(
		t,
		db,
		clientID,
		"Website Platform",
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

	endDate := time.Date(
		2026,
		time.September,
		30,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	valueCents := int64(250000)

	id, err := repository.Create(
		ctx,
		clientID,
		&projectID,
		"Website Development Agreement",
		"draft",
		&startDate,
		&endDate,
		&valueCents,
		"Initial contract notes.",
	)
	if err != nil {
		t.Fatalf("create contract: %v", err)
	}

	contract, err := repository.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("get contract: %v", err)
	}

	if contract.ID != id {
		t.Fatalf(
			"expected contract ID %d, got %d",
			id,
			contract.ID,
		)
	}

	if contract.ClientID != clientID {
		t.Fatalf(
			"expected client ID %d, got %d",
			clientID,
			contract.ClientID,
		)
	}

	if contract.ClientName != "Acme Studio" {
		t.Fatalf(
			"expected client name %q, got %q",
			"Acme Studio",
			contract.ClientName,
		)
	}

	if contract.ProjectID == nil {
		t.Fatal("expected project ID to be populated")
	}

	if *contract.ProjectID != projectID {
		t.Fatalf(
			"expected project ID %d, got %d",
			projectID,
			*contract.ProjectID,
		)
	}

	if contract.ProjectName != "Website Platform" {
		t.Fatalf(
			"expected project name %q, got %q",
			"Website Platform",
			contract.ProjectName,
		)
	}

	if contract.Title != "Website Development Agreement" {
		t.Fatalf(
			"unexpected title %q",
			contract.Title,
		)
	}

	if contract.Status != "draft" {
		t.Fatalf(
			"expected status %q, got %q",
			"draft",
			contract.Status,
		)
	}

	if contract.StartDate == nil ||
		!contract.StartDate.Equal(startDate) {
		t.Fatalf(
			"expected start date %v, got %v",
			startDate,
			contract.StartDate,
		)
	}

	if contract.EndDate == nil ||
		!contract.EndDate.Equal(endDate) {
		t.Fatalf(
			"expected end date %v, got %v",
			endDate,
			contract.EndDate,
		)
	}

	if contract.ValueCents == nil {
		t.Fatal("expected contract value to be populated")
	}

	if *contract.ValueCents != valueCents {
		t.Fatalf(
			"expected value %d, got %d",
			valueCents,
			*contract.ValueCents,
		)
	}

	if contract.Notes != "Initial contract notes." {
		t.Fatalf(
			"unexpected notes %q",
			contract.Notes,
		)
	}

	if contract.CreatedAt.IsZero() {
		t.Fatal("expected created_at to be populated")
	}

	if contract.UpdatedAt.IsZero() {
		t.Fatal("expected updated_at to be populated")
	}
}

func TestContractRepositoryAllowsContractWithoutProject(
	t *testing.T,
) {
	repository, db := newContractRepositoryTestDatabase(t)
	ctx := context.Background()

	clientID := createContractTestClient(
		t,
		db,
		"Acme Studio",
		"hello@acme.test",
	)

	id, err := repository.Create(
		ctx,
		clientID,
		nil,
		"General Services Agreement",
		"draft",
		nil,
		nil,
		nil,
		"",
	)
	if err != nil {
		t.Fatalf("create contract without project: %v", err)
	}

	contract, err := repository.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("get contract: %v", err)
	}

	if contract.ProjectID != nil {
		t.Fatalf(
			"expected nil project ID, got %v",
			contract.ProjectID,
		)
	}

	if contract.ProjectName != "" {
		t.Fatalf(
			"expected empty project name, got %q",
			contract.ProjectName,
		)
	}

	if contract.StartDate != nil {
		t.Fatalf(
			"expected nil start date, got %v",
			contract.StartDate,
		)
	}

	if contract.EndDate != nil {
		t.Fatalf(
			"expected nil end date, got %v",
			contract.EndDate,
		)
	}

	if contract.ValueCents != nil {
		t.Fatalf(
			"expected nil value, got %v",
			contract.ValueCents,
		)
	}
}

func TestContractRepositoryListAndListByClientID(
	t *testing.T,
) {
	repository, db := newContractRepositoryTestDatabase(t)
	ctx := context.Background()

	acmeID := createContractTestClient(
		t,
		db,
		"Acme Studio",
		"hello@acme.test",
	)

	orbitID := createContractTestClient(
		t,
		db,
		"Orbit Labs",
		"hello@orbit.test",
	)

	_, err := repository.Create(
		ctx,
		acmeID,
		nil,
		"Acme Agreement",
		"draft",
		nil,
		nil,
		nil,
		"",
	)
	if err != nil {
		t.Fatalf("create Acme contract: %v", err)
	}

	_, err = repository.Create(
		ctx,
		acmeID,
		nil,
		"Acme Support Agreement",
		"accepted",
		nil,
		nil,
		nil,
		"",
	)
	if err != nil {
		t.Fatalf("create second Acme contract: %v", err)
	}

	_, err = repository.Create(
		ctx,
		orbitID,
		nil,
		"Orbit Agreement",
		"completed",
		nil,
		nil,
		nil,
		"",
	)
	if err != nil {
		t.Fatalf("create Orbit contract: %v", err)
	}

	contracts, err := repository.List(ctx)
	if err != nil {
		t.Fatalf("list contracts: %v", err)
	}

	if len(contracts) != 3 {
		t.Fatalf(
			"expected 3 contracts, got %d",
			len(contracts),
		)
	}

	acmeContracts, err := repository.ListByClientID(
		ctx,
		acmeID,
	)
	if err != nil {
		t.Fatalf(
			"list contracts by client: %v",
			err,
		)
	}

	if len(acmeContracts) != 2 {
		t.Fatalf(
			"expected 2 Acme contracts, got %d",
			len(acmeContracts),
		)
	}

	for _, contract := range acmeContracts {
		if contract.ClientID != acmeID {
			t.Fatalf(
				"expected client ID %d, got %d",
				acmeID,
				contract.ClientID,
			)
		}
	}
}

func TestContractRepositoryGetByIDReturnsNotFound(
	t *testing.T,
) {
	repository, _ := newContractRepositoryTestDatabase(t)

	_, err := repository.GetByID(
		context.Background(),
		999999,
	)

	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf(
			"expected sql.ErrNoRows, got %v",
			err,
		)
	}
}

func TestContractRepositoryCreateRejectsMissingClient(
	t *testing.T,
) {
	repository, _ := newContractRepositoryTestDatabase(t)

	_, err := repository.Create(
		context.Background(),
		999999,
		nil,
		"Invalid Contract",
		"draft",
		nil,
		nil,
		nil,
		"",
	)

	if err == nil {
		t.Fatal("expected missing-client foreign-key error")
	}
}

func TestContractRepositoryCreateRejectsMissingProject(
	t *testing.T,
) {
	repository, db := newContractRepositoryTestDatabase(t)

	clientID := createContractTestClient(
		t,
		db,
		"Acme Studio",
		"hello@acme.test",
	)

	missingProjectID := int64(999999)

	_, err := repository.Create(
		context.Background(),
		clientID,
		&missingProjectID,
		"Invalid Contract",
		"draft",
		nil,
		nil,
		nil,
		"",
	)

	if err == nil {
		t.Fatal("expected missing-project foreign-key error")
	}
}

func TestContractRepositoryUpdate(
	t *testing.T,
) {
	repository, db := newContractRepositoryTestDatabase(t)
	ctx := context.Background()

	firstClientID := createContractTestClient(
		t,
		db,
		"Acme Studio",
		"hello@acme.test",
	)

	secondClientID := createContractTestClient(
		t,
		db,
		"Orbit Labs",
		"hello@orbit.test",
	)

	firstProjectID := createContractTestProject(
		t,
		db,
		firstClientID,
		"Original Project",
	)

	secondProjectID := createContractTestProject(
		t,
		db,
		secondClientID,
		"Updated Project",
	)

	id, err := repository.Create(
		ctx,
		firstClientID,
		&firstProjectID,
		"Original Contract",
		"draft",
		nil,
		nil,
		nil,
		"Original notes.",
	)
	if err != nil {
		t.Fatalf("create contract: %v", err)
	}

	startDate := time.Date(
		2026,
		time.September,
		1,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	endDate := time.Date(
		2026,
		time.December,
		31,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	valueCents := int64(500000)

	err = repository.Update(
		ctx,
		id,
		secondClientID,
		&secondProjectID,
		"Updated Contract",
		"accepted",
		&startDate,
		&endDate,
		&valueCents,
		"Updated notes.",
	)
	if err != nil {
		t.Fatalf("update contract: %v", err)
	}

	contract, err := repository.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("get updated contract: %v", err)
	}

	if contract.ClientID != secondClientID {
		t.Fatalf(
			"expected client ID %d, got %d",
			secondClientID,
			contract.ClientID,
		)
	}

	if contract.ProjectID == nil ||
		*contract.ProjectID != secondProjectID {
		t.Fatalf(
			"expected project ID %d, got %v",
			secondProjectID,
			contract.ProjectID,
		)
	}

	if contract.Title != "Updated Contract" {
		t.Fatalf(
			"unexpected title %q",
			contract.Title,
		)
	}

	if contract.Status != "accepted" {
		t.Fatalf(
			"expected accepted status, got %q",
			contract.Status,
		)
	}

	if contract.StartDate == nil ||
		!contract.StartDate.Equal(startDate) {
		t.Fatalf(
			"expected start date %v, got %v",
			startDate,
			contract.StartDate,
		)
	}

	if contract.EndDate == nil ||
		!contract.EndDate.Equal(endDate) {
		t.Fatalf(
			"expected end date %v, got %v",
			endDate,
			contract.EndDate,
		)
	}

	if contract.ValueCents == nil ||
		*contract.ValueCents != valueCents {
		t.Fatalf(
			"expected value %d, got %v",
			valueCents,
			contract.ValueCents,
		)
	}

	if contract.Notes != "Updated notes." {
		t.Fatalf(
			"unexpected notes %q",
			contract.Notes,
		)
	}
}

func TestContractRepositoryUpdateReturnsNotFound(
	t *testing.T,
) {
	repository, db := newContractRepositoryTestDatabase(t)

	clientID := createContractTestClient(
		t,
		db,
		"Acme Studio",
		"hello@acme.test",
	)

	err := repository.Update(
		context.Background(),
		999999,
		clientID,
		nil,
		"Missing Contract",
		"draft",
		nil,
		nil,
		nil,
		"",
	)

	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf(
			"expected sql.ErrNoRows, got %v",
			err,
		)
	}
}

func TestContractRepositoryCancel(
	t *testing.T,
) {
	repository, db := newContractRepositoryTestDatabase(t)
	ctx := context.Background()

	clientID := createContractTestClient(
		t,
		db,
		"Acme Studio",
		"hello@acme.test",
	)

	id, err := repository.Create(
		ctx,
		clientID,
		nil,
		"Contract to Cancel",
		"accepted",
		nil,
		nil,
		nil,
		"",
	)
	if err != nil {
		t.Fatalf("create contract: %v", err)
	}

	if err := repository.Cancel(ctx, id); err != nil {
		t.Fatalf("cancel contract: %v", err)
	}

	contract, err := repository.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("get cancelled contract: %v", err)
	}

	if contract.Status != "cancelled" {
		t.Fatalf(
			"expected status %q, got %q",
			"cancelled",
			contract.Status,
		)
	}
}

func TestContractRepositoryCancelReturnsNotFound(
	t *testing.T,
) {
	repository, _ := newContractRepositoryTestDatabase(t)

	err := repository.Cancel(
		context.Background(),
		999999,
	)

	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf(
			"expected sql.ErrNoRows, got %v",
			err,
		)
	}
}
