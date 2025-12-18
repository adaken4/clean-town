package app

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	"github.com/adaken4/clean-town/internal/auth"
	"github.com/adaken4/clean-town/internal/config"
	"github.com/adaken4/clean-town/internal/database"
	"github.com/adaken4/clean-town/internal/models"
	"github.com/adaken4/clean-town/internal/monitoring"
	"github.com/adaken4/clean-town/internal/services"
)

// App encapsulates core dependencies and services for the CleanTown application.
type App struct {
	Config    *config.Config          // Application configuration
	Logger    *slog.Logger            // Structured logger
	UserRepo  models.UserRepository   // Interface for user data access
	DB        *sql.DB                 // Database connection
	Blacklist *auth.InMemoryBlacklist // In-memory token blacklist for JWT revocation
	Auth      *services.AuthService
}

// New initializes the App with its dependencies and starts the blacklist cleanup routine.
func New() *App {
	// Initialize a new structured logger which writes log entries to the standard out
	// stream.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg, err := config.LoadConfig(".env.example")
	if err != nil {
		logger.Error("Failed to load configuration:", "error", err.Error())
		os.Exit(1)
	}

	db, err := database.ConnectWithRetry(cfg.Database)
	if err != nil {
		logger.Error("database connection failed", "error", err.Error())
		os.Exit(1)
	}

	// Defer a call to db.Close() so that the connection pool is closed before the
	// main() function exits.
	defer db.Close()

	// Also log a message to say that the connection pool has been successfully
	// established.
	logger.Info("database connection pool established")

	// Apply all pending "up" migrations from the ./migrations directory.
	// If there are no new migrations, it does nothing.
	// The application exits on migration errors to prevent running with an invalid schema.
	err = database.RunMigrations(db, "./migrations")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("database migrations applied")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	monitoring.StartMetricsMonitoring(ctx, db)

	userRepo := &models.PostgresUserRepository{DB: db} // Set up user repository
	blacklist := auth.NewInMemoryBlacklist()           // Create in-memory blacklist
	auth.InitBlacklist(blacklist)                      // Start periodic cleanup

	authService := services.AuthService{
		Config:    cfg,
		UserRepo:  userRepo,
		Blacklist: blacklist,
		Logger:    logger,
	}

	logger.Info("token blacklist initialized and periodic cleanup started")

	return &App{
		Config:    cfg,
		Logger:    logger,
		UserRepo:  userRepo,
		DB:        db,
		Blacklist: blacklist,
		Auth:      &authService,
	}
}

// Shutdown gracefully closes resources like the blacklist and database connection.
func (a *App) Shutdown(ctx context.Context) error {
	a.Blacklist.Close() // Stop blacklist cleanup and release resources
	return a.DB.Close() // Close database connection
}
