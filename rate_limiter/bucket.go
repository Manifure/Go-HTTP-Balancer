package rate_limiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	mu        sync.Mutex
	capacity  int
	tokens    float64
	rate      float64 // токенов в секунду
	lastCheck time.Time
}

func NewTokenBucket(capacity int, rate float64) *TokenBucket {
	return &TokenBucket{
		capacity:  capacity,
		tokens:    float64(capacity),
		rate:      rate,
		lastCheck: time.Now(),
	}
}

// Allow проверяет, можно ли пропустить запрос.
// Возвращает true, если в bucket есть токен, иначе false.
func (b *TokenBucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(b.lastCheck).Seconds()
	b.lastCheck = now

	// Пополнение токенов
	b.tokens += elapsed * b.rate
	if b.tokens > float64(b.capacity) {
		b.tokens = float64(b.capacity)
	}

	if b.tokens >= 1 {
		b.tokens--
		return true
	}
	return false
}
