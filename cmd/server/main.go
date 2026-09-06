// -----------------------------------------------------------------------------
// Daniel J. Manning
// https://danieljmanningdev.com
//
// Copyright © 2026 Daniel J. Manning.
// -----------------------------------------------------------------------------

package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/auth"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/blog"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/clients"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/config"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/contracts"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/database"
	apphttp "github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/http"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/logging"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/projects"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/repository"
)

func main() { os.Exit(run()) }

func run() int {
	cfg := config.Load()

	logger := logging.New(
		cfg.Environment,
		cfg.LogLevel,
	)

	slog.SetDefault(logger)

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
		return 1
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
		return 1
	}

	homeHandler, err := apphttp.NewHomeHandler(
		cfg.TemplateDir,
	)
	if err != nil {
		logger.Error(
			"failed to create home handler",
			"error", err,
		)
		return 1
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
		return 1
	}

	dashboardHandler, err := apphttp.NewDashboardHandler(
		db.SQL,
		cfg.TemplateDir,
	)
	if err != nil {
		logger.Error(
			"failed to create dashboard handler",
			"error", err,
		)
		return 1
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
		return 1
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
		return 1
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
		return 1
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
		return 1
	}

	blogAdminHandler, err := blog.NewAdminHandler(
		db.SQL,
		cfg.TemplateDir,
	)
	if err != nil {
		logger.Error(
			"failed to create blog admin handler",
			"error", err,
		)
		return 1
	}

	publicPageRoutes, err :=
		apphttp.BuildPublicPageRoutes(
			cfg.TemplateDir,
		)

	if err != nil {
		logger.Error(
			"failed to create public page routes",
			"error", err,
		)
		return 1
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
		return 1
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

	router := newRouter(routerDependencies{
		homeHandler:               homeHandler,
		portfolioCaseStudyHandler: portfolioCaseStudyHandler,

		publicPageRoutes: publicPageRoutes,

		adminAuthHandler: adminAuthHandler,
		dashboardHandler: dashboardHandler,
		clientsHandler:   clientsHandler,
		projectsHandler:  projectsHandler,
		contractsHandler: contractsHandler,
		blogAdminHandler: blogAdminHandler,

		sessionService: sessionService,
		blogHandler:    blogHandler,
	})

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

	logger.Info("server starting", "url", "http://localhost:"+strconv.Itoa(cfg.Port), "environment", cfg.Environment)
	if err := runHTTPServer(ctx, server); err != nil {
		logger.Error("server stopped with error", "error", err)
		return 1
	}
	logger.Info("server stopped")
	return 0
}

// runHTTPServer reports startup failures to the process and closes connections
// if the graceful-shutdown deadline expires. Tests can exercise it without a DB.
func runHTTPServer(ctx context.Context, server *http.Server) error {
	result := make(chan error, 1)
	go func() { result <- server.ListenAndServe() }()
	select {
	case err := <-result:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		_ = server.Close()
		return err
	}
	return nil
}
