package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/adaken4/clean-town/internal/app"
	"github.com/adaken4/clean-town/internal/router"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Declare an instance of the application struct, containing the config struct and
	// the logger.
	appInstance := app.New()
	defer appInstance.Shutdown(ctx)

	// Declare a HTTP server which listens on the port provided in the config struct,
	// uses the servemux we created above as the handler, has some sensible timeout
	// settings and writes any log messages to the structured logger at Error level.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", appInstance.Config.Server.Port),
		Handler:      router.New(appInstance),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(appInstance.Logger.Handler(), slog.LevelError),
	}

	// Start the HTTP server.
	appInstance.Logger.Info("starting server", "addr", srv.Addr, "env", appInstance.Config.Server.Env)

	err := srv.ListenAndServe()
	appInstance.Logger.Error(err.Error())
	os.Exit(1)
}
