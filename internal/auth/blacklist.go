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
