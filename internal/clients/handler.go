package clients

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/rendering"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/repository"
)

const clientsBasePath = "/dashboard/clients"

type ClientsHandler struct {
	repository *repository.ClientRepository

	clientsTemplates    *template.Template
	newClientTemplates  *template.Template
	clientTemplates     *template.Template
	editClientTemplates *template.Template
}

func NewClientsHandler(
	db *sql.DB,
	templateDir string,
) (*ClientsHandler, error) {
	repo := repository.NewClientRepository(db)

	clientsTemplates, err := rendering.LoadPageTemplate(
		templateDir,
		"clients.html",
	)
	if err != nil {
		return nil, err
	}

	newClientTemplates, err := rendering.LoadPageTemplate(
		templateDir,
		"client-new.html",
	)
	if err != nil {
		return nil, err
	}

	clientTemplates, err := rendering.LoadPageTemplate(
		templateDir,
		"client.html",
	)
	if err != nil {
		return nil, err
	}

	editClientTemplates, err := rendering.LoadPageTemplate(
		templateDir,
		"client-edit.html",
	)
	if err != nil {
		return nil, err
	}

	return &ClientsHandler{
		repository:          repo,
		clientsTemplates:    clientsTemplates,
		newClientTemplates:  newClientTemplates,
		clientTemplates:     clientTemplates,
		editClientTemplates: editClientTemplates,
	}, nil
}

func (h *ClientsHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	switch r.Method {
	case http.MethodGet:
		h.handleGET(w, r)

	case http.MethodPost:
		h.handlePOST(w, r)

	case http.MethodDelete:
		h.handleDELETE(w, r)

	default:
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
	}
}

func (h *ClientsHandler) handleGET(
	w http.ResponseWriter,
	r *http.Request,
) {
	path := strings.TrimPrefix(
		r.URL.Path,
		clientsBasePath,
	)

	// GET /dashboard/clients
	if path == "" || path == "/" {
		h.listClients(w, r)
		return
	}

	// GET /dashboard/clients/new
	if path == "/new" {
		h.renderNewClient(
			w,
			clientForm{
				Status: clientStatusActive,
			},
			clientFormErrors{},
			http.StatusOK,
		)
		return
	}

	// GET /dashboard/clients/:id/edit
	if strings.HasSuffix(path, "/edit") {
		idString := strings.TrimSuffix(
			strings.TrimPrefix(path, "/"),
			"/edit",
		)

		id, err := strconv.ParseInt(
			idString,
			10,
			64,
		)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		h.editClient(w, r, id)
		return
	}

	// GET /dashboard/clients/:id
	idString := strings.TrimPrefix(path, "/")

	if idString == "" || strings.Contains(idString, "/") {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.ParseInt(
		idString,
		10,
		64,
	)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	h.showClient(w, r, id)
}

func (h *ClientsHandler) handlePOST(
	w http.ResponseWriter,
	r *http.Request,
) {
	path := strings.TrimPrefix(
		r.URL.Path,
		clientsBasePath,
	)

	// POST /dashboard/clients/new
	if path == "/new" {
		h.createClient(w, r)
		return
	}

	// POST /dashboard/clients/:id
	idString := strings.TrimPrefix(path, "/")

	if idString == "" || strings.Contains(idString, "/") {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.ParseInt(
		idString,
		10,
		64,
	)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	h.updateClient(w, r, id)
}

func (h *ClientsHandler) handleDELETE(
	w http.ResponseWriter,
	r *http.Request,
) {
	path := strings.TrimPrefix(
		r.URL.Path,
		clientsBasePath,
	)

	idString := strings.TrimPrefix(path, "/")

	if idString == "" || strings.Contains(idString, "/") {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.ParseInt(
		idString,
		10,
		64,
	)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	h.handleDeleteClient(w, r, id)
}
