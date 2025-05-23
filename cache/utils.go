package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

// util to load value from cache
func Get[T any](redisConn *redis.Client, ctx context.Context, key string, dest *T) (bool, error) {
	if redisConn == nil {
		return false, nil
	}
	val, err := redisConn.Get(ctx, key).Result()
	if err != nil {
		return false, nil
	}
	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return false, err
	}
	return true, nil
}

// func to save data into cache
func Set[T any](redisConn *redis.Client, ctx context.Context, key string, value T, ttl time.Duration) error {
	if redisConn == nil {
		return nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return redisConn.Set(ctx, key, data, ttl).Err()
}
