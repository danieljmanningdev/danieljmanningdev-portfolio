package handlers

import (
	database "database/sql"
	"net/http"
	"time"

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

		// Insert user into SQLite (defaulting role to client if applicable)
		_, err = db.Exec("INSERT INTO users (email, password_hash, role) VALUES (?, ?, ?)", email, string(hashedPassword), "client")
		if err != nil {
			http.Error(w, "Email already registered or database error", http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User registered successfully!"))
	}
}

func HandleLogin(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.ServeFile(w, r, "backend/templates/login.html")
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		var id int64
		var hash string
		err := db.QueryRow("SELECT id, password_hash FROM users WHERE email = ?", email).Scan(&id, &hash)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// Set session token cookie
		cookie := http.Cookie{
			Name:     "session_token",
			Value:    email, // Replace with signed token/JWT for production
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Path:     "/",
		}
		http.SetCookie(w, &cookie)

		// Support HTMX redirect if requested, otherwise standard browser redirect
		if r.Header.Get("HX-Request") == "true" {
			w.Header().Set("HX-Redirect", "/portal")
			w.WriteHeader(http.StatusOK)
			return
		}

		http.Redirect(w, r, "/portal", http.StatusSeeOther)
	}
}
