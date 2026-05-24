package internal

import (
	"errors"
	"net/http"

	pkgerr "github.com/agamrai0123/wanderplan/pkg/errors"
	"github.com/agamrai0123/wanderplan/pkg/middleware"
	"github.com/agamrai0123/wanderplan/pkg/response"
	"github.com/gin-gonic/gin"
)

// Handlers wires HTTP requests to TripService methods.
type Handlers struct {
	svc *TripService
}

// NewHandlers constructs the handler set.
func NewHandlers(svc *TripService) *Handlers { return &Handlers{svc: svc} }

// ListTrips handles GET /trips
func (h *Handlers) ListTrips(c *gin.Context) {
	trips, err := h.svc.ListTrips(c.Request.Context(), middleware.UserID(c))
	if err != nil {
		response.Err(c, pkgerr.Internal(err.Error()))
		return
	}
	response.OK(c, trips)
}

// GetTrip handles GET /trips/:id
func (h *Handlers) GetTrip(c *gin.Context) {
	trip, err := h.svc.GetTrip(c.Request.Context(), c.Param("id"), middleware.UserID(c))
	if errors.Is(err, ErrNotFound) {
		response.Err(c, pkgerr.NotFound("trip not found"))
		return
	}
	if err != nil {
		response.Err(c, pkgerr.Internal(err.Error()))
		return
	}
	response.OK(c, trip)
}

// CreateTrip handles POST /trips
func (h *Handlers) CreateTrip(c *gin.Context) {
	var req CreateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, pkgerr.BadRequest(err.Error()))
		return
	}
	trip, err := h.svc.CreateTrip(c.Request.Context(), middleware.UserID(c), &req)
	if err != nil {
		response.Err(c, pkgerr.Internal(err.Error()))
		return
	}
	response.Created(c, trip)
}

// UpdateTrip handles PATCH /trips/:id
func (h *Handlers) UpdateTrip(c *gin.Context) {
	var req UpdateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, pkgerr.BadRequest(err.Error()))
		return
	}
	trip, err := h.svc.UpdateTrip(c.Request.Context(), c.Param("id"), middleware.UserID(c), &req)
	if errors.Is(err, ErrNotFound) {
		response.Err(c, pkgerr.NotFound("trip not found"))
		return
	}
	if err != nil {
		response.Err(c, pkgerr.Internal(err.Error()))
		return
	}
	response.OK(c, trip)
}

// DeleteTrip handles DELETE /trips/:id
func (h *Handlers) DeleteTrip(c *gin.Context) {
	err := h.svc.DeleteTrip(c.Request.Context(), c.Param("id"), middleware.UserID(c))
	if errors.Is(err, ErrNotFound) {
		response.Err(c, pkgerr.NotFound("trip not found"))
		return
	}
	if err != nil {
		response.Err(c, pkgerr.Internal(err.Error()))
		return
	}
	response.NoContent(c)
}

// DuplicateTrip handles POST /trips/:id/duplicate
func (h *Handlers) DuplicateTrip(c *gin.Context) {
	trip, err := h.svc.DuplicateTrip(c.Request.Context(), c.Param("id"), middleware.UserID(c))
	if errors.Is(err, ErrNotFound) {
		response.Err(c, pkgerr.NotFound("trip not found"))
		return
	}
	if err != nil {
		response.Err(c, pkgerr.Internal(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": trip})
}

// ── Day handlers ──────────────────────────────────────────────

// ListDays handles GET /trips/:id/days
func (h *Handlers) ListDays(c *gin.Context) {
	days, err := h.svc.ListDays(c.Request.Context(), c.Param("id"), middleware.UserID(c))
	if errors.Is(err, ErrNotFound) {
		response.Err(c, pkgerr.NotFound("trip not found"))
		return
	}
	if err != nil {
		response.Err(c, pkgerr.Internal(err.Error()))
		return
	}
	response.OK(c, days)
}

// CreateDay handles POST /trips/:id/days
func (h *Handlers) CreateDay(c *gin.Context) {
	var req CreateDayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, pkgerr.BadRequest(err.Error()))
		return
	}
	day, err := h.svc.CreateDay(c.Request.Context(), c.Param("id"), middleware.UserID(c), &req)
	if errors.Is(err, ErrNotFound) {
		response.Err(c, pkgerr.NotFound("trip not found"))
		return
	}
	if err != nil {
		response.Err(c, pkgerr.Internal(err.Error()))
		return
	}
	response.Created(c, day)
}

// UpdateDay handles PATCH /trips/:id/days/:dayId
func (h *Handlers) UpdateDay(c *gin.Context) {
	var req UpdateDayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, pkgerr.BadRequest(err.Error()))
		return
	}
	day, err := h.svc.UpdateDay(c.Request.Context(), c.Param("id"), c.Param("dayId"), middleware.UserID(c), &req)
	if errors.Is(err, ErrNotFound) {
		response.Err(c, pkgerr.NotFound("day not found"))
		return
	}
	if err != nil {
		response.Err(c, pkgerr.Internal(err.Error()))
		return
	}
	response.OK(c, day)
}

