package config

type Config struct {
	Host        string `yaml:"host" env:"HOST"`
	Port        int    `yaml:"port" env:"PORT"`
	DatabaseURL string `yaml:"database_url" env:"DATABASE_URL"`
	Worker      WorkerConfig
}

type WorkerConfig struct {
	BatchSize   int `yaml:"batch_size" env:"BATCH_SIZE"`
	WorkerCount int `yaml:"count" env:"WORKER_COUNT"`
}
