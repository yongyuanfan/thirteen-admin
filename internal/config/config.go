package config

import (
	"os"

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

	return nil
}

// GetApplication 获取应用配置
func GetApplication() *Application {
	return &config.Application
}
