package handlers

import (
	"encoding/json"
	"net/http"
)

// writeJSON marshals the given data to JSON and writes it to the response writer.
// It sets the provided status code and includes any optional HTTP headers.
// The JSON output is appended with a newline for readability.
func (h *Handlers) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Append newline for pretty formatting
	js = append(js, '\n')

	// Set additional headers if provided
	for key, value := range headers {
		w.Header()[key] = value
	}

	// Set content-type for JSON response
	w.Header().Set("Content-Type", "application/json")
	
	// Write status code and JSON response
	w.WriteHeader(status)
	w.Write(js)
	return nil
}