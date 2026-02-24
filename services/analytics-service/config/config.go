package config

type Config struct {
	Host        string `yaml:"host" env:"HOST"`
	Port        int    `yaml:"port" env:"PORT"`
	DatabaseURL string `yaml:"database_url" env:"DATABASE_URL"`
}
