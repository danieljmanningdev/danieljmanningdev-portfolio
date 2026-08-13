package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/auth"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/database"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/repository"
)

func newRouterTestSessionService(
	t *testing.T,
) *auth.SessionService {
	t.Helper()

	db, err := database.Open(
		context.Background(),
		":memory:",
	)
	if err != nil {
		t.Fatalf(
			"open test database: %v",
			err,
		)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve test path")
	}

	migrationsDir := filepath.Join(
		filepath.Dir(filename),
		"..",
		"..",
		"migrations",
	)

	if err := database.RunMigrations(
		db.SQL,
		migrationsDir,
	); err != nil {
		t.Fatalf(
			"run migrations: %v",
			err,
		)
	}

	return auth.NewSessionService(
		repository.NewAdminRepository(
			db.SQL,
		),
		repository.NewAdminSessionRepository(
			db.SQL,
		),
	)
}

func TestNewRouterPublicRoutes(
	t *testing.T,
) {
	sessionService :=
		newRouterTestSessionService(t)

	homeHandler := http.HandlerFunc(
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			_, _ = w.Write(
				[]byte("home"),
			)
		},
	)

	authHandler := http.HandlerFunc(
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			_, _ = w.Write(
				[]byte("auth"),
			)
		},
	)

	protectedHandler := http.HandlerFunc(
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			_, _ = w.Write(
				[]byte("protected"),
			)
		},
	)

	caseStudyHandler := http.HandlerFunc(
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			_, _ = w.Write(
				[]byte("case-study"),
			)
		},
	)

	router := newRouter(
		homeHandler,
		caseStudyHandler,
		authHandler,
		protectedHandler,
		protectedHandler,
		protectedHandler,
		protectedHandler,
		sessionService,
	)

	tests := []struct {
		path     string
		wantBody string
	}{
		{
			path:     "/",
			wantBody: "home",
		},
		{
			path:     "/work/portfolio",
			wantBody: "case-study",
		},
		{
			path:     "/login",
			wantBody: "auth",
		},
		{
			path:     "/logout",
			wantBody: "auth",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.path,
			func(t *testing.T) {
				req := httptest.NewRequest(
					http.MethodGet,
					tt.path,
					nil,
				)

				rec := httptest.NewRecorder()

				router.ServeHTTP(
					rec,
					req,
				)

				if rec.Code !=
					http.StatusOK {
					t.Fatalf(
						"expected 200, got %d",
						rec.Code,
					)
				}

				if rec.Body.String() !=
					tt.wantBody {
					t.Fatalf(
						"expected body %q, got %q",
						tt.wantBody,
						rec.Body.String(),
					)
				}
			},
		)
	}
}

func TestNewRouterProtectsDashboardRoutes(
	t *testing.T,
) {
	sessionService :=
		newRouterTestSessionService(t)

	handler := http.HandlerFunc(
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			_, _ = w.Write(
				[]byte("protected"),
			)
		},
	)

	router := newRouter(
		handler,
		handler,
		handler,
		handler,
		handler,
		handler,
		handler,
		sessionService,
	)

	protectedPaths := []string{
		"/dashboard/",
		"/dashboard/clients",
		"/dashboard/clients/1",
		"/dashboard/projects",
		"/dashboard/projects/1",
		"/dashboard/contracts",
		"/dashboard/contracts/1",
	}

	for _, path := range protectedPaths {
		t.Run(
			path,
			func(t *testing.T) {
				req := httptest.NewRequest(
					http.MethodGet,
					path,
					nil,
				)

				rec := httptest.NewRecorder()

				router.ServeHTTP(
					rec,
					req,
				)

				if rec.Code !=
					http.StatusSeeOther {
					t.Fatalf(
						"expected 303, got %d",
						rec.Code,
					)
				}

				if location :=
					rec.Header().Get(
						"Location",
					); location != "/login" {
					t.Fatalf(
						"expected /login redirect, got %q",
						location,
					)
				}
			},
		)
	}
}

func TestNewRouterRedirectsDashboardWithoutTrailingSlash(
	t *testing.T,
) {
	sessionService :=
		newRouterTestSessionService(t)

	handler := http.HandlerFunc(
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			w.WriteHeader(
				http.StatusOK,
			)
		},
	)

	router := newRouter(
		handler,
		handler,
		handler,
		handler,
		handler,
		handler,
		handler,
		sessionService,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/dashboard",
		nil,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(
		rec,
		req,
	)

	if rec.Code !=
		http.StatusPermanentRedirect {
		t.Fatalf(
			"expected 308, got %d",
			rec.Code,
		)
	}

	if location :=
		rec.Header().Get(
			"Location",
		); location != "/dashboard/" {
		t.Fatalf(
			"expected /dashboard/ redirect, got %q",
			location,
		)
	}
}
