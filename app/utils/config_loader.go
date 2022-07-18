package utils

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Database *DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
	Driver string `yaml:"driver"`
	Dsn    string `yaml:"dsn"`
}

func NewConfig(path string) (config *Config, err error) {
	dat, err := os.ReadFile(path)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(dat, config)
	if err != nil {
		return
	}

	return config, nil
}
