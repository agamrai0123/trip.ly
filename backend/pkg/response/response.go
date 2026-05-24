// Package response provides the standard JSON envelope for all WanderPlan API responses.
//
// Every endpoint must use OK, Created, NoContent, or Err rather than c.JSON directly.
package response

import (
	"net/http"
	"time"

	"github.com/agamrai0123/wanderplan/pkg/errors"

	"github.com/gin-gonic/gin"
)

// Envelope is the standard top-level JSON wrapper for every response.
type Envelope struct {
	Data  any      `json:"data"`
	Error *ErrBody `json:"error"`
	Meta  Meta     `json:"meta"`
}

// ErrBody is the error section of the envelope.
type ErrBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// Meta carries per-request metadata.
type Meta struct {
	RequestID string `json:"request_id"`
	Timestamp string `json:"timestamp"`
}

func buildMeta(c *gin.Context) Meta {
	rid := ""
	if v, ok := c.Get("request_id"); ok {
		if s, ok := v.(string); ok {
			rid = s
		}
	}
	return Meta{RequestID: rid, Timestamp: time.Now().UTC().Format(time.RFC3339)}
}

// OK sends HTTP 200 with the supplied data payload.
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Envelope{Data: data, Meta: buildMeta(c)})
}

// Created sends HTTP 201 with the supplied data payload.
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Envelope{Data: data, Meta: buildMeta(c)})
}

// NoContent sends HTTP 204 with no body.
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Err sends an error response derived from an *errors.AppError.
func Err(c *gin.Context, err *errors.AppError) {
	c.JSON(err.HTTPStatus, Envelope{
		Error: &ErrBody{Code: string(err.Code), Message: err.Message, Detail: err.Detail},
		Meta:  buildMeta(c),
	})
}

// ErrRaw sends an error response with explicit status / code / message.
func ErrRaw(c *gin.Context, status int, code, msg string) {
	c.JSON(status, Envelope{
		Error: &ErrBody{Code: code, Message: msg},
		Meta:  buildMeta(c),
	})
}
