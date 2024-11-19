package configs

import "os"

// RedisConfig
type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

func GetRedisConfig() *RedisConfig {
	return &RedisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}
}
