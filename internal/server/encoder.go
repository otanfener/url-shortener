package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Encoder handles encoding responses and errors.
type Encoder struct {
	logger logger
}

// NewEncoder initializes an encoder with logging.
func NewEncoder(l logger) *Encoder {
	return &Encoder{
		logger: l,
	}
}

// errorResponse encapsulates errors for HTTP responses.
type errorResponse struct {
	Message string `json:"message"`
}

// StatusResponse writes a JSON response with a given status code.
func (e *Encoder) StatusResponse(ctx context.Context, w http.ResponseWriter, response interface{}, status int) {
	if response == nil {
		w.WriteHeader(status)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		e.logger.Error(fmt.Errorf("error encoding response"), map[string]interface{}{"error": err.Error()})
	}
}

// Error writes an error response with an appropriate HTTP status code.
func (e *Encoder) Error(ctx context.Context, w http.ResponseWriter, err error, statusCode int) {
	e.logger.Error(fmt.Errorf("handling error"), map[string]interface{}{"error": err.Error(), "status": statusCode})

	resp := errorResponse{
		Message: err.Error(),
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		e.logger.Error(fmt.Errorf("error encoding response"), map[string]interface{}{"error": err.Error()})
	}
}

// RedirectResponse sends a redirect response with the appropriate status code.
func (e *Encoder) RedirectResponse(ctx context.Context, w http.ResponseWriter, location string, statusCode int) {
	if location == "" {
		e.logger.Error(fmt.Errorf("empty redirect location"), nil)
		http.Error(w, "Redirect location missing", http.StatusInternalServerError)
		return
	}

	e.logger.Info("Redirecting", map[string]interface{}{
		"location": location,
		"status":   statusCode,
	})

	w.Header().Set("Location", location)
	w.WriteHeader(statusCode)
}
