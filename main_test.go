package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"tier3-app/config"
	"tier3-app/database"
	"tier3-app/handlers"
	"tier3-app/models"
	"tier3-app/repositories"
	"tier3-app/services"
)

func clearQueue() {
	// Clear the queue before each test
	ctx := context.Background()
	database.RedisClient.FlushAll(ctx)
}

func TestGetQueue(t *testing.T) {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize Redis
	database.InitRedis(cfg)

	// Clear the queue
	clearQueue()

	// Set up repositories, services, and handlers
	queueRepo := repositories.NewQueueRepository(database.RedisClient)
	queueService := services.NewQueueService(queueRepo)
	queueHandler := handlers.NewQueueHandler(queueService)

	req, err := http.NewRequest("GET", "/api/v1/queue", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(queueHandler.GetQueue)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `[]` // Expecting an empty queue initially
	if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestAddToQueue(t *testing.T) {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize Redis
	database.InitRedis(cfg)

	// Clear the queue
	clearQueue()

	// Set up repositories, services, and handlers
	queueRepo := repositories.NewQueueRepository(database.RedisClient)
	queueService := services.NewQueueService(queueRepo)
	queueHandler := handlers.NewQueueHandler(queueService)

	newItem := models.QueueItem{Name: "John Doe", Email: "john.doe@example.com"}
	jsonItem, _ := json.Marshal(newItem)
	req, err := http.NewRequest("POST", "/api/v1/queue", bytes.NewBuffer(jsonItem))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(queueHandler.AddToQueue)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var responseItem models.QueueItem
	if err := json.NewDecoder(rr.Body).Decode(&responseItem); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if responseItem.ID == 0 {
		t.Errorf("handler returned unexpected ID: got %v want non-zero ID", responseItem.ID)
	}

	expected := models.QueueItem{ID: 1, Name: "John Doe", Email: "john.doe@example.com"}
	if responseItem.Name != expected.Name || responseItem.Email != expected.Email {
		t.Errorf("handler returned unexpected body: got %v want %v", responseItem, expected)
	}
}
