package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/adaken4/clean-town/internal/config"
	"github.com/adaken4/clean-town/internal/database"
)

// Declare a string containing the application version number
// TODO: Generate this automatically at build time
const version = "0.0.1"

// Define an application struct to hold the dependencies for our HTTP handlers, helpers,
// and middleware. At the moment this only contains a copy of the config struct and a
// logger, but it will grow to include a lot more as the application progresses.
type application struct {
	config *config.Config
	logger *slog.Logger
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	startMetricsMonitoring(ctx, db)

	// Declare an instance of the application struct, containing the config struct and
	// the logger.
	app := &application{
		config: cfg,
		logger: logger,
	}

	// Declare a new servemux and add a /v1/healthcheck route which dispatches requests
	// to the healthcheckHandler method.
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)

	// Declare a HTTP server which listens on the port provided in the config struct,
	// uses the servemux we created above as the handler, has some sensible timeout
	// settings and writes any log messages to the structured logger at Error level.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      mux,
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
