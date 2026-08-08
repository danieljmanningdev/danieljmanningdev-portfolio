package routes

import (
	"database/sql"
	"net/http"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/backend/handlers"
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

	// Handle portal route (TODO: Wrap with backend/middleware/auth.go later)
	mux.HandleFunc("GET /portal", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "backend/templates/portal.html")
	})

	return mux
}
