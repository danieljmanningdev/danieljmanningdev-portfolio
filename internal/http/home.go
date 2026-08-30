// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package http

import (
	"html/template"
	"net/http"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/rendering"
)

type HomeHandler struct {
	homeTemplate     *template.Template
	notFoundTemplate *template.Template
}

func NewHomeHandler(
	templateDir string,
) (*HomeHandler, error) {
	homeTemplate, err := rendering.LoadPageTemplate(
		templateDir,
		"public/home.html",
	)
	if err != nil {
		return nil, err
	}

	notFoundTemplate, err := rendering.LoadPageTemplate(
		templateDir,
		"public/404.html",
	)
	if err != nil {
		return nil, err
	}

	return &HomeHandler{
		homeTemplate:     homeTemplate,
		notFoundTemplate: notFoundTemplate,
	}, nil
}

func (h *HomeHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.URL.Path != "/" {
		renderNotFoundPage(
			w,
			h.notFoundTemplate,
			r.URL.Path,
		)
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

	data := newPublicPageData(
		"Web Designer, Developer & Product Designer | Daniel J. Manning",
		"Web designer, developer and digital product designer creating fast, accessible and maintainable websites and software for businesses in Leeds and beyond.",
		"/",
		"website",
		personStructuredData(),
	).withRequest(r)

	w.Header().Set(
		"Content-Type",
		"text/html; charset=utf-8",
	)
	w.Header().Set(
		"Link",
		`<https://danieljmanningdev.com/>; rel="canonical"`,
	)

	if err := h.homeTemplate.ExecuteTemplate(
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
