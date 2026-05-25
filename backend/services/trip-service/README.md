# trip-service

Core itinerary management service for WanderPlan. Manages trips, itinerary days, and itinerary items.

## Purpose

- Full CRUD for trips, itinerary days, and itinerary items
- Bulk reorder of itinerary items (powers dnd-kit drag-and-drop on the frontend)
- Trip duplication (deep-copy days + items, excludes collaborators)
- Exposes gRPC RPCs for user-service to fetch trip lists and stats
- Publishes Kafka events on trip lifecycle changes

## Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/healthz` | None | Health check |
| `GET` | `/metrics` | None | Prometheus metrics |
| `GET` | `/trips` | JWT | List authenticated user's trips |
| `POST` | `/trips` | JWT | Create a new trip |
| `GET` | `/trips/:id` | JWT | Get trip with days and items |
| `PATCH` | `/trips/:id` | JWT | Update trip metadata |
| `DELETE` | `/trips/:id` | JWT | Delete trip (owner only) |
| `POST` | `/trips/:id/duplicate` | JWT | Deep-copy trip under same user |
| `GET` | `/trips/:id/days` | JWT | List itinerary days |
| `POST` | `/trips/:id/days` | JWT | Add an itinerary day |
| `PATCH` | `/trips/:id/days/:dayId` | JWT | Update day notes/date |
| `DELETE` | `/trips/:id/days/:dayId` | JWT | Delete day and its items |
| `GET` | `/trips/:id/days/:dayId/items` | JWT | List items for a day |
| `POST` | `/trips/:id/days/:dayId/items` | JWT | Add itinerary item |
| `PATCH` | `/trips/:id/items/:itemId` | JWT | Update itinerary item |
| `DELETE` | `/trips/:id/items/:itemId` | JWT | Delete itinerary item |
| `PATCH` | `/trips/:id/items/reorder` | JWT | Bulk reorder — body: `[{id, order_index}]` |

## gRPC

- **Port:** `9082`
- **Service:** `TripService`
- **RPCs:**
  - `ListTripsByUser(user_id) → [TripSummary]` — called by user-service
  - `GetTripStats(user_id) → {total_trips, total_countries, total_days, total_budget}` — called by user-service

## Kafka Output

- Topic: `trip-events`
- Events: `trip.created`, `trip.updated`, `trip.deleted`

## Environment Variables

| Variable | Description |
|----------|-------------|
| `PORT` | HTTP port (default: `8082`) |
| `GRPC_PORT` | gRPC port (default: `9082`) |
| `DB_HOST` | PostgreSQL host |
| `DB_PORT` | PostgreSQL port |
| `DB_NAME` | Database name |
| `DB_USER` | Database user |
| `DB_PASSWORD` | Database password |
| `DB_SSL_MODE` | `disable` / `require` |
| `KAFKA_BROKERS` | Comma-separated Kafka broker addresses |
| `RATE_LIMIT_RPS` | Rate limit — requests per second |
| `RATE_LIMIT_BURST` | Rate limit — burst size |

## How to Run Locally

```bash
cp .env.example .env
cd backend/services/trip-service
air
```

## How to Run Tests

```bash
cd backend/services/trip-service
go test ./... -v -count=1

# Integration tests (requires Docker)
go test ./... -v -count=1 -tags integration
```

## Observability

- `GET /metrics` — `trip_operations_total`, `db_query_duration_seconds`, `go_goroutines`
- Structured JSON logs via `zerolog` + `lumberjack`

## Database Tables

- `trips` — trip metadata (owner, destination, dates, budget)
- `itinerary_days` — days ordered by `day_number`
- `itinerary_items` — ordered by `order_index`; bulk UPDATE in transaction for reorder
- `trip_tags` — batch insert/delete with each trip update
