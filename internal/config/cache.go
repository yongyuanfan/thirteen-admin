package config

// Cache 缓存配置，通过 driver 选择具体驱动。
type Cache struct {
	Driver string    `yaml:"driver"`
	File   FileCache `yaml:"file"`
	Redis  Redis     `yaml:"redis"`
}

type FileCache struct {
	Path string `yaml:"path"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}
