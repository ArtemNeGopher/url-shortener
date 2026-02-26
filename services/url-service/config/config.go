package config

type Config struct {
	Host             string `yaml:"host" env:"HOST"`
	Port             int    `yaml:"port" env:"PORT"`
	URLCreateRetries uint   `yaml:"url_create_retries" env:"URL_CREATE_RETRIES"`
	Database         DatabaseConfig
}

type DatabaseConfig struct {
	DatabaseURL        string `yaml:"database_url" env:"DATABASE_URL"`
	MaxOpenConnections int    `yaml:"max_open_connections" env:"MAX_OPEN_CONNECTIONS"`
	MaxIdleConnections int    `yaml:"max_idle_connections" env:"MAX_IDLE_CONNECTIONS"`
}
