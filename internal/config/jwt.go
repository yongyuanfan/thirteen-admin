package config

type JWT struct {
	Secret           string `yaml:"secret"`
	ExpiresIn        int    `yaml:"expires_in"`         // access token 有效期（秒）
	RefreshExpiresIn int    `yaml:"refresh_expires_in"` // refresh token 有效期（秒）
}
