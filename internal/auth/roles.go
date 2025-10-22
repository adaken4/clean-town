package auth

import (
	"net/http"
)

const (
	RoleAdmin     = "admin"
	RoleOrganizer = "organizer"
	RoleVolunteer = "volunteer"
)

// GetUserRole extracts the user's role from the request context.
func GetUserRole(r *http.Request) (string, bool) {
	// Retrieve claims from the context (set by AuthMiddleware after JWT verification).
	claims, ok := r.Context().Value(claimsKey).(*CustomClaims)
	if !ok {
		// No valid claims found in context; user not authenticated.
		return "", false
	}
	return claims.UserRole, true
}

// RequireRole is a role-based access control (RBAC) middleware generator.
// It restricts route access to handlers based on the user's role.
// Example usage:
//
//	http.HandleFunc("/admin", auth.RequireRole(auth.RoleAdmin)(adminHandler))
func RequireRole(allowedRoles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Extract the user's role from JWT claims.
			userRole, ok := GetUserRole(r)
			if !ok {
				http.Error(w, "missing or invalid user claims", http.StatusUnauthorized)
				return
			}

			// Check if the user's role is allowed to access this route.
			for _, role := range allowedRoles {
				if userRole == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Reject if user's role does not match any allowed role.
			http.Error(w, "forbidden: insufficient permissions", http.StatusForbidden)
		}
	}
}
