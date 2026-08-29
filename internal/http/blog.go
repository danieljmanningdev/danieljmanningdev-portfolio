package http

import (
	"bytes"
	"database/sql"
	"errors"
	"html/template"
	"net/http"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/blog"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/models"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/rendering"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/repository"
)

type BlogHandler struct {
	repository *repository.BlogRepository

	blogTemplates     *template.Template
	blogPostTemplates *template.Template
	notFoundTemplates *template.Template
}

func NewBlogHandler(
	db *sql.DB,
	templateDir string,
) (*BlogHandler, error) {
	repo := repository.NewBlogRepository(db)

	blogTemplates, err := rendering.LoadPageTemplate(
		templateDir,
		"blog.html",
	)
	if err != nil {
		return nil, err
	}

	blogPostTemplates, err := rendering.LoadPageTemplate(
		templateDir,
		"blog-post.html",
	)
	if err != nil {
		return nil, err
	}

	notFoundTemplates, err := rendering.LoadPageTemplate(
		templateDir,
		"public/404.html",
	)
	if err != nil {
		return nil, err
	}

	return &BlogHandler{
		repository:        repo,
		blogTemplates:     blogTemplates,
		blogPostTemplates: blogPostTemplates,
		notFoundTemplates: notFoundTemplates,
	}, nil
}

type blogPageData struct {
	publicPageData
	Posts []models.BlogPost
}

type blogPostPageData struct {
	publicPageData
	Post        models.BlogPost
	ContentHTML template.HTML
}

func (h *BlogHandler) List(
	w http.ResponseWriter,
	r *http.Request,
) {
	posts, err := h.repository.ListPublished(r.Context())
	if err != nil {
		http.Error(
			w,
			"failed to load blog posts",
			http.StatusInternalServerError,
		)
		return
	}

	h.render(
		w,
		h.blogTemplates,
		blogPageData{
			publicPageData: publicPageData{
				Title:       "Journal — Daniel J. Manning",
				Description: "Articles and technical notes on digital product design, Go, HTMX, security and building maintainable software.",
				OGTitle:     "Journal — Daniel J. Manning",
				OGType:      "website",
			},
			Posts: posts,
		},
	)
}

func (h *BlogHandler) Show(
	w http.ResponseWriter,
	r *http.Request,
) {
	slug := r.PathValue("slug")

	post, err := h.repository.GetBySlug(
		r.Context(),
		slug,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			renderNotFoundPage(
				w,
				h.notFoundTemplates,
				r.URL.Path,
			)
			return
		}

		http.Error(
			w,
			"failed to load blog post",
			http.StatusInternalServerError,
		)
		return
	}

	contentHTML, err := blog.RenderMarkdown(post.Content)
	if err != nil {
		http.Error(
			w,
			"failed to render blog post",
			http.StatusInternalServerError,
		)
		return
	}

	description := post.Excerpt
	if description == "" {
		description = "An article by Daniel J. Manning."
	}

	h.render(
		w,
		h.blogPostTemplates,
		blogPostPageData{
			publicPageData: publicPageData{
				Title:       post.Title + " — Daniel J. Manning",
				Description: description,
				OGTitle:     post.Title + " — Daniel J. Manning",
				OGType:      "article",
			},
			Post:        post,
			ContentHTML: contentHTML,
		},
	)
}

func (h *BlogHandler) render(
	w http.ResponseWriter,
	tmpl *template.Template,
	data any,
) {
	var body bytes.Buffer

	if err := tmpl.ExecuteTemplate(
		&body,
		"base",
		data,
	); err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"text/html; charset=utf-8",
	)

	w.WriteHeader(http.StatusOK)

	_, _ = body.WriteTo(w)
}
