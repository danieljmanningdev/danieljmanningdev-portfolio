// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package clients

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
)

func (h *ClientsHandler) createClient(
	w http.ResponseWriter,
	r *http.Request,
) {
	form, err := parseClientForm(r)
	if err != nil {
		h.renderNewClient(
			w,
			clientForm{
				Status: clientStatusActive,
			},
			clientFormErrors{
				General: "Invalid form submission.",
			},
			http.StatusBadRequest,
		)
		return
	}

	/*
		New clients always begin as active. The create form does
		not accept a client-controlled status value.
	*/
	form.Status = clientStatusActive

	formErrors := validateClientForm(form)

	if formErrors.Any() {
		h.renderNewClient(
			w,
			form,
			formErrors,
			http.StatusOK,
		)
		return
	}

	id, err := h.repository.Create(
		r.Context(),
		form.Name,
		form.Email,
		form.Company,
		form.Notes,
	)
	if err != nil {
		http.Error(
			w,
			"failed to create client",
			http.StatusInternalServerError,
		)
		return
	}

	http.Redirect(
		w,
		r,
		clientsBasePath+"/"+strconv.FormatInt(id, 10),
		http.StatusSeeOther,
	)
}

func (h *ClientsHandler) updateClient(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {
	form, err := parseClientForm(r)
	form.ID = id

	if err != nil {
		h.renderEditClient(
			w,
			form,
			clientFormErrors{
				General: "Invalid form submission.",
			},
			http.StatusBadRequest,
		)
		return
	}

	formErrors := validateClientForm(form)

	if formErrors.Any() {
		h.renderEditClient(
			w,
			form,
			formErrors,
			http.StatusBadRequest,
		)
		return
	}

	err = h.repository.Update(
		r.Context(),
		id,
		form.Name,
		form.Email,
		form.Company,
		form.Status,
		form.Notes,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}

		http.Error(
			w,
			"failed to update client",
			http.StatusInternalServerError,
		)
		return
	}

	http.Redirect(
		w,
		r,
		clientsBasePath+"/"+strconv.FormatInt(id, 10),
		http.StatusSeeOther,
	)
}

func (h *ClientsHandler) handleDeleteClient(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {
	err := h.repository.Delete(
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
			"failed to delete client",
			http.StatusInternalServerError,
		)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set(
			"HX-Redirect",
			clientsBasePath,
		)

		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(
		w,
		r,
		clientsBasePath,
		http.StatusSeeOther,
	)
}
