package http

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
	"unicode"
)

const requestIDHeader = "X-Request-ID"

type requestIDContextKey struct{}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	if r.status != 0 {
		return
	}

	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Write(body []byte) (int, error) {
	if r.status == 0 {
		r.WriteHeader(http.StatusOK)
	}

	return r.ResponseWriter.Write(body)
}

func (r *statusRecorder) Unwrap() http.ResponseWriter {
	return r.ResponseWriter
}

func RequestIDFromContext(ctx context.Context) string {
	requestID, _ := ctx.Value(requestIDContextKey{}).(string)
	return requestID
}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		requestID := strings.TrimSpace(
			r.Header.Get(requestIDHeader),
		)
		if !validRequestID(requestID) {
			requestID = newRequestID()
		}

		w.Header().Set(requestIDHeader, requestID)

		ctx := context.WithValue(
			r.Context(),
			requestIDContextKey{},
			requestID,
		)

		next.ServeHTTP(
			w,
			r.WithContext(ctx),
		)
	})
}

func RequestLogger(
	logger *slog.Logger,
	next http.Handler,
) http.Handler {
	recoveredHandler := RecoverPanics(logger, next)

	loggedHandler := http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		start := time.Now()

		recorder := &statusRecorder{
			ResponseWriter: w,
		}

		recoveredHandler.ServeHTTP(recorder, r)

		status := recorder.status
		if status == 0 {
			status = http.StatusOK
		}

		attributes := []any{
			"request_id", RequestIDFromContext(r.Context()),
			"method", r.Method,
			"path", r.URL.Path,
			"status", status,
			"duration", time.Since(start),
		}

		switch {
		case status >= http.StatusInternalServerError:
			logger.Error("http request failed", attributes...)
		case status >= http.StatusBadRequest:
			logger.Warn("http request rejected", attributes...)
		default:
			logger.Info("http request", attributes...)
		}
	})

	return RequestID(loggedHandler)
}

func RecoverPanics(
	logger *slog.Logger,
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		defer func() {
			recovered := recover()
			if recovered == nil {
				return
			}

			if recovered == http.ErrAbortHandler {
				panic(recovered)
			}

			logger.Error(
				"panic recovered",
				"request_id", RequestIDFromContext(r.Context()),
				"method", r.Method,
				"path", r.URL.Path,
				"panic", recovered,
				"stack", string(debug.Stack()),
			)

			w.Header().Set(
				"Connection",
				"close",
			)

			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
		}()

		next.ServeHTTP(w, r)
	})
}

func NoIndex(next http.Handler) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		w.Header().Set(
			"X-Robots-Tag",
			"noindex, nofollow",
		)

		next.ServeHTTP(w, r)
	})
}

func newRequestID() string {
	value := make([]byte, 16)
	if _, err := rand.Read(value); err != nil {
		return base64.RawURLEncoding.EncodeToString(
			[]byte(time.Now().UTC().Format(time.RFC3339Nano)),
		)
	}

	return base64.RawURLEncoding.EncodeToString(value)
}

func validRequestID(value string) bool {
	if len(value) < 8 || len(value) > 128 {
		return false
	}

	for _, character := range value {
		if unicode.IsLetter(character) ||
			unicode.IsDigit(character) ||
			strings.ContainsRune("-_.:", character) {
			continue
		}

		return false
	}

	return true
}
