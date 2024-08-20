package main

import (
	"log"
	"net/http"

	"tier3-app/config"
	"tier3-app/database"
	"tier3-app/handlers"
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
	queueHandler := handlers.NewQueueHandler(queueService)

	// Initialize the router
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/api/v1/queue", queueHandler.GetQueue).Methods("GET")
	r.HandleFunc("/api/v1/queue", queueHandler.AddToQueue).Methods("POST")

	// Start the server
	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
