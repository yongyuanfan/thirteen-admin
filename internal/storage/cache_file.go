package storage

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"thirteen-admin/internal/config"
	"thirteen-admin/internal/utils/pathx"
)

type fileCacheEntry struct {
	Value    string `json:"value"`
	ExpireAt int64  `json:"expire_at"` // unix 秒，0 表示不过期
}

type fileCache struct {
	root string
	mu   sync.Mutex
}

// InitFileCache 初始化文件缓存。
func InitFileCache(cfg config.FileCache) error {
	path := strings.TrimSpace(cfg.Path)
	if path == "" {
		path = "runtime/cache"
	}
	if !filepath.IsAbs(path) {
		path = filepath.Join(pathx.ResolveProjectRoot(), path)
	}
	if err := os.MkdirAll(path, 0o755); err != nil {
		return fmt.Errorf("create file cache dir: %w", err)
	}
	cacheClient = &fileCache{root: path}
	return nil
}

func (c *fileCache) filePath(key string) string {
	enc := base64.RawURLEncoding.EncodeToString([]byte(key))
	return filepath.Join(c.root, enc+".json")
}

func (c *fileCache) Get(_ context.Context, key string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, err := c.readEntry(key)
	if err != nil {
		return "", err
	}
	if c.expired(entry) {
		_ = os.Remove(c.filePath(key))
		return "", ErrCacheMiss
	}
	return entry.Value, nil
}

func (c *fileCache) Set(_ context.Context, key, value string, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry := fileCacheEntry{Value: value}
	if ttl > 0 {
		entry.ExpireAt = time.Now().Add(ttl).Unix()
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	return os.WriteFile(c.filePath(key), data, 0o644)
}

func (c *fileCache) Del(_ context.Context, keys ...string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, key := range keys {
		_ = os.Remove(c.filePath(key))
	}
	return nil
}

func (c *fileCache) Exists(_ context.Context, key string) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, err := c.readEntry(key)
	if err != nil {
		if IsCacheMiss(err) {
			return false, nil
		}
		return false, err
	}
	if c.expired(entry) {
		_ = os.Remove(c.filePath(key))
		return false, nil
	}
	return true, nil
}

func (c *fileCache) Scan(_ context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entries, err := os.ReadDir(c.root)
	if err != nil {
		return nil, 0, err
	}
	if count <= 0 {
		count = 100
	}

	all := make([]string, 0)
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		raw, err := base64.RawURLEncoding.DecodeString(strings.TrimSuffix(e.Name(), ".json"))
		if err != nil {
			continue
		}
		key := string(raw)
		entry, err := c.readEntryUnlocked(key)
		if err != nil {
			continue
		}
		if c.expired(entry) {
			_ = os.Remove(c.filePath(key))
			continue
		}
		if match == "" || matchGlob(match, key) {
			all = append(all, key)
		}
	}

	start := int(cursor)
	if start < 0 {
		start = 0
	}
	if start >= len(all) {
		return []string{}, 0, nil
	}
	end := start + int(count)
	if end >= len(all) {
		return all[start:], 0, nil
	}
	return all[start:end], uint64(end), nil
}

func (c *fileCache) Stats(_ context.Context) (CacheStats, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entries, err := os.ReadDir(c.root)
	if err != nil {
		return CacheStats{}, err
	}
	var size int64
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		raw, err := base64.RawURLEncoding.DecodeString(strings.TrimSuffix(e.Name(), ".json"))
		if err != nil {
			continue
		}
		entry, err := c.readEntryUnlocked(string(raw))
		if err != nil {
			continue
		}
		if c.expired(entry) {
			_ = os.Remove(filepath.Join(c.root, e.Name()))
			continue
		}
		size++
	}
	return CacheStats{
		Info: map[string]string{
			"redis_version":          "file",
			"redis_mode":             "file",
			"used_memory_human":      "-",
			"connected_clients":      "1",
			"tcp_port":               "-",
			"uptime_in_days":         "0",
			"maxmemory_human":        "-",
			"aof_enabled":            "0",
			"rdb_last_bgsave_status": "ok",
		},
		DBSize:       size,
		CommandStats: []map[string]string{},
	}, nil
}

func (c *fileCache) readEntry(key string) (fileCacheEntry, error) {
	return c.readEntryUnlocked(key)
}

func (c *fileCache) readEntryUnlocked(key string) (fileCacheEntry, error) {
	data, err := os.ReadFile(c.filePath(key))
	if err != nil {
		if os.IsNotExist(err) {
			return fileCacheEntry{}, ErrCacheMiss
		}
		return fileCacheEntry{}, err
	}
	var entry fileCacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return fileCacheEntry{}, err
	}
	return entry, nil
}

func (c *fileCache) expired(entry fileCacheEntry) bool {
	return entry.ExpireAt > 0 && time.Now().Unix() >= entry.ExpireAt
}

// matchGlob 支持 Redis SCAN 风格的 * 通配。
func matchGlob(pattern, s string) bool {
	if pattern == "*" || pattern == "" {
		return true
	}
	parts := strings.Split(pattern, "*")
	if len(parts) == 1 {
		return pattern == s
	}
	if !strings.HasPrefix(s, parts[0]) {
		return false
	}
	rest := s[len(parts[0]):]
	for i := 1; i < len(parts)-1; i++ {
		idx := strings.Index(rest, parts[i])
		if idx < 0 {
			return false
		}
		rest = rest[idx+len(parts[i]):]
	}
	return strings.HasSuffix(rest, parts[len(parts)-1])
}
