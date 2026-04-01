// Package errors provides custom error types for the Go Starter application.
package errors

import (
	"fmt"
	"net/http"
	"strings"
)

// AppError represents an application error with HTTP status code
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Err     error  `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Details)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// HTTPStatus returns the HTTP status code
func (e *AppError) HTTPStatus() int {
	return e.Code
}

// NewAppError creates a new application error
func NewAppError(code int, message, details string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
		Err:     err,
	}
}

// NewBadRequestError creates a new bad request error
func NewBadRequestError(message string, err error) *AppError {
	return NewAppError(http.StatusBadRequest, message, "", err)
}

// NewUnauthorizedError creates a new unauthorized error
func NewUnauthorizedError(message string) *AppError {
	return NewAppError(http.StatusUnauthorized, message, "", nil)
}

// NewForbiddenError creates a new forbidden error
func NewForbiddenError(message string) *AppError {
	return NewAppError(http.StatusForbidden, message, "", nil)
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string) *AppError {
	return NewAppError(http.StatusNotFound, message, "", nil)
}

// NewConflictError creates a new conflict error
func NewConflictError(message string, err error) *AppError {
	return NewAppError(http.StatusConflict, message, "", err)
}

// NewInternalServerError creates a new internal server error
func NewInternalServerError(message string, err error) *AppError {
	return NewAppError(http.StatusInternalServerError, message, "", err)
}

// NewServiceUnavailableError creates a new service unavailable error
func NewServiceUnavailableError(message string, err error) *AppError {
	return NewAppError(http.StatusServiceUnavailable, message, "", err)
}

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

func (ve *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", ve.Field, ve.Message)
}

// ValidationErrors represents a collection of validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// NewValidationErrors creates a new validation errors collection
func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		Errors: make([]ValidationError, 0),
	}
}

// Add adds a validation error to the collection
func (ve *ValidationErrors) Add(field, message string) {
	ve.Errors = append(ve.Errors, *NewValidationError(field, message))
}

// HasErrors returns true if there are validation errors
func (ve *ValidationErrors) HasErrors() bool {
	return len(ve.Errors) > 0
}

func (ve *ValidationErrors) Error() string {
	if !ve.HasErrors() {
		return ""
	}

	var messages []string
	for _, err := range ve.Errors {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

// ToAppError converts validation errors to an AppError
func (ve *ValidationErrors) ToAppError() *AppError {
	return NewValidationError("validation", ve.Error()).ToAppError()
}

// ToAppError converts a validation error to an AppError
func (ve *ValidationError) ToAppError() *AppError {
	return NewAppError(http.StatusBadRequest, "Validation failed", ve.Error(), nil)
}
