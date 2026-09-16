package http

import (
	"html/template"
	"net/http"

	casestudies "github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/case_studies"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/models"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/rendering"
)

type CaseStudyHandler struct {
	template *template.Template
}

type caseStudyPageData struct {
	publicPageData
	CaseStudy models.CaseStudy
}

func NewCaseStudyHandler(templateDir string) (*CaseStudyHandler, error) {
	tmpl, err := rendering.LoadPageTemplate(
		templateDir,
		"public/case-study.html",
	)
	if err != nil {
		return nil, err
	}

	return &CaseStudyHandler{
		template: tmpl,
	}, nil
}

func (h *CaseStudyHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	slug := r.PathValue("slug")

	study, ok := casestudies.Get(slug)
	if !ok {
		http.NotFound(w, r)
		return
	}

	path := "/work/" + study.Slug + "/"

	data := caseStudyPageData{
		publicPageData: newPublicPageData(
			study.Title+" "+study.Accent+" | Daniel J. Manning",
			study.Summary,
			path,
			"article",
			caseStudyStructuredData(
				study.Title+" "+study.Accent,
				path,
				study.Summary,
				study.Tags...,
			),
		).withRelatedLinks(
			relatedLinksForCaseStudy(study.Slug)...,
		).withRequest(r),

		CaseStudy: study,
	}

	w.Header().Set(
		"Link",
		"<"+absolutePublicURL(path)+">; rel=\"canonical\"",
	)

	renderPublicHTML(
		w,
		r,
		h.template,
		data,
	)
}
