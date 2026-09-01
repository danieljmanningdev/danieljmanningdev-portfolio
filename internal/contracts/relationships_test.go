// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package contracts

import (
	"context"
	"database/sql"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/database"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/repository"
)

func newContractRelationshipTestHandler(
	t *testing.T,
) (*ContractsHandler, *sql.DB) {
	t.Helper()

	db, err := database.Open(
		context.Background(),
		":memory:",
	)
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate contract relationship test file")
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

	formTemplate := template.Must(
		template.New("base").Parse(`
			{{define "base"}}
			{{.Errors.General}}
			{{.Errors.ClientID}}
			{{.Errors.ProjectID}}
			{{.Errors.Title}}
			{{.Form.ClientID}}
			{{.Form.ProjectID}}
			{{.Form.Title}}
			{{range .Clients}}{{.Name}}{{end}}
			{{range .Projects}}{{.Name}}{{end}}
			{{end}}
		`),
	)

	handler := &ContractsHandler{
		contractRepository: repository.NewContractRepository(
			db.SQL,
		),
		clientRepository: repository.NewClientRepository(
			db.SQL,
		),
		projectRepository: repository.NewProjectRepository(
			db.SQL,
		),
		newContractTemplates:  formTemplate,
		editContractTemplates: formTemplate,
	}

	return handler, db.SQL
}

func createContractRelationshipClient(
	t *testing.T,
	db *sql.DB,
	name string,
	email string,
) int64 {
	t.Helper()

	result, err := db.Exec(`
		INSERT INTO clients (
			name,
			email,
			company,
			notes
		)
		VALUES (?, ?, ?, ?)
	`,
		name,
		email,
		"",
		"",
	)
	if err != nil {
		t.Fatalf("insert test client: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("get test client ID: %v", err)
	}

	return id
}

func createContractRelationshipProject(
	t *testing.T,
	db *sql.DB,
	clientID int64,
	name string,
) int64 {
	t.Helper()

	result, err := db.Exec(`
		INSERT INTO projects (
			client_id,
			name,
			description,
			status
		)
		VALUES (?, ?, ?, ?)
	`,
		clientID,
		name,
		"",
		"active",
	)
	if err != nil {
		t.Fatalf("insert test project: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("get test project ID: %v", err)
	}

	return id
}

func createContractRelationshipContract(
	t *testing.T,
	db *sql.DB,
	clientID int64,
	projectID int64,
) int64 {
	t.Helper()

	result, err := db.Exec(`
		INSERT INTO contracts (
			client_id,
			project_id,
			title,
			status,
			notes
		)
		VALUES (?, ?, ?, ?, ?)
	`,
		clientID,
		projectID,
		"Existing Contract",
		contractStatusDraft,
		"",
	)
	if err != nil {
		t.Fatalf("insert test contract: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("get test contract ID: %v", err)
	}

	return id
}

func contractOwnershipFormRequest(
	target string,
	values url.Values,
) *http.Request {
	req := httptest.NewRequest(
		http.MethodPost,
		target,
		strings.NewReader(values.Encode()),
	)

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	return req
}

func validContractOwnershipForm(
	clientID int64,
	projectID int64,
) url.Values {
	values := url.Values{}

	values.Set(
		"client_id",
		strconv.FormatInt(clientID, 10),
	)

	if projectID > 0 {
		values.Set(
			"project_id",
			strconv.FormatInt(projectID, 10),
		)
	}

	values.Set("title", "Website Contract")
	values.Set("status", contractStatusDraft)

	return values
}

func TestContractRelationshipCreateAcceptsMatchingProject(
	t *testing.T,
) {
	handler, db := newContractRelationshipTestHandler(t)

	clientID := createContractRelationshipClient(
		t,
		db,
		"Acme Studio",
		"acme@example.test",
	)

	projectID := createContractRelationshipProject(
		t,
		db,
		clientID,
		"Website Platform",
	)

	req := contractOwnershipFormRequest(
		contractsBasePath+"/new",
		validContractOwnershipForm(
			clientID,
			projectID,
		),
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf(
			"expected status 303, got %d: %s",
			rec.Code,
			rec.Body.String(),
		)
	}

	var count int

	if err := db.QueryRow(`
		SELECT COUNT(*)
		FROM contracts
		WHERE client_id = ?
		AND project_id = ?
	`,
		clientID,
		projectID,
	).Scan(&count); err != nil {
		t.Fatalf("count created contract: %v", err)
	}

	if count != 1 {
		t.Fatalf(
			"expected one contract, found %d",
			count,
		)
	}
}

func TestContractRelationshipCreateAllowsNoProject(
	t *testing.T,
) {
	handler, db := newContractRelationshipTestHandler(t)

	clientID := createContractRelationshipClient(
		t,
		db,
		"Acme Studio",
		"acme@example.test",
	)

	req := contractOwnershipFormRequest(
		contractsBasePath+"/new",
		validContractOwnershipForm(
			clientID,
			0,
		),
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf(
			"expected status 303, got %d: %s",
			rec.Code,
			rec.Body.String(),
		)
	}

	var count int

	if err := db.QueryRow(`
		SELECT COUNT(*)
		FROM contracts
		WHERE client_id = ?
		AND project_id IS NULL
	`,
		clientID,
	).Scan(&count); err != nil {
		t.Fatalf(
			"count contract without project: %v",
			err,
		)
	}

	if count != 1 {
		t.Fatalf(
			"expected one contract without a project, found %d",
			count,
		)
	}
}

func TestContractRelationshipCreateRejectsDifferentClientProject(
	t *testing.T,
) {
	handler, db := newContractRelationshipTestHandler(t)

	firstClientID := createContractRelationshipClient(
		t,
		db,
		"First Client",
		"first@example.test",
	)

	secondClientID := createContractRelationshipClient(
		t,
		db,
		"Second Client",
		"second@example.test",
	)

	firstProjectID := createContractRelationshipProject(
		t,
		db,
		firstClientID,
		"First Client Project",
	)

	req := contractOwnershipFormRequest(
		contractsBasePath+"/new",
		validContractOwnershipForm(
			secondClientID,
			firstProjectID,
		),
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status 400, got %d",
			rec.Code,
		)
	}

	const expectedMessage = "Selected project does not belong to the selected client."

	if !strings.Contains(
		rec.Body.String(),
		expectedMessage,
	) {
		t.Fatalf(
			"expected response to contain %q, got %q",
			expectedMessage,
			rec.Body.String(),
		)
	}

	var count int

	if err := db.QueryRow(
		"SELECT COUNT(*) FROM contracts",
	).Scan(&count); err != nil {
		t.Fatalf("count contracts: %v", err)
	}

	if count != 0 {
		t.Fatalf(
			"expected no contract to be created, found %d",
			count,
		)
	}
}

func TestContractRelationshipCreateRejectsMissingProject(
	t *testing.T,
) {
	handler, db := newContractRelationshipTestHandler(t)

	clientID := createContractRelationshipClient(
		t,
		db,
		"Acme Studio",
		"acme@example.test",
	)

	req := contractOwnershipFormRequest(
		contractsBasePath+"/new",
		validContractOwnershipForm(
			clientID,
			999,
		),
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status 400, got %d",
			rec.Code,
		)
	}

	if !strings.Contains(
		rec.Body.String(),
		"Selected project does not exist.",
	) {
		t.Fatalf(
			"unexpected response body %q",
			rec.Body.String(),
		)
	}
}

func TestContractRelationshipCreateRejectsMissingClient(
	t *testing.T,
) {
	handler, _ := newContractRelationshipTestHandler(t)

	req := contractOwnershipFormRequest(
		contractsBasePath+"/new",
		validContractOwnershipForm(
			999,
			0,
		),
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status 400, got %d",
			rec.Code,
		)
	}

	if !strings.Contains(
		rec.Body.String(),
		"Selected client does not exist.",
	) {
		t.Fatalf(
			"unexpected response body %q",
			rec.Body.String(),
		)
	}
}

