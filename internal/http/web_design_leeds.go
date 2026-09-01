package http

import (
	"html/template"
	"net/http"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/rendering"
)

type WebDesignLeedsHandler struct {
	template *template.Template
}

func NewWebDesignLeedsHandler(templateDir string) (*WebDesignLeedsHandler, error) {
	tmpl, err := rendering.LoadPageTemplate(
		templateDir,
		"public/web-design-leeds.html",
	)
	if err != nil {
		return nil, err
	}

	return &WebDesignLeedsHandler{
		template: tmpl,
	}, nil
}

func (h *WebDesignLeedsHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	data := newPublicPageData(
		"Web Design & Development in Leeds | Daniel J. Manning",
		"Web design and development in Leeds for businesses that need fast, accessible and maintainable websites and digital products.",
		"/web-design-leeds/",
		"website",
		nil,
	).withRequest(r)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := h.template.ExecuteTemplate(
		w,
		"base",
		data,
	); err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
	}
}
