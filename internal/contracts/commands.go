// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package contracts

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
)

func (h *ContractsHandler) createContract(
	w http.ResponseWriter,
	r *http.Request,
) {
	form, err := parseContractForm(r)
	if err != nil {
		h.renderNewContract(
			w,
			r,
			contractForm{
				Status: contractStatusDraft,
			},
			contractFormErrors{
				General: "Invalid form submission.",
			},
			http.StatusBadRequest,
		)
		return
	}

	formErrors := validateContractForm(form)

	if !formErrors.Any() {
		formErrors, err = h.validateContractRelationships(
			r.Context(),
			form,
			formErrors,
		)
		if err != nil {
			http.Error(
				w,
				"failed to validate contract relationships",
				http.StatusInternalServerError,
			)
			return
		}
	}

	if formErrors.Any() {
		h.renderNewContract(
			w,
			r,
			form,
			formErrors,
			http.StatusBadRequest,
		)
		return
	}

	startDate, endDate, err := contractFormDateValues(form)
	if err != nil {
		h.renderNewContract(
			w,
			r,
			form,
			contractFormErrors{
				General: "Invalid contract dates.",
			},
			http.StatusBadRequest,
		)
		return
	}

	valueCents, err := contractValueCentsPointer(form)
	if err != nil {
		h.renderNewContract(
			w,
			r,
			form,
			contractFormErrors{
				General: "Invalid contract value.",
			},
			http.StatusBadRequest,
		)
		return
	}

	id, err := h.contractRepository.Create(
		r.Context(),
		form.ClientID,
		contractProjectIDValue(form),
		form.Title,
		form.Status,
		startDate,
		endDate,
		valueCents,
		form.Notes,
	)
	if err != nil {
		http.Error(
			w,
			"failed to create contract",
			http.StatusInternalServerError,
		)
		return
	}

	http.Redirect(
		w,
		r,
		contractsBasePath+"/"+strconv.FormatInt(id, 10),
		http.StatusSeeOther,
	)
}

func (h *ContractsHandler) updateContract(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {
	form, err := parseContractForm(r)
	form.ID = id

	if err != nil {
		h.renderEditContract(
			w,
			r,
			form,
			contractFormErrors{
				General: "Invalid form submission.",
			},
			http.StatusBadRequest,
		)
		return
	}

	formErrors := validateContractForm(form)

	if !formErrors.Any() {
		formErrors, err = h.validateContractRelationships(
			r.Context(),
			form,
			formErrors,
		)
		if err != nil {
			http.Error(
				w,
				"failed to validate contract relationships",
				http.StatusInternalServerError,
			)
			return
		}
	}

	if formErrors.Any() {
		h.renderEditContract(
			w,
			r,
			form,
			formErrors,
			http.StatusBadRequest,
		)
		return
	}

	startDate, endDate, err := contractFormDateValues(form)
	if err != nil {
		h.renderEditContract(
			w,
			r,
			form,
			contractFormErrors{
				General: "Invalid contract dates.",
			},
			http.StatusBadRequest,
		)
		return
	}

	valueCents, err := contractValueCentsPointer(form)
	if err != nil {
		h.renderEditContract(
			w,
			r,
			form,
			contractFormErrors{
				General: "Invalid contract value.",
			},
			http.StatusBadRequest,
		)
		return
	}

	err = h.contractRepository.Update(
		r.Context(),
		id,
		form.ClientID,
		contractProjectIDValue(form),
		form.Title,
		form.Status,
		startDate,
		endDate,
		valueCents,
		form.Notes,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}

		http.Error(
			w,
			"failed to update contract",
			http.StatusInternalServerError,
		)
		return
	}

	http.Redirect(
		w,
		r,
		contractsBasePath+"/"+strconv.FormatInt(id, 10),
		http.StatusSeeOther,
	)
}

func (h *ContractsHandler) cancelContract(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {
	err := h.contractRepository.Cancel(
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
			"failed to cancel contract",
			http.StatusInternalServerError,
		)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set(
			"HX-Redirect",
			contractsBasePath,
		)

		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(
		w,
		r,
		contractsBasePath,
		http.StatusSeeOther,
	)
}