func TestContractRelationshipUpdateRejectsDifferentClientProject(
	t *testing.T,
) {
	handler, db := newContractRelationshipTestHandler(t)

	firstClientID := createContractRelationshipClient(
		t,
		db,
		"First Client",
		"first@example.test",
	)

	secondClientID := createContractRelationshipClient(
		t,
		db,
		"Second Client",
		"second@example.test",
	)

	firstProjectID := createContractRelationshipProject(
		t,
		db,
		firstClientID,
		"First Client Project",
	)

	contractID := createContractRelationshipContract(
		t,
		db,
		firstClientID,
		firstProjectID,
	)

	req := contractOwnershipFormRequest(
		contractsBasePath+"/"+strconv.FormatInt(
			contractID,
			10,
		),
		validContractOwnershipForm(
			secondClientID,
			firstProjectID,
		),
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status 400, got %d",
			rec.Code,
		)
	}

	const expectedMessage = "Selected project does not belong to the selected client."

	if !strings.Contains(
		rec.Body.String(),
		expectedMessage,
	) {
		t.Fatalf(
			"expected response to contain %q, got %q",
			expectedMessage,
			rec.Body.String(),
		)
	}

	var (
		storedClientID  int64
		storedProjectID int64
	)

	if err := db.QueryRow(`
		SELECT
			client_id,
			project_id
		FROM contracts
		WHERE id = ?
	`,
		contractID,
	).Scan(
		&storedClientID,
		&storedProjectID,
	); err != nil {
		t.Fatalf(
			"query unchanged contract: %v",
			err,
		)
	}

	if storedClientID != firstClientID {
		t.Fatalf(
			"expected client ID %d, got %d",
			firstClientID,
			storedClientID,
		)
	}

	if storedProjectID != firstProjectID {
		t.Fatalf(
			"expected project ID %d, got %d",
			firstProjectID,
			storedProjectID,
		)
	}
}
