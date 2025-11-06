package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RedisCacheClient implements Cache interface using Redis
type RedisCacheClient struct {
	client *redis.Client
}

// NewRedisCacheClient creates a new Redis cache implementation
func NewRedisCacheClient(client *redis.Client) CacheClient {
	return &RedisCacheClient{
		client: client,
	}
}

// Set sets a key-value pair with expiration
// - If TTL is 0, the value will be stored forever
// - If TTL is negative, the value will be stored forever
// - If TTL is positive, the value will be stored for the given duration
func (c *RedisCacheClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	if value == nil {
		return false, ErrCacheKeyEmpty
	}

	if key == "" {
		return false, ErrCacheKeyEmpty
	}

	var writeData interface{}
	switch v := value.(type) {
	case string:
		// String: pass directly to Redis (no encoding needed)
		writeData = v
	case []byte:
		// []byte: pass directly to Redis (no encoding needed)
		writeData = v
	default:
		// Complex types (struct, map, slice, array, etc.): marshal to JSON
		jsonData, err := json.Marshal(v)
		if err != nil {
			zap.L().Error("failed to marshal cache value to JSON",
				zap.String("key", key),
				zap.Error(err),
			)
			return false, err
		}
		// Store JSON as []byte or string - Redis will handle it
		writeData = jsonData
	}

	// Set the value in the cache
	// - If TTL is 0, the value will be stored forever
	// - If TTL is negative, the value will be stored forever
	// - If TTL is positive, the value will be stored for the given duration
	// Redis client automatically handles string, []byte, and other types
	_, err := c.client.Set(ctx, key, writeData, ttl).Result()
	if err != nil {
		zap.L().Error("failed to set cache",
			zap.String("key", key),
			zap.Error(err),
		)
		return false, err
	}

	return true, nil
}

// Get gets a value by key
func (c *RedisCacheClient) Get(ctx context.Context, key string) (string, error) {
	// Validate the key
	if key == "" {
		return "", ErrCacheKeyEmpty
	}
	// Get the value from the cache
	value, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil // Cache miss - return nil, nil
		}
		zap.L().Error("failed to get cache",
			zap.String("key", key),
			zap.Error(err),
		)
		return "", err
	}

	return value, nil
}

// Del deletes a key from the cache
func (c *RedisCacheClient) Del(ctx context.Context, key string) (bool, error) {
	// Validate the key
	if key == "" {
		return false, ErrCacheKeyEmpty
	}

	// Delete the key from the cache - returns the number of keys deleted
	_, err := c.client.Del(ctx, key).Result()
	if err != nil {
		zap.L().Error("failed to delete cache",
			zap.String("key", key),
			zap.Error(err),
		)
		return false, err
	}

	return true, nil
}

// Exists checks if a key exists in the cache
func (c *RedisCacheClient) Exists(ctx context.Context, key string) (bool, error) {
	// Validate the key
	if key == "" {
		return false, ErrCacheKeyEmpty
	}

	// Check if the key exists in the cache - returns the number of keys found
	exists, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil // Cache miss - return false, nil
		}
		zap.L().Error("failed to check if cache exists",
			zap.String("key", key),
			zap.Error(err),
		)
		return false, err
	}
	return exists > 0, nil
}

// Incr increments the value of a key
func (c *RedisCacheClient) Incr(ctx context.Context, key string) (int64, error) {
	if key == "" {
		return 0, ErrCacheKeyEmpty
	}

	// Increment the value of the key - returns the new value
	// If key does not exist, it will be set to 0 and then incremented
	value, err := c.client.Incr(ctx, key).Result()
	if err != nil {
		zap.L().Error("failed to increment cache",
			zap.String("key", key),
			zap.Error(err),
		)
		return 0, err
	}
	return value, nil
}

// Decr decrements the value of a key
func (c *RedisCacheClient) Decr(ctx context.Context, key string) (int64, error) {
	if key == "" {
		return 0, ErrCacheKeyEmpty
	}

	// Decrement the value of the key - returns the new value
	// If key does not exist, it will be set to 0 and then decremented
	value, err := c.client.Decr(ctx, key).Result()
	if err != nil {
		zap.L().Error("failed to decrement cache",
			zap.String("key", key),
			zap.Error(err),
		)
		return 0, err
	}
	return value, nil
}
