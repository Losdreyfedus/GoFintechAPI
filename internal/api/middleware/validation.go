package middleware

import (
	"encoding/json"
	"net/http"

	"backend_path/pkg/errors"
	"backend_path/pkg/validator"
)

// ValidationMiddleware validates request body using go-playground/validator
func ValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only validate POST/PUT/PATCH requests with JSON content
		if r.Method == "GET" || r.Method == "DELETE" {
			next.ServeHTTP(w, r)
			return
		}

		// Check content type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			errors.WriteError(w, errors.BadRequest("Content-Type must be application/json"), r.Context())
			return
		}

		// Continue to next handler (validation will be done in handlers)
		next.ServeHTTP(w, r)
	})
}

// ValidateRequest validates a request body against a struct
func ValidateRequest(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	// Decode JSON body
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		errors.WriteError(w, errors.BadRequest("Invalid JSON format"), r.Context())
		return false
	}

	// Validate struct
	if validationErrors := validator.Validate(v); len(validationErrors) > 0 {
		// Convert validation errors to details map
		details := make(map[string]interface{})
		for _, err := range validationErrors {
			details[err.Field] = err.Message
		}

		errors.WriteError(w, errors.ValidationFailed("Validation failed", details), r.Context())
		return false
	}

	return true
}
