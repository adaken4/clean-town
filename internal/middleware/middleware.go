package middleware

import (
	"log/slog"
	"sync"

	"golang.org/x/time/rate"
)

// rateLimiter implements an IP-based request rate limiter middleware.
// It limits how frequently each client (by IP) can make requests to the server.
type rateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
	r        rate.Limit
	burst    int
	logger   *slog.Logger
}

// NewRateLimiter creates and returns a new rateLimiter instance.
// Parameters:
//   - r: rate of requests allowed per second (e.g., 1 means 1 request/sec)
//   - burst: number of requests allowed to burst beyond the rate
//   - logger: used for logging warnings and errors
func NewRateLimiter(r rate.Limit, burst int, logger *slog.Logger) *rateLimiter {
	return &rateLimiter{
		limiters: make(map[string]*rate.Limiter),
		r:        r,
		burst:    burst,
		logger:   logger,
	}
}

// getLimiter retrieves (or creates) a rate limiter for the given client IP.
// Each IP gets its own rate limiter instance.
func (rl *rateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Check if limiter for this IP already exists
	limiter, exists := rl.limiters[ip]
	if !exists {
		// Create a new limiter for this IP
		limiter = rate.NewLimiter(rl.r, rl.burst)
		rl.limiters[ip] = limiter
	}
	return limiter
}
