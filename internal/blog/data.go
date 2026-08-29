package blog

import "github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/models"

type ListPageData struct {
	Title string
	Posts []models.BlogPost
}

type FormPageData struct {
	Title  string
	Form   Form
	Errors FormErrors
}

type EditPageData struct {
	Title  string
	PostID int64
	Form   Form
	Errors FormErrors
}
