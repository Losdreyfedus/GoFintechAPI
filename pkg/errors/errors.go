package errors

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/trace"
)

// ErrorCode represents application error codes
type ErrorCode string

const (
	// Authentication errors
	ErrorCodeUnauthorized       ErrorCode = "UNAUTHORIZED"
	ErrorCodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
	ErrorCodeTokenExpired       ErrorCode = "TOKEN_EXPIRED"
	ErrorCodeInsufficientRole   ErrorCode = "INSUFFICIENT_ROLE"

	// Validation errors
	ErrorCodeValidationFailed ErrorCode = "VALIDATION_FAILED"
	ErrorCodeInvalidInput     ErrorCode = "INVALID_INPUT"
	ErrorCodeMissingRequired  ErrorCode = "MISSING_REQUIRED"

	// Business logic errors
	ErrorCodeInsufficientBalance ErrorCode = "INSUFFICIENT_BALANCE"
	ErrorCodeUserNotFound        ErrorCode = "USER_NOT_FOUND"
	ErrorCodeTransactionFailed   ErrorCode = "TRANSACTION_FAILED"
	ErrorCodeDuplicateResource   ErrorCode = "DUPLICATE_RESOURCE"

	// System errors
	ErrorCodeInternalError        ErrorCode = "INTERNAL_ERROR"
	ErrorCodeDatabaseError        ErrorCode = "DATABASE_ERROR"
	ErrorCodeExternalServiceError ErrorCode = "EXTERNAL_SERVICE_ERROR"
	ErrorCodeRateLimitExceeded    ErrorCode = "RATE_LIMIT_EXCEEDED"
)

// AppError represents an application error
type AppError struct {
	Code      ErrorCode              `json:"code"`
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details,omitempty"`
	TraceID   string                 `json:"trace_id,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	HTTPCode  int                    `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewAppError creates a new application error
func NewAppError(code ErrorCode, message string, httpCode int) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		HTTPCode:  httpCode,
		Timestamp: time.Now(),
	}
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}

// WithTraceID adds trace ID to the error
func (e *AppError) WithTraceID(traceID string) *AppError {
	e.TraceID = traceID
	return e
}

// WriteError writes an error response to HTTP response writer
func WriteError(w http.ResponseWriter, err error, ctx context.Context) {
	var appErr *AppError

	// Convert to AppError if it's not already
	if convertedErr, ok := err.(*AppError); ok {
		appErr = convertedErr
	} else {
		appErr = NewAppError(ErrorCodeInternalError, err.Error(), http.StatusInternalServerError)
	}

	// Add trace ID if available
	if span := trace.SpanFromContext(ctx); span != nil {
		appErr.TraceID = span.SpanContext().TraceID().String()
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.HTTPCode)

	// Write error response
	json.NewEncoder(w).Encode(appErr)
}

// Common error constructors
func Unauthorized(message string) *AppError {
	return NewAppError(ErrorCodeUnauthorized, message, http.StatusUnauthorized)
}

func ValidationFailed(message string, details map[string]interface{}) *AppError {
	return NewAppError(ErrorCodeValidationFailed, message, http.StatusBadRequest).WithDetails(details)
}

func NotFound(message string) *AppError {
	return NewAppError(ErrorCodeUserNotFound, message, http.StatusNotFound)
}

func BadRequest(message string) *AppError {
	return NewAppError(ErrorCodeInvalidInput, message, http.StatusBadRequest)
}

func InternalError(message string) *AppError {
	return NewAppError(ErrorCodeInternalError, message, http.StatusInternalServerError)
}

func InsufficientBalance(message string) *AppError {
	return NewAppError(ErrorCodeInsufficientBalance, message, http.StatusBadRequest)
}

func RateLimitExceeded(message string) *AppError {
	return NewAppError(ErrorCodeRateLimitExceeded, message, http.StatusTooManyRequests)
}
