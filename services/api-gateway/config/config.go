// Package config
// Содержит описание конфига сервиса
package config

type Config struct {
	Host            string                `yaml:"host"`
	Port            string                `yaml:"port"`
	URLClient       URLClientConfig       `yaml:"url_client"`
	CacheClient     CacheClientConfig     `yaml:"cache_client"`
	AnalyticsClient AnalyticsClientConfig `yaml:"analytics_client"`
}

type URLClientConfig struct {
	Addr string `yaml:"addr"`
}

type CacheClientConfig struct {
	Addr string `yaml:"addr"`
}

type AnalyticsClientConfig struct {
	Addr string `yaml:"addr"`
}
