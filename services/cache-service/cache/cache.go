package cache

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	// Максимально время жизни InMemory кэша
	localTTL time.Duration = 5 * time.Minute
)

type item struct {
	Data      string
	ExpiresAt time.Time // Когда запись умрёт
}

func newItem(data string, ttl time.Duration) *item {
	return &item{
		Data:      data,
		ExpiresAt: time.Now().Add(ttl),
	}
}

type Cache struct {
	client   *redis.Client
	localMap map[string]*item
	localMu  sync.RWMutex
	stopCh   chan struct{}
}

func New(client *redis.Client) *Cache {
	cache := &Cache{
		client:   client,
		localMap: make(map[string]*item),
		localMu:  sync.RWMutex{},
		stopCh:   make(chan struct{}),
	}

	// Запуск фоновой очистки каждую минуту
	go func() {
		for {
			select {
			case <-time.After(1 * time.Minute):
				cache.Clean()
			case <-cache.stopCh:
				return // Остановка фоновой задачи
			}
		}
	}()

	return cache
}

func (c *Cache) Clean() {
	c.localMu.Lock()
	// Фиксируем текущее время
	now := time.Now()
	for key, value := range c.localMap {
		// Удаляем, если время вышло
		if value.ExpiresAt.Before(now) {
			delete(c.localMap, key)
		}
	}
	c.localMu.Unlock()
}

func (c *Cache) Close() {
	c.client.Close()
	<-c.stopCh
}

func (c *Cache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	// Устанавливаем в Redis
	err := c.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return err
	}

	// Высчитываем сколько времени кэш будет жить локально
	localTTL := localTTL
	if ttl < localTTL {
		localTTL = ttl
	}

	// Утанавливаем локально
	c.localMu.Lock()
	c.localMap[key] = newItem(value, ttl)
	c.localMu.Unlock()

	return nil
}

func (c *Cache) Get(ctx context.Context, key string) (string, bool, error) {
	// Проверяем в локальном кэше
	c.localMu.RLock()
	value, exists := c.localMap[key]
	if exists && value.ExpiresAt.Before(time.Now()) {
		// Время жизни вышло, локально не найдено
		// Айтем удалит фоновая очистка
		exists = false
	}
	c.localMu.RUnlock()

	// Нашли в локальном кэше
	if exists {
		return value.Data, true, nil
	}

	// Не нашли в локальном кэше
	// Идём в редис
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}

	// Сохраняем локально
	c.localMu.Lock()
	c.localMap[key] = newItem(val, localTTL)
	c.localMu.Unlock()

	return val, true, nil
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	// Удаляем из локального кэша
	c.localMu.Lock()
	delete(c.localMap, key)
	c.localMu.Unlock()

	// Удаляем из редис
	err := c.client.Del(ctx, key).Err()
	return err
}
