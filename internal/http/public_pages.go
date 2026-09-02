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

type PublicPageDefinition struct {
	Path        string
	Template    string
	Name        string
	Title       string
	Description string

	ServiceType string
	AreaServed  []string
	Keywords    []string

	RelatedLinks []relatedLink
}

type PublicPageRoute struct {
	Path    string
	Handler http.Handler
}

var PublicPages = []PublicPageDefinition{
	{
		Path:        "/web-design-leeds/",
		Template:    "public/web-design-leeds.html",
		Name:        "Web Design & Development in Leeds",
		Title:       "Web Design & Development in Leeds | Daniel J. Manning",
		Description: "Web design and development in Leeds for businesses that need fast, accessible and maintainable websites and digital products.",
		ServiceType: "Web design and development",
		AreaServed: []string{
			"Leeds",
			"United Kingdom",
		},
		RelatedLinks: []relatedLink{
			salonCaseStudyRelatedLink(),
			webDesignRelatedLink(),
			uiUXDesignRelatedLink(),
		},
	},
	{
		Path:        "/web-development/",
		Template:    "public/web-development.html",
		Name:        "Web Development",
		Title:       "Web Development | Daniel J. Manning",
		Description: "Full-stack web development for fast, accessible and maintainable websites, web applications and digital products.",
		ServiceType: "Web development",
		AreaServed: []string{
			"United Kingdom",
		},
		RelatedLinks: []relatedLink{
			portfolioCaseStudyRelatedLink(),
			softwareDevelopmentRelatedLink(),
			journalRelatedLink(),
		},
	},
	{
		Path:        "/web-design/",
		Template:    "public/web-design.html",
		Name:        "Web Design",
		Title:       "Web Design | Daniel J. Manning",
		Description: "Web design for fast, accessible and maintainable websites, web applications and digital products.",
		ServiceType: "Web design",
		AreaServed: []string{
			"United Kingdom",
		},
		RelatedLinks: []relatedLink{
			salonCaseStudyRelatedLink(),
			uiUXDesignRelatedLink(),
			webDesignLeedsRelatedLink(),
		},
	},
	{
		Path:        "/software-development/",
		Template:    "public/software-development.html",
		Name:        "Software Development",
		Title:       "Software Development | Daniel J. Manning",
		Description: "Software development for fast, accessible and maintainable websites, web applications and digital products.",
		ServiceType: "Software development",
		AreaServed: []string{
			"United Kingdom",
		},
		RelatedLinks: []relatedLink{
			portfolioCaseStudyRelatedLink(),
			webDevelopmentRelatedLink(),
			journalRelatedLink(),
		},
	},
	{
		Path:        "/ui-ux-design/",
		Template:    "public/ui-ux-design.html",
		Name:        "UI & UX Design",
		Title:       "UI & UX Design | Daniel J. Manning",
		Description: "UI and UX design for fast, accessible and maintainable websites, web applications and digital products.",
		ServiceType: "UI and UX design",
		AreaServed: []string{
			"United Kingdom",
		},
		RelatedLinks: []relatedLink{
			salonCaseStudyRelatedLink(),
			webDesignRelatedLink(),
			portfolioCaseStudyRelatedLink(),
		},
	},
	{
		Path:        "/work/salon-rebuild/",
		Template:    "public/salon-rebuild.html",
		Name:        "Salon Rebuild",
		Title:       "Salon Rebuild | Daniel J. Manning",
		Description: "Revisiting an early freelance salon project with a modern UI/UX, responsive design and server-rendered Go implementation.",
		Keywords: []string{
			"web design",
			"UI/UX design",
			"responsive design",
			"Go",
			"server-rendered HTML",
		},
		RelatedLinks: []relatedLink{
			pricingLessonsRelatedLink(),
			webDesignRelatedLink(),
			uiUXDesignRelatedLink(),
		},
	},
}

type PublicPageHandler struct {
	template   *template.Template
	definition PublicPageDefinition
}

func NewPublicPageHandler(
	templateDir string,
	definition PublicPageDefinition,
) (*PublicPageHandler, error) {
	tmpl, err := rendering.LoadPageTemplate(
		templateDir,
		definition.Template,
	)
	if err != nil {
		return nil, err
	}

	return &PublicPageHandler{
		template:   tmpl,
		definition: definition,
	}, nil
}

func (h *PublicPageHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	structuredData := h.structuredData()

	data := newPublicPageData(
		h.definition.Title,
		h.definition.Description,
		h.definition.Path,
		"website",
		structuredData,
	).withRelatedLinks(
		h.definition.RelatedLinks...,
	).withRequest(r)

	w.Header().Set(
		"Content-Type",
		"text/html; charset=utf-8",
	)
	w.Header().Set(
		"Link",
		"<"+absolutePublicURL(h.definition.Path)+">; rel=\"canonical\"",
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

func (h *PublicPageHandler) structuredData() any {
	if h.definition.ServiceType != "" {
		return servicePageStructuredData(
			h.definition.Name,
			h.definition.Path,
			h.definition.Description,
			h.definition.ServiceType,
			h.definition.AreaServed...,
		)
	}

	return caseStudyStructuredData(
		h.definition.Name,
		h.definition.Path,
		h.definition.Description,
		h.definition.Keywords...,
	)
}

func BuildPublicPageRoutes(
	templateDir string,
) ([]PublicPageRoute, error) {
	routes := make(
		[]PublicPageRoute,
		0,
		len(PublicPages),
	)

	for _, page := range PublicPages {
		handler, err := NewPublicPageHandler(
			templateDir,
			page,
		)
		if err != nil {
			return nil, err
		}

		routes = append(
			routes,
			PublicPageRoute{
				Path:    page.Path,
				Handler: handler,
			},
		)
	}

	return routes, nil
}
