package http

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type HomeHandler struct {
	template *template.Template
}

func NewHomeHandler(
	templateDir string,
) (*HomeHandler, error) {
	tmpl, err := template.New("base").ParseFiles(
		filepath.Join(
			templateDir,
			"layouts",
			"base.html",
		),
		filepath.Join(
			templateDir,
			"components",
			"header.html",
		),
		filepath.Join(
			templateDir,
			"components",
			"footer.html",
		),
		filepath.Join(
			templateDir,
			"pages",
			"home.html",
		),
	)
	if err != nil {
		return nil, err
	}

	return &HomeHandler{
		template: tmpl,
	}, nil
}

func (h *HomeHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(
			w,
			http.StatusText(
				http.StatusMethodNotAllowed,
			),
			http.StatusMethodNotAllowed,
		)
		return
	}

	data := publicPageData{
		Title:       "Daniel J. Manning — Digital Product Designer & Engineer",
		Description: "Digital product design and full-stack engineering focused on fast, purposeful software, thoughtful user experiences and maintainable systems.",
		OGTitle:     "Daniel J. Manning — Digital Product Designer & Engineer",
		OGType:      "website",
	}

	w.Header().Set(
		"Content-Type",
		"text/html; charset=utf-8",
	)
	w.Header().Set(
		"Link",
		`<https://danieljmanningdev.com/>; rel="canonical"`,
	)

	if err := h.template.ExecuteTemplate(
		w,
		"base",
		data,
	); err != nil {
		http.Error(
			w,
			http.StatusText(
				http.StatusInternalServerError,
			),
			http.StatusInternalServerError,
		)
	}
}
