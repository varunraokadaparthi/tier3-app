package main

import (
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"tier3-app/config"
	"tier3-app/database"
	"tier3-app/handler"
	"tier3-app/repositories"
	"tier3-app/services"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize Redis
	database.InitRedis(cfg)

	// Initialize the database
	// (Assuming a similar function InitDB exists for your database setup)

	// Migrate the schema
	// db.AutoMigrate(&models.QueueItem{})

	// Set up repositories, services, and handlers
	queueRepo := repositories.NewQueueRepository(database.RedisClient) // or modify if using Redis only
	queueService := services.NewQueueService(queueRepo)
	queueHandler := handler.NewQueueHandler(queueService)

	// Initialize the router
	r := mux.NewRouter()

	// Handle preflight requests globally
	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
	})

	r.Use(handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	))

	// Define routes
	r.HandleFunc("/api/v1/queue", queueHandler.GetQueue).Methods("GET")
	r.HandleFunc("/api/v1/adduser", queueHandler.AddToQueue).Methods("POST")

	// Apply CORS middleware

	// Start the server
	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
