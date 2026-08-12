package http

import "net/http"

func SecurityHeaders(
	production bool,
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
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
				"frame-ancestors 'none'",
			)

			w.Header().Set(
				"Referrer-Policy",
				"strict-origin-when-cross-origin",
			)

			w.Header().Set(
				"Permissions-Policy",
				"camera=(), microphone=(), geolocation=()",
			)

			if production {
				w.Header().Set(
					"Strict-Transport-Security",
					"max-age=31536000",
				)
			}

			next.ServeHTTP(w, r)
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
