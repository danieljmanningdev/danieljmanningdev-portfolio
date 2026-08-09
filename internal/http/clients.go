package http

import (
	"database/sql"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/repository"
)

type ClientsHandler struct {
	repository *repository.ClientRepository

	clientsTemplates   *template.Template
	newClientTemplates *template.Template
	clientTemplates    *template.Template
}

func NewClientsHandler(
	db *sql.DB,
	templateDir string,
) (*ClientsHandler, error) {
	repo := repository.NewClientRepository(db)

	clientsTemplates, err := loadPageTemplate(
		templateDir,
		"clients.html",
	)
	if err != nil {
		return nil, err
	}

	newClientTemplates, err := loadPageTemplate(
		templateDir,
		"client-new.html",
	)
	if err != nil {
		return nil, err
	}

	clientTemplates, err := loadPageTemplate(
		templateDir,
		"client.html",
	)
	if err != nil {
		return nil, err
	}

	return &ClientsHandler{
		repository:         repo,
		clientsTemplates:   clientsTemplates,
		newClientTemplates: newClientTemplates,
		clientTemplates:    clientTemplates,
	}, nil
}

func loadPageTemplate(
	templateDir string,
	page string,
) (*template.Template, error) {
	return template.New("base").ParseFiles(
		filepath.Join(templateDir, "layouts", "base.html"),
		filepath.Join(templateDir, "components", "header.html"),
		filepath.Join(templateDir, "components", "footer.html"),
		filepath.Join(templateDir, "pages", page),
	)
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
		"/dashboard/clients",
	)

	// GET /dashboard/clients
	if path == "" || path == "/" {
		h.listClients(w, r)
		return
	}

	// GET /dashboard/clients/new
	if path == "/new" {
		h.renderNewClient(w, r, "")
		return
	}

	// GET /dashboard/clients/:id
	idString := strings.TrimPrefix(path, "/")

	if strings.Contains(idString, "/") {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.ParseInt(idString, 10, 64)
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
		"/dashboard/clients",
	)

	if path != "/new" {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.renderNewClient(
			w,
			r,
			"Invalid form submission.",
		)
		return
	}

	name := strings.TrimSpace(
		r.FormValue("name"),
	)
	email := strings.TrimSpace(
		r.FormValue("email"),
	)
	company := strings.TrimSpace(
		r.FormValue("company"),
	)
	notes := strings.TrimSpace(
		r.FormValue("notes"),
	)

	if name == "" {
		h.renderNewClient(
			w,
			r,
			"Name is required.",
		)
		return
	}

	if email == "" {
		h.renderNewClient(
			w,
			r,
			"Email is required.",
		)
		return
	}

	id, err := h.repository.Create(
		r.Context(),
		name,
		email,
		company,
		notes,
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
		"/dashboard/clients/"+strconv.FormatInt(id, 10),
		http.StatusSeeOther,
	)
}

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

	data := struct {
		Title   string
		Clients any
	}{
		Title:   "Clients — Daniel J. Manning",
		Clients: clients,
	}

	h.render(
		w,
		h.clientsTemplates,
		data,
	)
}

func (h *ClientsHandler) renderNewClient(
	w http.ResponseWriter,
	r *http.Request,
	errorMessage string,
) {
	data := struct {
		Title   string
		Error   string
		Name    string
		Email   string
		Company string
		Notes   string
	}{
		Title: "Add Client — Daniel J. Manning",
		Error: errorMessage,
		Name:  strings.TrimSpace(r.FormValue("name")),
		Email: strings.TrimSpace(r.FormValue("email")),
		Company: strings.TrimSpace(
			r.FormValue("company"),
		),
		Notes: strings.TrimSpace(
			r.FormValue("notes"),
		),
	}

	h.render(
		w,
		h.newClientTemplates,
		data,
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
		if err == sql.ErrNoRows {
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

	data := struct {
		Title  string
		Client any
	}{
		Title:  "Client — Daniel J. Manning",
		Client: client,
	}

	h.render(
		w,
		h.clientTemplates,
		data,
	)
}

func (h *ClientsHandler) render(
	w http.ResponseWriter,
	templates *template.Template,
	data any,
) {
	w.Header().Set(
		"Content-Type",
		"text/html; charset=utf-8",
	)

	if err := templates.ExecuteTemplate(
		w,
		"base",
		data,
	); err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
	}
}
