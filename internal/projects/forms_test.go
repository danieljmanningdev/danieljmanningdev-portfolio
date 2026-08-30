// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package projects

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/models"
)

func TestParseProjectFormTrimsValuesAndDefaultsStatus(
	t *testing.T,
) {
	values := url.Values{}

	values.Set("client_id", " 42 ")
	values.Set("name", "  Website Platform  ")
	values.Set(
		"description",
		"  Design and build the platform.  ",
	)
	values.Set("start_date", " 2026-08-10 ")
	values.Set("due_date", " 2026-09-30 ")

	req := httptest.NewRequest(
		http.MethodPost,
		"/dashboard/projects/new",
		strings.NewReader(values.Encode()),
	)

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	form, err := parseProjectForm(req)
	if err != nil {
		t.Fatalf("parse project form: %v", err)
	}

	if form.ClientID != 42 {
		t.Fatalf(
			"expected client ID 42, got %d",
			form.ClientID,
		)
	}

	if form.Name != "Website Platform" {
		t.Fatalf(
			"expected name %q, got %q",
			"Website Platform",
			form.Name,
		)
	}

	if form.Description !=
		"Design and build the platform." {
		t.Fatalf(
			"unexpected description %q",
			form.Description,
		)
	}

	if form.Status != projectStatusPlanned {
		t.Fatalf(
			"expected default status %q, got %q",
			projectStatusPlanned,
			form.Status,
		)
	}

	if form.StartDate != "2026-08-10" {
		t.Fatalf(
			"unexpected start date %q",
			form.StartDate,
		)
	}

	if form.DueDate != "2026-09-30" {
		t.Fatalf(
			"unexpected due date %q",
			form.DueDate,
		)
	}
}

func TestParseProjectFormInvalidClientIDBecomesUnselected(
	t *testing.T,
) {
	values := url.Values{}

	values.Set("client_id", "not-a-number")
	values.Set("name", "Website Platform")

	req := httptest.NewRequest(
		http.MethodPost,
		"/dashboard/projects/new",
		strings.NewReader(values.Encode()),
	)

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	form, err := parseProjectForm(req)
	if err != nil {
		t.Fatalf("parse project form: %v", err)
	}

	if form.ClientID != 0 {
		t.Fatalf(
			"expected invalid client ID to become 0, got %d",
			form.ClientID,
		)
	}

	formErrors := validateProjectForm(form)

	if formErrors.ClientID != "Please select a client." {
		t.Fatalf(
			"unexpected client error %q",
			formErrors.ClientID,
		)
	}
}

func TestValidateProjectFormAcceptsAllowedStatuses(
	t *testing.T,
) {
	statuses := []string{
		projectStatusPlanned,
		projectStatusActive,
		projectStatusCompleted,
		projectStatusArchived,
	}

	for _, status := range statuses {
		t.Run(status, func(t *testing.T) {
			formErrors := validateProjectForm(
				projectForm{
					ClientID: 1,
					Name:     "Test Project",
					Status:   status,
				},
			)

			if formErrors.Any() {
				t.Fatalf(
					"expected no validation errors, got %+v",
					formErrors,
				)
			}
		})
	}
}

func TestValidateProjectFormRequiresClientAndName(
	t *testing.T,
) {
	formErrors := validateProjectForm(
		projectForm{
			Status: projectStatusPlanned,
		},
	)

	if formErrors.ClientID != "Please select a client." {
		t.Fatalf(
			"unexpected client error %q",
			formErrors.ClientID,
		)
	}

	if formErrors.Name != "Project name is required." {
		t.Fatalf(
			"unexpected name error %q",
			formErrors.Name,
		)
	}
}

func TestValidateProjectFormRejectsInvalidStatus(
	t *testing.T,
) {
	formErrors := validateProjectForm(
		projectForm{
			ClientID: 1,
			Name:     "Test Project",
			Status:   "banana",
		},
	)

	if formErrors.Status != "Invalid project status." {
		t.Fatalf(
			"unexpected status error %q",
			formErrors.Status,
		)
	}
}

