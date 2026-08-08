package main

import (
	"log"
	"net/http"

	"github.com/danieljmanningdev/saas-master-boilerplate/backend/db"
	"github.com/danieljmanningdev/saas-master-boilerplate/backend/routes"
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
