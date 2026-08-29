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
	apphttp "github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/http"
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

	blogHandler := &apphttp.BlogHandler{}

	router := newRouter(
		homeHandler,
		caseStudyHandler,
		authHandler,
		protectedHandler,
		protectedHandler,
		protectedHandler,
		protectedHandler,
		protectedHandler,
		sessionService,
		blogHandler,
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
		handler,
		sessionService,
		&apphttp.BlogHandler{},
	)

	protectedPaths := []string{
		"/dashboard/",
		"/dashboard/clients",
		"/dashboard/clients/1",
		"/dashboard/projects",
		"/dashboard/projects/1",
		"/dashboard/contracts",
		"/dashboard/contracts/1",
		"/dashboard/blog",
		"/dashboard/blog/1/edit",
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
		handler,
		sessionService,
		&apphttp.BlogHandler{},
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

func TestNewRouterNoIndexPolicy(
	t *testing.T,
) {
	sessionService :=
		newRouterTestSessionService(t)

	handler := http.HandlerFunc(
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
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
		handler,
		sessionService,
		&apphttp.BlogHandler{},
	)

	tests := []struct {
		path        string
		wantNoIndex bool
	}{
		{
			path:        "/",
			wantNoIndex: false,
		},
		{
			path:        "/work/portfolio",
			wantNoIndex: false,
		},
		{
			path:        "/login",
			wantNoIndex: true,
		},
		{
			path:        "/logout",
			wantNoIndex: true,
		},
		{
			path:        "/dashboard",
			wantNoIndex: true,
		},
		{
			path:        "/dashboard/",
			wantNoIndex: true,
		},
		{
			path:        "/dashboard/clients",
			wantNoIndex: true,
		},
		{
			path:        "/dashboard/projects",
			wantNoIndex: true,
		},
		{
			path:        "/dashboard/contracts",
			wantNoIndex: true,
		},
		{
			path:        "/dashboard/blog",
			wantNoIndex: true,
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

				router.ServeHTTP(rec, req)

				got := rec.Header().Get(
					"X-Robots-Tag",
				)

				if tt.wantNoIndex {
					if got != "noindex, nofollow" {
						t.Fatalf(
							"expected noindex header, got %q",
							got,
						)
					}

					return
				}

				if got != "" {
					t.Fatalf(
						"expected no X-Robots-Tag, got %q",
						got,
					)
				}
			},
		)
	}
}
