package services

import (
	"context"
	"strings"

	"github.com/redis/go-redis/v9"
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

// CacheService Redis 缓存监控。
type CacheService struct {
	rdb *redis.Client
}

// NewCacheService 创建缓存服务。
func NewCacheService(rdb *redis.Client) *CacheService {
	return &CacheService{rdb: rdb}
}

func (s *CacheService) defByName(name string) (cacheDef, bool) {
	for _, d := range cacheDefs {
		if d.Name == name {
			return d, true
		}
	}
	return cacheDef{}, false
}

// Overview Redis 监控概览。
func (s *CacheService) Overview(ctx context.Context) (map[string]any, error) {
	if s.rdb == nil {
		return map[string]any{
			"info":         map[string]string{},
			"dbSize":       0,
			"commandStats": []map[string]string{},
		}, nil
	}
	infoStr, err := s.rdb.Info(ctx).Result()
	if err != nil {
		return nil, err
	}
	info := parseRedisInfo(infoStr)
	dbSize, err := s.rdb.DBSize(ctx).Result()
	if err != nil {
		return nil, err
	}
	commandStats := parseCommandStats(infoStr)
	return map[string]any{
		"info":         info,
		"dbSize":       dbSize,
		"commandStats": commandStats,
	}, nil
}

func parseRedisInfo(info string) map[string]string {
	out := map[string]string{}
	for _, line := range strings.Split(info, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		switch key {
		case "redis_version", "redis_mode", "tcp_port", "connected_clients",
			"uptime_in_days", "used_memory_human", "used_cpu_user_children",
			"maxmemory_human", "aof_enabled", "rdb_last_bgsave_status",
			"instantaneous_input_kbps", "instantaneous_output_kbps":
			out[key] = val
		}
	}
	return out
}

func parseCommandStats(info string) []map[string]string {
	out := make([]map[string]string, 0)
	inCmdstat := false
	for _, line := range strings.Split(info, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# Commandstats") {
			inCmdstat = true
			continue
		}
		if inCmdstat {
			if strings.HasPrefix(line, "#") {
				break
			}
			if !strings.HasPrefix(line, "cmdstat_") {
				continue
			}
			// cmdstat_get:calls=100,usec=50,...
			rest := strings.TrimPrefix(line, "cmdstat_")
			name, _, _ := strings.Cut(rest, ":")
			calls := ""
			if idx := strings.Index(rest, "calls="); idx >= 0 {
				seg := rest[idx+6:]
				if end := strings.Index(seg, ","); end >= 0 {
					calls = seg[:end]
				} else {
					calls = seg
				}
			}
			if name != "" {
				out = append(out, map[string]string{"name": name, "value": calls})
			}
		}
	}
	return out
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
	if s.rdb == nil {
		return []string{}, nil
	}
	var cursor uint64
	keys := make([]string, 0)
	for {
		batch, next, err := s.rdb.Scan(ctx, cursor, d.Prefix+"*", 100).Result()
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
	if s.rdb != nil {
		v, err := s.rdb.Get(ctx, fullKey).Result()
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
	if !ok || s.rdb == nil {
		return nil
	}
	return s.deleteByPrefix(ctx, d.Prefix)
}

// ClearCacheKey 按 key 删除（支持完整 key 或带前缀）。
func (s *CacheService) ClearCacheKey(ctx context.Context, cacheKey string) error {
	if s.rdb == nil || cacheKey == "" {
		return nil
	}
	return s.rdb.Del(ctx, cacheKey).Err()
}

// ClearCacheAll 清空全部缓存。
func (s *CacheService) ClearCacheAll(ctx context.Context) error {
	if s.rdb == nil {
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
		keys, next, err := s.rdb.Scan(ctx, cursor, prefix+"*", 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := s.rdb.Del(ctx, keys...).Err(); err != nil {
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
