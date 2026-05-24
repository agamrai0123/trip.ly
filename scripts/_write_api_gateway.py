"""Write api-gateway service files."""
import os

BASE = r'D:\Learn\trip.ly\backend\services\api-gateway'

files = {}

# ──────────────────────────────────────────────────────────────
# config/api-gateway-config.json
# ──────────────────────────────────────────────────────────────
files['config/api-gateway-config.json'] = '''{
    "version": "1.0",
    "server_port": 8080,
    "grpc_port": 9080,
    "metric_port": 7080,
    "logging": {
        "level": 1,
        "path": "./log/api-gateway.log",
        "max_size_mb": 256,
        "max_backups": 5,
        "max_age_days": 30
    },
    "services": {
        "auth_addr":          "localhost:9081",
        "trip_addr":          "localhost:8082",
        "user_addr":          "localhost:8083",
        "collaboration_addr": "localhost:8084",
        "notification_addr":  "localhost:8085",
        "search_addr":        "localhost:8086"
    },
    "cors": {
        "allowed_origins": ["http://localhost:5173"]
    },
    "rate_limit": {
        "rps": 200,
        "burst": 400
    }
}
'''

# ──────────────────────────────────────────────────────────────
# internal/config.go
# ──────────────────────────────────────────────────────────────
files['internal/config.go'] = '''package internal

// Config is the top-level api-gateway configuration.
type Config struct {
\tVersion    string     `mapstructure:"version"`
\tServerPort int        `mapstructure:"server_port"`
\tGRPCPort   int        `mapstructure:"grpc_port"`
\tMetricPort int        `mapstructure:"metric_port"`
\tLogging    LoggingCfg `mapstructure:"logging"`
\tServices   ServicesCfg `mapstructure:"services"`
\tCORS       CORSCfg    `mapstructure:"cors"`
\tRateLimit  RateLimitCfg `mapstructure:"rate_limit"`
}

type LoggingCfg struct {
\tLevel      int    `mapstructure:"level"`
\tPath       string `mapstructure:"path"`
\tMaxSizeMB  int    `mapstructure:"max_size_mb"`
\tMaxBackups int    `mapstructure:"max_backups"`
\tMaxAgeDays int    `mapstructure:"max_age_days"`
}

type ServicesCfg struct {
\tAuthAddr          string `mapstructure:"auth_addr"`
\tTripAddr          string `mapstructure:"trip_addr"`
\tUserAddr          string `mapstructure:"user_addr"`
\tCollaborationAddr string `mapstructure:"collaboration_addr"`
\tNotificationAddr  string `mapstructure:"notification_addr"`
\tSearchAddr        string `mapstructure:"search_addr"`
}

type CORSCfg struct {
\tAllowedOrigins []string `mapstructure:"allowed_origins"`
}

type RateLimitCfg struct {
\tRPS   float64 `mapstructure:"rps"`
\tBurst int     `mapstructure:"burst"`
}

// Validate returns an error if required fields are missing.
func (c *Config) Validate() error {
\tif c.Services.AuthAddr == "" {
\t\tc.Services.AuthAddr = "localhost:9081"
\t}
\tif c.ServerPort == 0 {
\t\tc.ServerPort = 8080
\t}
\tif c.MetricPort == 0 {
\t\tc.MetricPort = 7080
\t}
\tif c.RateLimit.RPS == 0 {
\t\tc.RateLimit.RPS = 100
\t\tc.RateLimit.Burst = 200
\t}
\tif len(c.CORS.AllowedOrigins) == 0 {
\t\tc.CORS.AllowedOrigins = []string{"http://localhost:5173"}
\t}
\treturn nil
}
'''

# ──────────────────────────────────────────────────────────────
# internal/errors.go
# ──────────────────────────────────────────────────────────────
files['internal/errors.go'] = '''package internal

import pkgerr "github.com/agamrai0123/wanderplan/pkg/errors"

var (
\tBadRequest   = pkgerr.BadRequest
\tUnauthorized = pkgerr.Unauthorized
\tForbidden    = pkgerr.Forbidden
\tNotFound     = pkgerr.NotFound
\tInternal     = pkgerr.Internal
)
'''

