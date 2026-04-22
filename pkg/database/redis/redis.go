package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(ctx context.Context, redisURL string) *RedisClient {
	if redisURL == "" {
		redisURL = "localhost:6379" // defualt fallback
	}

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		opts = &redis.Options{
			Addr: redisURL, // treat as direct host:port if url parse fails
		}
	}

	client := redis.NewClient(opts)

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("redis failed to connect: %v", err)
	} else {
		log.Println("redis connected successfully")
	}

	return &RedisClient{Client: client}
}
