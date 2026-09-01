package http

import (
	"html/template"
	"net/http"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/rendering"
)

type WebDevelopmentHandler struct {
	template *template.Template
}

func NewWebDevelopmentHandler(
	templateDir string,
) (*WebDevelopmentHandler, error) {
	tmpl, err := rendering.LoadPageTemplate(
		templateDir,
		"public/web-development.html",
	)
	if err != nil {
		return nil, err
	}

	return &WebDevelopmentHandler{
		template: tmpl,
	}, nil
}

func (h *WebDevelopmentHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	data := newPublicPageData(
		"Web Development | Daniel J. Manning",
		"Full-stack web development for fast, accessible and maintainable websites, web applications and digital products.",
		"/web-development/",
		"website",
		nil,
	).withRequest(r)

	w.Header().Set(
		"Content-Type",
		"text/html; charset=utf-8",
	)

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