# ──────────────────────────────────────────────────────────────
# internal/logger.go
# ──────────────────────────────────────────────────────────────
files['internal/logger.go'] = '''package internal

import (
\t"github.com/rs/zerolog"
\t"github.com/rs/zerolog/log"
\tpkglogger "github.com/agamrai0123/wanderplan/pkg/logger"
)

// InitLogger initialises the global zerolog logger for the api-gateway.
func InitLogger(cfg LoggingCfg, service string) zerolog.Logger {
\tlevel := zerolog.Level(cfg.Level)
\tif level < zerolog.TraceLevel || level > zerolog.Disabled {
\t\tlevel = zerolog.InfoLevel
\t}
\tmaxSizeMB := cfg.MaxSizeMB
\tif maxSizeMB <= 0 {
\t\tmaxSizeMB = 256
\t}
\tlogger := pkglogger.Init(pkglogger.Config{
\t\tLevel:      int(level),
\t\tFilePath:   cfg.Path,
\t\tMaxSizeMB:  maxSizeMB,
\t\tMaxBackups: cfg.MaxBackups,
\t\tMaxAgeDays: cfg.MaxAgeDays,
\t\tService:    service,
\t})
\tlog.Logger = logger
\tzerolog.DefaultContextLogger = &logger
\treturn logger
}
'''

# ──────────────────────────────────────────────────────────────
# internal/metrics.go
# ──────────────────────────────────────────────────────────────
files['internal/metrics.go'] = '''package internal

import (
\t"github.com/prometheus/client_golang/prometheus"
\t"github.com/prometheus/client_golang/prometheus/collectors"
)

// NewRegistry creates a Prometheus registry with standard process and Go metrics.
func NewRegistry() *prometheus.Registry {
\treg := prometheus.NewRegistry()
\treg.MustRegister(
\t\tcollectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
\t\tcollectors.NewGoCollector(),
\t)
\treturn reg
}
'''

# ──────────────────────────────────────────────────────────────
# internal/models.go
# ──────────────────────────────────────────────────────────────
files['internal/models.go'] = '''package internal

// ProxyTarget maps a route prefix to a downstream service address.
type ProxyTarget struct {
\tPrefix  string
\tAddress string
}
'''

# ──────────────────────────────────────────────────────────────
# internal/database.go
# ──────────────────────────────────────────────────────────────
files['internal/database.go'] = '''package internal
// The api-gateway does not own a database directly.
// It proxies requests to downstream services.
'''

# ──────────────────────────────────────────────────────────────
# internal/service.go
# ──────────────────────────────────────────────────────────────
files['internal/service.go'] = '''package internal

import (
\t"context"
\t"fmt"

\t"github.com/rs/zerolog/log"
\tpkgjwt "github.com/agamrai0123/wanderplan/pkg/jwt"
\tpbauth "github.com/agamrai0123/wanderplan/proto/gen/wanderplan/v1"
\t"google.golang.org/grpc"
\t"google.golang.org/grpc/credentials/insecure"
\t"google.golang.org/grpc/keepalive"
\t"time"
)

// AuthValidator calls auth-service over gRPC to validate a JWT.
type AuthValidator struct {
\tclient pbauth.AuthServiceClient
\tconn   *grpc.ClientConn
}

// NewAuthValidator dials the auth-service gRPC server.
func NewAuthValidator(authAddr string) (*AuthValidator, error) {
\tconn, err := grpc.NewClient(authAddr,
\t\tgrpc.WithTransportCredentials(insecure.NewCredentials()),
\t\tgrpc.WithKeepaliveParams(keepalive.ClientParameters{
\t\t\tTime:    30 * time.Second,
\t\t\tTimeout: 10 * time.Second,
\t\t}),
\t)
\tif err != nil {
\t\treturn nil, fmt.Errorf("dial auth-service: %w", err)
\t}
\tlog.Info().Str("addr", authAddr).Msg("connected to auth-service gRPC")
\treturn &AuthValidator{client: pbauth.NewAuthServiceClient(conn), conn: conn}, nil
}

// Validate calls AuthService.ValidateToken and returns the parsed claims.
func (v *AuthValidator) Validate(ctx context.Context, token string) (*pkgjwt.Claims, error) {
\tresp, err := v.client.ValidateToken(ctx, &pbauth.ValidateTokenRequest{Token: token})
\tif err != nil {
\t\treturn nil, fmt.Errorf("validate token: %w", err)
\t}
\treturn &pkgjwt.Claims{
\t\tUserID:    resp.UserId,
\t\tEmail:     resp.Email,
\t\tName:      resp.Name,
\t\tAvatarURL: resp.AvatarUrl,
\t}, nil
}

// Close releases the gRPC connection.
func (v *AuthValidator) Close() error { return v.conn.Close() }
'''

