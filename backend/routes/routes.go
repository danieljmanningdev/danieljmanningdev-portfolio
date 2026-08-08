package routes

import (
	"database/sql"
	"net/http"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/backend/handlers"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/backend/middleware"
)

func SetupRouter(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	// Serve local static assets (like htmx.min.js)
	fileServer := http.FileServer(http.Dir("backend/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Home page route
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		// Catch accidental sub-paths
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "backend/templates/index.html")
	})

	// Handle blog route
	mux.HandleFunc("GET /blog", handlers.BlogHandler)

	// Protected client portal route
	mux.Handle("GET /portal", middleware.RequireAuth(handlers.PortalHandler(db)))

	// Protected Stripe billing checkout route
	mux.Handle("GET /portal/billing", middleware.RequireAuth(http.HandlerFunc(handlers.HandleCreateCheckout)))

	// Handle Stripe Webhooks (Public route)
	mux.HandleFunc("POST /webhook", handlers.StripeWebhookHandler)

	// Handle calculate route
	mux.HandleFunc("POST /calculate", handlers.CalculateHandler)

	// Handle login
	mux.HandleFunc("GET /login", handlers.HandleLogin(db))
	mux.HandleFunc("POST /login", handlers.HandleLogin(db))

	// Register logout route to clear the session cookie and redirect back to login
	mux.HandleFunc("GET /logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "session_token",
			Value:  "",
			Path:   "/",
			MaxAge: -1, // Deletes the cookie immediately
		})
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})

	return mux
}
