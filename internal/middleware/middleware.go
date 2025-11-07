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
