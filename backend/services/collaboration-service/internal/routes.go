package internal

import (
	"github.com/agamrai0123/wanderplan/pkg/jwt"
	"github.com/agamrai0123/wanderplan/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// RegisterRoutes mounts all collaboration-service routes on the engine.
func RegisterRoutes(r *gin.Engine, h *Handlers, jwtMgr *jwt.Manager, reg *prometheus.Registry, cfg *Config) {
	metricsMW, metricsHandler := middleware.Metrics("collaboration_service", reg)

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
			trips.GET("/:id/collaborators", h.ListCollaborators)
			trips.POST("/:id/collaborators", h.InviteCollaborator)
			trips.PATCH("/:id/collaborators/:userId", h.UpdateCollaborator)
			trips.DELETE("/:id/collaborators/:userId", h.RemoveCollaborator)
		}
	}
}
