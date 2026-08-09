package http

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestClientsDeleteHTMXRedirect(t *testing.T) {
	handler, db := newTestClientsHandler(t)
	id := createTestClient(t, db)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/dashboard/clients/"+strconv.FormatInt(id, 10),
		nil,
	)
	req.Header.Set("HX-Request", "true")

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d",
			rec.Code,
		)
	}

	if redirect := rec.Header().Get(
		"HX-Redirect",
	); redirect != "/dashboard/clients" {
		t.Fatalf(
			"expected HX-Redirect %q, got %q",
			"/dashboard/clients",
			redirect,
		)
	}

	var count int

	if err := db.QueryRow(
		"SELECT COUNT(*) FROM clients WHERE id = ?",
		id,
	).Scan(&count); err != nil {
		t.Fatalf("check deleted client: %v", err)
	}

	if count != 0 {
		t.Fatalf(
			"expected client to be deleted, found %d rows",
			count,
		)
	}
}
