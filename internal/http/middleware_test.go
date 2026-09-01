// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package http

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestLoggerLogsRequestDetails(t *testing.T) {
	var output bytes.Buffer

	logger := slog.New(
		slog.NewTextHandler(
			&output,
			nil,
		),
	)

	next := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
		},
	)

	handler := RequestLogger(
		logger,
		next,
	)

	req := httptest.NewRequest(
		http.MethodPost,
		"/dashboard/projects",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	logOutput := output.String()

	expectedValues := []string{
		"msg=\"http request\"",
		"request_id=",
		"method=POST",
		"path=/dashboard/projects",
		"status=201",
		"duration=",
	}

	for _, expected := range expectedValues {
		if !strings.Contains(logOutput, expected) {
			t.Fatalf(
				"expected log output to contain %q, got %q",
				expected,
				logOutput,
			)
		}
	}

	if rec.Header().Get(requestIDHeader) == "" {
		t.Fatal("expected response request ID")
	}
}

func TestRequestLoggerDefaultsStatusToOK(t *testing.T) {
	var output bytes.Buffer

	logger := slog.New(
		slog.NewTextHandler(
			&output,
			nil,
		),
	)

	next := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("ok"))
		},
	)

	handler := RequestLogger(
		logger,
		next,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/health",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !strings.Contains(
		output.String(),
		"status=200",
	) {
		t.Fatalf(
			"expected status 200 in log output, got %q",
			output.String(),
		)
	}
}

func TestRequestLoggerPreservesResponseStatus(t *testing.T) {
	logger := slog.New(
		slog.NewTextHandler(
			&bytes.Buffer{},
			nil,
		),
	)

	next := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(
				w,
				"not found",
				http.StatusNotFound,
			)
		},
	)

	handler := RequestLogger(
		logger,
		next,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/missing",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNotFound,
			rec.Code,
		)
	}
}

func TestRequestLoggerUsesErrorLevelForServerErrors(
	t *testing.T,
) {
	var output bytes.Buffer

	logger := slog.New(
		slog.NewTextHandler(
			&output,
			nil,
		),
	)

	handler := RequestLogger(
		logger,
		http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			http.Error(
				w,
				"failed",
				http.StatusInternalServerError,
			)
		}),
	)

	handler.ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest(
			http.MethodGet,
			"/failure",
			nil,
		),
	)

	logOutput := output.String()
	if !strings.Contains(logOutput, "level=ERROR") ||
		!strings.Contains(logOutput, "status=500") ||
		!strings.Contains(logOutput, "msg=\"http request failed\"") {
		t.Fatalf(
			"expected structured server-error log, got %q",
			logOutput,
		)
	}
}

func TestRequestIDPreservesValidIncomingValue(
	t *testing.T,
) {
	const requestID = "external-request-123"

	handler := RequestID(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		if got := RequestIDFromContext(r.Context()); got != requestID {
			t.Fatalf(
				"expected context request ID %q, got %q",
				requestID,
				got,
			)
		}
	}))

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)
	req.Header.Set(requestIDHeader, requestID)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if got := rec.Header().Get(requestIDHeader); got != requestID {
		t.Fatalf(
			"expected response request ID %q, got %q",
			requestID,
			got,
		)
	}
}

func TestRequestIDReplacesInvalidIncomingValue(
	t *testing.T,
) {
	handler := RequestID(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		if RequestIDFromContext(r.Context()) == "" {
			t.Fatal("expected generated request ID in context")
		}
	}))

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)
	req.Header.Set(requestIDHeader, "bad\nrequest")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	got := rec.Header().Get(requestIDHeader)
	if got == "" || got == "bad\nrequest" {
		t.Fatalf("expected safe generated request ID, got %q", got)
	}
}

func TestRequestLoggerRecoversPanics(t *testing.T) {
	var output bytes.Buffer

	logger := slog.New(
		slog.NewTextHandler(
			&output,
			nil,
		),
	)

	handler := RequestLogger(
		logger,
		http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			panic("test panic")
		}),
	)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(
		rec,
		httptest.NewRequest(
			http.MethodGet,
			"/panic",
			nil,
		),
	)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status 500, got %d",
			rec.Code,
		)
	}

	logOutput := output.String()
	for _, expected := range []string{
		"msg=\"panic recovered\"",
		"panic=\"test panic\"",
		"request_id=",
		"msg=\"http request failed\"",
	} {
		if !strings.Contains(logOutput, expected) {
			t.Fatalf(
				"expected panic log to contain %q, got %q",
				expected,
				logOutput,
			)
		}
	}
}
