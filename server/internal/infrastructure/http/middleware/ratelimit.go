package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// ipEntry tracks request counts per IP
type ipEntry struct {
	count    int
	resetAt  time.Time
}

// RateLimiter implements per-IP rate limiting
type RateLimiter struct {
	mu        sync.Mutex
	ips       map[string]*ipEntry
	limit     int
	window    time.Duration
	stopCh    chan struct{}
	closeOnce sync.Once
}

// NewRateLimiter creates a rate limiter that allows `limit` requests per `window` per IP
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		ips:    make(map[string]*ipEntry),
		limit:  limit,
		window: window,
		stopCh: make(chan struct{}),
	}

	// Background cleanup every 5 minutes to prevent memory leak
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				rl.cleanup()
			case <-rl.stopCh:
				return
			}
		}
	}()

	return rl
}

// Close stops the background cleanup goroutine. Safe to call multiple times.
func (rl *RateLimiter) Close() {
	rl.closeOnce.Do(func() {
		close(rl.stopCh)
	})
}

// cleanup removes expired entries
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	for ip, entry := range rl.ips {
		if now.After(entry.resetAt) {
			delete(rl.ips, ip)
		}
	}
}

// allow checks if the IP is within limits
func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	entry, exists := rl.ips[ip]

	if !exists || now.After(entry.resetAt) {
		rl.ips[ip] = &ipEntry{count: 1, resetAt: now.Add(rl.window)}
		return true
	}

	if entry.count >= rl.limit {
		return false
	}

	entry.count++
	return true
}

// RateLimitMiddleware creates a Gin middleware for rate limiting
func RateLimitMiddleware(limiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.allow(c.ClientIP()) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "rate limit exceeded",
				"retry_after": limiter.window.String(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
