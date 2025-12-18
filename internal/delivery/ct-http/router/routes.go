package router

import (
	"net/http"

	"github.com/adaken4/clean-town/internal/app"
	"github.com/adaken4/clean-town/internal/delivery/ct-http/handlers"
	"github.com/adaken4/clean-town/internal/delivery/ct-http/middleware"
	"github.com/julienschmidt/httprouter"
)

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// New sets up the HTTP router and maps routes to their corresponding handlers.
func New(a *app.App) http.Handler {
	router := httprouter.New()

	h := handlers.New(a)

	// Rate limiter : 5 requests per minute per IP
	rl := middleware.NewRateLimiter(5.0/60.0, 5, a.Logger)

	// Health check endpoint to verify server status
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", h.HealthCheck)

	// User registration endpoint
	router.Handler(http.MethodPost, "/v1/auth/register", rl.Middleware(http.HandlerFunc(h.RegisterUser)))

	// User login endpoint
	router.Handler(http.MethodPost, "/v1/auth/login", rl.Middleware(http.HandlerFunc(h.LoginUser)))

	// Token refresh endpoint using refresh token cookie
	router.Handler(http.MethodPost, "/v1/auth/refresh", rl.Middleware(http.HandlerFunc(h.RefreshToken)))

	// User logout endpoint to clear tokens and revoke session
	router.Handler(http.MethodPost, "/v1/auth/logout", rl.Middleware(http.HandlerFunc(h.LogoutUser)))

	// Email verification endpoint (via token query param)
	router.Handler(http.MethodGet, "/v1/auth/verify-email", rl.Middleware(http.HandlerFunc(h.VerifyEmail)))

	return cors(router)
}
