package http

import (
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/models"
)

type projectsPageData struct {
	Title    string
	Projects []models.Project
}

type projectPageData struct {
	Title   string
	Project models.Project
}

type projectFormPageData struct {
	Title   string
	Form    projectForm
	Errors  projectFormErrors
	Clients []models.Client
}
