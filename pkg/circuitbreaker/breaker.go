package circuitbreaker

import (
	"errors"
	"time"
)

type State int

const (
	StateClosed State = 0
	StateOpen   State = 1
)

var ErrCircuitBreakerOpen = errors.New("circuit breaker is open")

type Breaker struct {
	state       State
	failures    int
	maxFailures int
	resetTime   time.Duration
	unlockIn    time.Time
}

func New(failures int, resetTime time.Duration) *Breaker {
	return &Breaker{
		state:       StateClosed,
		failures:    0,
		maxFailures: failures,
		resetTime:   resetTime,
	}
}

func (b *Breaker) Call(f func() error) error {
	// Если открыт, то проверяем не вышло ли время
	if b.state == StateOpen {
		if time.Now().Before(b.unlockIn) { // Время не вышло
			return ErrCircuitBreakerOpen
		}

		err := f()
		if err != nil {
			b.state = StateOpen
			return err
		}

		b.state = StateClosed
		b.failures = 0
	}

	err := f()
	if err != nil {
		b.failures++
		if b.failures >= b.maxFailures {
			b.state = StateOpen
			b.unlockIn = time.Now().Add(b.resetTime)
		}
		return err
	}

	// Успешный вызов сбрасывает счётчик
	b.failures = 0

	return nil
}
