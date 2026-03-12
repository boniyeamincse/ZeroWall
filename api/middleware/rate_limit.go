package middleware

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	requests map[string]*clientStats
	mu       sync.RWMutex
	rate     int
	burst    int
	timeout  time.Duration
}

type clientStats struct {
	tokens    int
	lastCheck time.Time
}

func NewRateLimiter(rate int, burst int, cleanupInterval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]*clientStats),
		rate:     rate,
		burst:    burst,
		timeout:  cleanupInterval,
	}

	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.timeout)
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, stats := range rl.requests {
			if now.Sub(stats.lastCheck) > rl.timeout*2 {
				delete(rl.requests, key)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	stats, exists := rl.requests[key]

	if !exists {
		rl.requests[key] = &clientStats{
			tokens:    rl.burst - 1,
			lastCheck: now,
		}
		return true
	}

	elapsed := now.Sub(stats.lastCheck)
	refill := int(elapsed.Seconds()) * rl.rate
	stats.tokens += refill
	if stats.tokens > rl.burst {
		stats.tokens = rl.burst
	}
	stats.lastCheck = now

	if stats.tokens > 0 {
		stats.tokens--
		return true
	}

	return false
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := getClientKey(r)

		if !rl.Allow(key) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error": "rate limit exceeded", "retry_after": 1}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getClientKey(r *http.Request) string {
	 forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func RateLimit(rate int, burst int) func(http.Handler) http.Handler {
	limiter := NewRateLimiter(rate, burst, 5*time.Minute)
	return limiter.Middleware
}
