package http

import (
	"bytes"
	"database/sql"
	"html/template"
	"net/http"
	"strings"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/models"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/rendering"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/repository"
)

type dashboardPageData struct {
	Title           string
	LogoutCSRFToken string
	Summary         models.DashboardSummary
}

type activityPageData struct {
	Title  string
	Filter string
	Events []models.AuditEvent
}

type DashboardHandler struct {
	dashboardRepository *repository.DashboardRepository
	auditRepository     *repository.AuditRepository
	dashboardTemplates  *template.Template
	activityTemplates   *template.Template
}

func NewDashboardHandler(
	db *sql.DB,
	templateDir string,
) (*DashboardHandler, error) {
	dashboardTemplates, err := rendering.LoadPageTemplate(
		templateDir,
		"admin/dashboard.html",
	)
	if err != nil {
		return nil, err
	}

	activityTemplates, err := rendering.LoadPageTemplate(
		templateDir,
		"admin/activity.html",
	)
	if err != nil {
		return nil, err
	}

	return &DashboardHandler{
		dashboardRepository: repository.NewDashboardRepository(db),
		auditRepository:     repository.NewAuditRepository(db),
		dashboardTemplates:  dashboardTemplates,
		activityTemplates:   activityTemplates,
	}, nil
}

func (h *DashboardHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
		return
	}

	switch r.URL.Path {
	case "/dashboard/":
		h.showDashboard(w, r)

	case "/dashboard/activity":
		h.showActivity(w, r)

	default:
		http.NotFound(w, r)
	}
}

func (h *DashboardHandler) showDashboard(
	w http.ResponseWriter,
	r *http.Request,
) {
	summary, err := h.dashboardRepository.Summary(
		r.Context(),
	)
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	logoutCSRFToken, _ :=
		AdminLogoutCSRFTokenFromContext(
			r.Context(),
		)

	h.renderDashboard(
		w,
		h.dashboardTemplates,
		dashboardPageData{
			Title:           "Dashboard — Daniel J. Manning",
			LogoutCSRFToken: logoutCSRFToken,
			Summary:         summary,
		},
	)
}

func (h *DashboardHandler) showActivity(
	w http.ResponseWriter,
	r *http.Request,
) {
	filter := normalizedActivityFilter(
		r.URL.Query().Get("type"),
	)

	events, err := h.auditRepository.ListRecent(
		r.Context(),
		100,
		filter,
	)
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	h.renderDashboard(
		w,
		h.activityTemplates,
		activityPageData{
			Title:  "Activity — Daniel J. Manning",
			Filter: filter,
			Events: events,
		},
	)
}

func (h *DashboardHandler) renderDashboard(
	w http.ResponseWriter,
	tmpl *template.Template,
	data any,
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
	w.WriteHeader(http.StatusOK)
	_, _ = body.WriteTo(w)
}

func normalizedActivityFilter(value string) string {
	switch strings.TrimSpace(value) {
	case "client", "project", "contract", "blog_post":
		return strings.TrimSpace(value)
	default:
		return ""
	}
}
