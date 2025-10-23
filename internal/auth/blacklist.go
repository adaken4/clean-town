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
