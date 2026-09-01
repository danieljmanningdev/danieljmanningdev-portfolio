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
	"strings"
	"testing"
)

func TestSecurityHeadersDevelopment(
	t *testing.T,
) {
	var nonceFromContext string

	handler := SecurityHeaders(
		false,
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				nonceFromContext = CSPNonceFromContext(
					r.Context(),
				)
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
		"X-Content-Type-Options":       "nosniff",
		"X-Frame-Options":              "DENY",
		"Referrer-Policy":              "strict-origin-when-cross-origin",
		"Permissions-Policy":           "camera=(), microphone=(), geolocation=()",
		"Cross-Origin-Opener-Policy":   "same-origin",
		"Cross-Origin-Resource-Policy": "same-origin",
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

	if nonceFromContext == "" {
		t.Fatal("expected CSP nonce in request context")
	}

	policy := headers.Get("Content-Security-Policy")
	for _, directive := range []string{
		"default-src 'self'",
		"base-uri 'self'",
		"form-action 'self'",
		"frame-ancestors 'none'",
		"object-src 'none'",
		"img-src 'self' data:",
		"script-src 'self' 'nonce-" + nonceFromContext + "'",
		"style-src 'self'",
	} {
		if !strings.Contains(policy, directive) {
			t.Fatalf(
				"expected CSP to contain %q, got %q",
				directive,
				policy,
			)
		}
	}

	if strings.Contains(policy, "upgrade-insecure-requests") {
		t.Fatalf(
			"expected development CSP without upgrade directive, got %q",
			policy,
		)
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
	); got != "max-age=31536000; includeSubDomains" {
		t.Fatalf(
			"unexpected HSTS header %q",
			got,
		)
	}

	policy := rec.Header().Get("Content-Security-Policy")
	if !strings.Contains(policy, "upgrade-insecure-requests") {
		t.Fatalf(
			"expected production upgrade directive, got %q",
			policy,
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
