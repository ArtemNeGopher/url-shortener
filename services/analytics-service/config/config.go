package config

type Config struct {
	Host     string `yaml:"host" env:"HOST"`
	Port     int    `yaml:"port" env:"PORT"`
	Database DatabaseConfig
	Worker   WorkerConfig
}

type WorkerConfig struct {
	BatchSize   int `yaml:"batch_size" env:"BATCH_SIZE"`
	WorkerCount int `yaml:"count" env:"WORKER_COUNT"`
}

type DatabaseConfig struct {
	DatabaseURL        string `yaml:"database_url" env:"DATABASE_URL"`
	MaxOpenConnections int    `yaml:"max_open_connections" env:"MAX_OPEN_CONNECTIONS"`
	MaxIdleConnections int    `yaml:"max_idle_connections" env:"MAX_IDLE_CONNECTIONS"`
}
