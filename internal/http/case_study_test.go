package http

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestCaseStudyHandler(t *testing.T) {
	handler, err := NewCaseStudyHandler(
		filepath.Join("..", "..", "web", "templates"),
	)
	if err != nil {
		t.Fatalf("create case study handler: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle(
		"GET /work/{slug}/{$}",
		handler,
	)

	tests := []struct {
		name       string
		path       string
		status     int
		contains   []string
		notContain []string
	}{
		{
			name:   "salon rebuild",
			path:   "/work/salon-rebuild/",
			status: http.StatusOK,
			contains: []string{
				"Salon",
				"Rebuild",
				"Full Site Rebuild",
				"Luxury without excess",
				"View live site",
			},
		},
		{
			name:   "wireframe kit",
			path:   "/work/wireframe-kit/",
			status: http.StatusOK,
			contains: []string{
				"Wireframe",
				"Kit",
				"40+",
				"Structure before styling",
				"View on Figma",
			},
		},
		{
			name:   "unknown case study",
			path:   "/work/does-not-exist/",
			status: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodGet,
				tt.path,
				nil,
			)

			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != tt.status {
				t.Fatalf(
					"expected status %d, got %d",
					tt.status,
					rec.Code,
				)
			}

			body := rec.Body.String()

			for _, expected := range tt.contains {
				if !strings.Contains(body, expected) {
					t.Errorf(
						"expected response to contain %q",
						expected,
					)
				}
			}

			for _, unexpected := range tt.notContain {
				if strings.Contains(body, unexpected) {
					t.Errorf(
						"response unexpectedly contained %q",
						unexpected,
					)
				}
			}
		})
	}
}

func TestCaseStudyCanonicalHeader(t *testing.T) {
	handler, err := NewCaseStudyHandler(
		filepath.Join("..", "..", "web", "templates"),
	)
	if err != nil {
		t.Fatalf("create case study handler: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle(
		"GET /work/{slug}/{$}",
		handler,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/work/wireframe-kit/",
		nil,
	)

	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	expected := `<https://danieljmanningdev.com/work/wireframe-kit/>; rel="canonical"`

	if got := rec.Header().Get("Link"); got != expected {
		t.Fatalf(
			"expected canonical header %q, got %q",
			expected,
			got,
		)
	}
}

func TestCaseStudyIncludesRelatedLinks(t *testing.T) {
	handler, err := NewCaseStudyHandler(
		filepath.Join("..", "..", "web", "templates"),
	)
	if err != nil {
		t.Fatalf("create case study handler: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle(
		"GET /work/{slug}/{$}",
		handler,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/work/wireframe-kit/",
		nil,
	)

	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	body := rec.Body.String()

	for _, expected := range []string{
		"Continue exploring",
		"UI &amp; UX Design",
		"Portfolio &amp; Client Workspace",
		"Web Design",
	} {
		if !strings.Contains(body, expected) {
			t.Errorf(
				"expected related content to contain %q",
				expected,
			)
		}
	}
}
