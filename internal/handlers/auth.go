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

// VerifyEmail handles user email verification requests using email verification token.
func (h *Handlers) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	// Create a context with timeout to prevent hanging requests
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	// Get token from query parameters
	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		http.Error(w, "missing token", http.StatusBadRequest)
		h.app.Logger.Warn("email verification attempted without token")
		return
	}

	// Attempt to verify the email using the token
	if err := h.app.Auth.VerifyEmail(ctx, tokenStr); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		h.app.Logger.Error("email verification failed: " + err.Error())
		return
	}

	// Respond with success message
	h.writeJSON(w, http.StatusOK, map[string]string{
		"message": "Email successfully verfied. You can now log in.",
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

// RefreshToken handles requests to refresh authentication tokens.
func (h *Handlers) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Create a context with timeout to prevent hanging requests
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	// Retrieve the refresh token from cookie
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "refresh token missing", http.StatusUnauthorized)
		h.app.Logger.Warn("refresh token cookie not found")
		return
	}

	// Attempt to refresh tokens using the provided refresh token
	access, refresh, err := h.app.Auth.RefreshTokens(ctx, cookie.Value)
	if err != nil {
		http.Error(w, "invalid or expired refresh token", http.StatusUnauthorized)
		h.app.Logger.Warn("invalid refresh token: " + err.Error())
		return
	}

	// Set new access token cookie (short-lived)
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    access,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.app.Config.Server.Env == "production",
		SameSite: http.SameSiteStrictMode,
		MaxAge:   900, // 15 minutes
	})

	// Set new refresh token cookie (long-lived)
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		Path:     "/v1/auth/refresh",
		HttpOnly: true,
		Secure:   h.app.Config.Server.Env == "production",
		SameSite: http.SameSiteStrictMode,
		MaxAge:   604800, // 7 days
	})

	// Respond with the new access_token
	h.writeJSON(w, http.StatusOK, map[string]string{
		"message":      "token refreshed successfully",
		"access_token": access,
	}, nil)
}

// LogoutUser handles user logout requests.
func (h *Handlers) LogoutUser(w http.ResponseWriter, r *http.Request) {
	// Create a context with timeout to prevent hanging requests
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	// Retrieve the refresh token from cookie
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "refresh token missing", http.StatusUnauthorized)
		return
	}

	// Attempt to revoke the refresh token
	if err := h.app.Auth.Logout(ctx, cookie.Value); err != nil {
		h.app.Logger.Error("logout failed", "error", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Clear access and refresh token cookies
	clearCookie := func(name string, path string) {
		http.SetCookie(w, &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     path,
			HttpOnly: true,
			Secure:   h.app.Config.Server.Env == "production",
			SameSite: http.SameSiteStrictMode,
			MaxAge:   -1, // Expire immediately
		})
	}

	clearCookie("access_token", "/")
	clearCookie("refresh_token", "/v1/auth/refresh")

	// Respond with success message
	h.writeJSON(w, http.StatusOK, map[string]string{
		"message": "logout successful",
	}, nil)
}
