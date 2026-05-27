# WanderPlan — GitHub Copilot Instructions

## Model Routing (read first)

Before starting any task, route it to the right model:

- **Default entry point**: type `/delegate` + your task. Qwen3.6 35B (local Ollama) runs first — zero cloud tokens. It either completes the task locally or emits a one-line escalation to a cloud prompt.
- **Simple tasks** (explain, refactor one function, boilerplate): `/explain-offline`, `/refactor-offline`, `/offline-task` → Qwen3.6 35B
- **Complex tasks** (multi-file feature, gRPC/proto, security audit): `/implement-feature`, `/online-task`, `/security-review-online` → Claude Sonnet 4.5

> If you receive a task without a `/prompt`, default to cloud (current model) but mention that `/delegate` would be more token-efficient for simple sub-tasks.

Full routing table: `.github/instructions/model-routing.instructions.md`

---

## Project Overview
WanderPlan is an AI-powered travel itinerary planning platform. It consists of a React 18 + TypeScript frontend and seven Go microservices communicating over HTTP, gRPC, and Kafka.

### Context Documents (read before working on an entity)
Before modifying any service, shared package, or frontend code, read the relevant context document:
- Architecture + relationship graph: `docs/ARCHITECTURE.md`
- Per-entity context + work state: `docs/context/<entity>.md` (api-gateway, auth-service, trip-service, user-service, collaboration-service, notification-service, search-service, frontend, shared-packages)

Each context doc is machine-readable and compact. It contains:
- `## state` block with `todo:`, `planned:`, `done:`, `errors:` — the work queue for that entity
- `## files`, `## routes`, `## env`, `## pkg`, `## db`, etc. — all facts about the entity

**After every code change**, update the `done:` (and clear `todo:`) in the affected context doc(s).

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

---

## Session Coordination Rules (concurrent sessions)

- **At the very start of every session**, before touching any source file, read `SESSIONS.md` and check the **Active Claims** table.
- **Claim your task** by adding a row to `SESSIONS.md` and committing + pushing it immediately. This is the distributed lock — the first push wins.
- **Never work on an entity that is already claimed** by another session. Pick a different unclaimed task from `WORKLOG.md`, or wait.
- **Release your claim** (remove your row from `SESSIONS.md`) as part of the final commit when the task is done.
- Full protocol (claim format, stale claim policy, overlap detection) is in `.github/instructions/session-coordination.instructions.md`.

---

## Session Handover & Task Tracking Rules

- **At the start of every session**, read `WORKLOG.md` first. It is a thin index showing which context doc has the next pending `todo:`.
- **Load only the context doc(s)** you need for the current task (`docs/context/<entity>.md`). Do not load all context docs — load only relevant ones to save tokens.
- **Per-entity work** (service implementation, frontend wiring, etc.) is tracked inside each `docs/context/<entity>.md` under `## state` with `todo:`, `planned:`, `done:`, `errors:`.
- **Cross-cutting work** (migrations, docker-compose, READMEs, tests) is tracked in `WORKLOG.md` under Cross-cutting TODO/PLANNED/DONE.
- **After completing each step**, update `todo:` → `done:` in the context doc. Update the status row in `WORKLOG.md`'s index table.
- **A request is only done** when all derived steps are marked done.
- Commit `WORKLOG.md` and the affected context doc(s) with every code change.

---

## Deployment Verification Rules

- **Code is only "done" when it is live and verified on Render.com.** Pushing to `main` is a prerequisite, not the finish line.
- After every `git push origin main`, wait for all affected Render services to finish building and deploying before declaring the task complete. A deploy that is still "In Progress" or "Failed" means the task is NOT done.
- **Before ending any session that touched backend code**, verify every affected service's `/healthz` endpoint returns HTTP 200:
  ```
  curl -sf https://<service-name>.onrender.com/healthz
  ```
  If any service returns non-200 or times out, diagnose the failure (check Render logs via the dashboard or API), fix it, and redeploy before closing the session.
- **Render deploy failure checklist** — if a service fails to deploy, check in this order:
  1. Docker build errors (missing files, wrong build context, import errors).
  2. Config / env var issues (missing required env vars, wrong key names).
  3. Database connection errors (wrong host/port, missing credentials, migrations not run).
  4. gRPC dial errors (wrong address or port for inter-service calls).
  5. Panic / runtime errors visible in Render's log stream.
- **Render service URLs** (use these for health checks):
  - api-gateway: `https://wanderplan-api-gateway.onrender.com/healthz`
  - auth-service: `https://wanderplan-auth-service.onrender.com/healthz`
  - trip-service: `https://wanderplan-trip-service.onrender.com/healthz`
  - user-service: `https://wanderplan-user-service.onrender.com/healthz`
  - collaboration-service: `https://wanderplan-collaboration-service.onrender.com/healthz`
  - notification-service: `https://wanderplan-notification-service.onrender.com/healthz`
  - search-service: `https://wanderplan-search-service.onrender.com/healthz`
- Use the Render API (`https://api.render.com/v1/`) with the project API key to check deploy status programmatically when needed. The key is stored in Render's dashboard under Account → API Keys.
- Do not call `task_complete` until all affected `/healthz` endpoints confirm HTTP 200.
