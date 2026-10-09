package apperror

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrorCode defines application-specific error codes for machine parsing.
type ErrorCode string

const (
	CodeBadRequest   ErrorCode = "BAD_REQUEST"
	CodeUnauthorized ErrorCode = "UNAUTHORIZED"
	CodeForbidden    ErrorCode = "FORBIDDEN"
	CodeNotFound     ErrorCode = "NOT_FOUND"
	CodeConflict     ErrorCode = "CONFLICT"
	CodeInternal     ErrorCode = "INTERNAL_SERVER_ERROR"
	CodeValidation   ErrorCode = "VALIDATION_ERROR"
)

// AppError represents a domain or application error with HTTP metadata.
type AppError struct {
	StatusCode int       `json:"-"`
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	Details    any       `json:"details,omitempty"`
	RawError   error     `json:"-"`
}

func (e *AppError) Error() string {
	if e.RawError != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.RawError)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.RawError
}

// New creates a generic AppError.
func New(statusCode int, code ErrorCode, message string, rawErr error) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		RawError:   rawErr,
	}
}

// NewBadRequest returns an HTTP 400 Bad Request error.
func NewBadRequest(message string, rawErr error) *AppError {
	return New(http.StatusBadRequest, CodeBadRequest, message, rawErr)
}

// NewUnauthorized returns an HTTP 401 Unauthorized error.
func NewUnauthorized(message string, rawErr error) *AppError {
	return New(http.StatusUnauthorized, CodeUnauthorized, message, rawErr)
}

// NewForbidden returns an HTTP 403 Forbidden error.
func NewForbidden(message string, rawErr error) *AppError {
	return New(http.StatusForbidden, CodeForbidden, message, rawErr)
}

// NewNotFound returns an HTTP 404 Not Found error.
func NewNotFound(message string, rawErr error) *AppError {
	return New(http.StatusNotFound, CodeNotFound, message, rawErr)
}

// NewConflict returns an HTTP 409 Conflict error.
func NewConflict(message string, rawErr error) *AppError {
	return New(http.StatusConflict, CodeConflict, message, rawErr)
}

// NewInternal returns an HTTP 500 Internal Server error.
func NewInternal(message string, rawErr error) *AppError {
	return New(http.StatusInternalServerError, CodeInternal, message, rawErr)
}

// NewValidationError returns an HTTP 422 Unprocessable Entity error with detailed field errors.
func NewValidationError(message string, details any) *AppError {
	return &AppError{
		StatusCode: http.StatusUnprocessableEntity,
		Code:       CodeValidation,
		Message:    message,
		Details:    details,
	}
}

// FromError extracts an AppError from an arbitrary error, or defaults to 500 Internal Server Error.
func FromError(err error) *AppError {
	if err == nil {
		return nil
	}

	if appErr, ok := errors.AsType[*AppError](err); ok {
		return appErr
	}

	return NewInternal("Internal server error", err)
}
