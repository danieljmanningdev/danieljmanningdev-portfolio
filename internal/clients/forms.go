package clients

import (
	"net/http"
	"net/mail"
	"strings"
	"unicode/utf8"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/models"
)

const (
	clientStatusActive   = "active"
	clientStatusInactive = "inactive"

	maxClientNameLength    = 100
	maxClientEmailLength   = 254
	maxClientCompanyLength = 150
	maxClientNotesLength   = 5000
)

type clientForm struct {
	ID      int64
	Name    string
	Email   string
	Company string
	Status  string
	Notes   string
}

type clientFormErrors struct {
	General string
	Name    string
	Email   string
	Company string
	Status  string
	Notes   string
}

func (e clientFormErrors) Any() bool {
	return e.General != "" ||
		e.Name != "" ||
		e.Email != "" ||
		e.Company != "" ||
		e.Status != "" ||
		e.Notes != ""
}

func parseClientForm(r *http.Request) (clientForm, error) {
	if err := r.ParseForm(); err != nil {
		return clientForm{}, err
	}

	form := clientForm{
		Name:    strings.TrimSpace(r.FormValue("name")),
		Email:   strings.TrimSpace(r.FormValue("email")),
		Company: strings.TrimSpace(r.FormValue("company")),
		Status:  strings.TrimSpace(r.FormValue("status")),
		Notes:   strings.TrimSpace(r.FormValue("notes")),
	}

	if form.Status == "" {
		form.Status = clientStatusActive
	}

	return form, nil
}

func validateClientForm(form clientForm) clientFormErrors {
	var errors clientFormErrors

	switch {
	case form.Name == "":
		errors.Name = "Name is required."

	case utf8.RuneCountInString(form.Name) > maxClientNameLength:
		errors.Name = "Name must be 100 characters or fewer."
	}

	switch {
	case form.Email == "":
		errors.Email = "Email is required."

	case utf8.RuneCountInString(form.Email) > maxClientEmailLength:
		errors.Email = "Email must be 254 characters or fewer."

	default:
		parsedEmail, err := mail.ParseAddress(form.Email)

		if err != nil || parsedEmail.Address != form.Email {
			errors.Email = "Please enter a valid email address."
		}
	}

	if utf8.RuneCountInString(form.Company) > maxClientCompanyLength {
		errors.Company = "Company must be 150 characters or fewer."
	}

	if form.Status != clientStatusActive &&
		form.Status != clientStatusInactive {
		errors.Status = "Invalid client status."
	}

	if utf8.RuneCountInString(form.Notes) > maxClientNotesLength {
		errors.Notes = "Notes must be 5,000 characters or fewer."
	}

	return errors
}

func clientFormFromModel(client models.Client) clientForm {
	status := client.Status

	if status == "" {
		status = clientStatusActive
	}

	return clientForm{
		ID:      client.ID,
		Name:    client.Name,
		Email:   client.Email,
		Company: client.Company,
		Status:  status,
		Notes:   client.Notes,
	}
}
