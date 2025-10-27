package router

import (
	"net/http"

	"github.com/adaken4/clean-town/internal/app"
	"github.com/adaken4/clean-town/internal/handlers"
	"github.com/julienschmidt/httprouter"
)

// New sets up the HTTP router and maps routes to their corresponding handlers.
func New(a *app.App) http.Handler {
	router := httprouter.New()

	h := handlers.New(a)

	// Health check endpoint to verify server status
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", h.HealthCheck)
	
	// User registration endpoint
	router.HandlerFunc(http.MethodPost, "/v1/auth/register", h.RegisterUser)
	
	// User login endpoint
	router.HandlerFunc(http.MethodPost, "/v1/auth/login", h.LoginUser)
	
	// Token refresh endpoint using refresh token cookie
	router.HandlerFunc(http.MethodPost, "/v1/auth/refresh", h.RefreshToken)
	
	// User logout endpoint to clear tokens and revoke session
	router.HandlerFunc(http.MethodPost, "/v1/auth/logout", h.LogoutUser)
	
	// Email verification endpoint (via token query param)
	router.HandlerFunc(http.MethodGet, "/v1/auth/verify-email", h.VerifyEmail)

	return router
}
