package auth

import (
	"context"
	"crypto/subtle"
	"log"
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
		if authHeader == "" {
			http.Error(w, "authentication failed", http.StatusUnauthorized)
			log.Printf("Auth attempt: %s %s", r.Method, r.URL.Path)
			return
		}
		
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "authentication failed", http.StatusUnauthorized)
			log.Printf("Auth attempt: %s %s", r.Method, r.URL.Path)
			return
		}
		headerToken = parts[1]
		
		// 2. Get cookie token
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "authentication failed", http.StatusUnauthorized)
			log.Printf("Auth attempt: %s %s", r.Method, r.URL.Path)
			return
		}
		cookieToken := cookie.Value
		
		// Compare header and cookie tokens
		if subtle.ConstantTimeCompare([]byte(headerToken), []byte(cookieToken)) != 1 {
			http.Error(w, "authentication failed", http.StatusUnauthorized)
			log.Printf("Auth attempt: %s %s", r.Method, r.URL.Path)
			return
		}
		
		// Verify token
		claims, err := VerifyToken(headerToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			log.Printf("Auth attempt: %s %s", r.Method, r.URL.Path)
			return
		}

		// Add claims to request context
		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