# ──────────────────────────────────────────────────────────────
# internal/handlers.go
# ──────────────────────────────────────────────────────────────
files['internal/handlers.go'] = '''package internal

import (
\t"net/http"
\t"net/http/httputil"
\t"net/url"
\t"strings"

\t"github.com/gin-gonic/gin"
\t"github.com/rs/zerolog/log"
\tpkgresp "github.com/agamrai0123/wanderplan/pkg/response"
)

// Handlers holds references shared across HTTP handlers.
type Handlers struct {
\tValidator *AuthValidator
\tTargets   map[string]string // prefix → base URL, e.g. "/api/v1/trips" → "http://localhost:8082"
}

// NewHandlers constructs Handlers.
func NewHandlers(validator *AuthValidator, cfg ServicesCfg) *Handlers {
\ttargets := map[string]string{
\t\t"/auth":               "http://" + serviceAddr(cfg.AuthAddr, 8081),
\t\t"/api/v1/trips":       "http://" + serviceAddr(cfg.TripAddr, 8082),
\t\t"/api/v1/users":       "http://" + serviceAddr(cfg.UserAddr, 8083),
\t\t"/api/v1/collaborators": "http://" + serviceAddr(cfg.CollaborationAddr, 8084),
\t\t"/api/v1/notifications": "http://" + serviceAddr(cfg.NotificationAddr, 8085),
\t\t"/api/v1/search":      "http://" + serviceAddr(cfg.SearchAddr, 8086),
\t}
\treturn &Handlers{Validator: validator, Targets: targets}
}

func serviceAddr(addr string, defaultPort int) string {
\tif addr != "" {
\t\treturn addr
\t}
\treturn fmt.Sprintf("localhost:%d", defaultPort)
}

// Proxy returns a Gin handler that reverse-proxies the request to the target.
func (h *Handlers) Proxy(targetBase string) gin.HandlerFunc {
\tbase, err := url.Parse(targetBase)
\tif err != nil {
\t\tpanic("bad proxy target: " + targetBase)
\t}
\tproxy := httputil.NewSingleHostReverseProxy(base)
\tproxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
\t\tlog.Error().Err(err).Str("target", targetBase).Msg("proxy error")
\t\tw.WriteHeader(http.StatusBadGateway)
\t}
\treturn func(c *gin.Context) {
\t\t// Strip gin wildcard prefix so that the full path is forwarded correctly.
\t\tc.Request.URL.Path = c.Param("path")
\t\tif c.Request.URL.Path == "" {
\t\t\tc.Request.URL.Path = "/"
\t\t}
\t\tproxy.ServeHTTP(c.Writer, c.Request)
\t}
}

// Health responds to liveness probes.
func Health(c *gin.Context) {
\tpkgresp.OK(c, gin.H{"status": "ok"})
}

// findTarget returns the longest prefix match for the given path.
func (h *Handlers) findTarget(path string) (string, bool) {
\tbest := ""
\tbestAddr := ""
\tfor prefix, addr := range h.Targets {
\t\tif strings.HasPrefix(path, prefix) && len(prefix) > len(best) {
\t\t\tbest = prefix
\t\t\tbestAddr = addr
\t\t}
\t}
\tif best == "" {
\t\treturn "", false
\t}
\treturn bestAddr, true
}
'''

