package projects

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
)

func (h *ProjectsHandler) createProject(
	w http.ResponseWriter,
	r *http.Request,
) {
	form, err := parseProjectForm(r)
	if err != nil {
		h.renderNewProject(
			w,
			r,
			projectForm{
				Status: projectStatusPlanned,
			},
			projectFormErrors{
				General: "Invalid form submission.",
			},
			http.StatusBadRequest,
		)
		return
	}

	formErrors := validateProjectForm(form)

	if formErrors.Any() {
		h.renderNewProject(
			w,
			r,
			form,
			formErrors,
			http.StatusBadRequest,
		)
		return
	}

	startDate, dueDate, err := projectFormDateValues(form)
	if err != nil {
		h.renderNewProject(
			w,
			r,
			form,
			projectFormErrors{
				General: "Invalid project dates.",
			},
			http.StatusBadRequest,
		)
		return
	}

	id, err := h.projectRepository.Create(
		r.Context(),
		form.ClientID,
		form.Name,
		form.Description,
		form.Status,
		startDate,
		dueDate,
	)
	if err != nil {
		http.Error(
			w,
			"failed to create project",
			http.StatusInternalServerError,
		)
		return
	}

	http.Redirect(
		w,
		r,
		projectsBasePath+"/"+strconv.FormatInt(id, 10),
		http.StatusSeeOther,
	)
}

func (h *ProjectsHandler) updateProject(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {
	form, err := parseProjectForm(r)
	form.ID = id

	if err != nil {
		h.renderEditProject(
			w,
			r,
			form,
			projectFormErrors{
				General: "Invalid form submission.",
			},
			http.StatusBadRequest,
		)
		return
	}

	formErrors := validateProjectForm(form)

	if formErrors.Any() {
		h.renderEditProject(
			w,
			r,
			form,
			formErrors,
			http.StatusBadRequest,
		)
		return
	}

	startDate, dueDate, err := projectFormDateValues(form)
	if err != nil {
		h.renderEditProject(
			w,
			r,
			form,
			projectFormErrors{
				General: "Invalid project dates.",
			},
			http.StatusBadRequest,
		)
		return
	}

	err = h.projectRepository.Update(
		r.Context(),
		id,
		form.ClientID,
		form.Name,
		form.Description,
		form.Status,
		startDate,
		dueDate,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}

		http.Error(
			w,
			"failed to update project",
			http.StatusInternalServerError,
		)
		return
	}

	http.Redirect(
		w,
		r,
		projectsBasePath+"/"+strconv.FormatInt(id, 10),
		http.StatusSeeOther,
	)
}

func (h *ProjectsHandler) archiveProject(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {
	err := h.projectRepository.Archive(
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
			"failed to archive project",
			http.StatusInternalServerError,
		)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set(
			"HX-Redirect",
			projectsBasePath,
		)

		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(
		w,
		r,
		projectsBasePath,
		http.StatusSeeOther,
	)
}
