package config

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	config *Config
)

// Config 全局配置
type Config struct {
	Application Application `yaml:"application"`
}

// Load 加载配置
func Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(data, &config); err != nil {
		return err
	}

	applyDefaults(&config.Application)
	return nil
}

func applyDefaults(app *Application) {
	if strings.TrimSpace(app.Database.Driver) == "" {
		app.Database.Driver = "mysql"
	}
	if strings.TrimSpace(app.Cache.Driver) == "" {
		app.Cache.Driver = "file"
	}
	if strings.TrimSpace(app.Cache.File.Path) == "" {
		app.Cache.File.Path = "runtime/cache"
	}
	if strings.TrimSpace(app.Database.MySQL.Charset) == "" {
		app.Database.MySQL.Charset = "utf8mb4"
	}
}

// GetApplication 获取应用配置
func GetApplication() *Application {
	return &config.Application
}
