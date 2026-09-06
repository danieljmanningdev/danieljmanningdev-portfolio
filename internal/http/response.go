package http

import (
	"bytes"
	"html/template"
	"net/http"
	"strconv"
)

// renderPublicHTML executes before committing a status or body. A template
// failure must not turn into a partial 200 response, including for HEAD.
func renderPublicHTML(w http.ResponseWriter, r *http.Request, tmpl *template.Template, data any) {
	var body bytes.Buffer
	if err := tmpl.ExecuteTemplate(&body, "base", data); err != nil {
		w.Header().Set("Cache-Control", "no-store")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(body.Len()))
	w.WriteHeader(http.StatusOK)
	if r.Method != http.MethodHead {
		_, _ = body.WriteTo(w)
	}
}