# ──────────────────────────────────────────────────────────────
# internal/routes.go
# ──────────────────────────────────────────────────────────────
files['internal/routes.go'] = '''package internal

import (
\t"fmt"
\t"net/http"

\t"github.com/gin-gonic/gin"
\t"github.com/prometheus/client_golang/prometheus"
\tpkgmw "github.com/agamrai0123/wanderplan/pkg/middleware"
\tpkgjwt "github.com/agamrai0123/wanderplan/pkg/jwt"
)

// RegisterRoutes attaches all routes to r.
func RegisterRoutes(
\tr *gin.Engine,
\th *Handlers,
\tjwtMgr *pkgjwt.Manager,
\treg *prometheus.Registry,
\tcfg *Config,
) {
\t_, metricsHandler := pkgmw.Metrics("api-gateway", reg)

\tr.Use(
\t\tpkgmw.RequestID(),
\t\tpkgmw.Logger(),
\t\tpkgmw.Recovery(),
\t\tpkgmw.CORS(cfg.CORS.AllowedOrigins),
\t\tpkgmw.RateLimit(cfg.RateLimit.RPS, cfg.RateLimit.Burst),
\t)

\t// Health & metrics (unauthenticated)
\tr.GET("/healthz", Health)
\tr.GET("/metrics", gin.WrapH(metricsHandler))

\t// Auth routes — proxied to auth-service (no JWT required)
\tauthBase := "http://" + cfg.Services.AuthAddr
\tif cfg.Services.AuthAddr == "" {
\t\tauthBase = "http://localhost:8081"
\t}
\tr.Any("/auth/*path", proxyTo(authBase))

\t// Protected API routes — validate JWT, then proxy
\tapi := r.Group("/api/v1")
\tapi.Use(pkgmw.Auth(jwtMgr))
\t{
\t\ttripBase := fmt.Sprintf("http://%s", coalesce(cfg.Services.TripAddr, "localhost:8082"))
\t\tuserBase := fmt.Sprintf("http://%s", coalesce(cfg.Services.UserAddr, "localhost:8083"))
\t\tcollabBase := fmt.Sprintf("http://%s", coalesce(cfg.Services.CollaborationAddr, "localhost:8084"))
\t\tnotifBase := fmt.Sprintf("http://%s", coalesce(cfg.Services.NotificationAddr, "localhost:8085"))
\t\tsearchBase := fmt.Sprintf("http://%s", coalesce(cfg.Services.SearchAddr, "localhost:8086"))

\t\tapi.Any("/trips/*path",          proxyTo(tripBase))
\t\tapi.Any("/users/*path",          proxyTo(userBase))
\t\tapi.Any("/collaborators/*path",  proxyTo(collabBase))
\t\tapi.Any("/notifications/*path",  proxyTo(notifBase))
\t\tapi.Any("/search/*path",         proxyTo(searchBase))
\t}
}

// proxyTo returns a Gin handler that reverse-proxies requests preserving the full path.
func proxyTo(baseURL string) gin.HandlerFunc {
\treturn func(c *gin.Context) {
\t\tpath := c.Param("path")
\t\tif path == "" {
\t\t\tpath = "/"
\t\t}
\t\thttp.StripPrefix("", nil) // ensure stdlib is imported
\t\tprxy := newReverseProxy(baseURL, path)
\t\tprxy.ServeHTTP(c.Writer, c.Request)
\t}
}

func coalesce(s, def string) string {
\tif s != "" {
\t\treturn s
\t}
\treturn def
}
'''

