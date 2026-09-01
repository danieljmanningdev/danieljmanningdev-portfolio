// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package contracts

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/models"
)

func TestParseContractFormTrimsValuesAndDefaultsStatus(
	t *testing.T,
) {
	values := url.Values{}

	values.Set("client_id", " 42 ")
	values.Set("project_id", " 7 ")
	values.Set("title", "  Website Development Agreement  ")
	values.Set("start_date", " 2026-08-10 ")
	values.Set("end_date", " 2026-09-30 ")
	values.Set("value", " 2500.50 ")
	values.Set("notes", "  Initial agreement notes.  ")

	req := httptest.NewRequest(
		http.MethodPost,
		"/dashboard/contracts/new",
		strings.NewReader(values.Encode()),
	)

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	form, err := parseContractForm(req)
	if err != nil {
		t.Fatalf("parse contract form: %v", err)
	}

	if form.ClientID != 42 {
		t.Fatalf(
			"expected client ID 42, got %d",
			form.ClientID,
		)
	}

	if form.ProjectID != 7 {
		t.Fatalf(
			"expected project ID 7, got %d",
			form.ProjectID,
		)
	}

	if form.Title != "Website Development Agreement" {
		t.Fatalf(
			"unexpected title %q",
			form.Title,
		)
	}

	if form.Status != contractStatusDraft {
		t.Fatalf(
			"expected default status %q, got %q",
			contractStatusDraft,
			form.Status,
		)
	}

	if form.StartDate != "2026-08-10" {
		t.Fatalf(
			"unexpected start date %q",
			form.StartDate,
		)
	}

	if form.EndDate != "2026-09-30" {
		t.Fatalf(
			"unexpected end date %q",
			form.EndDate,
		)
	}

	if form.Value != "2500.50" {
		t.Fatalf(
			"unexpected value %q",
			form.Value,
		)
	}

	if form.Notes != "Initial agreement notes." {
		t.Fatalf(
			"unexpected notes %q",
			form.Notes,
		)
	}
}

func TestParseContractFormInvalidIDsBecomeUnselected(
	t *testing.T,
) {
	values := url.Values{}

	values.Set("client_id", "not-a-number")
	values.Set("project_id", "-1")
	values.Set("title", "Agreement")

	req := httptest.NewRequest(
		http.MethodPost,
		"/dashboard/contracts/new",
		strings.NewReader(values.Encode()),
	)

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	form, err := parseContractForm(req)
	if err != nil {
		t.Fatalf("parse contract form: %v", err)
	}

	if form.ClientID != 0 {
		t.Fatalf(
			"expected invalid client ID to become 0, got %d",
			form.ClientID,
		)
	}

	if form.ProjectID != 0 {
		t.Fatalf(
			"expected invalid project ID to become 0, got %d",
			form.ProjectID,
		)
	}
}

