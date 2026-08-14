package clients

import (
	"bytes"
	"database/sql"
	"errors"
	"html/template"
	"net/http"
)

func (h *ClientsHandler) listClients(
	w http.ResponseWriter,
	r *http.Request,
) {
	clients, err := h.repository.List(
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

	h.render(
		w,
		h.clientsTemplates,
		clientsPageData{
			Title:   "Clients — Daniel J. Manning",
			Clients: clients,
		},
	)
}

func (h *ClientsHandler) showClient(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {
	client, err := h.repository.GetByID(
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
			"failed to load client",
			http.StatusInternalServerError,
		)
		return
	}

	h.render(
		w,
		h.clientTemplates,
		clientPageData{
			Title:  "Client — Daniel J. Manning",
			Client: client,
		},
	)
}

func (h *ClientsHandler) editClient(
	w http.ResponseWriter,
	r *http.Request,
	id int64,
) {
	client, err := h.repository.GetByID(
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
			"failed to load client",
			http.StatusInternalServerError,
		)
		return
	}

	h.renderEditClient(
		w,
		clientFormFromModel(client),
		clientFormErrors{},
		http.StatusOK,
	)
}

func (h *ClientsHandler) renderNewClient(
	w http.ResponseWriter,
	form clientForm,
	errors clientFormErrors,
	status int,
) {
	if form.Status == "" {
		form.Status = clientStatusActive
	}

	h.renderStatus(
		w,
		h.newClientTemplates,
		clientFormPageData{
			Title:  "Add Client — Daniel J. Manning",
			Form:   form,
			Errors: errors,
		},
		status,
	)
}

func (h *ClientsHandler) renderEditClient(
	w http.ResponseWriter,
	form clientForm,
	errors clientFormErrors,
	status int,
) {
	if form.Status == "" {
		form.Status = clientStatusActive
	}

	h.renderStatus(
		w,
		h.editClientTemplates,
		clientFormPageData{
			Title:  "Edit Client — Daniel J. Manning",
			Form:   form,
			Errors: errors,
		},
		status,
	)
}

func (h *ClientsHandler) render(
	w http.ResponseWriter,
	tmpl *template.Template,
	data any,
) {
	h.renderStatus(
		w,
		tmpl,
		data,
		http.StatusOK,
	)
}

func (h *ClientsHandler) renderStatus(
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
