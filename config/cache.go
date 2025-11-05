package config

import "github.com/redis/go-redis/v9"

type CacheConnection interface {
	Close() error
	GetRedisClient() *redis.Client
}
