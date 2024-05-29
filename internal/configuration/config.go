package configuration

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Network *NetworkConfig `yaml:"network"`
	Logger  *LoggerConfig  `yaml:"logger"`
}

func Load(filename string) (*Config, error) {
	if filename == "" {
		return &Config{}, nil
	}

	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %w", filename, err)
	}

	var config Config
	if err = yaml.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	return &config, nil
}

type NetworkConfig struct {
	Address        string        `yaml:"address"`
	MaxMessageSize int           `yaml:"max_message_size"`
	MaxConnections int           `yaml:"max_connections"`
	IdleTimeout    time.Duration `yaml:"idle_timeout"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}
