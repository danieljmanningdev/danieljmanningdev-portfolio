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

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/config"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/database"
	apphttp "github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/http"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/logging"
)

func newRouter(
	homeHandler http.Handler,
	dashboardHandler http.Handler,
	clientsHandler http.Handler,
	projectsHandler http.Handler,
) http.Handler {
	mux := http.NewServeMux()

	// Public application routes.
	mux.HandleFunc("/health", apphttp.HealthHandler)

	// Redirect the dashboard path to its canonical trailing-slash URL.
	mux.Handle(
		"/dashboard",
		http.RedirectHandler(
			"/dashboard/",
			http.StatusPermanentRedirect,
		),
	)

	// Match only the dashboard overview itself.
	mux.Handle(
		"/dashboard/{$}",
		dashboardHandler,
	)

	// Match the client list and every client sub-route.
	mux.Handle(
		"/dashboard/clients",
		clientsHandler,
	)

	mux.Handle(
		"/dashboard/clients/",
		clientsHandler,
	)

	// Match the project list and every project sub-route.
	mux.Handle(
		"/dashboard/projects",
		projectsHandler,
	)

	mux.Handle(
		"/dashboard/projects/",
		projectsHandler,
	)

	// Static files.
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

	// The public homepage remains the final fallback.
	// HomeHandler returns 404 for paths other than "/".
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

	clientsHandler, err := apphttp.NewClientsHandler(
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

	projectsHandler, err := apphttp.NewProjectsHandler(
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

	router := newRouter(
		homeHandler,
		dashboardHandler,
		clientsHandler,
		projectsHandler,
	)

	server := &http.Server{
		Addr:              ":" + strconv.Itoa(cfg.Port),
		Handler:           apphttp.RequestLogger(logger, router),
		ReadHeaderTimeout: 5 * time.Second,
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
