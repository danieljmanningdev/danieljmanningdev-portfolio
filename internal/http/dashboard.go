package http

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type dashboardPageData struct {
	Title           string
	LogoutCSRFToken string
}

type DashboardHandler struct {
	templates *template.Template
}

func NewDashboardHandler(
	templateDir string,
) (*DashboardHandler, error) {
	templates, err := template.New("base").ParseFiles(
		filepath.Join(
			templateDir,
			"layouts",
			"admin-base.html",
		),
		filepath.Join(
			templateDir,
			"components",
			"header.html",
		),
		filepath.Join(
			templateDir,
			"components",
			"footer.html",
		),
		filepath.Join(
			templateDir,
			"pages",
			"admin/dashboard.html",
		),
	)
	if err != nil {
		return nil, err
	}

	return &DashboardHandler{
		templates: templates,
	}, nil
}

func (h *DashboardHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.URL.Path != "/dashboard/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
		return
	}

	logoutCSRFToken, _ :=
		AdminLogoutCSRFTokenFromContext(
			r.Context(),
		)

	data := dashboardPageData{
		Title:           "Dashboard — Daniel J. Manning",
		LogoutCSRFToken: logoutCSRFToken,
	}

	w.Header().Set(
		"Content-Type",
		"text/html; charset=utf-8",
	)

	if err := h.templates.ExecuteTemplate(
		w,
		"base",
		data,
	); err != nil {
		http.Error(
			w,
			http.StatusText(
				http.StatusInternalServerError,
			),
			http.StatusInternalServerError,
		)
	}
}
