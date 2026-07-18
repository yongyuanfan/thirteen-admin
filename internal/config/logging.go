package config

type Logging struct {
	// Enabled 表示是否启用 JSON 文件日志。
	Enabled bool `yaml:"enabled"`
	// Level 日志级别：debug | info | warn | error，默认 info。
	Level string `yaml:"level"`
	// FilePath JSON 日志文件路径，默认 runtime/logs/thirteen-admin.log。
	FilePath string `yaml:"file_path"`
}
