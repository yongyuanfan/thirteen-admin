package storage

import (
	"errors"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

var ErrRedisNotInitialized = errors.New("redis 未初始化")

// GetRedis 返回全局 Redis 客户端实例。
func GetRedis() *redis.Client {
	return redisClient
}

// RequireRedis 返回全局 Redis 客户端实例，未初始化时返回错误。
func RequireRedis() (*redis.Client, error) {
	if redisClient == nil {
		return nil, ErrRedisNotInitialized
	}
	return redisClient, nil
}
