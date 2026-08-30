// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

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

	title := "Web Development & Product Design Journal | Daniel J. Manning"
	description := "Articles on web development, Go, HTMX, UI/UX, digital product design, security and building fast, maintainable software."

	w.Header().Set(
		"Link",
		`<https://danieljmanningdev.com/blog/>; rel="canonical"`,
	)

	h.render(
		w,
		h.blogTemplates,
		blogPageData{
			publicPageData: newPublicPageData(
				title,
				description,
				"/blog/",
				"website",
				map[string]any{
					"@context":    "https://schema.org",
					"@type":       "Blog",
					"name":        title,
					"description": description,
					"url":         publicSiteURL + "/blog/",
					"image":       defaultOGImage,
					"author": map[string]any{
						"@type": "Person",
						"name":  "Daniel J. Manning",
						"url":   publicSiteURL,
					},
				},
			).withRequest(r),
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

	canonicalPath := "/blog/" + post.Slug
	structuredData := map[string]any{
		"@context":         "https://schema.org",
		"@type":            "BlogPosting",
		"headline":         post.Title,
		"description":      description,
		"url":              publicSiteURL + canonicalPath,
		"mainEntityOfPage": publicSiteURL + canonicalPath,
		"image":            defaultOGImage,
		"dateModified":     post.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
		"author": map[string]any{
			"@type": "Person",
			"name":  "Daniel J. Manning",
			"url":   publicSiteURL,
		},
		"publisher": map[string]any{
			"@type": "Person",
			"name":  "Daniel J. Manning",
			"url":   publicSiteURL,
		},
	}

	if post.PublishedAt != nil {
		structuredData["datePublished"] = post.PublishedAt.UTC().Format(
			"2006-01-02T15:04:05Z07:00",
		)
	}

	w.Header().Set(
		"Link",
		"<"+publicSiteURL+canonicalPath+">; rel=\"canonical\"",
	)

	h.render(
		w,
		h.blogPostTemplates,
		blogPostPageData{
			publicPageData: newPublicPageData(
				post.Title+" — Daniel J. Manning",
				description,
				canonicalPath,
				"article",
				structuredData,
			).withRequest(r),
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
