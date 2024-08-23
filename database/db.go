package database

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"

	"tier3-app/config"
)

var (
	RedisClient *redis.Client
	ctx         = context.Background()
)

func InitRedis(cfg config.Config) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// Ping Redis to ensure connection is successful
	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}
