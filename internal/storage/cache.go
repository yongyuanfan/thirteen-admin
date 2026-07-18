package storage

import (
	"context"
	"errors"
	"time"
)

// ErrCacheMiss 表示缓存键不存在或已过期。
var ErrCacheMiss = errors.New("cache miss")

var ErrCacheNotInitialized = errors.New("cache 未初始化")

// CacheStats 缓存监控统计。
type CacheStats struct {
	Info         map[string]string
	DBSize       int64
	CommandStats []map[string]string
}

// Cache 统一缓存接口。
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, key string) (bool, error)
	Scan(ctx context.Context, cursor uint64, match string, count int64) (keys []string, next uint64, err error)
	Stats(ctx context.Context) (CacheStats, error)
}

var cacheClient Cache

// GetCache 返回全局缓存实例。
func GetCache() Cache {
	return cacheClient
}

// RequireCache 返回全局缓存实例，未初始化时返回错误。
func RequireCache() (Cache, error) {
	if cacheClient == nil {
		return nil, ErrCacheNotInitialized
	}
	return cacheClient, nil
}

// IsCacheMiss 判断是否为缓存未命中。
func IsCacheMiss(err error) bool {
	return errors.Is(err, ErrCacheMiss)
}
