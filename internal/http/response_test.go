package http

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func publicResponseHandlers(tmpl *template.Template) []struct {
	path    string
	handler http.Handler
} {
	return []struct {
		path    string
		handler http.Handler
	}{
		{"/", &HomeHandler{homeTemplate: tmpl}},
		{"/work/portfolio", &PortfolioCaseStudyHandler{template: tmpl}},
		{"/web-design/", &PublicPageHandler{template: tmpl, definition: PublicPageDefinition{Path: "/web-design/", Title: "Web design"}}},
	}
}

func TestPublicRenderingFailureDoesNotWritePartialSuccess(t *testing.T) {
	broken := template.Must(template.New("base").Parse(`partial-success{{.Missing.Field}}`))
	for _, tc := range publicResponseHandlers(broken) {
		t.Run(tc.path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			tc.handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))
			if rec.Code != http.StatusInternalServerError {
				t.Fatalf("status = %d", rec.Code)
			}
			if strings.Contains(rec.Body.String(), "partial-success") {
				t.Fatal("partial template escaped")
			}
			if rec.Header().Get("Cache-Control") != "no-store" {
				t.Fatal("failed response may be cached")
			}
		})
	}
}

func TestPublicResponsesSupportHEADAndRejectWrites(t *testing.T) {
	tmpl := template.Must(template.New("base").Parse(`<h1>{{.Title}}</h1>`))
	for _, tc := range publicResponseHandlers(tmpl) {
		t.Run(tc.path, func(t *testing.T) {
			get := httptest.NewRecorder()
			tc.handler.ServeHTTP(get, httptest.NewRequest(http.MethodGet, tc.path, nil))
			head := httptest.NewRecorder()
			tc.handler.ServeHTTP(head, httptest.NewRequest(http.MethodHead, tc.path, nil))
			if head.Code != http.StatusOK || head.Body.Len() != 0 {
				t.Fatalf("HEAD: %d %q", head.Code, head.Body.String())
			}
			if head.Header().Get("Content-Length") != strconv.Itoa(get.Body.Len()) {
				t.Fatal("HEAD representation length differs from GET")
			}
			post := httptest.NewRecorder()
			tc.handler.ServeHTTP(post, httptest.NewRequest(http.MethodPost, tc.path, nil))
			if post.Code != http.StatusMethodNotAllowed || post.Header().Get("Allow") != "GET, HEAD" {
				t.Fatalf("POST: %d, Allow=%q", post.Code, post.Header().Get("Allow"))
			}
		})
	}
}

func TestPublicPageRejectsUnknownSubpath(t *testing.T) {
	handler := &PublicPageHandler{definition: PublicPageDefinition{Path: "/web-design/"}}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/web-design/not-a-page", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("unknown subpath status = %d", rec.Code)
	}
}
