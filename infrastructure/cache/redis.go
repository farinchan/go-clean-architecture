package cache

import (
	"context"
	"fmt"
	"time"

	"kopeta-backend/config"
	"kopeta-backend/pkg/logger"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

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
		logger.Errorf("Failed to connect to Redis: %v", err)
		return nil, err
	}

	logger.Info("Successfully connected to Redis")
	return &RedisClient{Client: client}, nil
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *RedisClient) Delete(ctx context.Context, keys ...string) error {
	return r.Client.Del(ctx, keys...).Err()
}

func (r *RedisClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.Client.Exists(ctx, keys...).Result()
}

func (r *RedisClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return r.Client.SetNX(ctx, key, value, expiration).Result()
}

func (r *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.Client.Expire(ctx, key, expiration).Err()
}

func (r *RedisClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.Client.TTL(ctx, key).Result()
}

func (r *RedisClient) Incr(ctx context.Context, key string) (int64, error) {
	return r.Client.Incr(ctx, key).Result()
}

func (r *RedisClient) Decr(ctx context.Context, key string) (int64, error) {
	return r.Client.Decr(ctx, key).Result()
}

// Hash operations
func (r *RedisClient) HSet(ctx context.Context, key string, values ...interface{}) error {
	return r.Client.HSet(ctx, key, values...).Err()
}

func (r *RedisClient) HGet(ctx context.Context, key, field string) (string, error) {
	return r.Client.HGet(ctx, key, field).Result()
}

func (r *RedisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.Client.HGetAll(ctx, key).Result()
}

func (r *RedisClient) HDel(ctx context.Context, key string, fields ...string) error {
	return r.Client.HDel(ctx, key, fields...).Err()
}

// List operations
func (r *RedisClient) LPush(ctx context.Context, key string, values ...interface{}) error {
	return r.Client.LPush(ctx, key, values...).Err()
}

func (r *RedisClient) RPush(ctx context.Context, key string, values ...interface{}) error {
	return r.Client.RPush(ctx, key, values...).Err()
}

func (r *RedisClient) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.Client.LRange(ctx, key, start, stop).Result()
}

func (r *RedisClient) Close() error {
	return r.Client.Close()
}
