package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/auth"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/clients"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/config"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/contracts"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/database"
	apphttp "github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/http"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/logging"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/projects"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/repository"
)

func newRouter(
	homeHandler http.Handler,
	portfolioCaseStudyHandler http.Handler,
	adminAuthHandler http.Handler,
	dashboardHandler http.Handler,
	clientsHandler http.Handler,
	projectsHandler http.Handler,
	contractsHandler http.Handler,
	sessionService *auth.SessionService,
	blogHandler *apphttp.BlogHandler,
) http.Handler {
	mux := http.NewServeMux()

	// Public application routes.
	mux.HandleFunc("/health", apphttp.HealthHandler)

	mux.Handle(
		"/work/portfolio",
		portfolioCaseStudyHandler,
	)

	noIndexAdminAuth := apphttp.NoIndex(
		adminAuthHandler,
	)

	mux.Handle(
		"/login",
		noIndexAdminAuth,
	)

	mux.Handle(
		"/logout",
		noIndexAdminAuth,
	)

	// All dashboard routes require a valid admin session.
	protectedDashboard := apphttp.NoIndex(
		apphttp.NoStore(
			apphttp.RequireAdmin(
				sessionService,
				dashboardHandler,
			),
		),
	)

	protectedClients := apphttp.NoIndex(
		apphttp.NoStore(
			apphttp.RequireAdmin(
				sessionService,
				clientsHandler,
			),
		),
	)

	protectedProjects := apphttp.NoIndex(
		apphttp.NoStore(
			apphttp.RequireAdmin(
				sessionService,
				projectsHandler,
			),
		),
	)

	protectedContracts := apphttp.NoIndex(
		apphttp.NoStore(
			apphttp.RequireAdmin(
				sessionService,
				contractsHandler,
			),
		),
	)

	// Redirect the dashboard path to its canonical trailing-slash URL.
	mux.Handle(
		"/dashboard",
		apphttp.NoIndex(
			http.RedirectHandler(
				"/dashboard/",
				http.StatusPermanentRedirect,
			),
		),
	)

	// Match only the dashboard overview itself.
	mux.Handle(
		"/dashboard/{$}",
		protectedDashboard,
	)

	// Match the client list and every client sub-route.
	mux.Handle(
		"/dashboard/clients",
		protectedClients,
	)

	mux.Handle(
		"/dashboard/clients/",
		protectedClients,
	)

	// Match the project list and every project sub-route.
	mux.Handle(
		"/dashboard/projects",
		protectedProjects,
	)

	mux.Handle(
		"/dashboard/projects/",
		protectedProjects,
	)

	// Match the contract list and every contract sub-route.
	mux.Handle(
		"/dashboard/contracts",
		protectedContracts,
	)

	mux.Handle(
		"/dashboard/contracts/",
		protectedContracts,
	)

	// Static files remain public.
	fileServer := http.FileServer(
		http.Dir("web/static"),
	)

	mux.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			fileServer,
		),
	)

	mux.HandleFunc(
		"GET /blog/{$}",
		blogHandler.List,
	)

	mux.HandleFunc(
		"GET /blog/{slug}",
		blogHandler.Show,
	)

	// The public homepage remains the final fallback.
	mux.Handle("/", homeHandler)

	return mux
}

func main() {
	cfg := config.Load()

	logger := logging.New(
		cfg.Environment,
		cfg.LogLevel,
	)

	logger.Info(
		"starting application",
		"environment", cfg.Environment,
		"log_level", cfg.LogLevel,
	)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	db, err := database.Open(
		ctx,
		cfg.DatabasePath,
	)
	if err != nil {
		logger.Error(
			"failed to open database",
			"error", err,
		)
		os.Exit(1)
	}

	defer func() {
		if err := db.Close(); err != nil {
			logger.Error(
				"failed to close database",
				"error", err,
			)
		}
	}()

	if err := database.RunMigrations(
		db.SQL,
		"migrations",
	); err != nil {
		logger.Error(
			"failed to run migrations",
			"error", err,
		)
		os.Exit(1)
	}

	homeHandler, err := apphttp.NewHomeHandler(
		cfg.TemplateDir,
	)
	if err != nil {
		logger.Error(
			"failed to create home handler",
			"error", err,
		)
		os.Exit(1)
	}

	portfolioCaseStudyHandler, err :=
		apphttp.NewPortfolioCaseStudyHandler(
			cfg.TemplateDir,
		)
	if err != nil {
		logger.Error(
			"failed to create portfolio case study handler",
			"error", err,
		)
		os.Exit(1)
	}

	dashboardHandler, err := apphttp.NewDashboardHandler(
		cfg.TemplateDir,
	)
	if err != nil {
		logger.Error(
			"failed to create dashboard handler",
			"error", err,
		)
		os.Exit(1)
	}

	clientsHandler, err := clients.NewClientsHandler(
		db.SQL,
		cfg.TemplateDir,
	)
	if err != nil {
		logger.Error(
			"failed to create clients handler",
			"error", err,
		)
		os.Exit(1)
	}

	projectsHandler, err := projects.NewProjectsHandler(
		db.SQL,
		cfg.TemplateDir,
	)
	if err != nil {
		logger.Error(
			"failed to create projects handler",
			"error", err,
		)
		os.Exit(1)
	}

	contractsHandler, err := contracts.NewContractsHandler(
		db.SQL,
		cfg.TemplateDir,
	)
	if err != nil {
		logger.Error(
			"failed to create contracts handler",
			"error", err,
		)
		os.Exit(1)
	}

	blogHandler, err := apphttp.NewBlogHandler(
		db.SQL,
		cfg.TemplateDir,
	)
	if err != nil {
		logger.Error(
			"failed to create blog handler",
			"error", err,
		)
		os.Exit(1)
	}

	secureCookies :=
		cfg.Environment == "production"

	adminAuthHandler, err :=
		apphttp.NewAdminAuthHandler(
			db.SQL,
			cfg.TemplateDir,
			secureCookies,
		)
	if err != nil {
		logger.Error(
			"failed to create admin auth handler",
			"error", err,
		)
		os.Exit(1)
	}

	adminRepository :=
		repository.NewAdminRepository(
			db.SQL,
		)

	sessionRepository :=
		repository.NewAdminSessionRepository(
			db.SQL,
		)

	sessionService :=
		auth.NewSessionService(
			adminRepository,
			sessionRepository,
		)

	router := newRouter(
		homeHandler,
		portfolioCaseStudyHandler,
		adminAuthHandler,
		dashboardHandler,
		clientsHandler,
		projectsHandler,
		contractsHandler,
		sessionService,
		blogHandler,
	)

	crossOriginProtection :=
		http.NewCrossOriginProtection()

	protectedRouter :=
		crossOriginProtection.Handler(
			router,
		)

	securityHeaders := apphttp.SecurityHeaders(
		cfg.Environment == "production",
		protectedRouter,
	)

	applicationHandler :=
		apphttp.RequestLogger(
			logger,
			securityHeaders,
		)

	server := &http.Server{
		Addr:              ":" + strconv.Itoa(cfg.Port),
		Handler:           applicationHandler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	go func() {
		logger.Info(
			"server starting",
			"url", "http://localhost:"+strconv.Itoa(cfg.Port),
			"environment", cfg.Environment,
		)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			logger.Error(
				"server error",
				"error", err,
			)

			stop()
		}
	}()

	<-ctx.Done()

	logger.Info("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error(
			"server shutdown error",
			"error", err,
		)
	}
}
