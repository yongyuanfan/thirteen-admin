package urlx

import (
	"net/url"
	"strings"
)

// CanonicalStorageURLKey 生成用于去重比较的 storage URL 键。
// 去除 query/fragment，PathUnescape 路径，host 小写；不用于对外展示或写入存储。
func CanonicalStorageURLKey(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		if index := strings.Index(raw, "?"); index >= 0 {
			raw = raw[:index]
		}
		if index := strings.Index(raw, "#"); index >= 0 {
			raw = raw[:index]
		}
		return unescapePathLoop(raw)
	}
	parsed.RawQuery = ""
	parsed.Fragment = ""
	host := strings.ToLower(parsed.Host)
	path := unescapePathLoop(parsed.EscapedPath())
	if path == "" {
		path = "/"
	}
	return parsed.Scheme + "://" + host + path
}

func unescapePathLoop(path string) string {
	for i := 0; i < 3; i++ {
		decoded, err := url.PathUnescape(path)
		if err != nil || decoded == path {
			return path
		}
		path = decoded
	}
	return path
}
