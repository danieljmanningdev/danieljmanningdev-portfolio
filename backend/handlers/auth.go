package handlers

import (
	database "database/sql"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func HandleRegister(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			http.Error(w, "Email and password are required", http.StatusBadRequest)
			return
		}

		// Hash password securely
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Insert user into SQLite
		_, err = db.Exec("INSERT INTO users (email, password_hash) VALUES (?, ?)", email, string(hashedPassword))
		if err != nil {
			http.Error(w, "Email already registered or database error", http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User registered successfully!"))
	}
}
