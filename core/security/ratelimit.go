package security

import (
	"sync"
	"time"
)

type ipCounter struct {
	count int
	reset time.Time
}

type RateLimiter struct {
	mu     sync.Mutex
	window time.Duration
	max    int
	store  map[string]*ipCounter
}

func NewRateLimiter(max int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		max:    max,
		window: window,
		store:  make(map[string]*ipCounter),
	}
}

func (r *RateLimiter) Allow(ip string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	c, ok := r.store[ip]
	if !ok || time.Now().After(c.reset) {
		r.store[ip] = &ipCounter{
			count: 1,
			reset: time.Now().Add(r.window),
		}
		return true
	}

	if c.count >= r.max {
		return false
	}

	c.count++
	return true
}
