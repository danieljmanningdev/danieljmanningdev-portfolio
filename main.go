package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/backend/db"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/backend/routes"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 1. Initialize Database
	database, err := db.InitDB("app.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// 2. Setup Routes (passing database or services if needed)
	mux := routes.SetupRouter(database)

	// 3. Start Server
	port := ":8080"
	log.Printf("Server starting on http://localhost%s", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func CreateClient(db *sql.DB, email, rawPassword string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	_, err = db.Exec(
		"INSERT INTO users (email, password_hash, role, created_at) VALUES (?, ?, ?, datetime('now'))",
		email, string(hashedPassword), "client",
	)
	if err != nil {
		log.Fatalf("Failed to insert client: %v", err)
	}

	fmt.Printf("Successfully created client account for: %s\n", email)
}
