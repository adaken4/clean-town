package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/adaken4/clean-town/internal/auth"
	"github.com/adaken4/clean-town/internal/config"
	"github.com/adaken4/clean-town/internal/database"
	"github.com/adaken4/clean-town/internal/models"
)

// Declare a string containing the application version number
// TODO: Generate this automatically at build time
const version = "0.0.1"

// Define an application struct to hold the dependencies for our HTTP handlers, helpers,
// and middleware. At the moment this only contains a copy of the config struct and a
// logger, but it will grow to include a lot more as the application progresses.
type application struct {
	config   *config.Config
	logger   *slog.Logger
	userRepo models.UserRepository
}

func main() {
	// Initialize a new structured logger which writes log entries to the standard out
	// stream.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Declare an instance of the config struct.
	var cfg *config.Config

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

	startMetricsMonitoring(ctx, db)

	userRepo := &models.PostgresUserRepository{DB: db}

	blacklist := auth.NewInMemoryBlacklist()
	auth.InitBlacklist(blacklist)
	defer blacklist.Close()
	logger.Info("token blacklist initialized and periodic cleanup started")

	// Declare an instance of the application struct, containing the config struct and
	// the logger.
	app := &application{
		config:   cfg,
		logger:   logger,
		userRepo: userRepo,
	}

	// Declare a HTTP server which listens on the port provided in the config struct,
	// uses the servemux we created above as the handler, has some sensible timeout
	// settings and writes any log messages to the structured logger at Error level.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	// Start the HTTP server.
	logger.Info("starting server", "addr", srv.Addr, "env", cfg.Server.Env)

	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
