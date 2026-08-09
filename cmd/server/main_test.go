package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewRouterUsesExplicitDashboardRoutes(t *testing.T) {
	homeHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				http.NotFound(w, r)
				return
			}

			_, _ = w.Write([]byte("home"))
		},
	)

	dashboardHandler := http.HandlerFunc(
		func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte("dashboard"))
		},
	)

	clientsHandler := http.HandlerFunc(
		func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte("clients"))
		},
	)

	router := newRouter(
		homeHandler,
		dashboardHandler,
		clientsHandler,
	)

	tests := []struct {
		name             string
		target           string
		expectedStatus   int
		expectedBody     string
		expectedLocation string
	}{
		{
			name:           "public homepage",
			target:         "/",
			expectedStatus: http.StatusOK,
			expectedBody:   "home",
		},
		{
			name:             "dashboard redirects to canonical URL",
			target:           "/dashboard",
			expectedStatus:   http.StatusPermanentRedirect,
			expectedLocation: "/dashboard/",
		},
		{
			name:           "dashboard overview",
			target:         "/dashboard/",
			expectedStatus: http.StatusOK,
			expectedBody:   "dashboard",
		},
		{
			name:           "client list",
			target:         "/dashboard/clients",
			expectedStatus: http.StatusOK,
			expectedBody:   "clients",
		},
		{
			name:           "client list with trailing slash",
			target:         "/dashboard/clients/",
			expectedStatus: http.StatusOK,
			expectedBody:   "clients",
		},
		{
			name:           "client detail",
			target:         "/dashboard/clients/42",
			expectedStatus: http.StatusOK,
			expectedBody:   "clients",
		},
		{
			name:           "client edit",
			target:         "/dashboard/clients/42/edit",
			expectedStatus: http.StatusOK,
			expectedBody:   "clients",
		},
		{
			name:           "unknown dashboard domain",
			target:         "/dashboard/projects",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404 page not found",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodGet,
				test.target,
				nil,
			)

			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != test.expectedStatus {
				t.Fatalf(
					"expected status %d, got %d",
					test.expectedStatus,
					rec.Code,
				)
			}

			if test.expectedLocation != "" {
				location := rec.Header().Get("Location")

				if location != test.expectedLocation {
					t.Fatalf(
						"expected location %q, got %q",
						test.expectedLocation,
						location,
					)
				}
			}

			if test.expectedBody != "" &&
				!strings.Contains(
					rec.Body.String(),
					test.expectedBody,
				) {
				t.Fatalf(
					"expected body to contain %q, got %q",
					test.expectedBody,
					rec.Body.String(),
				)
			}
		})
	}
}
