package internal

import (
	"github.com/agamrai0123/wanderplan/pkg/jwt"
	"github.com/agamrai0123/wanderplan/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// RegisterRoutes mounts all notification-service routes.
func RegisterRoutes(r *gin.Engine, h *Handlers, jwtMgr *jwt.Manager, reg *prometheus.Registry, cfg *Config) {
	metricsMW, metricsHandler := middleware.Metrics("notification_service", reg)

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
		protected.GET("/notifications", h.ListNotifications)
		protected.PATCH("/notifications/:id/read", h.MarkRead)
		protected.PATCH("/notifications/read-all", h.MarkAllRead)
		protected.GET("/ws/notifications", h.WebSocket)
	}
}
