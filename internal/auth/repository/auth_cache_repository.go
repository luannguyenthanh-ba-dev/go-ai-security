package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/luannguyenthanh-ba-dev/go-ai-security/infrastructure/cache"
)

// Auth repository interface
var (
	AccessTokenWhiteListKey = "auth:access_token:whitelist:%s"
)

// AuthCacheRepository defines the interface for auth cache operations
type AuthCacheRepository interface {
	AddAccessTokenToWhiteList(ctx context.Context, uID string, accessToken string, ttl time.Duration) (bool, error)
}

type authCacheRepository struct {
	cacheClient cache.CacheClient
}

// NewAuthCacheRepository creates a new auth cache repository implementation
func NewAuthCacheRepository(cacheClient cache.CacheClient) AuthCacheRepository {
	return &authCacheRepository{cacheClient: cacheClient}
}

func (r *authCacheRepository) AddAccessTokenToWhiteList(ctx context.Context, uID string, accessToken string, ttl time.Duration) (bool, error) {
	key := fmt.Sprintf(AccessTokenWhiteListKey, uID)
	return r.cacheClient.Set(ctx, key, accessToken, ttl)
}
