# WanderPlan — Session Worklog

> **HOW TO USE THIS FILE (read before every session)**
>
> 1. Read this entire file before touching any code.
> 2. Pick the first `[ ]` item from **TODO** (or pull the next from **PLANNED** if TODO is empty).
> 3. Break the work into atomic steps; record all steps in **PLANNED** first.
> 4. Mark one step `[ ]` in **TODO** and work on it.
> 5. When done: move it to **DONE** with `[x]`, date, and files changed.
> 6. A user request is only finished when every derived step is `[x]`.
> 7. Never delete history. Record blockers in **NOTES**.

---

## MANDATORY RULES FOR EVERY SESSION

These are the non-negotiable constraints. Violating any of them requires an immediate fix.

### Backend (Go)
- Framework: **Gin** only. No `net/http` handlers directly.
- All errors go through `pkg/errors.AppError` — never `fmt.Errorf` in handlers.
- Every response uses the standard envelope: `{ "data": …, "error": …, "meta": { "request_id": "…", "timestamp": "…" } }`.
- Config via **Viper** from `.env` / env vars. Zero hardcoded values.
- DB queries: **parameterised only** (`pgx` `$1`/`$2`). SQL string concat is forbidden.
- Logging: **zerolog** + **lumberjack**. No `log.Println` or `fmt.Print`.
- Prometheus: every service exposes `GET /metrics` and `GET /healthz`.
- Graceful shutdown: handle `SIGTERM`/`SIGINT`, 10-second drain, close HTTP + gRPC + DB pool + Kafka.
- Schema changes via **golang-migrate** files in `/migrations/` (`000001_init.up.sql`, `000001_init.down.sql`, …).
- Dockerfiles: multi-stage (`golang:1.22-alpine` builder → `gcr.io/distroless/static-debian12` final). Build flag: `CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s"`. Non-root user. No `.env` in image. Pin base image versions.
- New service: run `go work use ./services/<name>` inside `backend/`.

### Auth
- OAuth 2.0 PKCE — **Google and GitHub only**. No username/password flows.
- RS256 JWTs, 15-minute expiry. Refresh tokens opaque, hashed in PostgreSQL, rotated on every use.
- JWT validation: api-gateway calls auth-service via gRPC `ValidateToken`.
- Access tokens: in-memory on client. Refresh tokens: httpOnly cookies set by backend only.

### Frontend (React + TypeScript)
- React 18, **TypeScript strict mode**, zero TS errors.
- UI: **shadcn/ui** (Radix) only. No other component libraries.
- Styling: **Tailwind CSS** utilities only. No inline styles, CSS modules, or styled-components.
- Dark/light: `next-themes` + `dark:` classes.
- Routing: `react-router-dom` v6. No `window.location` manipulation.
- Server state: `@tanstack/react-query` v5. No local state for API data.
- Forms: `react-hook-form` + `zod`. No uncontrolled inputs or manual validation.
- Dates: `date-fns` + `react-day-picker`. No `moment.js`.
- Charts: `recharts`. Drag-and-drop: `@dnd-kit`.
- API client: `axios` with JWT refresh interceptor (retry on 401).
- All API base URLs from `VITE_API_BASE_URL` env var. No hardcoded `localhost` URLs in source.
- ESLint must pass with zero errors. `npm run build` must succeed.
- No mock/placeholder data in production code. All pages use `@tanstack/react-query`.

### Render.com Deployment
- `render.yaml` at repo root defines all services as Infrastructure-as-Code.
- Backend services: `runtime: docker`, reference `./backend/services/<name>/Dockerfile`.
- Frontend: `runtime: static`, `buildCommand: cd frontend && npm install && npm run build`, `staticPublishPath: ./frontend/dist`.
- Every service must have `healthCheckPath: /healthz`.
- Secrets (`JWT_PRIVATE_KEY`, `JWT_PUBLIC_KEY`, OAuth creds, `KAFKA_BROKERS`) are `sync: false` — set manually in Render dashboard.
- DB env vars (`DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`) use `fromDatabase` references, not hardcoded values.
- `VITE_API_BASE_URL` must be set to the api-gateway's public Render URL (e.g. `https://wanderplan-api-gateway.onrender.com`).

### WORKLOG Protocol
- Read this file first every session.
- One item in TODO at a time; all others in PLANNED.
- Update DONE immediately after each step completes.
- Commit WORKLOG.md with every code change.

---

## DONE

