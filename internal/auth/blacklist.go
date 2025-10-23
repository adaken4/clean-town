package auth

import (
	"sync"
	"time"
)

// Blacklist defines the interface for managing revoked tokens.
type Blacklist interface {
	Revoke(token string, expiry time.Time) error // Marks a token as revoked until a specified expiry time.
	IsRevoked(token string) bool                 // Checks if a token is currently revoked.
	Cleanup()                                    // Removes expired tokens from the blacklist.
}

// InMemoryBlacklist is an in-memory implementation of the Blacklist interface.
type InMemoryBlacklist struct {
	revokedTokens map[string]time.Time // Stores revoked tokens and their expiry times.
	mutex         sync.RWMutex         // Ensures thread-safe access to revokedTokens.
	stopCleanup   chan struct{}        // Signals the cleanup goroutine to stop.
	wg            sync.WaitGroup       // Tracks the cleanup goroutine for graceful shutdown.
}

// NewInMemoryBlacklist initializes a new InMemoryBlacklist and starts periodic cleanup.
func NewInMemoryBlacklist() *InMemoryBlacklist {
	bl := &InMemoryBlacklist{
		revokedTokens: make(map[string]time.Time),
		stopCleanup:   make(chan struct{}),
	}

	// Start background cleanup goroutine.
	bl.wg.Add(1)
	go bl.periodicCleanup()

	return bl
}

// Revoke adds a token to the blacklist with its expiry time.
func (bl *InMemoryBlacklist) Revoke(token string, expiry time.Time) error {
	bl.mutex.Lock()
	defer bl.mutex.Unlock()

	bl.revokedTokens[token] = expiry
	return nil
}

// IsRevoked checks if a token is currently revoked and not yet expired.
func (bl *InMemoryBlacklist) IsRevoked(token string) bool {
	bl.mutex.Lock()
	defer bl.mutex.Unlock()

	if expiry, exists := bl.revokedTokens[token]; exists {
		// Token is considered revoked only if current time is before its expiry.
		return time.Now().Before(expiry)
	}
	return false
}

// Cleanup removes tokens that have expired from the blacklist.
func (bl *InMemoryBlacklist) Cleanup() {
	bl.mutex.Lock()
	defer bl.mutex.Unlock()

	now := time.Now()
	for token, expiry := range bl.revokedTokens {
		if now.After(expiry) {
			delete(bl.revokedTokens, token)
		}
	}
}

// periodicCleanup runs in the background and periodically invokes Cleanup.
func (bl *InMemoryBlacklist) periodicCleanup() {
	defer bl.wg.Done()

	ticker := time.NewTicker(1 * time.Hour) // Cleanup interval.
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			bl.Cleanup() // Perform cleanup every hour.
		case <-bl.stopCleanup:
			return // Exit the goroutine when stop signal is received.
		}
	}
}
