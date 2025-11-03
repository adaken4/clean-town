package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/adaken4/clean-town/internal/models"
)

// RegisterUser handles user registration requests.
func (h *Handlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Create a context with a timeout to avoid hanging requests
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	// Decode the incoming JSON request into the RegisterRequest struct
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		h.app.Logger.Error(err.Error())
		return
	}

	// Attempt to register the user via the Auth service
	_, err := h.app.Auth.Register(ctx, req)
	if err != nil {
		switch err {
		case models.ErrDuplicateEmail:
			http.Error(w, "email already registered", http.StatusConflict)
			h.app.Logger.Warn("duplicate email registration attempt: " + req.Email)
		default:
			http.Error(w, "server error", http.StatusInternalServerError)
			h.app.Logger.Error("registration error: " + err.Error())
		}
		return
	}

	// Respond with success message
	h.writeJSON(w, http.StatusCreated, map[string]any{
		"message": "registration successful, please verify your email",
	}, nil)
}

// LoginUser handles user login requests.
func (h *Handlers) LoginUser(w http.ResponseWriter, r *http.Request) {
	// Create a context with timeout to prevent hanging requests
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	// Decode the login request payload
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		h.app.Logger.Error("login request decode error: " + err.Error())
		return
	}

	// Authenticate user and generate tokens
	access, refresh, err := h.app.Auth.Login(ctx, req.Email, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		h.app.Logger.Warn("failed login attempt for email" + req.Email)
		return
	}

	// Set access token cookie (short-lived)
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    access,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.app.Config.Server.Env == "production",
		SameSite: http.SameSiteStrictMode,
		MaxAge:   900, // 15 minutes
	})

	// Set refresh token cookie (long-lived)
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		Path:     "/v1/auth/refresh",
		HttpOnly: true,
		Secure:   h.app.Config.Server.Env == "production",
		SameSite: http.SameSiteStrictMode,
		MaxAge:   604800, // 7 days
	})

	// Respond with success message
	h.writeJSON(w, http.StatusOK, map[string]string{
		"message":      "login successful",
		"access_token": access,
	}, nil)
}
