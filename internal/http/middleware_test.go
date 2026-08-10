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
