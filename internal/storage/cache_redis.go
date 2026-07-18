package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"thirteen-admin/internal/config"

	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	client *redis.Client
}

// InitRedis 初始化 Redis 缓存客户端。
func InitRedis(redisConfig config.Redis) error {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return fmt.Errorf("connect redis: %w", err)
	}

	redisClient = client
	cacheClient = &redisCache{client: client}
	return nil
}

func (c *redisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", ErrCacheMiss
	}
	return val, err
}

func (c *redisCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *redisCache) Del(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return c.client.Del(ctx, keys...).Err()
}

func (c *redisCache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (c *redisCache) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return c.client.Scan(ctx, cursor, match, count).Result()
}

func (c *redisCache) Stats(ctx context.Context) (CacheStats, error) {
	infoStr, err := c.client.Info(ctx).Result()
	if err != nil {
		return CacheStats{}, err
	}
	dbSize, err := c.client.DBSize(ctx).Result()
	if err != nil {
		return CacheStats{}, err
	}
	return CacheStats{
		Info:         parseRedisInfo(infoStr),
		DBSize:       dbSize,
		CommandStats: parseCommandStats(infoStr),
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
