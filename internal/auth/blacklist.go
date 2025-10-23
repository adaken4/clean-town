package auth

import (
	"time"
)

// Blacklist defines the interface for managing revoked tokens.
type Blacklist interface {
	Revoke(token string, expiry time.Time) error // Marks a token as revoked until a specified expiry time.
	IsRevoked(token string) bool                 // Checks if a token is currently revoked.
	Cleanup()                                    // Removes expired tokens from the blacklist.
}
