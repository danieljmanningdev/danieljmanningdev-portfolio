package handlers

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/backend/models"
)

func PortalHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Read the session cookie (set during login)
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		email := cookie.Value

		// 2. Fetch the user details from SQLite
		var user models.User
		err = db.QueryRow("SELECT id, email, role, created_at FROM users WHERE email = ?", email).
			Scan(&user.ID, &user.Email, &user.Role, &user.CreatedAt)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// 3. Parse and render the portal template with user data
		tmpl, err := template.ParseFiles("backend/templates/portal.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Pass user object to the HTML template
		tmpl.Execute(w, user)
	}
}
