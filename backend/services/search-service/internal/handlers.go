package internal

import (
	"strconv"

	pkgerr "github.com/agamrai0123/wanderplan/pkg/errors"
	"github.com/agamrai0123/wanderplan/pkg/response"
	"github.com/gin-gonic/gin"
)

// Handlers wires HTTP requests to SearchService.
type Handlers struct {
	svc *SearchService
}

// NewHandlers constructs the handler set.
func NewHandlers(svc *SearchService) *Handlers { return &Handlers{svc: svc} }

// SearchPlaces handles GET /search/places?q=&lat=&lng=
func (h *Handlers) SearchPlaces(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		response.Err(c, pkgerr.BadRequest("query parameter 'q' is required"))
		return
	}
	lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
	lng, _ := strconv.ParseFloat(c.Query("lng"), 64)

	places, err := h.svc.SearchPlaces(c.Request.Context(), q, lat, lng)
	if err != nil {
		response.Err(c, pkgerr.Internal(err.Error()))
		return
	}
	response.OK(c, places)
}

// SearchTrips handles GET /search/trips?q=
func (h *Handlers) SearchTrips(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		response.Err(c, pkgerr.BadRequest("query parameter 'q' is required"))
		return
	}
	trips, err := h.svc.SearchTrips(c.Request.Context(), q)
	if err != nil {
		response.Err(c, pkgerr.Internal(err.Error()))
		return
	}
	response.OK(c, trips)
}
