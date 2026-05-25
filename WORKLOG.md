# WanderPlan — Worklog Index

## Protocol
1. Read this file. Find the first context doc with a pending `todo:`.
2. Load that context doc (`docs/context/<entity>.md`). Work the `todo:` item.
3. Update `todo:/planned:/done:` in the context doc after each step. Commit together.
4. Cross-cutting items (no single owner) are tracked below.
5. Rules are in `.github/copilot-instructions.md` — do not duplicate here.

## Context Docs — Pending Work
| Context doc | Status |
|---|---|
| [api-gateway](docs/context/api-gateway.md) | 🟢 LIVE — P2 complete: README, .env.example, tests done |
| [auth-service](docs/context/auth-service.md) | 🟢 LIVE — P2 complete: README, .env.example, tests done |
| [trip-service](docs/context/trip-service.md) | 🟢 LIVE — P2 complete: README, .env.example, tests done |
| [user-service](docs/context/user-service.md) | 🟢 LIVE — P2 complete: README, .env.example, tests done |
| [collaboration-service](docs/context/collaboration-service.md) | 🟢 LIVE — P2 complete: README, .env.example, tests done |
| [notification-service](docs/context/notification-service.md) | 🟢 LIVE — P2 complete: README, .env.example, tests done |
| [search-service](docs/context/search-service.md) | 🟢 LIVE — P2 complete: README, .env.example, tests done |
| [frontend](docs/context/frontend.md) | � PARTIAL: Phase1+2+3 done (api.ts, auth, dashboard, trips, profile, collab, notif bell); todo: WS notifications, search bar, settings page |
| [shared-packages](docs/context/shared-packages.md) | 🟠 PLANNED: 000002_add_places_cache migration |

## Cross-cutting TODO / PLANNED / DONE

### DONE
- [x] copilot-instructions.md session handover rules @2026-05-24
- [x] WORKLOG.md created @2026-05-24
- [x] Full rule compliance audit @2026-05-24
- [x] Expanded PLANNED from spec (§1–§14), 28 tasks @2026-05-24
- [x] All 9 Render resources provisioned (7 services + frontend + postgres) @2026-05-24
- [x] Fix env var mismatch — flat DB_* overrides in all main.go; AUTH_SERVICE_GRPC_ADDR added @2026-05-24
- [x] docs/context/ created (9 files) + docs/ARCHITECTURE.md @2026-05-24
- [x] WORKLOG restructured — thin dispatcher, state moved into context docs @2026-05-24
- [x] Root .gitignore created — covers *.exe, .env, *.pem, nul, logs, coverage @2025-07-26
- [x] 4 committed exe files removed from git tracking (git rm --cached) @2025-07-26
- [x] Rate limits raised — all services: rps=1000, burst=2000 (was 20/40) @2025-07-26
- [x] docs/PROJECT_STATUS.md created — full audit: 47% complete, load test results, task delegation, 7-8 week estimate @2025-07-26

### TODO
- [x] **[P1] Frontend API wiring** — lib/api.ts created, Phases 1-3 complete; remaining: WS notifications, search bar, settings (tracked in frontend context doc)
- [x] **[P3] Missing protos** — collaboration.proto, notification.proto, search.proto (buf lint clean; Go code regenerated) @2026-05-25
- [x] **[P3] Grafana dashboard panels** — wanderplan-services.json has 7 real panels: RPS, P50/P95/P99 latency, error rate, goroutines, service health, active DB connections (db_pool_acquired/total_connections), Kafka consumer lag (notification-service) @2026-05-25. Metric instrumentation fixed @2026-05-26: pkg/middleware.Metrics() now exports unprefixed http_requests_total + http_request_duration_seconds (removed wanderplan_<svc>_ namespace), status label numeric (strconv.Itoa), goroutines gauge removed (go_goroutines from GoCollector); api-gateway routes.go wired metricsMW; notification-service consumer.go now records kafka_consumer_lag{topic,partition} GaugeVec.
- [x] **[P3] Migration 000002** — 000002_add_search_vector.up/down.sql: adds search_vector GENERATED STORED tsvector + GIN index on trips; GIN jsonb_path_ops on places_cache.results @2026-05-25

