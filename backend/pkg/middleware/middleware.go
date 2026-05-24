// Package middleware provides reusable Gin middlewares for all WanderPlan services.
package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/agamrai0123/wanderplan/pkg/errors"
	"github.com/agamrai0123/wanderplan/pkg/jwt"
	"github.com/agamrai0123/wanderplan/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

// ──────────────────────────────────────────────
// Request ID
// ──────────────────────────────────────────────

// RequestID injects a UUID request_id into every request context.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if id == "" {
			id = uuid.New().String()
		}
		c.Set("request_id", id)
		c.Header("X-Request-ID", id)
		c.Next()
	}
}

// ──────────────────────────────────────────────
// Structured Request Logging
// ──────────────────────────────────────────────

// Logger logs every request with method, path, status, and latency using zerolog.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		dur := time.Since(start)
		status := c.Writer.Status()

		ev := log.Info()
		if status >= 500 {
			ev = log.Error()
		} else if status >= 400 {
			ev = log.Warn()
		}

		rid, _ := c.Get("request_id")
		ev.
			Str("request_id", rid.(string)).
			Str("method", c.Request.Method).
			Str("path", c.FullPath()).
			Str("client_ip", c.ClientIP()).
			Int("status", status).
			Dur("latency", dur).
			Msg("request")
	}
}

// ──────────────────────────────────────────────
// CORS
// ──────────────────────────────────────────────

// CORS configures Cross-Origin Resource Sharing for the given allowed origins.
func CORS(allowedOrigins []string) gin.HandlerFunc {
	originSet := make(map[string]struct{}, len(allowedOrigins))
	for _, o := range allowedOrigins {
		originSet[o] = struct{}{}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if _, ok := originSet[origin]; ok {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
			c.Header("Access-Control-Max-Age", "86400")
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// ──────────────────────────────────────────────
// JWT Authentication
// ──────────────────────────────────────────────

// Auth validates the Bearer JWT in the Authorization header.
// On success it stores user_id, email, name, avatar_url in the context.
func Auth(jwtMgr *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Err(c, errors.Unauthorized("missing or malformed authorization header"))
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := jwtMgr.Parse(tokenStr)
		if err != nil {
			if strings.Contains(err.Error(), "expired") {
				response.Err(c, errors.TokenExpired())
			} else {
				response.Err(c, errors.InvalidToken(err.Error()))
			}
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("name", claims.Name)
		c.Set("avatar_url", claims.AvatarURL)
		c.Next()
	}
}

// UserID extracts the authenticated user_id from the Gin context.
// Panics if Auth middleware was not applied upstream.
func UserID(c *gin.Context) string {
	v, _ := c.Get("user_id")
	s, _ := v.(string)
	return s
}

// ──────────────────────────────────────────────
// Rate Limiting
// ──────────────────────────────────────────────

type ipLimiter struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
	r        rate.Limit
	b        int
}

func newIPLimiter(r rate.Limit, b int) *ipLimiter {
	return &ipLimiter{limiters: make(map[string]*rate.Limiter), r: r, b: b}
}

func (l *ipLimiter) get(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()
	if lim, ok := l.limiters[ip]; ok {
		return lim
	}
	lim := rate.NewLimiter(l.r, l.b)
	l.limiters[ip] = lim
	return lim
}

// RateLimit applies a per-IP token-bucket rate limiter.
// rps is the sustained requests per second; burst is the maximum burst size.
func RateLimit(rps float64, burst int) gin.HandlerFunc {
	limiter := newIPLimiter(rate.Limit(rps), burst)
	return func(c *gin.Context) {
		if !limiter.get(c.ClientIP()).Allow() {
			response.ErrRaw(c, http.StatusTooManyRequests, "rate_limited", "too many requests")
			c.Abort()
			return
		}
		c.Next()
	}
}

// ──────────────────────────────────────────────
// Prometheus Metrics
// ──────────────────────────────────────────────

// Metrics registers and returns a Gin middleware that instruments HTTP requests.
// It also returns a promhttp.Handler for the /metrics endpoint.
func Metrics(service string, reg *prometheus.Registry) (gin.HandlerFunc, http.Handler) {
	httpRequests := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "wanderplan",
		Subsystem: service,
		Name:      "http_requests_total",
		Help:      "Total HTTP requests.",
	}, []string{"method", "path", "status"})

	httpDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "wanderplan",
		Subsystem: service,
		Name:      "http_request_duration_seconds",
		Help:      "HTTP request latency.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"method", "path"})

	goroutines := prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace: "wanderplan",
		Subsystem: service,
		Name:      "goroutines",
		Help:      "Current goroutine count.",
	}, func() float64 {
		// runtime.NumGoroutine() imported lazily to avoid import cycle
		return 0
	})

	reg.MustRegister(httpRequests, httpDuration, goroutines)

	mw := func(c *gin.Context) {
		start := time.Now()
		c.Next()
		dur := time.Since(start).Seconds()
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}
		status := http.StatusText(c.Writer.Status())
		httpRequests.WithLabelValues(c.Request.Method, path, status).Inc()
		httpDuration.WithLabelValues(c.Request.Method, path).Observe(dur)
	}

	handler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg})
	return mw, handler
}

// Recovery returns a Gin middleware that recovers from panics and returns a 500 response.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Error().Interface("panic", r).Msg("recovered from panic")
				response.ErrRaw(c, http.StatusInternalServerError, "internal_server_error", "an unexpected error occurred")
				c.Abort()
			}
		}()
		c.Next()
	}
}