func TestValidateContractFormAcceptsAllowedStatuses(
	t *testing.T,
) {
	statuses := []string{
		contractStatusDraft,
		contractStatusSent,
		contractStatusAccepted,
		contractStatusCompleted,
		contractStatusCancelled,
	}

	for _, status := range statuses {
		t.Run(status, func(t *testing.T) {
			formErrors := validateContractForm(
				contractForm{
					ClientID: 1,
					Title:    "Test Contract",
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

func TestValidateContractFormRequiresClientAndTitle(
	t *testing.T,
) {
	formErrors := validateContractForm(
		contractForm{
			Status: contractStatusDraft,
		},
	)

	if formErrors.ClientID != "Please select a client." {
		t.Fatalf(
			"unexpected client error %q",
			formErrors.ClientID,
		)
	}

	if formErrors.Title != "Contract title is required." {
		t.Fatalf(
			"unexpected title error %q",
			formErrors.Title,
		)
	}
}

func TestValidateContractFormRejectsInvalidStatus(
	t *testing.T,
) {
	formErrors := validateContractForm(
		contractForm{
			ClientID: 1,
			Title:    "Test Contract",
			Status:   "banana",
		},
	)

	if formErrors.Status != "Invalid contract status." {
		t.Fatalf(
			"unexpected status error %q",
			formErrors.Status,
		)
	}
}

func TestValidateContractFormRejectsInvalidDates(
	t *testing.T,
) {
	formErrors := validateContractForm(
		contractForm{
			ClientID:  1,
			Title:     "Test Contract",
			Status:    contractStatusDraft,
			StartDate: "10/08/2026",
			EndDate:   "30/09/2026",
		},
	)

	if formErrors.StartDate !=
		"Start date must use YYYY-MM-DD." {
		t.Fatalf(
			"unexpected start-date error %q",
			formErrors.StartDate,
		)
	}

	if formErrors.EndDate !=
		"End date must use YYYY-MM-DD." {
		t.Fatalf(
			"unexpected end-date error %q",
			formErrors.EndDate,
		)
	}
}

func TestValidateContractFormRejectsEndDateBeforeStartDate(
	t *testing.T,
) {
	formErrors := validateContractForm(
		contractForm{
			ClientID:  1,
			Title:     "Test Contract",
			Status:    contractStatusAccepted,
			StartDate: "2026-09-30",
			EndDate:   "2026-08-10",
		},
	)

	if formErrors.EndDate !=
		"End date cannot be earlier than start date." {
		t.Fatalf(
			"unexpected end-date error %q",
			formErrors.EndDate,
		)
	}
}

func TestValidateContractFormRejectsMaximumLengths(
	t *testing.T,
) {
	formErrors := validateContractForm(
		contractForm{
			ClientID: 1,
			Title: strings.Repeat(
				"a",
				maxContractTitleLength+1,
			),
			Notes: strings.Repeat(
				"a",
				maxContractNotesLength+1,
			),
			Status: contractStatusDraft,
		},
	)

	if formErrors.Title == "" {
		t.Fatal("expected contract-title length error")
	}

	if formErrors.Notes == "" {
		t.Fatal("expected notes length error")
	}
}

func TestValidateContractFormAcceptsValidValue(
	t *testing.T,
) {
	formErrors := validateContractForm(
		contractForm{
			ClientID: 1,
			Title:    "Test Contract",
			Status:   contractStatusDraft,
			Value:    "2500.50",
		},
	)

	if formErrors.Value != "" {
		t.Fatalf(
			"unexpected value error %q",
			formErrors.Value,
		)
	}
}

func TestValidateContractFormRejectsInvalidValue(
	t *testing.T,
) {
	values := []string{
		"not-money",
		"-10.00",
	}

	for _, value := range values {
		t.Run(value, func(t *testing.T) {
			formErrors := validateContractForm(
				contractForm{
					ClientID: 1,
					Title:    "Test Contract",
					Status:   contractStatusDraft,
					Value:    value,
				},
			)

			if formErrors.Value == "" {
				t.Fatal("expected contract-value validation error")
			}
		})
	}
}

func TestParseContractValueCents(
	t *testing.T,
) {
	tests := []struct {
		name      string
		value     string
		wantCents int64
	}{
		{
			name:      "whole amount",
			value:     "2500",
			wantCents: 250000,
		},
		{
			name:      "decimal amount",
			value:     "2500.50",
			wantCents: 250050,
		},
		{
			name:      "single decimal",
			value:     "99.9",
			wantCents: 9990,
		},
		{
			name:      "zero",
			value:     "0",
			wantCents: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseContractValueCents(
				tt.value,
			)
			if err != nil {
				t.Fatalf(
					"parse contract value: %v",
					err,
				)
			}

			if got != tt.wantCents {
				t.Fatalf(
					"expected %d cents, got %d",
					tt.wantCents,
					got,
				)
			}
		})
	}
}

func TestContractFormDateValues(
	t *testing.T,
) {
	startDate, endDate, err := contractFormDateValues(
		contractForm{
			StartDate: "2026-08-10",
			EndDate:   "2026-09-30",
		},
	)
	if err != nil {
		t.Fatalf(
			"get contract date values: %v",
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

	expectedEndDate := time.Date(
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

	if endDate == nil ||
		!endDate.Equal(expectedEndDate) {
		t.Fatalf(
			"expected end date %v, got %v",
			expectedEndDate,
			endDate,
		)
	}
}

func TestContractFormDateValuesAllowsEmptyDates(
	t *testing.T,
) {
	startDate, endDate, err := contractFormDateValues(
		contractForm{},
	)
	if err != nil {
		t.Fatalf(
			"get empty contract dates: %v",
			err,
		)
	}

	if startDate != nil {
		t.Fatalf(
			"expected nil start date, got %v",
			startDate,
		)
	}

	if endDate != nil {
		t.Fatalf(
			"expected nil end date, got %v",
			endDate,
		)
	}
}

func TestContractProjectIDValueAllowsEmptyProject(
	t *testing.T,
) {
	value := contractProjectIDValue(
		contractForm{},
	)

	if value != nil {
		t.Fatalf(
			"expected nil project ID, got %v",
			value,
		)
	}
}

func TestContractProjectIDValueReturnsSelectedProject(
	t *testing.T,
) {
	value := contractProjectIDValue(
		contractForm{
			ProjectID: 42,
		},
	)

	if value == nil {
		t.Fatal("expected project ID")
	}

	if *value != 42 {
		t.Fatalf(
			"expected project ID 42, got %d",
			*value,
		)
	}
}

func TestContractValueCentsPointerAllowsEmptyValue(
	t *testing.T,
) {
	value, err := contractValueCentsPointer(
		contractForm{},
	)
	if err != nil {
		t.Fatalf(
			"get empty contract value: %v",
			err,
		)
	}

	if value != nil {
		t.Fatalf(
			"expected nil value, got %v",
			value,
		)
	}
}

func TestContractValueCentsPointerReturnsValue(
	t *testing.T,
) {
	value, err := contractValueCentsPointer(
		contractForm{
			Value: "1234.56",
		},
	)
	if err != nil {
		t.Fatalf(
			"get contract value: %v",
			err,
		)
	}

	if value == nil {
		t.Fatal("expected contract value")
	}

	if *value != 123456 {
		t.Fatalf(
			"expected 123456 cents, got %d",
			*value,
		)
	}
}

func TestContractFormFromModel(
	t *testing.T,
) {
	projectID := int64(7)
	valueCents := int64(250050)

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

	endDate := time.Date(
		2026,
		time.September,
		30,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	form := contractFormFromModel(
		models.Contract{
			ID:         9,
			ClientID:   42,
			ProjectID:  &projectID,
			Title:      "Website Development Agreement",
			Status:     contractStatusAccepted,
			StartDate:  &startDate,
			EndDate:    &endDate,
			ValueCents: &valueCents,
			Notes:      "Agreement notes.",
		},
	)

	if form.ID != 9 {
		t.Fatalf(
			"expected contract ID 9, got %d",
			form.ID,
		)
	}

	if form.ClientID != 42 {
		t.Fatalf(
			"expected client ID 42, got %d",
			form.ClientID,
		)
	}

	if form.ProjectID != 7 {
		t.Fatalf(
			"expected project ID 7, got %d",
			form.ProjectID,
		)
	}

	if form.Title != "Website Development Agreement" {
		t.Fatalf(
			"unexpected title %q",
			form.Title,
		)
	}

	if form.Status != contractStatusAccepted {
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

	if form.EndDate != "2026-09-30" {
		t.Fatalf(
			"unexpected end date %q",
			form.EndDate,
		)
	}

	if form.Value != "2500.50" {
		t.Fatalf(
			"unexpected value %q",
			form.Value,
		)
	}

	if form.Notes != "Agreement notes." {
		t.Fatalf(
			"unexpected notes %q",
			form.Notes,
		)
	}
}

func TestContractFormFromModelAllowsOptionalValues(
	t *testing.T,
) {
	form := contractFormFromModel(
		models.Contract{
			ID:       9,
			ClientID: 42,
			Title:    "General Agreement",
		},
	)

	if form.Status != contractStatusDraft {
		t.Fatalf(
			"expected default status %q, got %q",
			contractStatusDraft,
			form.Status,
		)
	}

	if form.ProjectID != 0 {
		t.Fatalf(
			"expected no project, got %d",
			form.ProjectID,
		)
	}

	if form.StartDate != "" {
		t.Fatalf(
			"expected empty start date, got %q",
			form.StartDate,
		)
	}

	if form.EndDate != "" {
		t.Fatalf(
			"expected empty end date, got %q",
			form.EndDate,
		)
	}

	if form.Value != "" {
		t.Fatalf(
			"expected empty value, got %q",
			form.Value,
		)
	}
}
