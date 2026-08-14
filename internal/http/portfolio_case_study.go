package http

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type PortfolioCaseStudyHandler struct {
	template *template.Template
}

func NewPortfolioCaseStudyHandler(
	templateDir string,
) (*PortfolioCaseStudyHandler, error) {
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
			"portfolio-case-study.html",
		),
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

	data := publicPageData{
		Title:       "Portfolio & Client Workspace — Daniel J. Manning",
		Description: "A case study covering the design and engineering of a secure Go portfolio and private client workspace for managing clients, projects and contracts.",
		OGTitle:     "Portfolio & Client Workspace — Daniel J. Manning",
		OGType:      "article",
	}

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
