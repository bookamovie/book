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

type Config struct {
	BookConfig   BookConfig   `yaml:"book"`
	SQLiteConfig SQLiteConfig `yaml:"sqlite"`
	KafkaConfig  KafkaConfig  `yaml:"kafka"`
}

type BookConfig struct {
	Network string `yaml:"network"`
	Address string `yaml:"address"`
}

type SQLiteConfig struct {
	Address string `yaml:"address"`
}

type KafkaConfig struct {
	Addresses []string `yaml:"addresses"`
	Topic     string   `yaml:"topic"`
	Offset    int64    `yaml:"offset"`
	Partition int32    `yaml:"partition"`
}

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
