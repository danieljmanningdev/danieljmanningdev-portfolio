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
		"Daniel J. Manning — Digital Product Designer & Engineer",
		"Digital product design and full-stack engineering focused on fast, purposeful software, thoughtful user experiences and maintainable systems.",
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
