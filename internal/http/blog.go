package http

import (
	"database/sql"
	"html/template"

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
