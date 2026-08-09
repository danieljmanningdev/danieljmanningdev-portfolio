package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/config"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/database"
	apphttp "github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/http"
)

func newRouter(
	homeHandler http.Handler,
	dashboardHandler http.Handler,
	clientsHandler http.Handler,
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
		log.Fatalf("open database: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("close database: %v", err)
		}
	}()

	if err := database.RunMigrations(
		db.SQL,
		"migrations",
	); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	homeHandler, err := apphttp.NewHomeHandler(
		cfg.TemplateDir,
	)
	if err != nil {
		log.Fatalf("create home handler: %v", err)
	}

	dashboardHandler, err := apphttp.NewDashboardHandler(
		cfg.TemplateDir,
	)
	if err != nil {
		log.Fatalf("create dashboard handler: %v", err)
	}

	clientsHandler, err := apphttp.NewClientsHandler(
		db.SQL,
		cfg.TemplateDir,
	)
	if err != nil {
		log.Fatalf("create clients handler: %v", err)
	}

	router := newRouter(
		homeHandler,
		dashboardHandler,
		clientsHandler,
	)

	server := &http.Server{
		Addr:              ":" + strconv.Itoa(cfg.Port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf(
			"server starting on http://localhost:%d (%s)",
			cfg.Port,
			cfg.Environment,
		)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Printf("server error: %v", err)
			stop()
		}
	}()

	<-ctx.Done()

	log.Println("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf(
			"server shutdown error: %v",
			err,
		)
	}
}
