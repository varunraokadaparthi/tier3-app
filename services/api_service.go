package services

import (
	"context"
	"encoding/json"

	"tier3-app/database"
	"tier3-app/models"
	"tier3-app/repositories"
)

type QueueService struct {
	Repo *repositories.QueueRepository
}

func NewQueueService(repo *repositories.QueueRepository) *QueueService {
	return &QueueService{Repo: repo}
}

func (s *QueueService) GetQueue() ([]models.QueueItem, error) {
	// Try fetching from Redis cache first
	cachedQueue, err := database.RedisClient.Get(context.Background(), "queue").Result()
	if err == nil && cachedQueue != "" {
		var queue []models.QueueItem
		if err := json.Unmarshal([]byte(cachedQueue), &queue); err == nil {
			return queue, nil
		}
	}

	// If cache miss, get data from the database
	queue, err := s.Repo.GetAll()
	if err != nil {
		return nil, err
	}
	if queue == nil {
		queue = []models.QueueItem{}
	}
	return queue, nil
}

func (s *QueueService) AddToQueue(name, email string) (models.QueueItem, error) {
	item := models.QueueItem{
		Name:  name,
		Email: email,
	}
	createdItem, err := s.Repo.Add(item)
	if err != nil {
		return models.QueueItem{}, err
	}

	// Optionally, update cache here if needed

	return createdItem, nil
}
