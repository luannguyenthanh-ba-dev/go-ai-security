package cache

import "errors"

var (
	ErrSetCacheValueNil = errors.New("cache value is nil")
	ErrCacheKeyEmpty    = errors.New("cache key is empty")
)
