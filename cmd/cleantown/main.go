package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/adaken4/clean-town/internal/app"
	"github.com/adaken4/clean-town/internal/config"
	"github.com/adaken4/clean-town/internal/database"
	"github.com/adaken4/clean-town/internal/monitoring"
	"github.com/adaken4/clean-town/internal/router"
)

func main() {
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

	// Declare an instance of the application struct, containing the config struct and
	// the logger.
	appInstance := app.New(cfg, logger, db)
	defer appInstance.Shutdown(ctx)

	// Declare a HTTP server which listens on the port provided in the config struct,
	// uses the servemux we created above as the handler, has some sensible timeout
	// settings and writes any log messages to the structured logger at Error level.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router.New(appInstance),
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
