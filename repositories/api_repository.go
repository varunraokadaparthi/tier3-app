package repositories

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"tier3-app/models"
)

type QueueRepository struct {
	redisClient *redis.Client
}

func NewQueueRepository(redisClient *redis.Client) *QueueRepository {
	return &QueueRepository{redisClient: redisClient}
}

func (r *QueueRepository) GetAll() ([]models.QueueItem, error) {
	ctx := context.Background()
	result, err := r.redisClient.LRange(ctx, "queue", 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var queue []models.QueueItem
	for _, item := range result {
		var queueItem models.QueueItem
		if err := json.Unmarshal([]byte(item), &queueItem); err != nil {
			return nil, err
		}
		queue = append(queue, queueItem)
	}

	return queue, nil
}

func (r *QueueRepository) Add(item models.QueueItem) (models.QueueItem, error) {
	ctx := context.Background()
	item.ID = r.getNextID(ctx)
	itemBytes, err := json.Marshal(item)
	if err != nil {
		return models.QueueItem{}, err
	}

	//// Check the type of the key before using it
	//keyType, err := r.redisClient.Type(ctx, "queue").Result()
	//if err != nil {
	//	return models.QueueItem{}, err
	//}
	//log.Printf("Type of 'queue' key: %s", keyType)
	//if keyType != "none" && keyType != "list" {
	//	return models.QueueItem{}, fmt.Errorf("key 'queue' is of type %s, expected list", keyType)
	//}

	if err := r.redisClient.RPush(ctx, "queue", itemBytes).Err(); err != nil {
		return models.QueueItem{}, err
	}

	return item, nil
}

func (r *QueueRepository) getNextID(ctx context.Context) uint {
	id, err := r.redisClient.Incr(ctx, "queue_id").Result()
	if err != nil {
		return 0
	}
	return uint(id)
}
