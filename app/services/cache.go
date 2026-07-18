package services

import (
	"context"
	"strings"

	"thirteen-admin/internal/storage"
)

type cacheDef struct {
	Name   string
	Prefix string
	Remark string
}

var cacheDefs = []cacheDef{
	{Name: "login_tokens", Prefix: "online:token:", Remark: "在线用户令牌"},
	{Name: "sys_config", Prefix: "sys_config:", Remark: "系统配置缓存"},
	{Name: "captcha", Prefix: "captcha:", Remark: "验证码缓存"},
	{Name: "login_lock", Prefix: "login:lock:", Remark: "登录失败锁定"},
}

// CacheService 缓存监控。
type CacheService struct {
	cache storage.Cache
}

// NewCacheService 创建缓存服务。
func NewCacheService(cache storage.Cache) *CacheService {
	return &CacheService{cache: cache}
}

func (s *CacheService) defByName(name string) (cacheDef, bool) {
	for _, d := range cacheDefs {
		if d.Name == name {
			return d, true
		}
	}
	return cacheDef{}, false
}

// Overview 缓存监控概览。
func (s *CacheService) Overview(ctx context.Context) (map[string]any, error) {
	if s.cache == nil {
		return map[string]any{
			"info":         map[string]string{},
			"dbSize":       0,
			"commandStats": []map[string]string{},
		}, nil
	}
	stats, err := s.cache.Stats(ctx)
	if err != nil {
		return nil, err
	}
	commandStats := stats.CommandStats
	if commandStats == nil {
		commandStats = []map[string]string{}
	}
	return map[string]any{
		"info":         stats.Info,
		"dbSize":       stats.DBSize,
		"commandStats": commandStats,
	}, nil
}

// GetNames 缓存名称列表。
func (s *CacheService) GetNames() []map[string]any {
	out := make([]map[string]any, 0, len(cacheDefs))
	for _, d := range cacheDefs {
		out = append(out, map[string]any{
			"cacheName": d.Name,
			"remark":    d.Remark,
		})
	}
	return out
}

// GetKeys 某缓存名下的 key 列表。
func (s *CacheService) GetKeys(ctx context.Context, cacheName string) ([]string, error) {
	d, ok := s.defByName(cacheName)
	if !ok {
		return []string{}, nil
	}
	if s.cache == nil {
		return []string{}, nil
	}
	var cursor uint64
	keys := make([]string, 0)
	for {
		batch, next, err := s.cache.Scan(ctx, cursor, d.Prefix+"*", 100)
		if err != nil {
			return nil, err
		}
		for _, k := range batch {
			keys = append(keys, strings.TrimPrefix(k, d.Prefix))
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return keys, nil
}

// GetValue 查询缓存值。
func (s *CacheService) GetValue(ctx context.Context, cacheName, cacheKey string) (map[string]any, error) {
	d, ok := s.defByName(cacheName)
	if !ok {
		return map[string]any{
			"cacheName":  cacheName,
			"cacheKey":   cacheKey,
			"cacheValue": "",
			"remark":     "",
		}, nil
	}
	fullKey := d.Prefix + cacheKey
	val := ""
	if s.cache != nil {
		v, err := s.cache.Get(ctx, fullKey)
		if err == nil {
			val = v
		}
	}
	return map[string]any{
		"cacheName":  cacheName,
		"cacheKey":   cacheKey,
		"cacheValue": val,
		"remark":     d.Remark,
	}, nil
}

// ClearCacheName 按缓存名清空。
func (s *CacheService) ClearCacheName(ctx context.Context, cacheName string) error {
	d, ok := s.defByName(cacheName)
	if !ok || s.cache == nil {
		return nil
	}
	return s.deleteByPrefix(ctx, d.Prefix)
}

// ClearCacheKey 按 key 删除（支持完整 key 或带前缀）。
func (s *CacheService) ClearCacheKey(ctx context.Context, cacheKey string) error {
	if s.cache == nil || cacheKey == "" {
		return nil
	}
	return s.cache.Del(ctx, cacheKey)
}

// ClearCacheAll 清空全部缓存。
func (s *CacheService) ClearCacheAll(ctx context.Context) error {
	if s.cache == nil {
		return nil
	}
	for _, d := range cacheDefs {
		if err := s.deleteByPrefix(ctx, d.Prefix); err != nil {
			return err
		}
	}
	return nil
}

func (s *CacheService) deleteByPrefix(ctx context.Context, prefix string) error {
	var cursor uint64
	for {
		keys, next, err := s.cache.Scan(ctx, cursor, prefix+"*", 100)
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := s.cache.Del(ctx, keys...); err != nil {
				return err
			}
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return nil
}
