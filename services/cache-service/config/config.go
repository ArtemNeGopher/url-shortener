// Package config описывает структуру конфигурации
package config

import "time"

type Config struct {
	Host      string        `yaml:"host" env:"HOST"`
	Port      int           `yaml:"port" env:"PORT"`
	LocalTTL  time.Duration `yaml:"local_ttl" env:"LOCAL_TTL"`
	RedisAddr string        `yaml:"redis_addr" env:"REDIS_ADDR"`
	RedisTTL  time.Duration `yaml:"redis_ttl" env:"REDIS_TTL"`
}
