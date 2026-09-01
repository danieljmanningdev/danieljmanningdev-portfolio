// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package projects

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/models"
)

const (
	projectStatusPlanned   = "planned"
	projectStatusActive    = "active"
	projectStatusCompleted = "completed"
	projectStatusArchived  = "archived"

	maxProjectNameLength        = 150
	maxProjectDescriptionLength = 5000
)

type projectForm struct {
	ID          int64
	ClientID    int64
	Name        string
	Description string
	Status      string
	StartDate   string
	DueDate     string
}

type projectFormErrors struct {
	General     string
	ClientID    string
	Name        string
	Description string
	Status      string
	StartDate   string
	DueDate     string
}

func (e projectFormErrors) Any() bool {
	return e.General != "" ||
		e.ClientID != "" ||
		e.Name != "" ||
		e.Description != "" ||
		e.Status != "" ||
		e.StartDate != "" ||
		e.DueDate != ""
}

func parseProjectForm(
	r *http.Request,
) (projectForm, error) {
	if err := r.ParseForm(); err != nil {
		return projectForm{}, err
	}

	form := projectForm{
		Name: strings.TrimSpace(
			r.FormValue("name"),
		),
		Description: strings.TrimSpace(
			r.FormValue("description"),
		),
		Status: strings.TrimSpace(
			r.FormValue("status"),
		),
		StartDate: strings.TrimSpace(
			r.FormValue("start_date"),
		),
		DueDate: strings.TrimSpace(
			r.FormValue("due_date"),
		),
	}

	clientIDValue := strings.TrimSpace(
		r.FormValue("client_id"),
	)

	if clientIDValue != "" {
		clientID, err := strconv.ParseInt(
			clientIDValue,
			10,
			64,
		)

		if err == nil && clientID > 0 {
			form.ClientID = clientID
		}
	}

	if form.Status == "" {
		form.Status = projectStatusPlanned
	}

	return form, nil
}

func validateProjectForm(
	form projectForm,
) projectFormErrors {
	var formErrors projectFormErrors

	if form.ClientID <= 0 {
		formErrors.ClientID = "Please select a client."
	}

	switch {
	case form.Name == "":
		formErrors.Name = "Project name is required."

	case utf8.RuneCountInString(
		form.Name,
	) > maxProjectNameLength:
		formErrors.Name =
			"Project name must be 150 characters or fewer."
	}

	if utf8.RuneCountInString(
		form.Description,
	) > maxProjectDescriptionLength {
		formErrors.Description =
			"Description must be 5,000 characters or fewer."
	}

	if !validProjectStatus(form.Status) {
		formErrors.Status = "Invalid project status."
	}

	startDate, startDateErr := parseOptionalProjectDate(
		form.StartDate,
	)

	if startDateErr != nil {
		formErrors.StartDate =
			"Start date must use YYYY-MM-DD."
	}

	dueDate, dueDateErr := parseOptionalProjectDate(
		form.DueDate,
	)

	if dueDateErr != nil {
		formErrors.DueDate =
			"Due date must use YYYY-MM-DD."
	}

	if startDateErr == nil &&
		dueDateErr == nil &&
		startDate != nil &&
		dueDate != nil &&
		dueDate.Before(*startDate) {
		formErrors.DueDate =
			"Due date cannot be earlier than start date."
	}

	return formErrors
}

func validProjectStatus(
	status string,
) bool {
	switch status {
	case projectStatusPlanned,
		projectStatusActive,
		projectStatusCompleted,
		projectStatusArchived:
		return true

	default:
		return false
	}
}

func parseOptionalProjectDate(
	value string,
) (*time.Time, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return nil, nil
	}

	parsed, err := time.Parse(
		time.DateOnly,
		value,
	)
	if err != nil {
		return nil, err
	}

	date := time.Date(
		parsed.Year(),
		parsed.Month(),
		parsed.Day(),
		0,
		0,
		0,
		0,
		time.UTC,
	)

	return &date, nil
}

func projectFormDateValues(
	form projectForm,
) (*time.Time, *time.Time, error) {
	startDate, err := parseOptionalProjectDate(
		form.StartDate,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"parse project start date: %w",
			err,
		)
	}

	dueDate, err := parseOptionalProjectDate(
		form.DueDate,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"parse project due date: %w",
			err,
		)
	}

	return startDate, dueDate, nil
}

func projectFormFromModel(
	project models.Project,
) projectForm {
	status := project.Status

	if status == "" {
		status = projectStatusPlanned
	}

	form := projectForm{
		ID:          project.ID,
		ClientID:    project.ClientID,
		Name:        project.Name,
		Description: project.Description,
		Status:      status,
	}

	if project.StartDate != nil {
		form.StartDate = project.StartDate.Format(
			time.DateOnly,
		)
	}

	if project.DueDate != nil {
		form.DueDate = project.DueDate.Format(
			time.DateOnly,
		)
	}

	return form
}
