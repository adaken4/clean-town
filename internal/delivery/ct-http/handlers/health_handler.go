package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adaken4/clean-town/internal/database"
)

// Declare a string containing the application version number
// TODO: Generate this automatically at build time
// const version = "0.0.1"
const version = "0.0.1"

// Declare a handler which writes a plain-text response with information about the
// application status, operating environment and version.
func (h *Handlers) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "environment: %s\n", h.app.Config.Server.Env)
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
