# collaboration-service

Trip collaboration management service for WanderPlan. Handles inviting users, managing roles, and Kafka event publishing.

## Purpose

- Invite users to trips by email address (resolves email → user_id via `users` table)
- List, update role (viewer/editor/admin), and remove collaborators
- Enforces ownership: only the trip owner can invite, change roles, or remove collaborators
- Publishes `collaboration.invited` events to Kafka topic `collab-events`

## Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/healthz` | None | Health check |
| `GET` | `/metrics` | None | Prometheus metrics |
| `GET` | `/trips/:id/collaborators` | JWT | List collaborators for a trip |
| `POST` | `/trips/:id/collaborators` | JWT (owner) | Invite by email |
| `PATCH` | `/trips/:id/collaborators/:userId` | JWT (owner) | Change role |
| `DELETE` | `/trips/:id/collaborators/:userId` | JWT (owner) | Remove collaborator |

### Invite Request Body

```json
{ "email": "user@example.com", "role": "editor" }
```

### Roles

`viewer` < `editor` < `admin` < `owner`

- Only owner can invite/change-role/remove
- Cannot elevate your own role
- Cannot remove the owner (returns 400)

## Kafka Output

- Topic: `collab-events`
- Events:
  - `collaboration.invited` — `{trip_id, inviter_user_id, invitee_user_id, invitee_email, role}`

## Environment Variables

| Variable | Description |
|----------|-------------|
| `PORT` | HTTP port (default: `8084`) |
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
cd backend/services/collaboration-service
air
```

## How to Run Tests

```bash
cd backend/services/collaboration-service
go test ./... -v -count=1

# Integration tests (requires Docker)
go test ./... -v -count=1 -tags integration
```

## Observability

- `GET /metrics` — `collab_operations_total`, `db_query_duration_seconds`, `go_goroutines`
- Structured JSON logs via `zerolog` + `lumberjack`

## Database Tables

- `collaborators` — `SELECT` by trip_id, `INSERT`, `UPDATE role`, `DELETE`
- `users` — `SELECT` by email to resolve invitee
- `trips` — `SELECT owner_id` for ownership check
