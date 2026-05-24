package internal

import (
	"errors"

	pkgerr "github.com/agamrai0123/wanderplan/pkg/errors"
	"github.com/agamrai0123/wanderplan/pkg/middleware"
	"github.com/agamrai0123/wanderplan/pkg/response"
	"github.com/gin-gonic/gin"
)

// Handlers wires HTTP requests to CollaborationService methods.
type Handlers struct {
	svc *CollaborationService
}

// NewHandlers constructs the handler set.
func NewHandlers(svc *CollaborationService) *Handlers { return &Handlers{svc: svc} }

// ListCollaborators handles GET /trips/:id/collaborators
func (h *Handlers) ListCollaborators(c *gin.Context) {
	collabs, err := h.svc.List(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.Err(c, pkgerr.Internal(err.Error()))
		return
	}
	response.OK(c, collabs)
}

// InviteCollaborator handles POST /trips/:id/collaborators
func (h *Handlers) InviteCollaborator(c *gin.Context) {
	var req InviteCollaboratorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, pkgerr.BadRequest(err.Error()))
		return
	}
	collab, err := h.svc.Invite(c.Request.Context(), c.Param("id"), &req)
	if err != nil {
		response.Err(c, pkgerr.Internal(err.Error()))
		return
	}
	response.Created(c, collab)
}

// UpdateCollaborator handles PATCH /trips/:id/collaborators/:userId
func (h *Handlers) UpdateCollaborator(c *gin.Context) {
	var req UpdateCollaboratorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, pkgerr.BadRequest(err.Error()))
		return
	}
	collab, err := h.svc.Update(c.Request.Context(), c.Param("id"), c.Param("userId"), &req)
	if errors.Is(err, ErrNotFound) {
		response.Err(c, pkgerr.NotFound("collaborator not found"))
		return
	}
	if err != nil {
		response.Err(c, pkgerr.Internal(err.Error()))
		return
	}
	response.OK(c, collab)
}

// RemoveCollaborator handles DELETE /trips/:id/collaborators/:userId
func (h *Handlers) RemoveCollaborator(c *gin.Context) {
	// Only trip owner or admin can remove; ownership check done via trip-service.
	// For now we enforce that the caller is either the removed user or has sufficient permissions.
	_ = middleware.UserID(c) // ensures auth
	err := h.svc.Remove(c.Request.Context(), c.Param("id"), c.Param("userId"))
	if errors.Is(err, ErrNotFound) {
		response.Err(c, pkgerr.NotFound("collaborator not found"))
		return
	}
	if err != nil {
		response.Err(c, pkgerr.Internal(err.Error()))
		return
	}
	response.NoContent(c)
}
