package app

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/adaken4/clean-town/internal/auth"
	"github.com/adaken4/clean-town/internal/config"
	"github.com/adaken4/clean-town/internal/models"
)

// App encapsulates core dependencies and services for the CleanTown application.
type App struct {
	Config    *config.Config          // Application configuration
	Logger    *slog.Logger            // Structured logger
	UserRepo  models.UserRepository   // Interface for user data access
	DB        *sql.DB                 // Database connection
	Blacklist *auth.InMemoryBlacklist // In-memory token blacklist for JWT revocation

}

// New initializes the App with its dependencies and starts the blacklist cleanup routine.
func New(cfg *config.Config, logger *slog.Logger, db *sql.DB) *App {
	userRepo := &models.PostgresUserRepository{DB: db} // Set up user repository
	blacklist := auth.NewInMemoryBlacklist()           // Create in-memory blacklist
	auth.InitBlacklist(blacklist)                      // Start periodic cleanup

	logger.Info("token blacklist initialized and periodic cleanup started")

	return &App{
		Config:    cfg,
		Logger:    logger,
		UserRepo:  userRepo,
		DB:        db,
		Blacklist: blacklist,
	}
}

// Shutdown gracefully closes resources like the blacklist and database connection.
func (a *App) Shutdown(ctx context.Context) error {
	a.Blacklist.Close() // Stop blacklist cleanup and release resources
	return a.DB.Close() // Close database connection
}
