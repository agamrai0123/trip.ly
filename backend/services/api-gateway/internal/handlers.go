package internal

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	pkgresp "github.com/agamrai0123/wanderplan/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Handlers holds references shared across HTTP handlers.
type Handlers struct {
	Validator *AuthValidator
	Targets   map[string]string // prefix → base URL, e.g. "/api/v1/trips" → "http://localhost:8082"
}

// NewHandlers constructs Handlers.
func NewHandlers(validator *AuthValidator, cfg ServicesCfg) *Handlers {
	targets := map[string]string{
		"/auth":                 "http://" + serviceAddr(cfg.AuthAddr, 8081),
		"/api/v1/trips":         "http://" + serviceAddr(cfg.TripAddr, 8082),
		"/api/v1/users":         "http://" + serviceAddr(cfg.UserAddr, 8083),
		"/api/v1/collaborators": "http://" + serviceAddr(cfg.CollaborationAddr, 8084),
		"/api/v1/notifications": "http://" + serviceAddr(cfg.NotificationAddr, 8085),
		"/api/v1/search":        "http://" + serviceAddr(cfg.SearchAddr, 8086),
	}
	return &Handlers{Validator: validator, Targets: targets}
}

func serviceAddr(addr string, defaultPort int) string {
	if addr != "" {
		return addr
	}
	return fmt.Sprintf("localhost:%d", defaultPort)
}

// Proxy returns a Gin handler that reverse-proxies the request to the target.
func (h *Handlers) Proxy(targetBase string) gin.HandlerFunc {
	base, err := url.Parse(targetBase)
	if err != nil {
		panic("bad proxy target: " + targetBase)
	}
	proxy := httputil.NewSingleHostReverseProxy(base)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Error().Err(err).Str("target", targetBase).Msg("proxy error")
		w.WriteHeader(http.StatusBadGateway)
	}
	return func(c *gin.Context) {
		// Strip gin wildcard prefix so that the full path is forwarded correctly.
		c.Request.URL.Path = c.Param("path")
		if c.Request.URL.Path == "" {
			c.Request.URL.Path = "/"
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// Health responds to liveness probes.
func Health(c *gin.Context) {
	pkgresp.OK(c, gin.H{"status": "ok"})
}

// findTarget returns the longest prefix match for the given path.
func (h *Handlers) findTarget(path string) (string, bool) {
	best := ""
	bestAddr := ""
	for prefix, addr := range h.Targets {
		if strings.HasPrefix(path, prefix) && len(prefix) > len(best) {
			best = prefix
			bestAddr = addr
		}
	}
	if best == "" {
		return "", false
	}
	return bestAddr, true
}
