package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const onlineKeyPref = "online:token:"

// OnlineUser 在线用户信息。
type OnlineUser struct {
	TokenID       string `json:"tokenId"`
	DeptName      string `json:"deptName"`
	UserName      string `json:"userName"`
	IPAddr        string `json:"ipaddr"`
	LoginLocation string `json:"loginLocation"`
	Browser       string `json:"browser"`
	OS            string `json:"os"`
	LoginTime     int64  `json:"loginTime"`
}

// OnlineService 在线用户（Redis）。
type OnlineService struct {
	rdb *redis.Client
	ttl time.Duration
}

// NewOnlineService 创建在线用户服务。
func NewOnlineService(rdb *redis.Client, ttl time.Duration) *OnlineService {
	if ttl <= 0 {
		ttl = 2 * time.Hour
	}
	return &OnlineService{rdb: rdb, ttl: ttl}
}

func (s *OnlineService) key(token string) string {
	return onlineKeyPref + token
}

// Put 写入在线会话。
func (s *OnlineService) Put(ctx context.Context, u OnlineUser) error {
	if s.rdb == nil {
		return nil
	}
	b, err := json.Marshal(u)
	if err != nil {
		return err
	}
	return s.rdb.Set(ctx, s.key(u.TokenID), b, s.ttl).Err()
}

// Remove 删除会话。
func (s *OnlineService) Remove(ctx context.Context, tokenID string) error {
	if s.rdb == nil {
		return nil
	}
	return s.rdb.Del(ctx, s.key(tokenID)).Err()
}

// List 列出在线用户。
func (s *OnlineService) List(ctx context.Context, ipaddr, userName string) ([]OnlineUser, error) {
	if s.rdb == nil {
		return []OnlineUser{}, nil
	}
	var cursor uint64
	out := make([]OnlineUser, 0)
	for {
		keys, next, err := s.rdb.Scan(ctx, cursor, onlineKeyPref+"*", 100).Result()
		if err != nil {
			return nil, err
		}
		for _, k := range keys {
			val, err := s.rdb.Get(ctx, k).Bytes()
			if err != nil {
				continue
			}
			var u OnlineUser
			if err := json.Unmarshal(val, &u); err != nil {
				continue
			}
			if ipaddr != "" && u.IPAddr != ipaddr {
				continue
			}
			if userName != "" && u.UserName != userName {
				continue
			}
			out = append(out, u)
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return out, nil
}

// TokenExists 会话是否存在（用于强退后校验）。
func (s *OnlineService) TokenExists(ctx context.Context, token string) (bool, error) {
	if s.rdb == nil {
		return true, nil
	}
	n, err := s.rdb.Exists(ctx, s.key(token)).Result()
	if err != nil {
		return false, fmt.Errorf("redis exists: %w", err)
	}
	return n > 0, nil
}
