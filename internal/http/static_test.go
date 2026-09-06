package http

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStaticFilesDoNotExposeDirectoriesOrDotfiles(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "images"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"style.css", ".secret", "images/.secret"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte("body {}"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	handler := StaticFiles(root)
	for _, path := range []string{"/", "/images/", "/.secret", "/images/.secret", "/%2esecret", "/missing.css"} {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		if rec.Code != http.StatusNotFound {
			t.Errorf("%s: status %d", path, rec.Code)
		}
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/style.css", nil))
	if rec.Code != http.StatusOK || rec.Body.String() != "body {}" || !strings.HasPrefix(rec.Header().Get("Content-Type"), "text/css") {
		t.Fatalf("asset response: %d %v %q", rec.Code, rec.Header(), rec.Body.String())
	}
}
