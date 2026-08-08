package http

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	templateDir := t.TempDir()

	for _, dir := range []string{
		"layouts",
		"components",
		"pages",
	} {
		if err := os.MkdirAll(filepath.Join(templateDir, dir), 0755); err != nil {
			t.Fatal(err)
		}
	}

	files := map[string]string{
		"layouts/base.html":          `{{define "base"}}{{template "home" .}}{{end}}`,
		"components/navigation.html": `{{define "navigation"}}nav{{end}}`,
		"components/footer.html":     `{{define "footer"}}footer{{end}}`,
		"pages/home.html":            `{{define "home"}}Daniel Manning{{end}}`,
	}

	for name, contents := range files {
		path := filepath.Join(templateDir, name)

		if err := os.WriteFile(path, []byte(contents), 0644); err != nil {
			t.Fatal(err)
		}
	}

	handler, err := NewHomeHandler(templateDir)
	if err != nil {
		t.Fatalf("create home handler: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if contentType := rec.Header().Get("Content-Type"); contentType != "text/html; charset=utf-8" {
		t.Fatalf("expected HTML content type, got %q", contentType)
	}

	if body := rec.Body.String(); body != "Daniel Manning" {
		t.Fatalf("expected homepage content, got %q", body)
	}
}

func TestHomeHandlerRejectsNonGet(t *testing.T) {
	templateDir := t.TempDir()

	for _, dir := range []string{
		"layouts",
		"components",
		"pages",
	} {
		if err := os.MkdirAll(filepath.Join(templateDir, dir), 0755); err != nil {
			t.Fatal(err)
		}
	}

	files := map[string]string{
		"layouts/base.html":          `{{define "base"}}{{template "home" .}}{{end}}`,
		"components/navigation.html": `{{define "navigation"}}nav{{end}}`,
		"components/footer.html":     `{{define "footer"}}footer{{end}}`,
		"pages/home.html":            `{{define "home"}}home{{end}}`,
	}

	for name, contents := range files {
		path := filepath.Join(templateDir, name)

		if err := os.WriteFile(path, []byte(contents), 0644); err != nil {
			t.Fatal(err)
		}
	}

	handler, err := NewHomeHandler(templateDir)
	if err != nil {
		t.Fatalf("create home handler: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rec.Code)
	}
}
