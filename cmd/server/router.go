package main

import (
	"net/http"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/auth"
	apphttp "github.com/danieljmanningdev/danieljmanningdev-portfolio/internal/http"
)

type routerDependencies struct {
	homeHandler               http.Handler
	portfolioCaseStudyHandler http.Handler
	webDesignLeedsHandler     http.Handler

	adminAuthHandler http.Handler
	dashboardHandler http.Handler
	clientsHandler   http.Handler
	projectsHandler  http.Handler
	contractsHandler http.Handler
	blogAdminHandler http.Handler

	sessionService        *auth.SessionService
	blogHandler           *apphttp.BlogHandler
	webDevelopmentHandler *apphttp.WebDevelopmentHandler
}

func newRouter(deps routerDependencies) http.Handler {
	mux := http.NewServeMux()

	// Public application routes.
	mux.HandleFunc("/health", apphttp.HealthHandler)
	mux.HandleFunc("GET /robots.txt", apphttp.RobotsHandler)
	mux.HandleFunc("GET /sitemap.xml", deps.blogHandler.Sitemap)

	mux.Handle(
		"/work/portfolio",
		deps.portfolioCaseStudyHandler,
	)

	mux.Handle(
		"/web-design-leeds/",
		deps.webDesignLeedsHandler,
	)

	mux.Handle(
		"/web-development/",
		deps.webDevelopmentHandler,
	)

	noIndexAdminAuth := apphttp.NoIndex(
		deps.adminAuthHandler,
	)

	mux.Handle(
		"/login",
		noIndexAdminAuth,
	)

	mux.Handle(
		"/logout",
		noIndexAdminAuth,
	)

	protectedDashboard := apphttp.NoIndex(
		apphttp.NoStore(
			apphttp.RequireAdmin(
				deps.sessionService,
				deps.dashboardHandler,
			),
		),
	)

	protectedClients := apphttp.NoIndex(
		apphttp.NoStore(
			apphttp.RequireAdmin(
				deps.sessionService,
				deps.clientsHandler,
			),
		),
	)

	protectedProjects := apphttp.NoIndex(
		apphttp.NoStore(
			apphttp.RequireAdmin(
				deps.sessionService,
				deps.projectsHandler,
			),
		),
	)

	protectedContracts := apphttp.NoIndex(
		apphttp.NoStore(
			apphttp.RequireAdmin(
				deps.sessionService,
				deps.contractsHandler,
			),
		),
	)

	protectedBlogAdmin := apphttp.NoIndex(
		apphttp.NoStore(
			apphttp.RequireAdmin(
				deps.sessionService,
				deps.blogAdminHandler,
			),
		),
	)

	mux.Handle(
		"/dashboard",
		apphttp.NoIndex(
			http.RedirectHandler(
				"/dashboard/",
				http.StatusPermanentRedirect,
			),
		),
	)

	mux.Handle("/dashboard/{$}", protectedDashboard)
	mux.Handle("/dashboard/activity", protectedDashboard)

	mux.Handle("/dashboard/clients", protectedClients)
	mux.Handle("/dashboard/clients/", protectedClients)

	mux.Handle("/dashboard/projects", protectedProjects)
	mux.Handle("/dashboard/projects/", protectedProjects)

	mux.Handle("/dashboard/contracts", protectedContracts)
	mux.Handle("/dashboard/contracts/", protectedContracts)

	mux.Handle("/dashboard/blog", protectedBlogAdmin)
	mux.Handle("/dashboard/blog/", protectedBlogAdmin)

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
		deps.blogHandler.List,
	)

	mux.HandleFunc(
		"GET /blog/{slug}",
		deps.blogHandler.Show,
	)

	mux.Handle(
		"/",
		deps.homeHandler,
	)

	return mux
}
