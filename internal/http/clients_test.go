package http

import (
	"context"
	"database/sql"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/database"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/repository"
)

func newTestClientsHandler(t *testing.T) (*ClientsHandler, *sql.DB) {
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

	_, err = db.SQL.Exec(`
		CREATE TABLE clients (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			company TEXT,
			status TEXT NOT NULL DEFAULT 'active',
			notes TEXT,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("create clients table: %v", err)
	}

	errorTemplate := template.Must(
		template.New("base").Parse(`
			{{define "base"}}
			{{.Error}}
			{{end}}
		`),
	)

	handler := &ClientsHandler{
		repository:         repository.NewClientRepository(db.SQL),
		newClientTemplates: errorTemplate,
	}

	return handler, db.SQL
}

func createTestClient(t *testing.T, db *sql.DB) int64 {
	t.Helper()

	result, err := db.Exec(`
		INSERT INTO clients (
			name,
			email,
			company,
			status,
			notes
		)
		VALUES (?, ?, ?, ?, ?)
	`,
		"Test Client",
		"test@example.com",
		"Test Company",
		"active",
		"Test notes",
	)
	if err != nil {
		t.Fatalf("insert test client: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("get test client id: %v", err)
	}

	return id
}

func formRequest(
	method string,
	target string,
	values url.Values,
) *http.Request {
	req := httptest.NewRequest(
		method,
		target,
		strings.NewReader(values.Encode()),
	)

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	return req
}

func TestClientsCreate(t *testing.T) {
	handler, db := newTestClientsHandler(t)

	form := url.Values{}
	form.Set("name", "Alice Smith")
	form.Set("email", "alice@example.com")
	form.Set("company", "Acme Ltd")
	form.Set("notes", "Important client")

	req := formRequest(
		http.MethodPost,
		"/dashboard/clients/new",
		form,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf(
			"expected status 303, got %d",
			rec.Code,
		)
	}

	location := rec.Header().Get("Location")
	if location == "" {
		t.Fatal("expected redirect location")
	}

	var (
		id      int64
		name    string
		email   string
		company string
		status  string
		notes   string
	)

	err := db.QueryRow(`
		SELECT
			id,
			name,
			email,
			company,
			status,
			notes
		FROM clients
		WHERE name = ?
	`, "Alice Smith").Scan(
		&id,
		&name,
		&email,
		&company,
		&status,
		&notes,
	)
	if err != nil {
		t.Fatalf("query created client: %v", err)
	}

	expectedLocation := "/dashboard/clients/" +
		strconv.FormatInt(id, 10)

	if location != expectedLocation {
		t.Fatalf(
			"expected location %q, got %q",
			expectedLocation,
			location,
		)
	}

	if name != "Alice Smith" {
		t.Fatalf(
			"expected name %q, got %q",
			"Alice Smith",
			name,
		)
	}

	if email != "alice@example.com" {
		t.Fatalf(
			"expected email %q, got %q",
			"alice@example.com",
			email,
		)
	}

	if company != "Acme Ltd" {
		t.Fatalf(
			"expected company %q, got %q",
			"Acme Ltd",
			company,
		)
	}

	if status != "active" {
		t.Fatalf(
			"expected status %q, got %q",
			"active",
			status,
		)
	}

	if notes != "Important client" {
		t.Fatalf(
			"expected notes %q, got %q",
			"Important client",
			notes,
		)
	}
}

func TestClientsCreateRejectsMissingName(t *testing.T) {
	handler, _ := newTestClientsHandler(t)

	form := url.Values{}
	form.Set("name", "")
	form.Set("email", "test@example.com")

	req := formRequest(
		http.MethodPost,
		"/dashboard/clients/new",
		form,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d",
			rec.Code,
		)
	}

	if !strings.Contains(
		rec.Body.String(),
		"Name is required.",
	) {
		t.Fatalf(
			"expected name validation message, got %q",
			rec.Body.String(),
		)
	}
}

func TestClientsCreateRejectsMissingEmail(t *testing.T) {
	handler, _ := newTestClientsHandler(t)

	form := url.Values{}
	form.Set("name", "Test Client")
	form.Set("email", "")

	req := formRequest(
		http.MethodPost,
		"/dashboard/clients/new",
		form,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d",
			rec.Code,
		)
	}

	if !strings.Contains(
		rec.Body.String(),
		"Email is required.",
	) {
		t.Fatalf(
			"expected email validation message, got %q",
			rec.Body.String(),
		)
	}
}

func TestClientsCreateRejectsInvalidEmail(t *testing.T) {
	handler, _ := newTestClientsHandler(t)

	form := url.Values{}
	form.Set("name", "Test Client")
	form.Set("email", "not-an-email")

	req := formRequest(
		http.MethodPost,
		"/dashboard/clients/new",
		form,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d",
			rec.Code,
		)
	}

	if !strings.Contains(
		rec.Body.String(),
		"Please enter a valid email address.",
	) {
		t.Fatalf(
			"expected invalid email message, got %q",
			rec.Body.String(),
		)
	}
}

func TestClientsUpdate(t *testing.T) {
	handler, db := newTestClientsHandler(t)

	id := createTestClient(t, db)

	form := url.Values{}
	form.Set("name", "Updated Client")
	form.Set("email", "updated@example.com")
	form.Set("company", "Updated Company")
	form.Set("status", "inactive")
	form.Set("notes", "Updated notes")

	req := formRequest(
		http.MethodPost,
		"/dashboard/clients/"+strconv.FormatInt(id, 10),
		form,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf(
			"expected status 303, got %d",
			rec.Code,
		)
	}

	expectedLocation := "/dashboard/clients/" +
		strconv.FormatInt(id, 10)

	if location := rec.Header().Get("Location"); location != expectedLocation {
		t.Fatalf(
			"expected location %q, got %q",
			expectedLocation,
			location,
		)
	}

	var (
		name    string
		email   string
		company string
		status  string
		notes   string
	)

	err := db.QueryRow(`
		SELECT
			name,
			email,
			company,
			status,
			notes
		FROM clients
		WHERE id = ?
	`, id).Scan(
		&name,
		&email,
		&company,
		&status,
		&notes,
	)
	if err != nil {
		t.Fatalf("query updated client: %v", err)
	}

	if name != "Updated Client" {
		t.Fatalf(
			"expected name %q, got %q",
			"Updated Client",
			name,
		)
	}

	if email != "updated@example.com" {
		t.Fatalf(
			"expected email %q, got %q",
			"updated@example.com",
			email,
		)
	}

	if company != "Updated Company" {
		t.Fatalf(
			"expected company %q, got %q",
			"Updated Company",
			company,
		)
	}

	if status != "inactive" {
		t.Fatalf(
			"expected status %q, got %q",
			"inactive",
			status,
		)
	}

	if notes != "Updated notes" {
		t.Fatalf(
			"expected notes %q, got %q",
			"Updated notes",
			notes,
		)
	}
}

func TestClientsUpdateRejectsMissingName(t *testing.T) {
	handler, db := newTestClientsHandler(t)

	id := createTestClient(t, db)

	form := url.Values{}
	form.Set("name", "")
	form.Set("email", "updated@example.com")
	form.Set("status", "active")

	req := formRequest(
		http.MethodPost,
		"/dashboard/clients/"+strconv.FormatInt(id, 10),
		form,
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
		"Name is required.",
	) {
		t.Fatalf(
			"expected name validation message, got %q",
			rec.Body.String(),
		)
	}
}

func TestClientsUpdateRejectsInvalidStatus(t *testing.T) {
	handler, db := newTestClientsHandler(t)

	id := createTestClient(t, db)

	form := url.Values{}
	form.Set("name", "Updated Client")
	form.Set("email", "updated@example.com")
	form.Set("status", "banana")

	req := formRequest(
		http.MethodPost,
		"/dashboard/clients/"+strconv.FormatInt(id, 10),
		form,
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
		"Invalid client status.",
	) {
		t.Fatalf(
			"expected status validation message, got %q",
			rec.Body.String(),
		)
	}
}

func TestClientsUpdateRejectsInvalidEmail(t *testing.T) {
	handler, db := newTestClientsHandler(t)

	id := createTestClient(t, db)

	form := url.Values{}
	form.Set("name", "Updated Client")
	form.Set("email", "not-an-email")
	form.Set("status", "active")

	req := formRequest(
		http.MethodPost,
		"/dashboard/clients/"+strconv.FormatInt(id, 10),
		form,
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
		"Please enter a valid email address.",
	) {
		t.Fatalf(
			"expected invalid email message, got %q",
			rec.Body.String(),
		)
	}
}

func TestClientsUpdateReturnsNotFoundForMissingClient(t *testing.T) {
	handler, _ := newTestClientsHandler(t)

	form := url.Values{}
	form.Set("name", "Missing Client")
	form.Set("email", "missing@example.com")
	form.Set("status", "active")

	req := formRequest(
		http.MethodPost,
		"/dashboard/clients/999",
		form,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status 404, got %d",
			rec.Code,
		)
	}
}

func TestClientsDelete(t *testing.T) {
	handler, db := newTestClientsHandler(t)

	id := createTestClient(t, db)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/dashboard/clients/"+strconv.FormatInt(id, 10),
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf(
			"expected status 303, got %d",
			rec.Code,
		)
	}

	if location := rec.Header().Get("Location"); location != "/dashboard/clients" {
		t.Fatalf(
			"expected location %q, got %q",
			"/dashboard/clients",
			location,
		)
	}

	var count int

	err := db.QueryRow(
		"SELECT COUNT(*) FROM clients WHERE id = ?",
		id,
	).Scan(&count)
	if err != nil {
		t.Fatalf("check deleted client: %v", err)
	}

	if count != 0 {
		t.Fatalf(
			"expected client to be deleted, found %d rows",
			count,
		)
	}
}

func TestClientsDeleteReturnsNotFoundForMissingClient(t *testing.T) {
	handler, _ := newTestClientsHandler(t)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/dashboard/clients/999",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status 404, got %d",
			rec.Code,
		)
	}
}