// DeleteDay handles DELETE /trips/:id/days/:dayId
func (h *Handlers) DeleteDay(c *gin.Context) {
	err := h.svc.DeleteDay(c.Request.Context(), c.Param("id"), c.Param("dayId"), middleware.UserID(c))
	if errors.Is(err, ErrNotFound) {
		response.Err(c, pkgerr.NotFound("day not found"))
		return
	}
	if err != nil {
		response.Err(c, pkgerr.Internal(err.Error()))
		return
	}
	response.NoContent(c)
}

// ── Item handlers ─────────────────────────────────────────────

// ListItems handles GET /trips/:id/days/:dayId/items
func (h *Handlers) ListItems(c *gin.Context) {
	items, err := h.svc.ListItems(c.Request.Context(), c.Param("id"), c.Param("dayId"), middleware.UserID(c))
	if errors.Is(err, ErrNotFound) {
		response.Err(c, pkgerr.NotFound("trip or day not found"))
		return
	}
	if err != nil {
		response.Err(c, pkgerr.Internal(err.Error()))
		return
	}
	response.OK(c, items)
}

// CreateItem handles POST /trips/:id/days/:dayId/items
func (h *Handlers) CreateItem(c *gin.Context) {
	var req CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, pkgerr.BadRequest(err.Error()))
		return
	}
	item, err := h.svc.CreateItem(c.Request.Context(), c.Param("id"), c.Param("dayId"), middleware.UserID(c), &req)
	if errors.Is(err, ErrNotFound) {
		response.Err(c, pkgerr.NotFound("trip or day not found"))
		return
	}
	if err != nil {
		response.Err(c, pkgerr.Internal(err.Error()))
		return
	}
	response.Created(c, item)
}

// UpdateItem handles PATCH /trips/:id/items/:itemId
func (h *Handlers) UpdateItem(c *gin.Context) {
	var req UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, pkgerr.BadRequest(err.Error()))
		return
	}
	item, err := h.svc.UpdateItem(c.Request.Context(), c.Param("id"), c.Param("itemId"), middleware.UserID(c), &req)
	if errors.Is(err, ErrNotFound) {
		response.Err(c, pkgerr.NotFound("item not found"))
		return
	}
	if err != nil {
		response.Err(c, pkgerr.Internal(err.Error()))
		return
	}
	response.OK(c, item)
}

// DeleteItem handles DELETE /trips/:id/items/:itemId
func (h *Handlers) DeleteItem(c *gin.Context) {
	err := h.svc.DeleteItem(c.Request.Context(), c.Param("id"), c.Param("itemId"), middleware.UserID(c))
	if errors.Is(err, ErrNotFound) {
		response.Err(c, pkgerr.NotFound("item not found"))
		return
	}
	if err != nil {
		response.Err(c, pkgerr.Internal(err.Error()))
		return
	}
	response.NoContent(c)
}

// ReorderItems handles PATCH /trips/:id/items/reorder
func (h *Handlers) ReorderItems(c *gin.Context) {
	var req ReorderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, pkgerr.BadRequest(err.Error()))
		return
	}
	if err := h.svc.ReorderItems(c.Request.Context(), c.Param("id"), middleware.UserID(c), &req); err != nil {
		if errors.Is(err, ErrNotFound) {
			response.Err(c, pkgerr.NotFound("trip not found"))
			return
		}
		response.Err(c, pkgerr.Internal(err.Error()))
		return
	}
	response.NoContent(c)
}
