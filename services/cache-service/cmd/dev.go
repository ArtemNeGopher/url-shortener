//go:build !prod

package main

import "github.com/go-redis/redis/v8"

func InitRedis(addr string, poolSize int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       1, // Обращаемся к тестовым данным
		PoolSize: poolSize,
	})
}
