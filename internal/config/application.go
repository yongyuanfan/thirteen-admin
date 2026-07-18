package config

type Application struct {
	Mode     string   `yaml:"mode"`
	Name     string   `yaml:"name"`
	HTTP     HTTP     `yaml:"http"`
	JWT      JWT      `yaml:"jwt"`
	Logging  Logging  `yaml:"logging"`
	Database Database `yaml:"database"`
	Cache    Cache    `yaml:"cache"`
}
