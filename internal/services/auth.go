package services

import (
	"log/slog"
	"time"

	"github.com/adaken4/clean-town/internal/auth"
	"github.com/adaken4/clean-town/internal/config"
	"github.com/adaken4/clean-town/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService provides all authentication and authorization operations
// such as registration, login, logout, token refresh, and email verification.
// It depends on a UserRepository for persistence and a JWT-based token system.
type AuthService struct {
	Config    *config.Config
	UserRepo  models.UserRepository
	Blacklist *auth.InMemoryBlacklist
	Logger    *slog.Logger
}

// HashPassword hashes a plain password
func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// CheckPassword compares a plain password with the stored hash
func CheckPassword(hash []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(password))
}

// GenerateEmailVerificationToken creates a signed JWT for verifying user emails.
// The token contains the email address, token type, and expiry time (30 minutes).
func GenerateEmailVerificationToken(signingKey []byte, email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"type":  "email_verification",
		"exp":   time.Now().Add(30 * time.Minute).Unix(),
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}
