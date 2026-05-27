package internal

import (
	"github.com/agamrai0123/wanderplan/pkg/jwt"
	"github.com/agamrai0123/wanderplan/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// RegisterRoutes mounts all user-service routes on the engine.
func RegisterRoutes(r *gin.Engine, h *Handlers, jwtMgr *jwt.Manager, reg *prometheus.Registry, cfg *Config) {
	metricsMW, metricsHandler := middleware.Metrics("user_service", reg)

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
		users := protected.Group("/users")
		{
			users.GET("/me", h.GetMe)
			users.PATCH("/me", h.UpdateMe)
			users.GET("/me/trips", h.GetMyTrips)
			users.GET("/me/stats", h.GetMyStats)
		}
	}
}
