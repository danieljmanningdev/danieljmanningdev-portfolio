package http

import (
	"bytes"
	"database/sql"
	"errors"
	"html/template"
	"net/http"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/models"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/rendering"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/repository"
)

type BlogHandler struct {
	repository *repository.BlogRepository

	blogTemplates     *template.Template
	blogPostTemplates *template.Template
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

	return &BlogHandler{
		repository:        repo,
		blogTemplates:     blogTemplates,
		blogPostTemplates: blogPostTemplates,
	}, nil
}

type blogPageData struct {
	Title string
	Posts []models.BlogPost
}

type blogPostPageData struct {
	Title string
	Post  models.BlogPost
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
			Title: "Blog — Daniel J. Manning",
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
			http.NotFound(w, r)
			return
		}

		http.Error(
			w,
			"failed to load blog post",
			http.StatusInternalServerError,
		)
		return
	}

	h.render(
		w,
		h.blogPostTemplates,
		blogPostPageData{
			Title: post.Title + " — Daniel J. Manning",
			Post:  post,
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
