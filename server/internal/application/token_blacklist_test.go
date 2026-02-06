package application

import (
	"sync"
	"testing"
	"time"
)

func TestNewTokenBlacklist(t *testing.T) {
	bl := NewTokenBlacklist()
	defer bl.Close()
	if bl == nil {
		t.Fatal("Expected non-nil blacklist")
	}
	if bl.tokens == nil {
		t.Fatal("Expected initialized tokens map")
	}
}

func TestTokenBlacklist_Revoke(t *testing.T) {
	bl := &TokenBlacklist{tokens: make(map[string]time.Time)}

	expiry := time.Now().Add(time.Hour)
	bl.Revoke("token1", expiry)

	bl.mu.RLock()
	storedExpiry, exists := bl.tokens[hashToken("token1")]
	bl.mu.RUnlock()

	if !exists {
		t.Fatal("Expected token hash to exist in blacklist")
	}
	if !storedExpiry.Equal(expiry) {
		t.Errorf("Expected expiry %v, got %v", expiry, storedExpiry)
	}
}

func TestTokenBlacklist_Revoke_OverwritesExisting(t *testing.T) {
	bl := &TokenBlacklist{tokens: make(map[string]time.Time)}

	bl.Revoke("token1", time.Now().Add(time.Hour))
	newExpiry := time.Now().Add(2 * time.Hour)
	bl.Revoke("token1", newExpiry)

	bl.mu.RLock()
	storedExpiry := bl.tokens[hashToken("token1")]
	bl.mu.RUnlock()

	if !storedExpiry.Equal(newExpiry) {
		t.Errorf("Expected overwritten expiry %v, got %v", newExpiry, storedExpiry)
	}
}

func TestTokenBlacklist_IsRevoked_NotRevoked(t *testing.T) {
	bl := &TokenBlacklist{tokens: make(map[string]time.Time)}

	if bl.IsRevoked("nonexistent") {
		t.Error("Expected non-existent token to not be revoked")
	}
}

func TestTokenBlacklist_IsRevoked_ActiveToken(t *testing.T) {
	bl := &TokenBlacklist{tokens: make(map[string]time.Time)}

	bl.Revoke("token1", time.Now().Add(time.Hour))

	if !bl.IsRevoked("token1") {
		t.Error("Expected token with future expiry to be revoked")
	}
}

func TestTokenBlacklist_IsRevoked_ExpiredToken(t *testing.T) {
	bl := &TokenBlacklist{tokens: make(map[string]time.Time)}

	// Token that expired in the past - store using hash
	bl.tokens[hashToken("token1")] = time.Now().Add(-time.Hour)

	if bl.IsRevoked("token1") {
		t.Error("Expected expired token to not be considered revoked")
	}
}

func TestTokenBlacklist_Cleanup(t *testing.T) {
	bl := &TokenBlacklist{tokens: make(map[string]time.Time)}

	// Add expired and active tokens using hashed keys (as Revoke does)
	bl.tokens[hashToken("expired1")] = time.Now().Add(-time.Hour)
	bl.tokens[hashToken("expired2")] = time.Now().Add(-2 * time.Hour)
	bl.tokens[hashToken("active1")] = time.Now().Add(time.Hour)

	bl.cleanup()

	bl.mu.RLock()
	defer bl.mu.RUnlock()

	if _, exists := bl.tokens[hashToken("expired1")]; exists {
		t.Error("Expected expired1 to be cleaned up")
	}
	if _, exists := bl.tokens[hashToken("expired2")]; exists {
		t.Error("Expected expired2 to be cleaned up")
	}
	if _, exists := bl.tokens[hashToken("active1")]; !exists {
		t.Error("Expected active1 to remain")
	}
}

func TestTokenBlacklist_Cleanup_EmptyMap(t *testing.T) {
	bl := &TokenBlacklist{tokens: make(map[string]time.Time)}
	bl.cleanup() // Should not panic
}

func TestTokenBlacklist_Close_Idempotent(t *testing.T) {
	bl := NewTokenBlacklist()
	bl.Close()
	bl.Close() // must not panic
}

func TestTokenBlacklist_ConcurrentAccess(t *testing.T) {
	bl := NewTokenBlacklist()
	defer bl.Close()

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		token := "token" + string(rune('A'+i%26))
		go func() {
			defer wg.Done()
			bl.Revoke(token, time.Now().Add(time.Hour))
		}()
		go func() {
			defer wg.Done()
			bl.IsRevoked(token)
		}()
	}
	wg.Wait()
}
