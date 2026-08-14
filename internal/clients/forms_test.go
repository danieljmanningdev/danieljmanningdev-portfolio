package clients

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestParseClientFormTrimsValuesAndDefaultsStatus(
	t *testing.T,
) {
	values := url.Values{}
	values.Set("name", "  Test Client  ")
	values.Set("email", "  test@example.com  ")
	values.Set("company", "  Test Company  ")
	values.Set("notes", "  Test notes  ")

	req := httptest.NewRequest(
		http.MethodPost,
		"/dashboard/clients/new",
		strings.NewReader(values.Encode()),
	)
	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	form, err := parseClientForm(req)
	if err != nil {
		t.Fatalf("parse client form: %v", err)
	}

	if form.Name != "Test Client" {
		t.Fatalf(
			"expected trimmed name %q, got %q",
			"Test Client",
			form.Name,
		)
	}

	if form.Email != "test@example.com" {
		t.Fatalf(
			"expected trimmed email %q, got %q",
			"test@example.com",
			form.Email,
		)
	}

	if form.Company != "Test Company" {
		t.Fatalf(
			"expected trimmed company %q, got %q",
			"Test Company",
			form.Company,
		)
	}

	if form.Notes != "Test notes" {
		t.Fatalf(
			"expected trimmed notes %q, got %q",
			"Test notes",
			form.Notes,
		)
	}

	if form.Status != clientStatusActive {
		t.Fatalf(
			"expected default status %q, got %q",
			clientStatusActive,
			form.Status,
		)
	}
}

func TestValidateClientForm(t *testing.T) {
	tests := []struct {
		name      string
		form      clientForm
		assertion func(*testing.T, clientFormErrors)
	}{
		{
			name: "valid form",
			form: clientForm{
				Name:    "Test Client",
				Email:   "test@example.com",
				Company: "Test Company",
				Status:  clientStatusActive,
				Notes:   "Test notes",
			},
			assertion: func(
				t *testing.T,
				errors clientFormErrors,
			) {
				t.Helper()

				if errors.Any() {
					t.Fatalf(
						"expected no validation errors, got %+v",
						errors,
					)
				}
			},
		},
		{
			name: "missing name",
			form: clientForm{
				Email:  "test@example.com",
				Status: clientStatusActive,
			},
			assertion: func(
				t *testing.T,
				errors clientFormErrors,
			) {
				t.Helper()

				if errors.Name != "Name is required." {
					t.Fatalf(
						"unexpected name error %q",
						errors.Name,
					)
				}
			},
		},
		{
			name: "invalid email",
			form: clientForm{
				Name:   "Test Client",
				Email:  "not-an-email",
				Status: clientStatusActive,
			},
			assertion: func(
				t *testing.T,
				errors clientFormErrors,
			) {
				t.Helper()

				if errors.Email !=
					"Please enter a valid email address." {
					t.Fatalf(
						"unexpected email error %q",
						errors.Email,
					)
				}
			},
		},
		{
			name: "invalid status",
			form: clientForm{
				Name:   "Test Client",
				Email:  "test@example.com",
				Status: "banana",
			},
			assertion: func(
				t *testing.T,
				errors clientFormErrors,
			) {
				t.Helper()

				if errors.Status != "Invalid client status." {
					t.Fatalf(
						"unexpected status error %q",
						errors.Status,
					)
				}
			},
		},
		{
			name: "fields exceed maximum lengths",
			form: clientForm{
				Name: strings.Repeat(
					"a",
					maxClientNameLength+1,
				),
				Email: strings.Repeat(
					"a",
					maxClientEmailLength+1,
				),
				Company: strings.Repeat(
					"a",
					maxClientCompanyLength+1,
				),
				Status: clientStatusActive,
				Notes: strings.Repeat(
					"a",
					maxClientNotesLength+1,
				),
			},
			assertion: func(
				t *testing.T,
				errors clientFormErrors,
			) {
				t.Helper()

				if errors.Name == "" {
					t.Fatal("expected a name length error")
				}

				if errors.Email == "" {
					t.Fatal("expected an email length error")
				}

				if errors.Company == "" {
					t.Fatal("expected a company length error")
				}

				if errors.Notes == "" {
					t.Fatal("expected a notes length error")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.assertion(
				t,
				validateClientForm(test.form),
			)
		})
	}
}

func TestClientsCreateValidationPreservesSubmittedValues(
	t *testing.T,
) {
	handler, _ := newTestClientsHandler(t)

	form := url.Values{}
	form.Set("name", "Submitted Client")
	form.Set("email", "not-an-email")
	form.Set("company", "Submitted Company")
	form.Set("notes", "Submitted notes")

	req := formRequest(
		http.MethodPost,
		"/dashboard/clients/new",
		form,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d",
			rec.Code,
		)
	}

	for _, expected := range []string{
		"Please enter a valid email address.",
		"Submitted Client",
		"not-an-email",
		"Submitted Company",
		"Submitted notes",
	} {
		if !strings.Contains(
			rec.Body.String(),
			expected,
		) {
			t.Fatalf(
				"expected response to contain %q, got %q",
				expected,
				rec.Body.String(),
			)
		}
	}
}

func TestClientsUpdateValidationPreservesSubmittedValues(
	t *testing.T,
) {
	handler, db := newTestClientsHandler(t)

	id := createTestClient(t, db)

	form := url.Values{}
	form.Set("name", "Submitted Update")
	form.Set("email", "not-an-email")
	form.Set("company", "Updated Company")
	form.Set("status", clientStatusInactive)
	form.Set("notes", "Updated notes")

	req := formRequest(
		http.MethodPost,
		clientsBasePath+"/"+strconv.FormatInt(id, 10),
		form,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status 400, got %d",
			rec.Code,
		)
	}

	for _, expected := range []string{
		"Please enter a valid email address.",
		"Submitted Update",
		"not-an-email",
		"Updated Company",
		clientStatusInactive,
		"Updated notes",
	} {
		if !strings.Contains(
			rec.Body.String(),
			expected,
		) {
			t.Fatalf(
				"expected response to contain %q, got %q",
				expected,
				rec.Body.String(),
			)
		}
	}
}
