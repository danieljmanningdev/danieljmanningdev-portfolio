package blog

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/rendering"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/repository"
)

const blogAdminBasePath = "/dashboard/blog"

type AdminHandler struct {
	repository *repository.BlogRepository

	indexTemplates *template.Template
	newTemplates   *template.Template
	editTemplates  *template.Template
}

func NewAdminHandler(
	db *sql.DB,
	templateDir string,
) (*AdminHandler, error) {
	repo := repository.NewBlogRepository(db)

	indexTemplates, err := rendering.LoadPageTemplate(
		templateDir,
		"admin/blog/index.html",
	)
	if err != nil {
		return nil, err
	}

	newTemplates, err := rendering.LoadPageTemplate(
		templateDir,
		"admin/blog/new.html",
	)
	if err != nil {
		return nil, err
	}

	editTemplates, err := rendering.LoadPageTemplate(
		templateDir,
		"admin/blog/edit.html",
	)
	if err != nil {
		return nil, err
	}

	return &AdminHandler{
		repository:     repo,
		indexTemplates: indexTemplates,
		newTemplates:   newTemplates,
		editTemplates:  editTemplates,
	}, nil
}

func (h *AdminHandler) ServeHTTP(
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

func (h *AdminHandler) handleGET(
	w http.ResponseWriter,
	r *http.Request,
) {
	path := strings.TrimPrefix(
		r.URL.Path,
		blogAdminBasePath,
	)

	// GET /dashboard/blog
	if path == "" || path == "/" {
		h.listPosts(w, r)
		return
	}

	// GET /dashboard/blog/new
	if path == "/new" {
		h.renderNewPost(
			w,
			Form{
				Status: "draft",
			},
			FormErrors{},
			http.StatusOK,
		)
		return
	}

	// GET /dashboard/blog/:id/edit
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

		h.editPost(w, r, id)
		return
	}

	http.NotFound(w, r)
}

func (h *AdminHandler) handlePOST(
	w http.ResponseWriter,
	r *http.Request,
) {
	path := strings.TrimPrefix(
		r.URL.Path,
		blogAdminBasePath,
	)

	// POST /dashboard/blog/new
	if path == "/new" {
		h.createPost(w, r)
		return
	}

	// POST /dashboard/blog/:id/publish
	if strings.HasSuffix(path, "/publish") {
		idString := strings.TrimSuffix(
			strings.TrimPrefix(path, "/"),
			"/publish",
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

		h.publishPost(w, r, id)
		return
	}

	// POST /dashboard/blog/:id/unpublish
	if strings.HasSuffix(path, "/unpublish") {
		idString := strings.TrimSuffix(
			strings.TrimPrefix(path, "/"),
			"/unpublish",
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

		h.unpublishPost(w, r, id)
		return
	}

	// POST /dashboard/blog/:id
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

	h.updatePost(w, r, id)
}

func (h *AdminHandler) handleDELETE(
	w http.ResponseWriter,
	r *http.Request,
) {
	path := strings.TrimPrefix(
		r.URL.Path,
		blogAdminBasePath,
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

	h.deletePost(w, r, id)
}
