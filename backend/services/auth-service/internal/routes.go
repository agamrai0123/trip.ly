package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/agamrai0123/wanderplan/pkg/jwt"
	"github.com/agamrai0123/wanderplan/pkg/middleware"
)

// RegisterRoutes mounts all auth-service routes on the engine.
func RegisterRoutes(r *gin.Engine, h *Handlers, jwtMgr *jwt.Manager, reg *prometheus.Registry, cfg *Config) {
	metricsMW, metricsHandler := middleware.Metrics("auth_service", reg)

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

	auth := r.Group("/auth")
	{
		auth.GET("/:provider/login", h.Login)
		auth.GET("/:provider/callback", h.Callback)
		auth.POST("/refresh", h.Refresh)
		auth.POST("/logout", middleware.Auth(jwtMgr), h.Logout)
		auth.GET("/me", middleware.Auth(jwtMgr), h.Me)
	}
}
