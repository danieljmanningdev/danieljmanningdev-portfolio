package http

import "net/http"

// LimitRequestBody bounds form/request memory use. Known oversized bodies are
// rejected before handlers run; streaming bodies fail when read beyond the cap.
func LimitRequestBody(next http.Handler, maximum int64) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength > maximum {
			http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
			return
		}
		if r.Body != nil {
			r.Body = http.MaxBytesReader(w, r.Body, maximum)
		}
		next.ServeHTTP(w, r)
	})
}
