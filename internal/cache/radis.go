package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(addr, pass string, db int) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})

	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)

	defer close()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connnect to redis %w", err)
	}

	slog.Info("redis client connnected", "addr", addr)
	return &RedisClient{client: client}, nil
}

func (r *RedisClient) Get(ctx context.Context, key string, target interface{}) error {
	val, err := r.client.Get(ctx, key).Result()

	if err == redis.Nil {
		return fmt.Errorf("cache miss")
	}

	if err != nil {
		slog.Error("error getting the cache", "key", key, "error", err)
		return err
	}

	if err := json.Unmarshal([]byte(val), target); err != nil {
		slog.Error("failed to Unmarshal data", "key", key, "error", err)
		return err
	}

	slog.Debug("cache hit", "key", key)
	return nil
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		slog.Error("error Marshaling data", "key", key, "error", err)
		return err
	}

	sixHours := 6 * time.Hour

	if err := r.client.Set(ctx, key, jsonData, sixHours).Err(); err != nil {
		slog.Error("error setting data to redis", "key", key, "error", err)
	}

	slog.Debug("setting was successfull", "key", key, "ttl", sixHours)
	return nil
}

func (r *RedisClient) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}
