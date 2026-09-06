package http

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestBodyLimit(t *testing.T) {
	for _, tc := range []struct {
		name  string
		size  int
		known bool
		want  int
	}{
		{"within limit", 64, true, 204}, {"declared excess", 65, true, 413}, {"streamed excess", 65, false, 413},
	} {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			handler := LimitRequestBody(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				body, err := io.ReadAll(r.Body)
				if len(body) > 64 {
					t.Fatal("body grew past limit")
				}
				var tooLarge *http.MaxBytesError
				if errors.As(err, &tooLarge) {
					w.WriteHeader(413)
					return
				}
				if err != nil {
					t.Fatal(err)
				}
				w.WriteHeader(204)
			}), 64)
			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(strings.Repeat("a", tc.size)))
			if !tc.known {
				req.ContentLength = -1
			}
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			if rec.Code != tc.want {
				t.Fatalf("status = %d, want %d", rec.Code, tc.want)
			}
			if tc.known && tc.size > 64 && called {
				t.Fatal("oversized request reached handler")
			}
		})
	}
}
