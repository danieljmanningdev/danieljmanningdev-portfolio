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
		"pages/public",
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
{{template "content" .}}
{{template "footer" .}}
{{end}}`,

		"components/header.html": `{{define "header"}}header{{end}}`,

		"components/footer.html": `{{define "footer"}}footer{{end}}`,

		"pages/public/home.html": `{{define "content"}}Daniel Manning{{end}}`,

		"pages/public/404.html": `{{define "content"}}Page not found: {{.Path}}{{end}}`,
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
		"Web Designer, Developer &amp; Product Designer | Daniel J. Manning",
		"Web designer, developer and digital product designer creating fast, accessible and maintainable websites and software for businesses in Leeds and beyond.",
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

func TestHomeHandlerRendersNotFoundPage(
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
		"/missing-page",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(
		rec,
		req,
	)

	if rec.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status 404, got %d",
			rec.Code,
		)
	}

	if !strings.Contains(
		rec.Body.String(),
		"Page not found: /missing-page",
	) {
		t.Fatalf(
			"unexpected not-found page %q",
			rec.Body.String(),
		)
	}

	if robots := rec.Header().Get(
		"X-Robots-Tag",
	); robots != "noindex, nofollow" {
		t.Fatalf(
			"expected noindex header, got %q",
			robots,
		)
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
