// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package projects

import (
	"bytes"
	"database/sql"
	"errors"
	"html/template"
	"net/http"
)

func (h *ProjectsHandler) listProjects(
	w http.ResponseWriter,
	r *http.Request,
) {
	projects, err := h.projectRepository.List(
		r.Context(),
	)
	if err != nil {
		http.Error(
			w,
			"failed to load projects",
			http.StatusInternalServerError,
		)
		return
	}

	h.renderProjectPage(
		w,
		h.projectsTemplates,
		projectsPageData{
			Title:    "Projects — Daniel J. Manning",
			Projects: projects,
		},
	)
}

func (h *ProjectsHandler) showProject(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {
	project, err := h.projectRepository.GetByID(
		r.Context(),
		id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}

		http.Error(
			w,
			"failed to load project",
			http.StatusInternalServerError,
		)
		return
	}

	h.renderProjectPage(
		w,
		h.projectTemplates,
		projectPageData{
			Title:   "Project — Daniel J. Manning",
			Project: project,
		},
	)
}

func (h *ProjectsHandler) editProject(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {
	project, err := h.projectRepository.GetByID(
		r.Context(),
		id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}

		http.Error(
			w,
			"failed to load project",
			http.StatusInternalServerError,
		)
		return
	}

	h.renderEditProject(
		w,
		r,
		projectFormFromModel(project),
		projectFormErrors{},
		http.StatusOK,
	)
}

func (h *ProjectsHandler) renderNewProject(
	w http.ResponseWriter,
	r *http.Request,
	form projectForm,
	formErrors projectFormErrors,
	status int,
) {
	if form.Status == "" {
		form.Status = projectStatusPlanned
	}

	clients, err := h.clientRepository.List(
		r.Context(),
	)
	if err != nil {
		http.Error(
			w,
			"failed to load clients",
			http.StatusInternalServerError,
		)
		return
	}

	h.renderProjectStatus(
		w,
		h.newProjectTemplates,
		projectFormPageData{
			Title:   "Add Project — Daniel J. Manning",
			Form:    form,
			Errors:  formErrors,
			Clients: clients,
		},
		status,
	)
}

func (h *ProjectsHandler) renderEditProject(
	w http.ResponseWriter,
	r *http.Request,
	form projectForm,
	formErrors projectFormErrors,
	status int,
) {
	if form.Status == "" {
		form.Status = projectStatusPlanned
	}

	clients, err := h.clientRepository.List(
		r.Context(),
	)
	if err != nil {
		http.Error(
			w,
			"failed to load clients",
			http.StatusInternalServerError,
		)
		return
	}

	h.renderProjectStatus(
		w,
		h.editProjectTemplates,
		projectFormPageData{
			Title:   "Edit Project — Daniel J. Manning",
			Form:    form,
			Errors:  formErrors,
			Clients: clients,
		},
		status,
	)
}

func (h *ProjectsHandler) renderProjectPage(
	w http.ResponseWriter,
	tmpl *template.Template,
	data any,
) {
	h.renderProjectStatus(
		w,
		tmpl,
		data,
		http.StatusOK,
	)
}

func (h *ProjectsHandler) renderProjectStatus(
	w http.ResponseWriter,
	tmpl *template.Template,
	data any,
	status int,
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

	w.WriteHeader(status)

	_, _ = body.WriteTo(w)
}
