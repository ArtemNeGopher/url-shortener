package ratelimit

import (
	"sync"
	"time"
)

type entry struct {
	refillIn time.Time // Время когда нужно сбросить rate
	rate     int       // Количество запросов
}

type Limiter struct {
	limit   int
	period  time.Duration
	entries map[string]*entry
	mu      sync.Mutex
}

func New(limit int, period time.Duration) *Limiter {
	return &Limiter{
		limit:   limit,
		period:  period,
		entries: make(map[string]*entry),
		mu:      sync.Mutex{},
	}
}

func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	en, ok := l.entries[key]
	if !ok {
		en = &entry{
			refillIn: now.Add(l.period),
			rate:     1,
		}
		l.entries[key] = en
	}

	// Обнуление rate, если время вышло
	if now.After(en.refillIn) {
		en.refillIn = now.Add(l.period)
		en.rate = 0
	}

	if en.rate <= l.limit {
		en.rate++
		return true
	}

	return false
}
