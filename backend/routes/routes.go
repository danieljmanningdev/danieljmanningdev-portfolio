package routes

import (
	"database/sql"
	"net/http"
)

func SetupRouter(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	// Serve local static assets (like htmx.min.js)
	fileServer := http.FileServer(http.Dir("backend/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Home page route
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "backend/templates/index.html")
	})

	return mux
}
