package auth

import (
	"context"
	"crypto/subtle"
	"net/http"
	"strings"
)

type contextKey string

const claimsKey = contextKey("claims")

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var headerToken string
		// 1. Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			headerToken = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			http.Error(w, "missing Authorization header", http.StatusUnauthorized)
			return
		}

		// 2. Get cookie token
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "missing access token cookie", http.StatusUnauthorized)
			return
		}
		cookieToken := cookie.Value

		// Compare header and cookie tokens
		if subtle.ConstantTimeCompare([]byte(headerToken), []byte(cookieToken)) != 1 {
			http.Error(w, "token mismatch", http.StatusUnauthorized)
			return
		}

		// Verify token
		claims, err := VerifyToken(headerToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Add claims to request context
		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
