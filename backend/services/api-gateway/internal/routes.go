package internal

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	pkgjwt "github.com/agamrai0123/wanderplan/pkg/jwt"
	pkgmw "github.com/agamrai0123/wanderplan/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

// RegisterRoutes attaches all routes to r.
func RegisterRoutes(
	r *gin.Engine,
	h *Handlers,
	jwtMgr *pkgjwt.Manager,
	reg *prometheus.Registry,
	cfg *Config,
) {
	metricsMW, metricsHandler := pkgmw.Metrics("api-gateway", reg)

	r.Use(
		pkgmw.RequestID(),
		pkgmw.Logger(),
		pkgmw.Recovery(),
		pkgmw.CORS(cfg.CORS.AllowedOrigins),
		pkgmw.RateLimit(cfg.RateLimit.RPS, cfg.RateLimit.Burst),
		metricsMW,
	)

	// Health & metrics (unauthenticated)
	r.GET("/healthz", Health)
	r.GET("/metrics", gin.WrapH(metricsHandler))

	// Auth routes — proxied to auth-service (no JWT required)
	authBase := "http://" + cfg.Services.AuthAddr
	if cfg.Services.AuthAddr == "" {
		authBase = "http://localhost:8081"
	}
	r.Any("/auth/*path", proxyTo(authBase))

	// Protected API routes — validate JWT, then proxy
	api := r.Group("/api/v1")
	api.Use(pkgmw.Auth(jwtMgr))
	{
		tripBase := fmt.Sprintf("http://%s", coalesce(cfg.Services.TripAddr, "localhost:8082"))
		userBase := fmt.Sprintf("http://%s", coalesce(cfg.Services.UserAddr, "localhost:8083"))
		collabBase := fmt.Sprintf("http://%s", coalesce(cfg.Services.CollaborationAddr, "localhost:8084"))
		notifBase := fmt.Sprintf("http://%s", coalesce(cfg.Services.NotificationAddr, "localhost:8085"))
		searchBase := fmt.Sprintf("http://%s", coalesce(cfg.Services.SearchAddr, "localhost:8086"))

		api.Any("/trips/*path", proxyTo(tripBase))
		api.Any("/users/*path", proxyTo(userBase))
		api.Any("/collaborators/*path", proxyTo(collabBase))
		api.Any("/notifications/*path", proxyTo(notifBase))
		api.Any("/search/*path", proxyTo(searchBase))
	}
}

// proxyTo returns a Gin handler that reverse-proxies requests preserving the full path.
func proxyTo(baseURL string) gin.HandlerFunc {
	base, err := url.Parse(baseURL)
	if err != nil {
		panic("bad proxy target: " + baseURL)
	}
	proxy := httputil.NewSingleHostReverseProxy(base)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
		log.Error().Err(e).Str("target", baseURL).Msg("proxy error")
		w.WriteHeader(http.StatusBadGateway)
	}
	return func(c *gin.Context) {
		path := c.Param("path")
		if path == "" {
			path = "/"
		}
		c.Request.URL.Path = path
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func coalesce(s, def string) string {
	if s != "" {
		return s
	}
	return def
}
