package application

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

// TokenBlacklist maintains a set of revoked JWT tokens in memory.
// Tokens are stored as SHA-256 hashes to prevent attacker-controlled strings
// from being used as map keys and to reduce memory usage.
type TokenBlacklist struct {
	mu        sync.RWMutex
	tokens    map[string]time.Time // hash(token) -> expiry time
	stopCh    chan struct{}
	closeOnce sync.Once
}

// hashToken returns the SHA-256 hex digest of a token string.
func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

// NewTokenBlacklist creates a new token blacklist with background cleanup.
func NewTokenBlacklist() *TokenBlacklist {
	bl := &TokenBlacklist{
		tokens: make(map[string]time.Time),
		stopCh: make(chan struct{}),
	}
	go bl.cleanupLoop()
	return bl
}

// Close stops the background cleanup goroutine. Safe to call multiple times.
func (bl *TokenBlacklist) Close() {
	bl.closeOnce.Do(func() {
		close(bl.stopCh)
	})
}

// Revoke adds a token to the blacklist until its expiry time.
// The token is stored as a SHA-256 hash.
func (bl *TokenBlacklist) Revoke(token string, expiresAt time.Time) {
	bl.mu.Lock()
	defer bl.mu.Unlock()
	bl.tokens[hashToken(token)] = expiresAt
}

// IsRevoked checks if a token has been revoked and is still within its expiry window.
func (bl *TokenBlacklist) IsRevoked(token string) bool {
	bl.mu.RLock()
	defer bl.mu.RUnlock()
	expiresAt, exists := bl.tokens[hashToken(token)]
	if !exists {
		return false
	}
	return time.Now().Before(expiresAt)
}

// cleanupLoop removes expired tokens every minute to prevent memory growth.
func (bl *TokenBlacklist) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			bl.cleanup()
		case <-bl.stopCh:
			return
		}
	}
}

func (bl *TokenBlacklist) cleanup() {
	bl.mu.Lock()
	defer bl.mu.Unlock()
	now := time.Now()
	for token, expiresAt := range bl.tokens {
		if now.After(expiresAt) {
			delete(bl.tokens, token)
		}
	}
}