### DONE (P2 — completed this session)
- [x] **[P2] .env.example** — root + 7 service files created @2026-05-25
- [x] **[P2] Service READMEs** — 7 service README.md files created @2026-05-25
- [x] **[P2] Tests** — unit tests for all 7 services (handlers_test.go + service_test.go, all PASS) + integration test skeleton (//go:build integration) for auth-service with testcontainers @2026-05-25

## NOTES
- Render service URLs: wanderplan-{api-gateway,auth-service,trip-service,user-service,collaboration-service,notification-service,search-service}.onrender.com
- DB owned by Render postgres; all services use fromDatabase env refs in render.yaml
- Kafka: external broker; KAFKA_BROKERS set manually in Render dashboard (sync:false)
- Topics: auth-events, trip-events, collab-events
- DB tables: users, refresh_tokens, trips, itinerary_days, itinerary_items, collaborators, trip_tags, notifications, audit_log, places_cache

- [x] **Create architecture relationship graph + per-entity context documents** — Created `docs/ARCHITECTURE.md` (system Mermaid diagram, inter-service communication map, DB table ownership, shared package dependencies, proto file map, frontend→backend API map). Created `docs/context/` with 8 context docs: api-gateway, auth-service, trip-service, user-service, collaboration-service, notification-service, search-service, frontend, shared-packages. Added `.github/instructions/context-maintenance.instructions.md` with auto-update rules. Updated `copilot-instructions.md` to reference context system. _(2026-05-24)_



---

## PLANNED

> Priority order. Move one at a time into TODO.
> Source of truth for all items: `wanderplan_copilot_prompt.docx` (§1–§14).

### 🔴 BLOCKING — Render Deploy

- [x] **[BLOCKER] Dockerfile — api-gateway** — `backend/services/api-gateway/Dockerfile`. Multi-stage alpine→alpine, LABEL, HEALTHCHECK /healthz, EXPOSE 8080. @2026-05-25
- [x] **[BLOCKER] Dockerfile — auth-service** — `backend/services/auth-service/Dockerfile`. Same pattern. EXPOSE 8081 9081, HEALTHCHECK, LABEL. @2026-05-25
- [x] **[BLOCKER] Dockerfiles — trip/user/collaboration/notification/search services** — All 5 Dockerfiles complete with LABEL, HEALTHCHECK, correct EXPOSE (trip:8082/9082, user:8083/9083, collab:8084, notif:8085, search:8086). @2026-05-25
- [ ] **[BLOCKER] Copy render.yaml into develop branch** — File exists only in `main`; cherry-pick commit or recreate at repo root in `develop`. Required for Render blueprint to apply.

### 🟠 FOUNDATION — Migrations & Proto

- [x] **Migration 000001_init** — `000001_init.up/down.sql` fixed: `trips.user_id` (was owner_id), `trips.budget_total` (was budget), `trips.visibility` added, `itinerary_items.trip_id` added, `itinerary_items.order_index` (was position). Notifications keep title/body/read_at (matches service code). `000003_fix_trips_and_items.up/down.sql` created for live DB that already applied old 000001. @2025-07-26
- [x] **Migration 000002_add_places_cache** — `000002_add_search_vector.up/down.sql`. Adds `search_vector` GENERATED tsvector + GIN index on `trips`; GIN `jsonb_path_ops` on `places_cache.results`. @2026-05-25
- [x] **Proto — collaboration.proto** — RPCs: `ListCollaborators`, `InviteCollaborator`, `UpdateCollaborator`, `RemoveCollaborator`. All fields commented. Go code generated. @2026-05-25
- [x] **Proto — notification.proto** — RPCs: `ListNotifications`, `MarkRead`, `MarkAllRead`. All fields commented. Go code generated. @2026-05-25
- [x] **Proto — search.proto** — RPCs: `SearchPlaces`, `SearchTrips`. All fields commented. Go code generated. @2026-05-25

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
