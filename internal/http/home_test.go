package http

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func createTestTemplates(
	t *testing.T,
) string {
	t.Helper()

	templateDir := t.TempDir()

	for _, dir := range []string{
		"layouts",
		"components",
		"pages",
	} {
		if err := os.MkdirAll(
			filepath.Join(
				templateDir,
				dir,
			),
			0755,
		); err != nil {
			t.Fatal(err)
		}
	}

	files := map[string]string{
		"layouts/base.html": `{{define "base"}}
<title>{{.Title}}</title>
<meta name="description" content="{{.Description}}">
{{template "header" .}}
{{template "home" .}}
{{template "footer" .}}
{{end}}`,

		"components/header.html": `{{define "header"}}header{{end}}`,

		"components/footer.html": `{{define "footer"}}footer{{end}}`,

		"pages/home.html": `{{define "home"}}Daniel Manning{{end}}`,
	}

	for name, contents := range files {
		path := filepath.Join(
			templateDir,
			name,
		)

		if err := os.WriteFile(
			path,
			[]byte(contents),
			0644,
		); err != nil {
			t.Fatal(err)
		}
	}

	return templateDir
}

func TestHomeHandler(
	t *testing.T,
) {
	templateDir := createTestTemplates(t)

	handler, err := NewHomeHandler(
		templateDir,
	)
	if err != nil {
		t.Fatalf(
			"create home handler: %v",
			err,
		)
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(
		rec,
		req,
	)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d",
			rec.Code,
		)
	}

	if contentType :=
		rec.Header().Get(
			"Content-Type",
		); contentType !=
		"text/html; charset=utf-8" {
		t.Fatalf(
			"expected HTML content type, got %q",
			contentType,
		)
	}

	body := rec.Body.String()

	for _, expected := range []string{
		"Daniel Manning",
		"Digital Product Designer &amp; Engineer",
		"Digital product design and full-stack engineering",
	} {
		if !strings.Contains(
			body,
			expected,
		) {
			t.Fatalf(
				"expected homepage output to contain %q, got %q",
				expected,
				body,
			)
		}
	}
}

func TestHomeHandlerRejectsNonGet(
	t *testing.T,
) {
	templateDir := createTestTemplates(t)

	handler, err := NewHomeHandler(
		templateDir,
	)
	if err != nil {
		t.Fatalf(
			"create home handler: %v",
			err,
		)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(
		rec,
		req,
	)

	if rec.Code !=
		http.StatusMethodNotAllowed {
		t.Fatalf(
			"expected status 405, got %d",
			rec.Code,
		)
	}
}
