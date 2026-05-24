// Package errors provides typed error handling for all WanderPlan services.
// Every handler must return an *AppError rather than a raw error.
package errors

import (
	"fmt"
	"net/http"
)

// Code is a machine-readable error identifier sent to clients.
type Code string

const (
	CodeBadRequest   Code = "bad_request"
	CodeUnauthorized Code = "unauthorized"
	CodeForbidden    Code = "forbidden"
	CodeNotFound     Code = "not_found"
	CodeConflict     Code = "conflict"
	CodeValidation   Code = "validation_failed"
	CodeInternal     Code = "internal_server_error"
	CodeDatabase     Code = "database_error"
	CodeUnavailable  Code = "service_unavailable"
	CodeInvalidToken Code = "invalid_token"
	CodeTokenExpired Code = "token_expired"
)

// AppError is the standard error type returned by all WanderPlan handlers.
type AppError struct {
	// Code is the machine-readable error code.
	Code Code `json:"code"`
	// Message is a human-readable description.
	Message string `json:"message"`
	// HTTPStatus is the HTTP status code (not serialised).
	HTTPStatus int `json:"-"`
	// Detail carries optional context shown to the caller.
	Detail string `json:"detail,omitempty"`
	cause  error
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap supports errors.Is / errors.As traversal.
func (e *AppError) Unwrap() error { return e.cause }

// Cause returns the underlying cause, if any.
func (e *AppError) Cause() error { return e.cause }

// WithCause attaches an underlying error for logging purposes.
func (e *AppError) WithCause(err error) *AppError {
	cp := *e
	cp.cause = err
	return &cp
}

// WithDetail adds a caller-visible detail string.
func (e *AppError) WithDetail(d string) *AppError {
	cp := *e
	cp.Detail = d
	return &cp
}

// New creates an AppError with the given code, message, and HTTP status.
func New(code Code, msg string, status int) *AppError {
	return &AppError{Code: code, Message: msg, HTTPStatus: status}
}

// BadRequest returns a 400 AppError.
func BadRequest(msg string) *AppError {
	return New(CodeBadRequest, msg, http.StatusBadRequest)
}

// Unauthorized returns a 401 AppError.
func Unauthorized(msg string) *AppError {
	return New(CodeUnauthorized, msg, http.StatusUnauthorized)
}

// Forbidden returns a 403 AppError.
func Forbidden(msg string) *AppError {
	return New(CodeForbidden, msg, http.StatusForbidden)
}

// NotFound returns a 404 AppError.
func NotFound(msg string) *AppError {
	return New(CodeNotFound, msg, http.StatusNotFound)
}

// Conflict returns a 409 AppError.
func Conflict(msg string) *AppError {
	return New(CodeConflict, msg, http.StatusConflict)
}

// Validation returns a 422 AppError.
func Validation(msg string) *AppError {
	return New(CodeValidation, msg, http.StatusUnprocessableEntity)
}

// Internal returns a 500 AppError.
func Internal(msg string) *AppError {
	return New(CodeInternal, msg, http.StatusInternalServerError)
}

// Database wraps a DB error into a 500 AppError.
func Database(err error) *AppError {
	return New(CodeDatabase, "a database error occurred", http.StatusInternalServerError).WithCause(err)
}

// Unavailable returns a 503 AppError.
func Unavailable(msg string) *AppError {
	return New(CodeUnavailable, msg, http.StatusServiceUnavailable)
}

// InvalidToken returns a 401 AppError for bad JWTs.
func InvalidToken(msg string) *AppError {
	return New(CodeInvalidToken, msg, http.StatusUnauthorized)
}

// TokenExpired returns a 401 AppError for expired JWTs.
func TokenExpired() *AppError {
	return New(CodeTokenExpired, "access token has expired", http.StatusUnauthorized)
}
