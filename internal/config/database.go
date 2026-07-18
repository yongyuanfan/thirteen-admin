package config

// Database 数据库配置，通过 driver 选择具体驱动。
type Database struct {
	Driver   string   `yaml:"driver"`
	MySQL    MySQL    `yaml:"mysql"`
	Postgres Postgres `yaml:"postgres"`
}

type MySQL struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Charset  string `yaml:"charset"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}
