package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/danieljmanningdev/danieljmanningdev-portfolio/backend/db"
	"github.com/danieljmanningdev/danieljmanningdev-portfolio/backend/routes"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 1. Initialize Database & Run Migrations
	database, err := db.InitDB("app.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Ensure database migrations run so tables (like users) exist
	if err := db.RunMigrations(database); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// 2. Check for CLI command: `go run main.go create-client <email> <password>`
	if len(os.Args) > 1 && os.Args[1] == "create-client" {
		if len(os.Args) < 4 {
			log.Fatalf("Usage: go run main.go create-client <email> <password>")
		}
		email := os.Args[2]
		password := os.Args[3]
		CreateClient(database, email, password)
		return // Exit after creating client so it doesn't start the web server
	}

	// 3. Setup Routes
	mux := routes.SetupRouter(database)

	// 4. Start Server
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