func TestValidateProjectFormRejectsInvalidDates(
	t *testing.T,
) {
	formErrors := validateProjectForm(
		projectForm{
			ClientID:  1,
			Name:      "Test Project",
			Status:    projectStatusPlanned,
			StartDate: "10/08/2026",
			DueDate:   "30/09/2026",
		},
	)

	if formErrors.StartDate !=
		"Start date must use YYYY-MM-DD." {
		t.Fatalf(
			"unexpected start-date error %q",
			formErrors.StartDate,
		)
	}

	if formErrors.DueDate !=
		"Due date must use YYYY-MM-DD." {
		t.Fatalf(
			"unexpected due-date error %q",
			formErrors.DueDate,
		)
	}
}

func TestValidateProjectFormRejectsDueDateBeforeStartDate(
	t *testing.T,
) {
	formErrors := validateProjectForm(
		projectForm{
			ClientID:  1,
			Name:      "Test Project",
			Status:    projectStatusActive,
			StartDate: "2026-09-30",
			DueDate:   "2026-08-10",
		},
	)

	if formErrors.DueDate !=
		"Due date cannot be earlier than start date." {
		t.Fatalf(
			"unexpected due-date error %q",
			formErrors.DueDate,
		)
	}
}

func TestValidateProjectFormRejectsMaximumLengths(
	t *testing.T,
) {
	formErrors := validateProjectForm(
		projectForm{
			ClientID: 1,
			Name: strings.Repeat(
				"a",
				maxProjectNameLength+1,
			),
			Description: strings.Repeat(
				"a",
				maxProjectDescriptionLength+1,
			),
			Status: projectStatusPlanned,
		},
	)

	if formErrors.Name == "" {
		t.Fatal("expected project-name length error")
	}

	if formErrors.Description == "" {
		t.Fatal("expected description length error")
	}
}

func TestProjectFormDateValues(t *testing.T) {
	startDate, dueDate, err := projectFormDateValues(
		projectForm{
			StartDate: "2026-08-10",
			DueDate:   "2026-09-30",
		},
	)
	if err != nil {
		t.Fatalf(
			"get project date values: %v",
			err,
		)
	}

	expectedStartDate := time.Date(
		2026,
		time.August,
		10,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	expectedDueDate := time.Date(
		2026,
		time.September,
		30,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	if startDate == nil ||
		!startDate.Equal(expectedStartDate) {
		t.Fatalf(
			"expected start date %v, got %v",
			expectedStartDate,
			startDate,
		)
	}

	if dueDate == nil ||
		!dueDate.Equal(expectedDueDate) {
		t.Fatalf(
			"expected due date %v, got %v",
			expectedDueDate,
			dueDate,
		)
	}
}

func TestProjectFormDateValuesAllowsEmptyDates(
	t *testing.T,
) {
	startDate, dueDate, err := projectFormDateValues(
		projectForm{},
	)
	if err != nil {
		t.Fatalf(
			"get empty project dates: %v",
			err,
		)
	}

	if startDate != nil {
		t.Fatalf(
			"expected nil start date, got %v",
			startDate,
		)
	}

	if dueDate != nil {
		t.Fatalf(
			"expected nil due date, got %v",
			dueDate,
		)
	}
}

func TestProjectFormFromModel(t *testing.T) {
	startDate := time.Date(
		2026,
		time.August,
		10,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	dueDate := time.Date(
		2026,
		time.September,
		30,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	form := projectFormFromModel(
		models.Project{
			ID:          7,
			ClientID:    42,
			Name:        "Website Platform",
			Description: "Build the platform.",
			Status:      projectStatusActive,
			StartDate:   &startDate,
			DueDate:     &dueDate,
		},
	)

	if form.ID != 7 {
		t.Fatalf(
			"expected project ID 7, got %d",
			form.ID,
		)
	}

	if form.ClientID != 42 {
		t.Fatalf(
			"expected client ID 42, got %d",
			form.ClientID,
		)
	}

	if form.Name != "Website Platform" {
		t.Fatalf(
			"unexpected name %q",
			form.Name,
		)
	}

	if form.Description != "Build the platform." {
		t.Fatalf(
			"unexpected description %q",
			form.Description,
		)
	}

	if form.Status != projectStatusActive {
		t.Fatalf(
			"unexpected status %q",
			form.Status,
		)
	}

	if form.StartDate != "2026-08-10" {
		t.Fatalf(
			"unexpected start date %q",
			form.StartDate,
		)
	}

	if form.DueDate != "2026-09-30" {
		t.Fatalf(
			"unexpected due date %q",
			form.DueDate,
		)
	}
}
