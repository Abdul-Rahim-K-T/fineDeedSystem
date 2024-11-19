package database

import (
	"fmt"
	"log"
	"time"

	"fineDeedSystem/employer-service/configs"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

// RedisClient wraps the Redis client
type RedisClient struct {
	client *redis.Client
}

// InitRedis initializes the Redis client
func InitRedis(config configs.RedisConfig) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: config.Password, // no password set
		DB:       0,               // use default DB
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return &RedisClient{client: rdb}
}

// Set is a wrapper to set key-value pairs in Redis
func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.client.Set(context.Background(), key, value, expiration)
}

// Get is a wrapper to get values from Redis by key
func (r *RedisClient) Get(key string) *redis.StringCmd {
	return r.client.Get(context.Background(), key)
}

// Close closes the Redis client
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// Del is a wrapper to delete a key from Redis
func (r *RedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return r.client.Del(ctx, keys...)
}
