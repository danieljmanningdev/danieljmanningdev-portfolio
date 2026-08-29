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

	title := "Portfolio & Client Workspace — Daniel J. Manning"
	description := "A case study covering the design and engineering of a secure Go portfolio and private workspace for managing clients, projects, contracts and journal publishing."

	data := newPublicPageData(
		title,
		description,
		"/work/portfolio",
		"article",
		map[string]any{
			"@context":    "https://schema.org",
			"@type":       "CreativeWork",
			"name":        title,
			"description": description,
			"url":         publicSiteURL + "/work/portfolio",
			"image":       defaultOGImage,
			"author": map[string]any{
				"@type": "Person",
				"name":  "Daniel J. Manning",
				"url":   publicSiteURL,
			},
			"keywords": []string{
				"digital product design",
				"Go",
				"HTMX",
				"SQLite",
				"web application security",
			},
		},
	).withRequest(r)

	w.Header().Set(
		"Content-Type",
		"text/html; charset=utf-8",
	)

	w.Header().Set(
		"Link",
		`<https://danieljmanningdev.com/work/portfolio>; rel="canonical"`,
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
