// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package projects

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/rendering"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/repository"
)

const projectsBasePath = "/dashboard/projects"

type ProjectsHandler struct {
	projectRepository *repository.ProjectRepository
	clientRepository  *repository.ClientRepository

	projectsTemplates    *template.Template
	newProjectTemplates  *template.Template
	projectTemplates     *template.Template
	editProjectTemplates *template.Template
}

func NewProjectsHandler(
	db *sql.DB,
	templateDir string,
) (*ProjectsHandler, error) {
	projectRepo := repository.NewProjectRepository(db)
	clientRepo := repository.NewClientRepository(db)

	projectsTemplates, err := rendering.LoadPageTemplate(
		templateDir,
		"admin/projects/index.html",
	)
	if err != nil {
		return nil, err
	}

	newProjectTemplates, err := rendering.LoadPageTemplate(
		templateDir,
		"admin/projects/new.html",
	)
	if err != nil {
		return nil, err
	}

	projectTemplates, err := rendering.LoadPageTemplate(
		templateDir,
		"admin/projects/show.html",
	)
	if err != nil {
		return nil, err
	}

	editProjectTemplates, err := rendering.LoadPageTemplate(
		templateDir,
		"admin/projects/edit.html",
	)
	if err != nil {
		return nil, err
	}

	return &ProjectsHandler{
		projectRepository:    projectRepo,
		clientRepository:     clientRepo,
		projectsTemplates:    projectsTemplates,
		newProjectTemplates:  newProjectTemplates,
		projectTemplates:     projectTemplates,
		editProjectTemplates: editProjectTemplates,
	}, nil
}

func (h *ProjectsHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	switch r.Method {
	case http.MethodGet:
		h.handleProjectGET(w, r)

	case http.MethodPost:
		h.handleProjectPOST(w, r)

	default:
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
	}
}

func (h *ProjectsHandler) handleProjectGET(
	w http.ResponseWriter,
	r *http.Request,
) {
	path := strings.TrimPrefix(
		r.URL.Path,
		projectsBasePath,
	)

	// GET /dashboard/projects
	if path == "" || path == "/" {
		h.listProjects(w, r)
		return
	}

	// GET /dashboard/projects/new
	if path == "/new" {
		h.renderNewProject(
			w,
			r,
			projectForm{
				Status: projectStatusPlanned,
			},
			projectFormErrors{},
			http.StatusOK,
		)
		return
	}

	// GET /dashboard/projects/:id/edit
	if strings.HasSuffix(path, "/edit") {
		idString := strings.TrimSuffix(
			strings.TrimPrefix(path, "/"),
			"/edit",
		)

		id, err := strconv.ParseInt(
			idString,
			10,
			64,
		)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		h.editProject(w, r, id)
		return
	}

	// GET /dashboard/projects/:id
	idString := strings.TrimPrefix(path, "/")

	if idString == "" || strings.Contains(idString, "/") {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.ParseInt(
		idString,
		10,
		64,
	)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	h.showProject(w, r, id)
}

func (h *ProjectsHandler) handleProjectPOST(
	w http.ResponseWriter,
	r *http.Request,
) {
	path := strings.TrimPrefix(
		r.URL.Path,
		projectsBasePath,
	)

	// POST /dashboard/projects/new
	if path == "/new" {
		h.createProject(w, r)
		return
	}

	// POST /dashboard/projects/:id/archive
	if strings.HasSuffix(path, "/archive") {
		idString := strings.TrimSuffix(
			strings.TrimPrefix(path, "/"),
			"/archive",
		)

		id, err := strconv.ParseInt(
			idString,
			10,
			64,
		)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		h.archiveProject(w, r, id)
		return
	}

	// POST /dashboard/projects/:id
	idString := strings.TrimPrefix(path, "/")

	if idString == "" || strings.Contains(idString, "/") {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.ParseInt(
		idString,
		10,
		64,
	)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	h.updateProject(w, r, id)
}
