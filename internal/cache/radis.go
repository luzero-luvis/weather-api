package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

// creating our own struct to create a connect and perform set and get Delete operation

type RedisClient struct {
	client *redis.Client
}

// this is where we connect to redis

func NewRedisClient(addr, pass string, db int) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})

	// if redis take more than 5 sec stop trying

	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)

	// cut the connection once the work done

	defer close()

	// ping redis

	// commented out because i don't want to termicate code if redis was not there

	if err := client.Ping(ctx).Err(); err != nil {
		// return nil, fmt.Errorf("failed to connnect to redis %w", err)
	}

	// return the clinet connnection

	slog.Info("redis client connnected", "addr", addr)
	return &RedisClient{client: client}, nil
}

// this where you will get the cached data

func (r *RedisClient) Get(ctx context.Context, key string, target interface{}) error {
	val, err := r.client.Get(ctx, key).Result()

	if err == redis.Nil {
		return fmt.Errorf("cache miss")
	}

	if err != nil {
		slog.Error("error getting the cache", "key", key, "error", err)
		return err
	}

	// convert json to go struct fetch

	if err := json.Unmarshal([]byte(val), target); err != nil {
		slog.Error("failed to Unmarshal data", "key", key, "error", err)
		return err
	}

	slog.Debug("cache hit", "key", key)
	return nil
}

// this where you will put cache to redis DB

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}) error {
	// convertin go data to json to put redis
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

func (r *RedisClient) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// deleting the data after 6 hour

func (r *RedisClient) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Close the connection

func (r *RedisClient) Close() error {
	return r.client.Close()
}
