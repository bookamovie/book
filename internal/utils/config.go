package utils

type Config struct {
	BookaMovieConfig *BookaMovieConfig
	KafkaConfig      *KafkaConfig
}

type BookaMovieConfig struct {
	Network string `yaml:"network"`
	Address string `yaml:"address"`
}

type KafkaConfig struct{}

func LoadConfig() *Config {
	var cfg *Config

	return cfg
}
