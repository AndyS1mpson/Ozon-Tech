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
	Token    string `yaml:"token"`
	Services struct {
		Loms     string `yaml:"loms"`
		Products string `yaml:"products"`
	} `yaml:"services"`
	Postgres struct {
		ConnectionString       string `yaml:"connection_string"`
		TestDBConnectionString string `yaml:"test_db_connection_string"`
	} `yaml:"postgres"`
	Jaeger struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"jaeger"`
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
