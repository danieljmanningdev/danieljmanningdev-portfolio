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

type PortfolioCaseStudyHandler struct {
	template *template.Template
}

func NewPortfolioCaseStudyHandler(
	templateDir string,
) (*PortfolioCaseStudyHandler, error) {
	tmpl, err := rendering.LoadPageTemplate(
		templateDir,
		"public/portfolio.html",
	)
	if err != nil {
		return nil, err
	}

	return &PortfolioCaseStudyHandler{
		template: tmpl,
	}, nil
}

func (h *PortfolioCaseStudyHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.URL.Path != "/work/portfolio" {
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

	const (
		name        = "Portfolio & Client Workspace"
		path        = "/work/portfolio"
		title       = "Go Portfolio & Client Workspace Case Study | Daniel J. Manning"
		description = "Case study of a secure portfolio and client workspace built with Go, HTMX and SQLite, covering product design, authentication, security and full-stack engineering."
	)

	data := newPublicPageData(
		title,
		description,
		path,
		"article",
		caseStudyStructuredData(
			name,
			path,
			description,
			"digital product design",
			"Go",
			"HTMX",
			"SQLite",
			"web application security",
		),
	).withRelatedLinks(
		webDevelopmentRelatedLink(),
		softwareDevelopmentRelatedLink(),
		journalRelatedLink(),
	).withRequest(r)

	w.Header().Set(
		"Content-Type",
		"text/html; charset=utf-8",
	)
	w.Header().Set(
		"Link",
		"<"+absolutePublicURL(path)+">; rel=\"canonical\"",
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
