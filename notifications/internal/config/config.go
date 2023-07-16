// Project Configuration
package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const pathToConfig = "config.yaml"

// Service configuraton
type Config struct {
	Postgres struct {
		ConnectionString string `yaml:"connection_string"`
	} `yaml:"postgres"`
	Brokers  []string `yaml:"brokers"`
	Telegram struct {
		APIKey string `yaml:"api_key"`
		ChatID int64  `yaml:"chat_id"`
	} `yaml:"telegram"`
}

// Create a new instance of the config
func New() (*Config, error) {
	cfg := &Config{}

	rawYaml, err := os.ReadFile(pathToConfig)
	if err != nil {
		return nil, errors.Wrap(err, "read config file")
	}

	err = yaml.Unmarshal(rawYaml, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "parse config file")
	}

	return cfg, nil
}
