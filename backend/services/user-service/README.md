# user-service

User profile service for WanderPlan. Serves authenticated user profile data and aggregates trip statistics via gRPC.

## Purpose

- Read and update the authenticated user's profile (name, avatar)
- Proxy trip list requests to trip-service via gRPC (`ListTripsByUser`)
- Aggregate trip statistics for dashboard charts via gRPC (`GetTripStats`)
- Exposes gRPC `GetUser` RPC (reserved for future inter-service calls)

## Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/healthz` | None | Health check |
| `GET` | `/metrics` | None | Prometheus metrics |
| `GET` | `/users/me` | JWT | Return authenticated user profile |
| `PATCH` | `/users/me` | JWT | Update name and/or avatar URL |
| `GET` | `/users/me/trips` | JWT | List trips (delegates to trip-service gRPC) |
| `GET` | `/users/me/stats` | JWT | Trip stats for recharts dashboard |

### Stats Response Shape

```json
{
  "total_trips": 12,
  "total_countries": 8,
  "total_days": 45,
  "total_budget": 15000.00
}
```

## gRPC

- **Port:** `9083`
- **Outbound:** trip-service `TripService.ListTripsByUser` + `GetTripStats` @ `TRIP_SERVICE_GRPC_ADDR`
- **Inbound:** `UserService.GetUser(user_id) → {id, email, name, avatar_url}` (reserved)

## Environment Variables

| Variable | Description |
|----------|-------------|
| `PORT` | HTTP port (default: `8083`) |
| `GRPC_PORT` | gRPC port (default: `9083`) |
| `DB_HOST` | PostgreSQL host |
| `DB_PORT` | PostgreSQL port |
| `DB_NAME` | Database name |
| `DB_USER` | Database user |
| `DB_PASSWORD` | Database password |
| `DB_SSL_MODE` | `disable` / `require` |
| `TRIP_SERVICE_GRPC_ADDR` | trip-service gRPC address (default: `localhost:9082`) |
| `RATE_LIMIT_RPS` | Rate limit — requests per second |
| `RATE_LIMIT_BURST` | Rate limit — burst size |

## How to Run Locally

```bash
cp .env.example .env
cd backend/services/user-service
air
```

Requires trip-service to be running at `TRIP_SERVICE_GRPC_ADDR` for `/me/trips` and `/me/stats` to respond.

## How to Run Tests

```bash
cd backend/services/user-service
go test ./... -v -count=1

# Integration tests (requires Docker)
go test ./... -v -count=1 -tags integration
```

## Observability

- `GET /metrics` — `user_operations_total`, `grpc_call_duration_seconds`, `db_query_duration_seconds`, `go_goroutines`
- Structured JSON logs via `zerolog` + `lumberjack`

## Database Tables

- `users` — `SELECT` by JWT user ID for `GetMe`; `UPDATE` name/avatar_url for `UpdateMe`
