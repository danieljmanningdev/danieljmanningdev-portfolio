package http

import (
	"html/template"
	"net/http"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/rendering"
)

type PublicPageDefinition struct {
	Path            string
	Template        string
	Title           string
	Description     string
	ChangeFrequency string
	Priority        string
}

type PublicPageRoute struct {
	Path    string
	Handler http.Handler
}

var PublicPages = []PublicPageDefinition{
	{
		Path:            "/web-design-leeds/",
		Template:        "public/web-design-leeds.html",
		Title:           "Web Design & Development in Leeds | Daniel J. Manning",
		Description:     "Web design and development in Leeds for businesses that need fast, accessible and maintainable websites and digital products.",
		ChangeFrequency: "monthly",
		Priority:        "0.9",
	},
	{
		Path:            "/web-development/",
		Template:        "public/web-development.html",
		Title:           "Web Development | Daniel J. Manning",
		Description:     "Full-stack web development for fast, accessible and maintainable websites, web applications and digital products.",
		ChangeFrequency: "monthly",
		Priority:        "0.9",
	},
	{
		Path:            "/web-design/",
		Template:        "public/web-design.html",
		Title:           "Web Design | Daniel J. Manning",
		Description:     "Web design for fast, accessible and maintainable websites, web applications and digital products.",
		ChangeFrequency: "monthly",
		Priority:        "0.9",
	},
	{
		Path:            "/software-development/",
		Template:        "public/software-development.html",
		Title:           "Software Development | Daniel J. Manning",
		Description:     "Software development for fast, accessible and maintainable websites, web applications and digital products.",
		ChangeFrequency: "monthly",
		Priority:        "0.9",
	},
	{
		Path:            "/ui-ux-design/",
		Template:        "public/ui-ux-design.html",
		Title:           "UI & UX Design | Daniel J. Manning",
		Description:     "UI and UX design for fast, accessible and maintainable websites, web applications and digital products.",
		ChangeFrequency: "monthly",
		Priority:        "0.9",
	},
	{
		Path:            "/salon-rebuild/",
		Template:        "public/salon-rebuild.html",
		Title:           "Portfolio & Salon Rebuild | Daniel J. Manning",
		Description:     "A portfolio and secure personal workspace designed as one coherent product — clear in public, focused in private and deliberately lightweight throughout.",
		ChangeFrequency: "monthly",
		Priority:        "0.9",
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
	data := newPublicPageData(
		h.definition.Title,
		h.definition.Description,
		h.definition.Path,
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
