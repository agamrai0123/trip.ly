# ctx:api-gateway | 2026-05-24 | HTTP:8080
## state
todo: -
planned: Dockerfile — multi-stage(golang:1.22-alpine→gcr.io/distroless/static-debian12)
done: scaffold(cmd+internal/*)@2026-05-24 | env-var-fix(AUTH_SERVICE_GRPC_ADDR split from HTTP addr)@2026-05-24 | .env.example@2026-05-25 | README.md@2026-05-25 | tests(handlers_test.go)@2026-05-25
errors: -
_update when: endpoints/proxy rules/gRPC calls/middleware/config keys change_

## files
cmd/main.go: bootstrap — gin,prometheus,grpc-auth-client,shutdown(SIGTERM/SIGINT 10s)
config/api-gateway-config.toml: TOML defaults (overridden by env vars)
internal/config.go: Viper config struct; reads all env vars
internal/routes.go: RegisterRoutes; path-prefix→upstream; middleware chain
internal/handlers.go: ReverseProxyHandler — httputil.ReverseProxy per upstream, strip prefix
internal/service.go: AuthValidator.ValidateToken → gRPC(auth-service) → UserClaims
internal/models.go: Config, UserClaims structs
internal/metrics.go: http_requests_total,http_request_duration_seconds,http_errors_total,go_goroutines; /metrics
internal/logger.go: delegates pkg/logger | internal/errors.go: delegates pkg/errors
internal/database.go: unused | certs/: gRPC TLS certs
## routes
GET /healthz,/metrics [local no-auth]
ANY /auth/* → AUTH:8081 [no-auth]
ANY /api/v1/trips/* → TRIP:8082 [JWT]
ANY /api/v1/users/* → USER:8083 [JWT]
ANY /api/v1/collaborators/* → COLLAB:8084 [JWT]
ANY /api/v1/notifications/* → NOTIF:8085 [JWT]
ANY /api/v1/search/* → SEARCH:8086 [JWT]
## chain: RequestID→Logger→Recovery→CORS→RateLimit→Auth→ReverseProxy
## grpc-out: auth-service AuthService.ValidateToken @AUTH_SERVICE_GRPC_ADDR (every JWT req; injects UserClaims into gin ctx)
## env: PORT=8080 AUTH_SERVICE_GRPC_ADDR=localhost:9081 AUTH_SERVICE_ADDR=localhost:8081 TRIP_SERVICE_ADDR=localhost:8082 USER_SERVICE_ADDR=localhost:8083 COLLAB_SERVICE_ADDR=localhost:8084 NOTIF_SERVICE_ADDR=localhost:8085 SEARCH_SERVICE_ADDR=localhost:8086 ALLOWED_ORIGINS=* RATE_LIMIT_RPS=100
## pkg: config errors grpc jwt middleware response
## db: none | kafka: none
