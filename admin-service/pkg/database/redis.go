package database

import (
	"context"
	"log"

	"fineDeedSystem/admin-service/configs"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func InitRedis(config configs.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password,
		DB:       0,
	})

	// Test connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	return client
}

type RedisClient struct {
	Client *redis.Client
}

func (r *RedisClient) InvalidateCache(key string) error {
	ctx := context.Background()
	err := r.Client.Del(ctx, key).Err()
	if err != nil {
		log.Printf("Failed to invalidate cache for key %s: %v", key, err)
		return err
	}
	return nil
}
