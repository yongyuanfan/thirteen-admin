package config

type JWT struct {
	Secret           string `yaml:"secret"`
	ExpiresIn        int    `yaml:"expires_in"`         // access token 有效期（秒）
	RefreshExpiresIn int    `yaml:"refresh_expires_in"` // refresh token 窗口（秒）
	AutoRefreshToken *bool  `yaml:"auto_refresh_token"` // 滑动续期，nil 默认 true
	AutoRevokeToken  *bool  `yaml:"auto_revoke_token"`  // 退出/强退吊销，nil 默认 true
}
