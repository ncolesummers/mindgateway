package errors

import (
	"fmt"
	"net/http"
)

// Error represents a domain error with a status code and user-friendly message
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"error,omitempty"`
}

// Error returns the error message
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *Error) Unwrap() error {
	return e.Err
}

// Common errors
var (
	// Authentication errors
	ErrUnauthorized = &Error{Code: http.StatusUnauthorized, Message: "Unauthorized"}
	ErrForbidden    = &Error{Code: http.StatusForbidden, Message: "Forbidden"}
	
	// Input validation errors
	ErrBadRequest     = &Error{Code: http.StatusBadRequest, Message: "Bad request"}
	ErrInvalidInput   = &Error{Code: http.StatusBadRequest, Message: "Invalid input"}
	ErrMissingField   = &Error{Code: http.StatusBadRequest, Message: "Missing required field"}
	ErrMalformedInput = &Error{Code: http.StatusBadRequest, Message: "Malformed input"}
	
	// Resource errors
	ErrNotFound     = &Error{Code: http.StatusNotFound, Message: "Resource not found"}
	ErrAlreadyExists = &Error{Code: http.StatusConflict, Message: "Resource already exists"}
	
	// Service errors
	ErrInternal          = &Error{Code: http.StatusInternalServerError, Message: "Internal server error"}
	ErrServiceUnavailable = &Error{Code: http.StatusServiceUnavailable, Message: "Service unavailable"}
	ErrTimeout            = &Error{Code: http.StatusGatewayTimeout, Message: "Request timeout"}
	ErrRateLimited        = &Error{Code: http.StatusTooManyRequests, Message: "Rate limit exceeded"}
	
	// Worker errors
	ErrNoWorkersAvailable = &Error{Code: http.StatusServiceUnavailable, Message: "No workers available"}
	ErrWorkerNotFound     = &Error{Code: http.StatusNotFound, Message: "Worker not found"}
	ErrWorkerTimeout      = &Error{Code: http.StatusGatewayTimeout, Message: "Worker request timeout"}
)

// WithMessage adds context to a standard error
func WithMessage(err error, message string) *Error {
	if err == nil {
		return nil
	}
	
	// If the error is already an Error, use its code
	if e, ok := err.(*Error); ok {
		return &Error{
			Code:    e.Code,
			Message: message,
			Err:     e.Err,
		}
	}
	
	// Default to internal server error
	return &Error{
		Code:    http.StatusInternalServerError,
		Message: message,
		Err:     err,
	}
}

// New creates a new error with the given code and message
func New(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Wrap wraps an error with additional context
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}