# api-gateway

HTTP reverse proxy and authentication gateway for WanderPlan. Validates JWTs via gRPC to auth-service and proxies requests to downstream services.

## Purpose

- Single ingress for all client traffic on port **8080**
- JWT validation (RS256) delegated to auth-service via gRPC `ValidateToken`
- Reverse-proxies path prefixes to the correct upstream service
- Applies rate limiting, CORS, request-ID injection, and structured logging to every request

## Endpoints

| Method | Path | Auth | Upstream |
|--------|------|------|----------|
| `GET` | `/healthz` | None | Local |
| `GET` | `/metrics` | None | Local (Prometheus) |
| `ANY` | `/auth/*` | None | auth-service :8081 |
| `ANY` | `/api/v1/trips/*` | JWT | trip-service :8082 |
| `ANY` | `/api/v1/users/*` | JWT | user-service :8083 |
| `ANY` | `/api/v1/collaborators/*` | JWT | collaboration-service :8084 |
| `ANY` | `/api/v1/notifications/*` | JWT | notification-service :8085 |
| `ANY` | `/api/v1/search/*` | JWT | search-service :8086 |

## Middleware Chain

`RequestID → Logger → Recovery → CORS → RateLimit → Auth → ReverseProxy`

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP listen port |
| `AUTH_SERVICE_GRPC_ADDR` | `localhost:9081` | auth-service gRPC address for JWT validation |
| `AUTH_SERVICE_ADDR` | `localhost:8081` | auth-service HTTP upstream |
| `TRIP_SERVICE_ADDR` | `localhost:8082` | trip-service HTTP upstream |
| `USER_SERVICE_ADDR` | `localhost:8083` | user-service HTTP upstream |
| `COLLAB_SERVICE_ADDR` | `localhost:8084` | collaboration-service HTTP upstream |
| `NOTIF_SERVICE_ADDR` | `localhost:8085` | notification-service HTTP upstream |
| `SEARCH_SERVICE_ADDR` | `localhost:8086` | search-service HTTP upstream |
| `ALLOWED_ORIGINS` | `*` | CORS allowed origins (comma-separated) |
| `RATE_LIMIT_RPS` | `1000` | Rate limit — requests per second |
| `RATE_LIMIT_BURST` | `2000` | Rate limit — burst size |

## How to Run Locally

```bash
# From repo root — requires all downstream services running
cp .env.example .env   # edit values
cd backend/services/api-gateway
air                    # hot-reload via .air.toml
```

Or without hot-reload:

```bash
cd backend/services/api-gateway
go run ./cmd/...
```

## How to Run Tests

```bash
cd backend/services/api-gateway
go test ./... -v -count=1
```

Integration tests spin up a real PostgreSQL container via `testcontainers-go`. Requires Docker.

```bash
go test ./... -v -count=1 -tags integration
```

## Build

```bash
CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ../../bin/api-gateway ./cmd/...
```

## Observability

- `GET /metrics` — Prometheus metrics: `http_requests_total`, `http_request_duration_seconds`, `http_errors_total`, `go_goroutines`
- Structured JSON logs via `zerolog`, rotated by `lumberjack`
- Every response includes `meta.request_id` and `meta.timestamp`
