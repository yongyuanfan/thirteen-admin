package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"thirteen-admin/internal/storage"
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

// OnlineService 在线用户。
type OnlineService struct {
	cache storage.Cache
	ttl   time.Duration
}

// NewOnlineService 创建在线用户服务。
func NewOnlineService(cache storage.Cache, ttl time.Duration) *OnlineService {
	if ttl <= 0 {
		ttl = 7 * 24 * time.Hour
	}
	return &OnlineService{cache: cache, ttl: ttl}
}

func (s *OnlineService) key(tokenID string) string {
	return onlineKeyPref + tokenID
}

// Put 写入在线会话。可选 ttl 覆盖默认（绝对过期剩余时间）。
func (s *OnlineService) Put(ctx context.Context, u OnlineUser, ttl ...time.Duration) error {
	if s.cache == nil {
		return nil
	}
	b, err := json.Marshal(u)
	if err != nil {
		return err
	}
	d := s.ttl
	if len(ttl) > 0 && ttl[0] > 0 {
		d = ttl[0]
	}
	return s.cache.Set(ctx, s.key(u.TokenID), string(b), d)
}

// Get 读取在线会话。
func (s *OnlineService) Get(ctx context.Context, tokenID string) (OnlineUser, error) {
	if s.cache == nil {
		return OnlineUser{}, storage.ErrCacheMiss
	}
	val, err := s.cache.Get(ctx, s.key(tokenID))
	if err != nil {
		return OnlineUser{}, err
	}
	var u OnlineUser
	if err := json.Unmarshal([]byte(val), &u); err != nil {
		return OnlineUser{}, err
	}
	return u, nil
}

// Remove 删除会话。
func (s *OnlineService) Remove(ctx context.Context, tokenID string) error {
	if s.cache == nil {
		return nil
	}
	return s.cache.Del(ctx, s.key(tokenID))
}

// List 列出在线用户。
func (s *OnlineService) List(ctx context.Context, ipaddr, userName string) ([]OnlineUser, error) {
	if s.cache == nil {
		return []OnlineUser{}, nil
	}
	var cursor uint64
	out := make([]OnlineUser, 0)
	for {
		keys, next, err := s.cache.Scan(ctx, cursor, onlineKeyPref+"*", 100)
		if err != nil {
			return nil, err
		}
		for _, k := range keys {
			val, err := s.cache.Get(ctx, k)
			if err != nil {
				continue
			}
			var u OnlineUser
			if err := json.Unmarshal([]byte(val), &u); err != nil {
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

// TokenExists 会话是否存在（用于强退/退出后校验）。
func (s *OnlineService) TokenExists(ctx context.Context, tokenID string) (bool, error) {
	if s.cache == nil {
		return true, nil
	}
	ok, err := s.cache.Exists(ctx, s.key(tokenID))
	if err != nil {
		return false, fmt.Errorf("cache exists: %w", err)
	}
	return ok, nil
}
