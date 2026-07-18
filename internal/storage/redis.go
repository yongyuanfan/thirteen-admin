package storage

import (
	"context"
	"errors"
	"fmt"

	"clash-admin/internal/config"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

var ErrRedisNotInitialized = errors.New("redis 未初始化")

// InitRedis 初始化 Redis 客户端连接。
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
	return nil
}

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
