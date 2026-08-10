package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRouterRoutesRequests(t *testing.T) {
	homeHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("home"))
		},
	)

	dashboardHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("dashboard"))
		},
	)

	clientsHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("clients"))
		},
	)

	projectsHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("projects"))
		},
	)

	contractsHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("contracts"))
		},
	)

	router := newRouter(
		homeHandler,
		dashboardHandler,
		clientsHandler,
		projectsHandler,
		contractsHandler,
	)

	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "homepage",
			path:       "/",
			wantStatus: http.StatusOK,
			wantBody:   "home",
		},
		{
			name:       "dashboard",
			path:       "/dashboard/",
			wantStatus: http.StatusOK,
			wantBody:   "dashboard",
		},
		{
			name:       "clients",
			path:       "/dashboard/clients",
			wantStatus: http.StatusOK,
			wantBody:   "clients",
		},
		{
			name:       "client sub-route",
			path:       "/dashboard/clients/1",
			wantStatus: http.StatusOK,
			wantBody:   "clients",
		},
		{
			name:       "projects",
			path:       "/dashboard/projects",
			wantStatus: http.StatusOK,
			wantBody:   "projects",
		},
		{
			name:       "project sub-route",
			path:       "/dashboard/projects/1",
			wantStatus: http.StatusOK,
			wantBody:   "projects",
		},
		{
			name:       "contracts",
			path:       "/dashboard/contracts",
			wantStatus: http.StatusOK,
			wantBody:   "contracts",
		},
		{
			name:       "contract sub-route",
			path:       "/dashboard/contracts/1",
			wantStatus: http.StatusOK,
			wantBody:   "contracts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodGet,
				tt.path,
				nil,
			)

			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf(
					"expected status %d, got %d",
					tt.wantStatus,
					rec.Code,
				)
			}

			if rec.Body.String() != tt.wantBody {
				t.Fatalf(
					"expected body %q, got %q",
					tt.wantBody,
					rec.Body.String(),
				)
			}
		})
	}
}

func TestNewRouterRedirectsDashboardWithoutTrailingSlash(
	t *testing.T,
) {
	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	)

	router := newRouter(
		handler,
		handler,
		handler,
		handler,
		handler,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/dashboard",
		nil,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusPermanentRedirect {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusPermanentRedirect,
			rec.Code,
		)
	}

	if location := rec.Header().Get("Location"); location != "/dashboard/" {
		t.Fatalf(
			"expected redirect location %q, got %q",
			"/dashboard/",
			location,
		)
	}
}
