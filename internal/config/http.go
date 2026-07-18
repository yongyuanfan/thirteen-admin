package config

type HTTP struct {
	// Enabled 表示是否启用 http 服务。
	Enabled bool   `yaml:"enabled"`
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
}
