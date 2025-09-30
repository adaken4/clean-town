package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adaken4/clean-town/internal/database"
)

// Declare a handler which writes a plain-text response with information about the
// application status, operating environment and version.
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "environment: %s\n", app.config.Server.Env)
	fmt.Fprintf(w, "version: %s\n", version)
}

func healthHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics := database.CollectMetrics(db)

		w.Header().Set("Content-Type", "application/json")
		if !metrics.Healthy {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		json.NewEncoder(w).Encode(metrics)
	}
}
