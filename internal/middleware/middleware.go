package middleware

import (
	"log/slog"
	"net"
	"net/http"
	"strings"
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

// getClientIP extracts the client's IP address from the HTTP request.
// It prioritizes the X-Forwarded-For and X-Real-IP headers, commonly set by proxies.
func getClientIP(r *http.Request) (string, error) {
	// Check X-Forwarded-For header (may contain multiple IPs)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// The first IP in the list is the original client
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0]), nil
	}

	// Check X-Real-IP header (used by some reverse proxies)
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return xrip, nil
	}

	// Fallback: extract IP directly from the request’s remote address
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	return ip, nil
}
