package services

import (
	"log/slog"

	"github.com/adaken4/clean-town/internal/auth"
	"github.com/adaken4/clean-town/internal/config"
	"github.com/adaken4/clean-town/internal/models"
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
