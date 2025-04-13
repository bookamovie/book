package utils

type Config struct {
	LogMode string

	BookaMovieConfig *BookaMovieConfig `yaml:"bookamovie"`
	SQLiteConfig     *SQLiteConfig     `yaml:"sqlite"`
	KafkaConfig      *KafkaConfig      `yaml:"kafka"`
}

type BookaMovieConfig struct {
	Network string `yaml:"network"`
	Address string `yaml:"address"`
}

type SQLiteConfig struct {
	Address string `yaml:"address"`
}

type KafkaConfig struct {
	Addresses []string `yaml:"addresses"`
}

func LoadConfig() *Config {
	var cfg *Config

	return cfg
}