# ──────────────────────────────────────────────────────────────
# cmd/main.go
# ──────────────────────────────────────────────────────────────
files['cmd/main.go'] = '''package main

import (
\t"context"
\t"fmt"
\t"net/http"
\t"os"
\t"os/signal"
\t"syscall"
\t"time"

\t"github.com/gin-gonic/gin"
\t"github.com/rs/zerolog/log"
\tpkgcfg "github.com/agamrai0123/wanderplan/pkg/config"
\tpkgjwt "github.com/agamrai0123/wanderplan/pkg/jwt"
\tinternal "github.com/agamrai0123/wanderplan/services/api-gateway/internal"
)

func main() {
\t// ── Config ────────────────────────────────────────────────
\tvar cfg internal.Config
\tif err := pkgcfg.Load("api-gateway-config", "./config", "GW", &cfg); err != nil {
\t\tfmt.Fprintf(os.Stderr, "load config: %v\\n", err)
\t}
\tif err := cfg.Validate(); err != nil {
\t\tfmt.Fprintf(os.Stderr, "config validation: %v\\n", err)
\t\tos.Exit(1)
\t}

\t// ── Logger ────────────────────────────────────────────────
\tinternal.InitLogger(cfg.Logging, "api-gateway")

\t// ── JWT Manager (for Auth middleware on protected routes) ─
\tprivB64 := os.Getenv("JWT_PRIVATE_KEY")
\tpubB64  := os.Getenv("JWT_PUBLIC_KEY")
\tvar jwtMgr *pkgjwt.Manager
\tif privB64 != "" && pubB64 != "" {
\t\tvar err error
\t\tjwtMgr, err = pkgjwt.NewManagerFromBase64(privB64, pubB64, 15*time.Minute)
\t\tif err != nil {
\t\t\tlog.Fatal().Err(err).Msg("init jwt manager")
\t\t}
\t} else {
\t\tlog.Warn().Msg("JWT_PRIVATE_KEY / JWT_PUBLIC_KEY not set — Auth middleware disabled")
\t}

\t// ── Auth-service gRPC validator ───────────────────────────
\tauthValidator, err := internal.NewAuthValidator(cfg.Services.AuthAddr)
\tif err != nil {
\t\tlog.Fatal().Err(err).Msg("dial auth-service")
\t}
\tdefer authValidator.Close()

\t// ── HTTP server ───────────────────────────────────────────
\tgin.SetMode(gin.ReleaseMode)
\trouter := gin.New()

\treg := internal.NewRegistry()
\tinternal.RegisterRoutes(router, internal.NewHandlers(authValidator, cfg.Services), jwtMgr, reg, &cfg)

\tsrv := &http.Server{
\t\tAddr:         fmt.Sprintf(":%d", cfg.ServerPort),
\t\tHandler:      router,
\t\tReadTimeout:  30 * time.Second,
\t\tWriteTimeout: 60 * time.Second,
\t\tIdleTimeout:  120 * time.Second,
\t}

\tgo func() {
\t\tlog.Info().Int("port", cfg.ServerPort).Msg("api-gateway HTTP listening")
\t\tif err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
\t\t\tlog.Fatal().Err(err).Msg("http server")
\t\t}
\t}()

\t// ── Graceful shutdown ─────────────────────────────────────
\tquit := make(chan os.Signal, 1)
\tsignal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
\t<-quit

\tlog.Info().Msg("shutting down api-gateway…")
\tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
\tdefer cancel()
\tif err := srv.Shutdown(ctx); err != nil {
\t\tlog.Error().Err(err).Msg("graceful shutdown failed")
\t}
\tlog.Info().Msg("api-gateway stopped")
}
'''

# Write all files
for relpath, content in files.items():
    fpath = os.path.join(BASE, relpath)
    os.makedirs(os.path.dirname(fpath), exist_ok=True)
    with open(fpath, 'w', encoding='utf-8') as f:
        f.write(content)
    print(f'wrote {relpath}')

print('\nDone.')
