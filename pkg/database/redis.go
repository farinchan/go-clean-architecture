package database

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/your-username/go-clean-architecture/config"
	"github.com/your-username/go-clean-architecture/pkg/logger"
)

// RedisClient holds the redis client
type RedisClient struct {
	Client *redis.Client
}

// NewRedisClient creates a new redis client
func NewRedisClient(cfg *config.RedisConfig) (*RedisClient, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	logger.Info("Redis connected successfully")

	return &RedisClient{Client: client}, nil
}

// Close closes the redis connection
func (r *RedisClient) Close() error {
	return r.Client.Close()
}

// Set sets a key-value pair with expiration
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

// Get gets a value by key
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

// Delete deletes a key
func (r *RedisClient) Delete(ctx context.Context, keys ...string) error {
	return r.Client.Del(ctx, keys...).Err()
}

// Exists checks if a key exists
func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}
