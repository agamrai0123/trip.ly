# WanderPlan — GitHub Copilot Instructions

## Project Overview
WanderPlan is an AI-powered travel itinerary planning platform. It consists of a React 18 + TypeScript frontend and seven Go microservices communicating over HTTP, gRPC, and Kafka.

---

## Architecture Rules

- There are seven backend services: api-gateway (8080), auth-service (8081), trip-service (8082), user-service (8083), collaboration-service (8084), notification-service (8085), search-service (8086).
- Synchronous inter-service calls use gRPC. All `.proto` files live in `/backend/proto/`; regenerate with `make proto` (buf).
- Async inter-service events use Apache Kafka via Sarama. Never use direct HTTP between services for events.
- All services share packages from `/backend/pkg/`. Import from there; do not duplicate logic in individual services.

---

## Backend (Go) Rules

- Go version: 1.22+. HTTP framework: Gin. No net/http handlers directly in service code.
- Never use `fmt.Errorf` in Gin handlers. All errors must go through `pkg/errors.AppError` with `Code`, `Message`, and `HTTPStatus` set.
- Every JSON response must use the standard envelope:
  ```json
  { "data": <payload or null>, "error": <error object or null>, "meta": { "request_id": "<uuid>", "timestamp": "<RFC3339>" } }
  ```
- All configuration (ports, DSNs, secrets, API keys) must be loaded via Viper from `.env` or environment variables. No hardcoded values anywhere in source.
- All database queries must use parameterised statements (`pgx` named args or positional `$1`, `$2`). SQL string concatenation is forbidden.
- Use `zerolog` for all structured logging. Use `lumberjack` for rotating file logs. Never use `log.Println` or `fmt.Print` for application logs.
- Every service must instrument the following with Prometheus: HTTP request count, latency histograms, error rates, DB query duration, Kafka consumer lag (notification-service only), and goroutine count. Expose `GET /metrics`.
- Every service must expose `GET /healthz`.
- Every service must handle `SIGTERM` and `SIGINT` with a 10-second drain timeout, closing the HTTP server, gRPC server, PostgreSQL pool, and Kafka producer/consumer cleanly.
- Every service needs its own `.air.toml` for hot-reload with `air`.
- Docker images must use multi-stage builds with a distroless or alpine final stage.
- All exported Go functions and types must have godoc comments.
- Code must compile cleanly with `go vet` and pass `golangci-lint` with the default ruleset.
- Use `golang-migrate` for all schema changes. Migration files go in `/migrations/` with sequential naming: `000001_init.up.sql`, `000001_init.down.sql`, etc.

---

## Authentication Rules

- Auth uses OAuth 2.0 PKCE with Google and GitHub only. No username/password flows.
- Issue RS256 JWTs with 15-minute expiry. Refresh tokens are opaque, stored hashed in PostgreSQL, and rotated on every refresh.
- JWT validation in the api-gateway delegates to auth-service via gRPC `ValidateToken`.
- Access tokens are stored in-memory on the frontend client. Refresh tokens are set in httpOnly cookies by the backend only.

---

## Frontend (React + TypeScript) Rules

- Use React 18 + TypeScript in strict mode. Zero TypeScript errors permitted.
- All UI components must come from shadcn/ui (Radix UI primitives). Do not introduce other component libraries.
- All styling uses Tailwind CSS utility classes. No inline styles, no CSS modules, no styled-components.
- Dark/light mode is handled exclusively by `next-themes` + Tailwind. Use `dark:` prefixed classes throughout.
- All routing is via `react-router-dom` v6. No manual `window.location` manipulation.
- All server state (API data) is managed by `@tanstack/react-query` v5. No local state for data that comes from the API.
- All forms use `react-hook-form` + `zod` for validation. No uncontrolled inputs or manual validation.
- All date handling uses `date-fns` and `react-day-picker`. Do not use `moment.js` (it is deprecated and increases bundle size).
- All data visualisations (trip stats, budget charts) use `recharts`.
- Drag-and-drop itinerary reordering uses `@dnd-kit`. On reorder, immediately call `PATCH /trips/:id/items/reorder`.
- The API client is an `axios` instance with a JWT refresh interceptor that retries on 401.
- All API base URLs must come from the `VITE_API_BASE_URL` environment variable. No hardcoded localhost URLs.
- The frontend must pass `eslint` with the existing `eslint.config.js` with zero errors.
- Do not remove or regress any existing component. Only extend and wire to real backend endpoints.
- Replace all mock and hardcoded data with real `@tanstack/react-query` calls. No placeholder data in production code.

---

## Testing Rules

### Go
- Use table-driven tests with the standard `testing` package + `testify/assert` and `testify/require`.
- Integration tests must use `testcontainers-go` to spin up real PostgreSQL and Kafka containers.
- Mock gRPC clients with `mockery`.
- Every service must have `*_test.go` files covering: HTTP handlers, service/business-logic, and repository/DB layers.
- Minimum 80% coverage on all business logic packages.

### Frontend
- Use `vitest` + `@testing-library/react` for component tests.
- Test all form validation logic, `@tanstack/react-query` hook behaviours, and route guards.

---

## Code Quality Rules

- Remove all placeholder files, empty components, TODO stubs, and unused imports left from the Lovable scaffold.
- All `.proto` files must have comment annotations on every RPC and message field.
- Every environment variable must have a corresponding entry in `.env.example`.
- Each service must have its own `README.md` covering: purpose, endpoints, environment variables, how to run locally, and how to run its tests.
- The root `README.md` must include a Mermaid architecture diagram, full setup guide, environment variable reference table, and API reference index.

---

## Observability Rules

- `docker-compose.yml` must start Prometheus (scraping all services) and Grafana with a pre-provisioned datasource and dashboards in `/deployments/grafana/dashboards/`.
- Each service Grafana dashboard must show: RPS, P50/P95/P99 latency, error rate, active DB connections, Kafka consumer lag (notification-service), and goroutine count.
- All HTTP routes, DB queries, Kafka produce/consume operations, and gRPC calls must be individually instrumented.
