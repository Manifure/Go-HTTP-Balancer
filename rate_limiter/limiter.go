package rate_limiter

import (
	"net/http"
	"sync"
)

type Limiter struct {
	mu        sync.RWMutex
	buckets   map[string]*TokenBucket
	defaultFn func(key string) (int, float64) // лимиты: capacity, rate
}

func NewLimiter(defaultFn func(string) (int, float64)) *Limiter {
	return &Limiter{
		buckets:   make(map[string]*TokenBucket),
		defaultFn: defaultFn,
	}
}

func (l *Limiter) getBucket(key string) *TokenBucket {
	l.mu.RLock()
	bucket, exists := l.buckets[key]
	l.mu.RUnlock()
	if exists {
		return bucket
	}

	// Создание нового bucket по ключу (IP или API-ключ)
	capacity, rate := l.defaultFn(key)
	bucket = NewTokenBucket(capacity, rate)

	l.mu.Lock()
	l.buckets[key] = bucket
	l.mu.Unlock()

	return bucket
}

// Middleware для HTTP
func (l *Limiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.RemoteAddr // Можно заменить на API-ключ, если он есть
		if !l.getBucket(key).Allow() {
			http.Error(w, "429 Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
