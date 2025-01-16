package main

import (
	"forum/handlers"
	"forum/middleware"
	"forum/routes"
	"forum/utils"
	"log"
	"net/http"
)

func main() {
	// Initialize database
	db := utils.InitialiseDB()
	defer db.Close()

	// Initialize handlers with database connection
	handlers.InitDB(db)

	// Create ServeMux for routing
	mux := http.NewServeMux()

	// Set up routes
	routes.SetupRoutes(mux)

	// Apply middleware
	handler := middleware.LoggingMiddleware(
		middleware.AuthMiddleware(db)(mux),
	)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
