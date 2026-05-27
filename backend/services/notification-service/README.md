# notification-service

Kafka-driven notification fanout service for WanderPlan. Persists notifications and delivers them in real-time via WebSocket.

## Purpose

- Consumes Kafka events from `auth-events`, `trip-events`, and `collab-events`
- Persists notifications to PostgreSQL `notifications` table
- Delivers real-time notifications to connected browser tabs via WebSocket
- Exposes REST endpoints for listing, marking read, and marking all read
- Exposes Prometheus `kafka_consumer_lag` metric (required by spec)

## Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/healthz` | None | Health check |
| `GET` | `/metrics` | None | Prometheus metrics |
| `GET` | `/notifications` | JWT | List notifications (paginated, newest first) |
| `PATCH` | `/notifications/:id/read` | JWT | Mark a notification as read |
| `PATCH` | `/notifications/read-all` | JWT | Mark all notifications as read |
| `GET` | `/ws/notifications` | JWT (query param or header) | WebSocket — real-time push |

### WebSocket Message Format

```json
{
  "type": "notification",
  "data": {
    "id": "uuid",
    "type": "collaboration.invited",
    "payload": {},
    "read": false,
    "created_at": "2026-01-01T00:00:00Z"
  }
}
```

## Kafka Input

| Topic | Event | Action |
|-------|-------|--------|
| `auth-events` | `user.login` | Create "new device login" notification |
| `trip-events` | `trip.created` / `trip.updated` / `trip.deleted` | Notify collaborators |
| `collab-events` | `collaboration.invited` | Notify invitee |

- Consumer group: `wanderplan-notifications`
- Balance strategy: `BalanceStrategyRange`

## WebSocket Architecture

One `Hub` per process; one `Client` per browser tab. Hub broadcasts to all clients registered to the same `user_id`. PostgreSQL is the source of truth — WebSocket is delivery-only.

## Environment Variables

| Variable | Description |
|----------|-------------|
| `PORT` | HTTP port (default: `8085`) |
| `DB_HOST` | PostgreSQL host |
| `DB_PORT` | PostgreSQL port |
| `DB_NAME` | Database name |
| `DB_USER` | Database user |
| `DB_PASSWORD` | Database password |
| `DB_SSL_MODE` | `disable` / `require` |
| `KAFKA_BROKERS` | Comma-separated Kafka broker addresses |
| `KAFKA_CONSUMER_GROUP` | Consumer group ID (default: `wanderplan-notifications`) |
| `RATE_LIMIT_RPS` | Rate limit — requests per second |
| `RATE_LIMIT_BURST` | Rate limit — burst size |

## How to Run Locally

```bash
cp .env.example .env
cd backend/services/notification-service
air
```

Requires Kafka broker(s) at `KAFKA_BROKERS` and PostgreSQL.

## How to Run Tests

```bash
cd backend/services/notification-service
go test ./... -v -count=1

# Integration tests (requires Docker)
go test ./... -v -count=1 -tags integration
```

## Observability

- `GET /metrics` — `notifications_created_total`, `kafka_consumer_lag` (**GaugeVec{topic,partition}**), `websocket_active_connections`, `db_query_duration_seconds`, `go_goroutines`
- Consumer lag updated after each `MarkOffset` call

## Database Tables

- `notifications` — `INSERT` from Kafka consumer; `SELECT` paginated by user_id; `UPDATE read_at`
