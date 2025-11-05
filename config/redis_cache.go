package config

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisConfig struct {
	Host       string
	Port       int
	DB         int
	MaxRetries int
	RetryDelay time.Duration
}

type RedisCacheConnection struct {
	Client *redis.Client
}

// NewCache creates a new Redis client connection with retry mechanism
// Note: MaxRetries in redis.Options is for operations (GET, SET, etc.), not for initial connection
// We use separate retry logic with Ping() to verify initial connection (similar to MongoDB)
func NewRedisCacheConnection(config RedisConfig) (CacheConnection, error) {
	// Set default values if not provided
	maxRetries := config.MaxRetries
	if maxRetries == 0 {
		maxRetries = 3 // Default 3 retries
	}
	retryDelay := config.RetryDelay
	if retryDelay == 0 {
		retryDelay = 2 * time.Second // Default 2 seconds
	}

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	// Create Redis client (this doesn't actually connect to server)
	// Connection happens lazily on first operation (e.g., Ping)
	client := redis.NewClient(&redis.Options{
		Addr:       addr,
		DB:         config.DB,
		MaxRetries: maxRetries, // This is for operations retry, not initial connection
	})

	zap.L().Debug("Redis client object created",
		zap.String("addr", addr),
		zap.Int("db", config.DB),
	)

	// Retry Ping to verify actual connection (this is where real connection happens)
	currentDelay := retryDelay
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			zap.L().Info("retrying Redis connection (ping)",
				zap.Int("attempt", attempt),
				zap.Int("max_retries", maxRetries),
				zap.Duration("delay", currentDelay),
			)
			time.Sleep(currentDelay)
			// Exponential backoff: double the delay for each retry
			currentDelay *= 2
		}

		pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := client.Ping(pingCtx).Err()
		pingCancel()

		if err == nil {
			zap.L().Info("successfully connected to Redis",
				zap.String("addr", addr),
				zap.Int("db", config.DB),
				zap.Int("total_attempts", attempt+1),
			)
			break
		}

		zap.L().Warn("failed to connect to Redis (ping failed)",
			zap.Int("attempt", attempt+1),
			zap.Int("max_retries", maxRetries),
			zap.String("addr", addr),
			zap.Error(err),
		)

		// If this is the last attempt, cleanup and return error
		if attempt == maxRetries {
			zap.L().Error("failed to connect to Redis after all retries",
				zap.Int("total_attempts", maxRetries+1),
				zap.String("addr", addr),
				zap.Error(err),
			)
			// Cleanup: close client before returning error
			_ = client.Close()
			return nil, err
		}
	}

	return &RedisCacheConnection{Client: client}, nil
}

// Close closes the Redis client connection
func (c *RedisCacheConnection) Close() error {
	if err := c.Client.Close(); err != nil {
		zap.L().Error("failed to close Redis connection", zap.Error(err))
		return err
	}
	zap.L().Info("disconnected from Redis")
	return nil
}

func (c *RedisCacheConnection) GetRedisClient() *redis.Client {
	return c.Client
}
