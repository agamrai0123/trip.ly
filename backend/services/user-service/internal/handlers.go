package internal

import (
	"errors"

	pkgerr "github.com/agamrai0123/wanderplan/pkg/errors"
	"github.com/agamrai0123/wanderplan/pkg/middleware"
	"github.com/agamrai0123/wanderplan/pkg/response"
	"github.com/gin-gonic/gin"
)

// Handlers wires HTTP requests to UserService methods.
type Handlers struct {
	svc *UserService
}

// NewHandlers constructs the handler set.
func NewHandlers(svc *UserService) *Handlers { return &Handlers{svc: svc} }

// GetMe handles GET /users/me
func (h *Handlers) GetMe(c *gin.Context) {
	user, err := h.svc.GetMe(c.Request.Context(), middleware.UserID(c))
	if errors.Is(err, ErrNotFound) {
		response.Err(c, pkgerr.NotFound("user not found"))
		return
	}
	if err != nil {
		response.Err(c, pkgerr.Internal(err.Error()))
		return
	}
	response.OK(c, user)
}

// UpdateMe handles PATCH /users/me
func (h *Handlers) UpdateMe(c *gin.Context) {
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, pkgerr.BadRequest(err.Error()))
		return
	}
	user, err := h.svc.UpdateMe(c.Request.Context(), middleware.UserID(c), &req)
	if errors.Is(err, ErrNotFound) {
		response.Err(c, pkgerr.NotFound("user not found"))
		return
	}
	if err != nil {
		response.Err(c, pkgerr.Internal(err.Error()))
		return
	}
	response.OK(c, user)
}

// GetMyTrips handles GET /users/me/trips — delegates to trip-service via gRPC
func (h *Handlers) GetMyTrips(c *gin.Context) {
	resp, err := h.svc.GetMyTrips(c.Request.Context(), middleware.UserID(c))
	if err != nil {
		response.Err(c, pkgerr.Internal(err.Error()))
		return
	}
	response.OK(c, resp.Trips)
}

// GetMyStats handles GET /users/me/stats
func (h *Handlers) GetMyStats(c *gin.Context) {
	resp, err := h.svc.GetMyStats(c.Request.Context(), middleware.UserID(c))
	if err != nil {
		response.Err(c, pkgerr.Internal(err.Error()))
		return
	}
	response.OK(c, resp)
}
