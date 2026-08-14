package contracts

import (
	"bytes"
	"database/sql"
	"errors"
	"html/template"
	"net/http"
)

func (h *ContractsHandler) listContracts(
	w http.ResponseWriter,
	r *http.Request,
) {
	contracts, err := h.contractRepository.List(
		r.Context(),
	)
	if err != nil {
		http.Error(
			w,
			"failed to load contracts",
			http.StatusInternalServerError,
		)
		return
	}

	items := make(
		[]contractListItem,
		0,
		len(contracts),
	)

	for _, contract := range contracts {
		items = append(
			items,
			contractListItem{
				Contract:     contract,
				DisplayValue: formatContractValue(contract.ValueCents),
			},
		)
	}

	h.renderContractPage(
		w,
		h.contractsTemplates,
		contractsPageData{
			Title:     "Contracts — Daniel J. Manning",
			Contracts: items,
		},
	)
}

func (h *ContractsHandler) showContract(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {
	contract, err := h.contractRepository.GetByID(
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
			"failed to load contract",
			http.StatusInternalServerError,
		)
		return
	}

	h.renderContractPage(
		w,
		h.contractTemplates,
		contractPageData{
			Title:        "Contract — Daniel J. Manning",
			Contract:     contract,
			DisplayValue: formatContractValue(contract.ValueCents),
		},
	)
}

func (h *ContractsHandler) editContract(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {
	contract, err := h.contractRepository.GetByID(
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
			"failed to load contract",
			http.StatusInternalServerError,
		)
		return
	}

	h.renderEditContract(
		w,
		r,
		contractFormFromModel(contract),
		contractFormErrors{},
		http.StatusOK,
	)
}

func (h *ContractsHandler) renderNewContract(
	w http.ResponseWriter,
	r *http.Request,
	form contractForm,
	formErrors contractFormErrors,
	status int,
) {
	if form.Status == "" {
		form.Status = contractStatusDraft
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

	h.renderContractStatus(
		w,
		h.newContractTemplates,
		contractFormPageData{
			Title:    "Add Contract — Daniel J. Manning",
			Form:     form,
			Errors:   formErrors,
			Clients:  clients,
			Projects: projects,
		},
		status,
	)
}

func (h *ContractsHandler) renderEditContract(
	w http.ResponseWriter,
	r *http.Request,
	form contractForm,
	formErrors contractFormErrors,
	status int,
) {
	if form.Status == "" {
		form.Status = contractStatusDraft
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

	h.renderContractStatus(
		w,
		h.editContractTemplates,
		contractFormPageData{
			Title:    "Edit Contract — Daniel J. Manning",
			Form:     form,
			Errors:   formErrors,
			Clients:  clients,
			Projects: projects,
		},
		status,
	)
}

func (h *ContractsHandler) renderContractPage(
	w http.ResponseWriter,
	tmpl *template.Template,
	data any,
) {
	h.renderContractStatus(
		w,
		tmpl,
		data,
		http.StatusOK,
	)
}

func (h *ContractsHandler) renderContractStatus(
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
