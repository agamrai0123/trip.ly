package internal

import (
	"github.com/agamrai0123/wanderplan/pkg/jwt"
	"github.com/agamrai0123/wanderplan/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// RegisterRoutes mounts all trip-service routes on the engine.
func RegisterRoutes(r *gin.Engine, h *Handlers, jwtMgr *jwt.Manager, reg *prometheus.Registry, cfg *Config) {
	metricsMW, metricsHandler := middleware.Metrics("trip_service", reg)

	r.Use(
		middleware.Recovery(),
		middleware.RequestID(),
		middleware.Logger(),
		middleware.CORS(cfg.CORS.AllowedOrigins),
		middleware.RateLimit(cfg.RateLimit.RPS, cfg.RateLimit.Burst),
		metricsMW,
	)

	r.GET("/metrics", gin.WrapH(metricsHandler))
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	protected := r.Group("/", middleware.Auth(jwtMgr))
	{
		trips := protected.Group("/trips")
		{
			trips.GET("", h.ListTrips)
			trips.POST("", h.CreateTrip)
			trips.GET("/:id", h.GetTrip)
			trips.PATCH("/:id", h.UpdateTrip)
			trips.DELETE("/:id", h.DeleteTrip)
			trips.POST("/:id/duplicate", h.DuplicateTrip)

			// Days
			trips.GET("/:id/days", h.ListDays)
			trips.POST("/:id/days", h.CreateDay)
			trips.PATCH("/:id/days/:dayId", h.UpdateDay)
			trips.DELETE("/:id/days/:dayId", h.DeleteDay)

			// Items
			trips.GET("/:id/days/:dayId/items", h.ListItems)
			trips.POST("/:id/days/:dayId/items", h.CreateItem)
			trips.PATCH("/:id/items/reorder", h.ReorderItems)
			trips.PATCH("/:id/items/:itemId", h.UpdateItem)
			trips.DELETE("/:id/items/:itemId", h.DeleteItem)
		}
	}
}
