package handlers

import "github.com/adaken4/clean-town/internal/app"

// Handlers groups all HTTP handler functions and provides access to shared application resources.
type Handlers struct {
	app *app.App // Reference to the core application struct containing config, logger, DB, etc.
}

// New initializes and returns a Handler instance with access to the app's dependencies
func New(app *app.App) *Handlers {
	return &Handlers{app}
}
