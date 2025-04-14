package utils

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const cpEnvName = "CONFIG_PATH"

var (
	ErrConfigPathNotSpecified = fmt.Errorf("%s env variable must be specified", cpEnvName)
	ErrConfigNotFound         = fmt.Errorf("config not found")
)

// Config{} defines the structure of the entire application configuration.
//
// It is typically populated from a YAML file.
type Config struct {
	BookConfig   BookConfig   `yaml:"book"`
	SQLiteConfig SQLiteConfig `yaml:"sqlite"`
	KafkaConfig  KafkaConfig  `yaml:"kafka"`
}

// BookConfig{} contains network settings for the gRPC book service.
type BookConfig struct {
	Network string `yaml:"network"`
	Address string `yaml:"address"`
}

// SQLiteConfig{} holds database configuration for SQLite.
type SQLiteConfig struct {
	Address string `yaml:"address"`
}

// KafkaConfig{} includes all configuration needed for Kafka producers.
type KafkaConfig struct {
	Addresses []string `yaml:"addresses"`
	Topic     string   `yaml:"topic"`
	Offset    int64    `yaml:"offset"`
	Partition int32    `yaml:"partition"`
}

// LoadConfig() loads and validates configuration from a YAML file specified by the CONFIG_PATH environment variable. Only known paths are accepted.
func LoadConfig() (Config, error) {
	configPath := os.Getenv(cpEnvName)
	if configPath == "" {
		return Config{}, ErrConfigPathNotSpecified
	}
	defer os.Unsetenv(cpEnvName)

	var cfg Config

	switch configPath {
	case "config/local.yaml":
	case "config/dev.yaml":
	case "config/prod.yaml":
	case "config/custom.yaml":
	default:
		return Config{}, ErrConfigNotFound
	}

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return Config{}, ErrConfigNotFound
	}

	return cfg, nil
}
