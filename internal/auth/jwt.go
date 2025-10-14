package auth

import (
	"fmt"
	"time"

	"github.com/adaken4/clean-town/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	UserID   int64
	UserRole string `json:"role"`
}

const TokenIssuer = "clean-town-auth"

func GenerateAccessToken(signingKey []byte, user models.User) (string, error) {
	claims := CustomClaims{
		UserID:   user.ID,
		UserRole: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    TokenIssuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("user-%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(signingKey)
}

func GenerateRefreshToken(signingKey []byte, user models.User) (string, error) {
	claims := CustomClaims{
		UserID:   user.ID,
		UserRole: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    TokenIssuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // longer expiry (7 days)
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("user-%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}

func ValidateClaims(claims *CustomClaims) error {
	// Validate expiration
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return fmt.Errorf("token has expired")
	}

	// Validate issuer (if needed)
	if claims.Issuer != TokenIssuer {
		return fmt.Errorf("invalid issuer")
	}

	// Add any other custom validation logic

	return nil
}
