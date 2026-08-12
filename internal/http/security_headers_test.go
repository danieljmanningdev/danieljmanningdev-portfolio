package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecurityHeadersDevelopment(
	t *testing.T,
) {
	handler := SecurityHeaders(
		false,
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				w.WriteHeader(
					http.StatusOK,
				)
			},
		),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	headers := rec.Header()

	expected := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
		"Content-Security-Policy": "frame-ancestors 'none'",
		"Referrer-Policy": "strict-origin-when-cross-origin",
		"Permissions-Policy": "camera=(), microphone=(), geolocation=()",
	}

	for name, want := range expected {
		if got := headers.Get(name); got != want {
			t.Fatalf(
				"expected %s=%q, got %q",
				name,
				want,
				got,
			)
		}
	}

	if got := headers.Get(
		"Strict-Transport-Security",
	); got != "" {
		t.Fatalf(
			"expected no HSTS in development, got %q",
			got,
		)
	}
}

func TestSecurityHeadersProductionAddsHSTS(
	t *testing.T,
) {
	handler := SecurityHeaders(
		true,
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				w.WriteHeader(
					http.StatusOK,
				)
			},
		),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if got := rec.Header().Get(
		"Strict-Transport-Security",
	); got != "max-age=31536000" {
		t.Fatalf(
			"unexpected HSTS header %q",
			got,
		)
	}
}

func TestNoStore(
	t *testing.T,
) {
	handler := NoStore(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				w.WriteHeader(
					http.StatusOK,
				)
			},
		),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/dashboard/",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if got := rec.Header().Get(
		"Cache-Control",
	); got != "no-store" {
		t.Fatalf(
			"expected Cache-Control no-store, got %q",
			got,
		)
	}

	if got := rec.Header().Get(
		"Pragma",
	); got != "no-cache" {
		t.Fatalf(
			"expected Pragma no-cache, got %q",
			got,
		)
	}
}
