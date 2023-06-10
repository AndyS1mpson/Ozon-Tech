// Project Configuration
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const pathToConfig = "config.yaml"

// Service configuraton
type Config struct {
	Token    string `yaml:"token"`
	Services struct {
		Loms     string `yaml:"loms"`
		Products string `yaml:"products"`
	} `yaml:"services"`
	Postgres struct {
		ConnectionString string `yaml:"connection_string"`
	} `yaml:"postgres"`
}

// Create a new instance of the config
func New() (*Config, error) {

	cfg := &Config{}

	rawYaml, err := os.ReadFile(pathToConfig)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	err = yaml.Unmarshal(rawYaml, &cfg)
	if err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	return cfg, nil
}
