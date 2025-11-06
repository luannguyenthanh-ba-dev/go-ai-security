package cache

import (
	"context"
	"time"
)

// CacheClient defines the interface for cache operations
// This abstraction allows different cache implementations (Redis, Memcached, etc.)
// In Clean Architecture, this is an infrastructure abstraction (not Application Service)
type CacheClient interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error)
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) (bool, error)
	Exists(ctx context.Context, key string) (bool, error)
	Incr(ctx context.Context, key string) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)
}