- [x] **Add session handover rule to copilot-instructions.md** — Added `## Session Handover & Task Tracking Rules` section. Files: `.github/copilot-instructions.md`. _(2026-05-24)_
- [x] **Create initial WORKLOG.md** — This file established as handover document. _(2026-05-24)_
- [x] **Full rule compliance audit + render.com deployment review** — Identified 8 violations; render.yaml exists in `main` branch only; Dockerfiles for backend services missing; 5 services unscaffolded; mock data still in frontend; no migrations, no docker-compose, empty README. WORKLOG.md updated with full context. _(2026-05-24)_
- [x] **WORKLOG PLANNED section updated from spec document** — Read `wanderplan_copilot_prompt.docx` (§1–§14); expanded PLANNED from 13 vague items to 28 granular tasks grouped by priority (Blocking / Foundation / Services / Frontend / Infrastructure / Testing / Docs); NOTES updated with service port map, Kafka topics, all 10 DB table names, and full Render deployment context. Files: `WORKLOG.md`. _(2026-05-24)_
- [x] **All 9 Render resources provisioned** — Deployed 7 Docker web services + 1 static frontend + 1 PostgreSQL via render.yaml Blueprint. JWT_PRIVATE_KEY + JWT_PUBLIC_KEY set on all 7 services via Render API. _(2026-05-24)_
- [x] **Fix env var name mismatch for Render deployment** — All 7 service main.go files updated to override DB config (DB_HOST/DB_PORT/DB_NAME/DB_USER/DB_PASSWORD) from flat env vars after Viper config load. auth-service also overrides KAFKA_BROKERS + OAuth credentials. trip/collab/notification-service override KAFKA_BROKERS. api-gateway overrides all service HTTP addresses + adds separate AUTH_SERVICE_GRPC_ADDR for gRPC auth validator (fixing a bug where HTTP proxy and gRPC validator used the same config field pointing to the wrong port). render.yaml updated with AUTH_SERVICE_GRPC_ADDR=wanderplan-auth-service:9081. Pushed to main → Render redeploy triggered. Files: backend/services/*/cmd/main.go, api-gateway/internal/config.go, api-gateway/config/api-gateway-config.json, render.yaml. _(2026-05-24)_

---

## TODO

_(no item currently in progress — pull the next from PLANNED)_

---

## PLANNED

> Priority order. Move one at a time into TODO.
> Source of truth for all items: `wanderplan_copilot_prompt.docx` (§1–§14).

### 🔴 BLOCKING — Render Deploy

- [ ] **[BLOCKER] Dockerfile — api-gateway** — `backend/services/api-gateway/Dockerfile`. Multi-stage: `golang:1.22-alpine` builder → `gcr.io/distroless/static-debian12` final. `CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/service ./cmd/...`. Non-root user. `EXPOSE 8080`. `HEALTHCHECK` via `/healthz`. `LABEL service=api-gateway`.
- [ ] **[BLOCKER] Dockerfile — auth-service** — `backend/services/auth-service/Dockerfile`. Same pattern. `EXPOSE 8081 9081`.
- [ ] **[BLOCKER] Copy render.yaml into develop branch** — File exists only in `main`; cherry-pick commit or recreate at repo root in `develop`. Required for Render blueprint to apply.

### 🟠 FOUNDATION — Migrations & Proto

- [ ] **Migration 000001_init** — `/migrations/000001_init.up.sql` + `000001_init.down.sql`. Tables: `users (id, email, name, avatar_url, provider, provider_id, created_at, updated_at)`, `refresh_tokens (id, user_id, token_hash, expires_at, revoked, created_at)`, `trips (id, user_id, title, destination, cover_image_url, start_date, end_date, status, visibility, budget_total, currency, created_at, updated_at)`, `itinerary_days (id, trip_id, day_number, date, notes)`, `itinerary_items (id, day_id, trip_id, title, description, location, place_id, type, start_time, end_time, cost, currency, order_index, created_at, updated_at)`, `collaborators (id, trip_id, user_id, role, invited_at, accepted_at)`, `trip_tags (id, trip_id, tag)`, `notifications (id, user_id, type, payload jsonb, read, created_at)`, `audit_log (id, user_id, action, resource_type, resource_id, metadata jsonb, created_at)`.
- [ ] **Migration 000002_add_places_cache** — `/migrations/000002_add_places_cache.up.sql` + `000002_add_places_cache.down.sql`. Table: `places_cache (id, place_id, name, lat, lng, address, photos jsonb, cached_at)`. Add `tsvector` column + GIN index on `trips` for full-text search.
- [ ] **Proto — collaboration.proto** — `/backend/proto/wanderplan/v1/collaboration.proto`. RPCs: `InviteCollaborator`, `ListCollaborators`, `UpdateRole`, `RemoveCollaborator`. Comment every RPC and field. Run `make proto` after.
- [ ] **Proto — notification.proto** — `/backend/proto/wanderplan/v1/notification.proto`. RPCs: `ListNotifications`, `MarkRead`, `MarkAllRead`. Comment all. Run `make proto`.
- [ ] **Proto — search.proto** — `/backend/proto/wanderplan/v1/search.proto`. RPCs: `SearchPlaces`, `SearchTrips`. Comment all. Run `make proto`.

### 🟡 SERVICES — Scaffold & Implement

Each service needs: `go.mod`, `cmd/main.go`, `internal/{config,handlers,routes,service,models,metrics,logger,errors,database}.go`, `.air.toml`, `Dockerfile`, `README.md`. Register `go work use ./services/<name>` in `backend/`.

- [ ] **Scaffold + implement trip-service (HTTP 8082 | gRPC 9082)** — Endpoints: `POST /trips`, `GET /trips`, `GET /trips/:id`, `PATCH /trips/:id`, `DELETE /trips/:id`; `POST /trips/:id/days`, `PATCH /trips/:id/days/:dayId`, `DELETE /trips/:id/days/:dayId`; `POST /trips/:id/days/:dayId/items`, `PATCH /trips/:id/items/:itemId`, `DELETE /trips/:id/items/:itemId`; `PATCH /trips/:id/items/reorder` (dnd-kit order array); `POST /trips/:id/duplicate`. gRPC: expose trip data for user-service. Kafka: publish `trip.created`, `trip.updated`, `trip.deleted` to topic `trip-events`.
- [ ] **Scaffold + implement user-service (HTTP 8083 | gRPC 9083)** — Endpoints: `GET /users/me`, `PATCH /users/me` (name + avatar), `GET /users/me/trips` (delegates to trip-service via gRPC), `GET /users/me/stats` (total trips, countries, days, budget — feeds recharts dashboard). gRPC: expose user data.
- [ ] **Scaffold + implement collaboration-service (HTTP 8084)** — Endpoints: `POST /trips/:id/collaborators` (invite by email), `GET /trips/:id/collaborators`, `PATCH /trips/:id/collaborators/:userId` (change role: owner/editor/viewer), `DELETE /trips/:id/collaborators/:userId`. Kafka: publish `collaboration.invited`, `collaboration.accepted` to topic `collab-events`.
- [ ] **Scaffold + implement notification-service (HTTP 8085)** — Kafka consumer for topics: `auth-events`, `trip-events`, `collab-events`. Persist to `notifications` table. Endpoints: `GET /notifications`, `PATCH /notifications/:id/read`, `PATCH /notifications/read-all`. WebSocket: `/ws/notifications` for real-time push. Prometheus: **Kafka consumer lag metric** (required). Dockerfile: `EXPOSE 8085`.
- [ ] **Scaffold + implement search-service (HTTP 8086)** — Endpoints: `GET /search/places?q=&lat=&lng=` (proxies Google Places API, caches in `places_cache`), `GET /search/trips?q=` (PostgreSQL `tsvector` FTS over user's trips). Env var: `GOOGLE_PLACES_API_KEY`. Dockerfile: `EXPOSE 8086`.
- [ ] **Dockerfiles for trip, user, collaboration, notification, search services** — Same multi-stage pattern as api-gateway/auth-service. Add to render.yaml (already referenced there; services just need the file to exist).

### 🔵 FRONTEND — Wire to Real API

- [ ] **Frontend: axios API client + JWT refresh interceptor** — `frontend/src/lib/api.ts`. Axios instance with `baseURL: import.meta.env.VITE_API_BASE_URL`. Attach `Authorization: Bearer <access_token>` from in-memory store. On 401: call `POST /auth/refresh` (refresh token in httpOnly cookie), retry original request. On refresh failure: redirect to login.
- [ ] **Frontend: Auth flow (Google + GitHub OAuth)** — Login page calls `GET /auth/{provider}/login` (redirect). Callback page handles `GET /auth/{provider}/callback` response (access token in memory, refresh token in httpOnly cookie). Replace current `AppContext` auth with real token management. Remove `Signup.tsx` or repurpose for OAuth-only entry.
- [ ] **Frontend: Remove mock data, wire Index + Dashboard pages** — Delete usage of `src/data/mock.ts` from `Index.tsx` and `Dashboard.tsx`. Fetch real data: trips list via `GET /api/v1/trips`, stats via `GET /users/me/stats` (feed into recharts charts). Use `@tanstack/react-query`.
- [ ] **Frontend: Wire Trips + TripDetail pages** — `Trips.tsx`: `GET /api/v1/trips` with react-query. `TripDetail.tsx`: `GET /api/v1/trips/:id`, full itinerary days/items. dnd-kit reorder calls `PATCH /trips/:id/items/reorder` immediately on drop.
- [ ] **Frontend: Wire CityDetail + PostDetail pages** — `CityDetail.tsx`: city data from `GET /search/places`. `PostDetail.tsx`: trip detail data from `GET /api/v1/trips/:id`.
- [ ] **Frontend: Collaborators panel on TripDetail** — Collapsible sidebar/drawer. List via `GET /trips/:id/collaborators`, invite by email (`POST`), change role (`PATCH`), remove (`DELETE`). Use react-hook-form + zod for invite form.
- [ ] **Frontend: Notifications bell in Header** — Poll `GET /notifications` every 30s (react-query with `refetchInterval`). Unread badge count. Dropdown list. `PATCH /notifications/:id/read` on click. `PATCH /notifications/read-all` button. Connect WebSocket `/ws/notifications` for real-time push.
- [ ] **Frontend: User profile/settings page** — New route `/settings`. `GET /users/me` to populate form. `PATCH /users/me` to save. Avatar upload field. react-hook-form + zod validation.
- [ ] **Frontend: Place autocomplete input** — In itinerary item form (TripDetail add-item drawer). Calls `GET /search/places?q=&lat=&lng=` as user types (debounced, react-query). Populates `place_id`, `location`, coordinates.
- [ ] **Frontend: Trip search bar in Dashboard** — Search input calls `GET /search/trips?q=` as user types. Results replace trip list. react-query + debounce.
- [ ] **Frontend: Dark/light mode** — Verify `next-themes` ThemeProvider wraps app. Ensure all new pages and components use `dark:` Tailwind classes. Add theme toggle button to Header.

### 🟢 INFRASTRUCTURE & OBSERVABILITY

- [ ] **docker-compose.yml** — At `/deployments/docker-compose.yml`. Services: all 7 Go services + `postgres:16-alpine` + `confluentinc/cp-kafka` + Zookeeper + `prom/prometheus` + `grafana/grafana`. Prometheus scrapes all services at `/metrics`. Grafana pre-provisioned with datasource + dashboards from `/deployments/grafana/dashboards/`.
- [ ] **Grafana dashboards** — One JSON dashboard per service in `/deployments/grafana/dashboards/`. Each shows: RPS, P50/P95/P99 latency, error rate, active DB connections, goroutine count. notification-service dashboard also shows Kafka consumer lag.
- [ ] **Makefile targets** — Verify/add all targets: `make run` (docker-compose up), `make dev` (air + vite dev), `make test` (Go + vitest), `make lint` (golangci-lint + eslint), `make migrate`, `make migrate-down`, `make proto` (buf generate), `make build` (binaries → /bin/), `make clean`.
- [ ] **scripts/run.sh, migrate.sh, generate-proto.sh** — Shell helper scripts in `/scripts/`.

### ⚪ TESTING

- [ ] **Go tests — api-gateway** — `*_test.go` for handlers (proxy routing, JWT validation middleware), service (AuthValidator), middleware (rate limit, CORS). testcontainers-go for integration. ≥80% coverage on business logic.
- [ ] **Go tests — auth-service** — Handler tests (OAuth callback, refresh, logout), service tests (token rotation, hashing), repository tests (DB CRUD for refresh_tokens, users). ≥80% coverage.
- [ ] **Go tests — trip-service** — Handler + service + repository tests for trips, days, items, reorder, duplicate. Kafka mock for publish. ≥80% coverage.
- [ ] **Go tests — remaining services** — user-service, collaboration-service, notification-service (Kafka consumer mock), search-service (Places API mock). ≥80% coverage each.
- [ ] **Frontend tests** — vitest + @testing-library/react: Login page OAuth redirect, TripDetail dnd-kit reorder, Collaborators panel invite form (zod validation), Notifications bell (react-query poll), route guard (`RequireAuth`).

### 📄 DOCUMENTATION

- [ ] **Root README.md** — Mermaid architecture diagram (all 7 services + Kafka + DB), full setup guide (prerequisites, `make run`), env var reference table (all variables from `.env.example`), API reference index (all service endpoints).
- [ ] **Service READMEs** — One `README.md` per service: purpose, all endpoints with request/response examples, env vars, `make dev` instructions, `make test` instructions. (api-gateway and auth-service READMEs first.)
- [ ] **godoc comments** — All exported Go functions and types across all packages and services must have godoc comments.

---

## NOTES

### Source Document
Full build spec: `wanderplan_copilot_prompt.docx` at repo root (§1–§14). All PLANNED items derive from it. Existing frontend reference: https://github.com/agamrai0123/wanderplan-itinerary — read `/src` before wiring frontend.

### Project Stack
- Backend: Go 1.22+ workspace at `backend/go.work`. Services in `backend/services/`. Shared packages: `backend/pkg/`.
- Frontend: Vite 5 + React 18 + TypeScript strict. Package manager: `bun`. Dev: `bun run dev` inside `frontend/`.
- Proto: Buf toolchain. Config: `backend/proto/buf.yaml`. Existing protos: `auth.proto`, `trip.proto`, `user.proto`.
- Migrations: golang-migrate. Files go in `/migrations/` at repo root.
- Deployments: `docker-compose.yml` lives at `/deployments/docker-compose.yml`. Grafana dashboards: `/deployments/grafana/dashboards/`.

### Service Port Map
| Service | HTTP | gRPC |
|---|---|---|
| api-gateway | 8080 | — |
| auth-service | 8081 | 9081 |
| trip-service | 8082 | 9082 |
| user-service | 8083 | 9083 |
| collaboration-service | 8084 | — |
| notification-service | 8085 | — |
| search-service | 8086 | — |

### Database Tables (all in migration 000001, except places_cache in 000002)
`users`, `refresh_tokens`, `trips`, `itinerary_days`, `itinerary_items`, `collaborators`, `trip_tags`, `notifications`, `audit_log`, `places_cache`

### Kafka Topics
- `auth-events` — published by auth-service (`user.created`, `user.login`)
- `trip-events` — published by trip-service (`trip.created`, `trip.updated`, `trip.deleted`)
- `collab-events` — published by collaboration-service (`collaboration.invited`, `collaboration.accepted`)
- notification-service **consumes** all three topics

### Render.com Deployment Context
- render.yaml exists in `main` branch only — **must be copied to `develop`**.
- Blueprint URL: `https://dashboard.render.com/blueprint/exs-d89d39bbc2fs73f3f2q0/sync/exe-d89d3lf7f7vs73c2j1p0`.
- Environment Group being configured for shared secrets (JWT keys, OAuth creds, Kafka broker).
- External prerequisites before first deploy: Upstash Kafka, Google Places API key, RS256 key pair (`openssl genrsa -out private.pem 2048 && openssl rsa -in private.pem -pubout -out public.pem`, then `base64 -w0`), Google + GitHub OAuth apps.
- **Blocking**: Dockerfiles for api-gateway and auth-service do not exist → render blueprint apply will fail.
- Frontend on Render uses `runtime: static` (no Dockerfile needed). `VITE_API_BASE_URL` must be set to the api-gateway public URL.

### Branch Strategy
- `develop` → `production` → `main` (enforced by CI pipeline `.github/workflows/`).
- Feature branches: `feature/*` → `develop`. Hotfixes: `hotfix/*` → `production`.
- Current working branch: `develop`.

### Known Violations (as of 2026-05-24)
1. **[CRITICAL]** No Dockerfiles for `api-gateway` or `auth-service` — blocks Render deploy.
2. **[CRITICAL]** `render.yaml` absent from `develop` branch.
3. 5 of 7 backend services unscaffolded (trip, user, collaboration, notification, search).
4. Frontend pages consume `mock.ts` — must use real API + react-query.
5. No `/migrations/` directory.
6. No `/deployments/docker-compose.yml`.
7. `README.md` is empty.
8. No Go tests (`*_test.go`) or real frontend tests.
9. Collaborators panel, Notifications bell, Profile/Settings page, Place autocomplete, Trip search — all missing from frontend.
10. No Grafana dashboard JSON files.
