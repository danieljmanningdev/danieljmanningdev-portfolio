// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package http

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func createPortfolioCaseStudyTestTemplates(
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

		"pages/public/portfolio.html": `{{define "content"}}
Portfolio & Client Workspace
{{end}}`,
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

func TestPortfolioCaseStudyHandler(
	t *testing.T,
) {
	handler, err :=
		NewPortfolioCaseStudyHandler(
			createPortfolioCaseStudyTestTemplates(t),
		)
	if err != nil {
		t.Fatalf(
			"create handler: %v",
			err,
		)
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/work/portfolio",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected 200, got %d",
			rec.Code,
		)
	}

	body := rec.Body.String()

	for _, expected := range []string{
		"Portfolio &amp; Client Workspace",
		"secure portfolio and client workspace",
	} {
		if !strings.Contains(
			body,
			expected,
		) {
			t.Fatalf(
				"expected body to contain %q, got %q",
				expected,
				body,
			)
		}
	}
}

func TestPortfolioCaseStudyRejectsNonGet(
	t *testing.T,
) {
	handler, err :=
		NewPortfolioCaseStudyHandler(
			createPortfolioCaseStudyTestTemplates(t),
		)
	if err != nil {
		t.Fatalf(
			"create handler: %v",
			err,
		)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/work/portfolio",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code !=
		http.StatusMethodNotAllowed {
		t.Fatalf(
			"expected 405, got %d",
			rec.Code,
		)
	}
}
