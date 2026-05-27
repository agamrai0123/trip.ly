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
| [frontend](docs/context/frontend.md) | 🟢 COMPLETE — all todos done: api.ts, auth, dashboard, trips, profile, collab, notif bell, WS notifications, settings, place autocomplete, trip search, dark/light mode |
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

- [x] **Migration 000001_init** — fixed schema: `trips.user_id/budget_total/visibility`; `itinerary_items.trip_id/order_index`; `000003_fix_trips_and_items.up/down.sql` for live DB. @2025-07-26
- [x] **Migration 000002_add_places_cache** — `000002_add_search_vector.up/down.sql`. Adds `search_vector` GENERATED tsvector + GIN index on `trips`; GIN `jsonb_path_ops` on `places_cache.results`. @2026-05-25
- [x] **Proto — collaboration.proto** — RPCs: `ListCollaborators`, `InviteCollaborator`, `UpdateCollaborator`, `RemoveCollaborator`. All fields commented. Go code generated. @2026-05-25
- [x] **Proto — notification.proto** — RPCs: `ListNotifications`, `MarkRead`, `MarkAllRead`. All fields commented. Go code generated. @2026-05-25
- [x] **Proto — search.proto** — RPCs: `SearchPlaces`, `SearchTrips`. All fields commented. Go code generated. @2026-05-25

### 🟡 SERVICES — Scaffold & Implement

- [x] **Scaffold + implement trip-service** — all CRUD, reorder, duplicate, Kafka, gRPC. @2026-05-24
- [x] **Scaffold + implement user-service** — profile CRUD, gRPC→trip-service stats. @2026-05-24
- [x] **Scaffold + implement collaboration-service** — invite/list/update/remove, Kafka. @2026-05-24
- [x] **Scaffold + implement notification-service** — Kafka consumer, WS push, REST, consumer-lag metric. @2026-05-24
- [x] **Scaffold + implement search-service** — Google Places proxy, places_cache, tsvector FTS. @2026-05-24
- [x] **Dockerfiles for trip, user, collaboration, notification, search services** — multi-stage, HEALTHCHECK, LABEL, EXPOSE. @2026-05-25

### 🔵 FRONTEND — Wire to Real API

- [x] **Frontend: axios API client + JWT refresh interceptor** — `frontend/src/lib/api.ts`. @session2
- [x] **Frontend: Auth flow (Google + GitHub OAuth)** — Login+AuthCallback, token in-memory, httpOnly cookie. @session2
- [x] **Frontend: Remove mock data, wire Index + Dashboard pages** — real react-query, mock.ts cleared. @session2
- [x] **Frontend: Wire Trips + TripDetail pages** — dnd-kit reorder → PATCH /trips/:id/items/reorder. @session2
- [x] **Frontend: Wire CityDetail + PostDetail pages** — @session2
- [x] **Frontend: Collaborators panel on TripDetail** — invite/role/remove, react-hook-form+zod. @session3
- [x] **Frontend: Notifications bell in Header** — poll 30s, WS push, useNotificationsWS.ts. @2026-05-25
- [x] **Frontend: User profile/settings page** — /settings route, GET/PATCH /users/me. @2026-05-25
- [x] **Frontend: Place autocomplete input** — debounced searchPlaces, 6-suggestion dropdown. @2026-05-25
- [x] **Frontend: Trip search bar in Dashboard** — debounced searchTrips, result cards. @2026-05-25
- [x] **Frontend: Dark/light mode** — ThemeProvider(next-themes), Sun/Moon toggle in Header. @2026-05-25

### 🟢 INFRASTRUCTURE & OBSERVABILITY

- [x] **docker-compose.yml** — at repo root; 7 services + postgres + kafka + zookeeper + prometheus + grafana; scripts/postgres-init.sh applies migrations. @2026-05-24
- [x] **Grafana per-service dashboards** — 7 JSON dashboards: api-gateway, auth-service, trip-service, user-service, collaboration-service, notification-service(+Kafka lag), search-service. @2025-07-26
- [x] **Makefile targets** — run/migrate/migrate-down/proto/build/test/lint/clean/dev all present. @2025-07-26
- [x] **scripts/run.sh, migrate.sh, generate-proto.sh** — created in `/scripts/`. @2025-07-26

### ⚪ TESTING

- [x] **Go tests — all 7 services** — handlers_test.go + service_test.go; auth integration test with testcontainers. @2026-05-25
- [x] **Frontend tests** — vitest + @testing-library/react: Login OAuth navigation, RequireAuth route guard, TripDetail dnd-kit reorder, CollaboratorsPanel invite form, Header notifications bell. 23/23 tests pass. @2025-07-27

### 📄 DOCUMENTATION

- [x] **Root README.md** — added Mermaid architecture diagram + API reference index (all endpoints) + env var reference table. @2025-07-26
- [x] **Service READMEs** — all 7 services have README.md. @2026-05-25
- [x] **godoc comments** — all exported Go functions and types have godoc comments; golint reports zero missing-comment issues. @2026-05-27

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
