package http

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type HomeHandler struct {
	template *template.Template
}

func NewHomeHandler(templateDir string) (*HomeHandler, error) {
	tmpl, err := template.ParseFiles(
		filepath.Join(templateDir, "layouts", "base.html"),
		filepath.Join(templateDir, "components", "header.html"),
		filepath.Join(templateDir, "components", "footer.html"),
		filepath.Join(templateDir, "pages", "home.html"),
	)
	if err != nil {
		return nil, err
	}

	return &HomeHandler{
		template: tmpl,
	}, nil
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := h.template.ExecuteTemplate(w, "base", nil); err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
	}
}
