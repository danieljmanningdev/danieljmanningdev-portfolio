package http

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

type cspNonceContextKey struct{}

func CSPNonceFromContext(ctx context.Context) string {
	nonce, _ := ctx.Value(cspNonceContextKey{}).(string)
	return nonce
}

func SecurityHeaders(
	production bool,
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			nonce, err := generateCSPNonce()
			if err != nil {
				http.Error(
					w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError,
				)
				return
			}

			w.Header().Set(
				"X-Content-Type-Options",
				"nosniff",
			)

			w.Header().Set(
				"X-Frame-Options",
				"DENY",
			)

			w.Header().Set(
				"Content-Security-Policy",
				contentSecurityPolicy(nonce, production),
			)

			w.Header().Set(
				"Referrer-Policy",
				"strict-origin-when-cross-origin",
			)

			w.Header().Set(
				"Permissions-Policy",
				"camera=(), microphone=(), geolocation=()",
			)

			w.Header().Set(
				"Cross-Origin-Opener-Policy",
				"same-origin",
			)

			w.Header().Set(
				"Cross-Origin-Resource-Policy",
				"same-origin",
			)

			if production {
				w.Header().Set(
					"Strict-Transport-Security",
					"max-age=31536000; includeSubDomains",
				)
			}

			ctx := context.WithValue(
				r.Context(),
				cspNonceContextKey{},
				nonce,
			)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		},
	)
}

func NoStore(
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			w.Header().Set(
				"Cache-Control",
				"no-store",
			)

			w.Header().Set(
				"Pragma",
				"no-cache",
			)

			next.ServeHTTP(w, r)
		},
	)
}

func contentSecurityPolicy(
	nonce string,
	production bool,
) string {
	directives := []string{
		"default-src 'self'",
		"base-uri 'self'",
		"connect-src 'self'",
		"font-src 'self'",
		"form-action 'self'",
		"frame-ancestors 'none'",
		"frame-src 'none'",
		"img-src 'self' data:",
		"manifest-src 'self'",
		"object-src 'none'",
		fmt.Sprintf("script-src 'self' 'nonce-%s'", nonce),
		"style-src 'self'",
	}

	if production {
		directives = append(
			directives,
			"upgrade-insecure-requests",
		)
	}

	return strings.Join(directives, "; ")
}

func generateCSPNonce() (string, error) {
	value := make([]byte, 16)
	if _, err := rand.Read(value); err != nil {
		return "", fmt.Errorf("generate CSP nonce: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(value), nil
}
