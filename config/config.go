package config

import (
	"os"
)

type Config struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

func LoadConfig() Config {
	return Config{
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       0, // Default to 0
	}
}
