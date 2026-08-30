// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package contracts

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
	contractStatusDraft     = "draft"
	contractStatusSent      = "sent"
	contractStatusAccepted  = "accepted"
	contractStatusCompleted = "completed"
	contractStatusCancelled = "cancelled"

	maxContractTitleLength = 150
	maxContractNotesLength = 5000
)

type contractForm struct {
	ID        int64
	ClientID  int64
	ProjectID int64

	Title     string
	Status    string
	StartDate string
	EndDate   string
	Value     string
	Notes     string
}

type contractFormErrors struct {
	General   string
	ClientID  string
	ProjectID string
	Title     string
	Status    string
	StartDate string
	EndDate   string
	Value     string
	Notes     string
}

func (e contractFormErrors) Any() bool {
	return e.General != "" ||
		e.ClientID != "" ||
		e.ProjectID != "" ||
		e.Title != "" ||
		e.Status != "" ||
		e.StartDate != "" ||
		e.EndDate != "" ||
		e.Value != "" ||
		e.Notes != ""
}

func parseContractForm(
	r *http.Request,
) (contractForm, error) {
	if err := r.ParseForm(); err != nil {
		return contractForm{}, err
	}

	form := contractForm{
		Title: strings.TrimSpace(
			r.FormValue("title"),
		),
		Status: strings.TrimSpace(
			r.FormValue("status"),
		),
		StartDate: strings.TrimSpace(
			r.FormValue("start_date"),
		),
		EndDate: strings.TrimSpace(
			r.FormValue("end_date"),
		),
		Value: strings.TrimSpace(
			r.FormValue("value"),
		),
		Notes: strings.TrimSpace(
			r.FormValue("notes"),
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

	projectIDValue := strings.TrimSpace(
		r.FormValue("project_id"),
	)

	if projectIDValue != "" {
		projectID, err := strconv.ParseInt(
			projectIDValue,
			10,
			64,
		)

		if err == nil && projectID > 0 {
			form.ProjectID = projectID
		}
	}

	if form.Status == "" {
		form.Status = contractStatusDraft
	}

	return form, nil
}

func validateContractForm(
	form contractForm,
) contractFormErrors {
	var formErrors contractFormErrors

	if form.ClientID <= 0 {
		formErrors.ClientID = "Please select a client."
	}

	switch {
	case form.Title == "":
		formErrors.Title = "Contract title is required."

	case utf8.RuneCountInString(
		form.Title,
	) > maxContractTitleLength:
		formErrors.Title =
			"Contract title must be 150 characters or fewer."
	}

	if utf8.RuneCountInString(
		form.Notes,
	) > maxContractNotesLength {
		formErrors.Notes =
			"Notes must be 5,000 characters or fewer."
	}

	if !validContractStatus(form.Status) {
		formErrors.Status = "Invalid contract status."
	}

	startDate, startDateErr := parseOptionalContractDate(
		form.StartDate,
	)

	if startDateErr != nil {
		formErrors.StartDate =
			"Start date must use YYYY-MM-DD."
	}

	endDate, endDateErr := parseOptionalContractDate(
		form.EndDate,
	)

	if endDateErr != nil {
		formErrors.EndDate =
			"End date must use YYYY-MM-DD."
	}

	if startDateErr == nil &&
		endDateErr == nil &&
		startDate != nil &&
		endDate != nil &&
		endDate.Before(*startDate) {
		formErrors.EndDate =
			"End date cannot be earlier than start date."
	}

	if form.Value != "" {
		valueCents, err := parseContractValueCents(
			form.Value,
		)

		if err != nil || valueCents < 0 {
			formErrors.Value =
				"Contract value must be a valid non-negative amount."
		}
	}

	return formErrors
}

func validContractStatus(
	status string,
) bool {
	switch status {
	case contractStatusDraft,
		contractStatusSent,
		contractStatusAccepted,
		contractStatusCompleted,
		contractStatusCancelled:
		return true

	default:
		return false
	}
}

func parseOptionalContractDate(
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

func contractFormDateValues(
	form contractForm,
) (*time.Time, *time.Time, error) {
	startDate, err := parseOptionalContractDate(
		form.StartDate,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"parse contract start date: %w",
			err,
		)
	}

	endDate, err := parseOptionalContractDate(
		form.EndDate,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"parse contract end date: %w",
			err,
		)
	}

	return startDate, endDate, nil
}

func contractProjectIDValue(
	form contractForm,
) *int64 {
	if form.ProjectID <= 0 {
		return nil
	}

	projectID := form.ProjectID

	return &projectID
}

func parseContractValueCents(
	value string,
) (int64, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return 0, nil
	}

	amount, err := strconv.ParseFloat(
		value,
		64,
	)
	if err != nil {
		return 0, err
	}

	if amount < 0 {
		return 0, fmt.Errorf(
			"contract value cannot be negative",
		)
	}

	return int64(amount*100 + 0.5), nil
}

func contractValueCentsPointer(
	form contractForm,
) (*int64, error) {
	if strings.TrimSpace(form.Value) == "" {
		return nil, nil
	}

	valueCents, err := parseContractValueCents(
		form.Value,
	)
	if err != nil {
		return nil, err
	}

	return &valueCents, nil
}

func contractFormFromModel(
	contract models.Contract,
) contractForm {
	status := contract.Status

	if status == "" {
		status = contractStatusDraft
	}

	form := contractForm{
		ID:       contract.ID,
		ClientID: contract.ClientID,
		Title:    contract.Title,
		Status:   status,
		Notes:    contract.Notes,
	}

	if contract.ProjectID != nil {
		form.ProjectID = *contract.ProjectID
	}

	if contract.StartDate != nil {
		form.StartDate = contract.StartDate.Format(
			time.DateOnly,
		)
	}

	if contract.EndDate != nil {
		form.EndDate = contract.EndDate.Format(
			time.DateOnly,
		)
	}

	if contract.ValueCents != nil {
		form.Value = fmt.Sprintf(
			"%.2f",
			float64(*contract.ValueCents)/100,
		)
	}

	return form
}
