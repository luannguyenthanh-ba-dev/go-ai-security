package cache

import (
	"github.com/redis/go-redis/v9"
)

// NewCacheClient creates a Cache directly from Redis client
// This factory function abstracts the creation of cache
func NewCacheClient(client *redis.Client) CacheClient {
	if client == nil {
		return nil
	}

	return NewRedisCacheClient(client)
}
