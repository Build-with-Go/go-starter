package errors

import (
	"fmt"
	"net/http"
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

// Common error constructors
func NewBadRequestError(message string, err error) *AppError {
	return NewAppError(http.StatusBadRequest, message, "", err)
}

func NewUnauthorizedError(message string) *AppError {
	return NewAppError(http.StatusUnauthorized, message, "", nil)
}

func NewForbiddenError(message string) *AppError {
	return NewAppError(http.StatusForbidden, message, "", nil)
}

func NewNotFoundError(message string) *AppError {
	return NewAppError(http.StatusNotFound, message, "", nil)
}

func NewConflictError(message string, err error) *AppError {
	return NewAppError(http.StatusConflict, message, "", err)
}

func NewInternalServerError(message string, err error) *AppError {
	return NewAppError(http.StatusInternalServerError, message, "", err)
}

func NewServiceUnavailableError(message string, err error) *AppError {
	return NewAppError(http.StatusServiceUnavailable, message, "", err)
}

// Validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// Error collection for multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		Errors: make([]ValidationError, 0),
	}
}

func (ve *ValidationErrors) Add(field, message string) {
	ve.Errors = append(ve.Errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

func (ve *ValidationErrors) HasErrors() bool {
	return len(ve.Errors) > 0
}

func (ve *ValidationErrors) Error() string {
	if len(ve.Errors) == 0 {
		return "no validation errors"
	}
	return fmt.Sprintf("validation errors: %d issues found", len(ve.Errors))
}

func (ve *ValidationErrors) ToAppError() *AppError {
	return NewValidationError("validation", ve.Error()).ToAppError()
}

func (ve *ValidationError) ToAppError() *AppError {
	return NewAppError(http.StatusBadRequest, "Validation failed", ve.Error(), nil)
}
