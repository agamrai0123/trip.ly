# WanderPlan — Architecture & Entity Relationship Graph

> **Maintenance rule**: Update this file whenever you add, remove, or rename a service, endpoint, Kafka topic, proto RPC, database table, or shared package. The section headings and Mermaid diagrams are the single source of truth for how every entity connects.

---

## System Overview

```mermaid
graph TB
    subgraph "Client"
        FE["Frontend<br/>React 18 + TypeScript<br/>:3000 (dev)"]
    end

    subgraph "API Layer"
        GW["api-gateway<br/>:8080"]
    end

    subgraph "Business Services"
        AUTH["auth-service<br/>:8081 / gRPC :9081"]
        TRIP["trip-service<br/>:8082 / gRPC :9082"]
        USER["user-service<br/>:8083 / gRPC :9083"]
        COLLAB["collaboration-service<br/>:8084"]
        NOTIF["notification-service<br/>:8085"]
        SEARCH["search-service<br/>:8086"]
    end

    subgraph "Message Bus"
        KAFKA["Apache Kafka<br/>auth-events | trip-events | collab-events"]
    end

    subgraph "Data"
        PG["PostgreSQL<br/>wanderplan schema"]
        PLACES["Google Places API"]
    end

    FE -->|"HTTP/REST"| GW
    GW -->|"HTTP proxy (unauthenticated)"| AUTH
    GW -->|"HTTP proxy (JWT required)"| TRIP
    GW -->|"HTTP proxy (JWT required)"| USER
    GW -->|"HTTP proxy (JWT required)"| COLLAB
    GW -->|"HTTP proxy (JWT required)"| NOTIF
    GW -->|"HTTP proxy (JWT required)"| SEARCH
    GW -->|"gRPC ValidateToken"| AUTH

    AUTH -->|"produces: user.login"| KAFKA
    TRIP -->|"produces: trip.created/updated/deleted"| KAFKA
    COLLAB -->|"produces: collaboration.invited"| KAFKA
    KAFKA -->|"consumes: all 3 topics"| NOTIF

    USER -->|"gRPC ListTripsByUser / GetTripStats"| TRIP

    AUTH --- PG
    TRIP --- PG
    USER --- PG
    COLLAB --- PG
    NOTIF --- PG
    SEARCH --- PG
    SEARCH -->|"REST"| PLACES

    FE -->|"WebSocket /ws/notifications"| NOTIF
```

---

## Inter-Service Communication Map

### gRPC Calls

| Caller | Callee | Proto File | RPC | Purpose |
|--------|--------|------------|-----|---------|
| `api-gateway` | `auth-service` | `auth.proto` | `ValidateToken` | JWT validation on every request |
| `user-service` | `trip-service` | `trip.proto` | `ListTripsByUser` | Return trips for `/users/me/trips` |
| `user-service` | `trip-service` | `trip.proto` | `GetTripStats` | Dashboard stats for `/users/me/stats` |

### Kafka Event Flow

| Topic | Producer | Events Published | Consumer | Action on Consume |
|-------|----------|-----------------|----------|-------------------|
| `auth-events` | `auth-service` | `user.login` | `notification-service` | Create Notification record + WebSocket push |
| `trip-events` | `trip-service` | `trip.created`, `trip.updated`, `trip.deleted` | `notification-service` | Create Notification + push |
| `collab-events` | `collaboration-service` | `collaboration.invited` | `notification-service` | Create Notification + push |

### HTTP Proxy Rules (api-gateway)

| Path Pattern | Upstream | Auth Required |
|-------------|----------|---------------|
| `/auth/*` | `auth-service:8081` | No |
| `/api/v1/trips/*` | `trip-service:8082` | Yes (JWT) |
| `/api/v1/users/*` | `user-service:8083` | Yes (JWT) |
| `/api/v1/collaborators/*` | `collaboration-service:8084` | Yes (JWT) |
| `/api/v1/notifications/*` | `notification-service:8085` | Yes (JWT) |
| `/api/v1/search/*` | `search-service:8086` | Yes (JWT) |

---

## Database Table Ownership

| Table | Owner Service | Referenced By |
|-------|--------------|---------------|
| `users` | `auth-service` (write) | `user-service` (read), `collaboration-service` (read) |
| `refresh_tokens` | `auth-service` | — |
| `trips` | `trip-service` | `collaboration-service`, `search-service`, `user-service` (via gRPC) |
| `itinerary_days` | `trip-service` | — |
| `itinerary_items` | `trip-service` | — |
| `collaborators` | `collaboration-service` | — |
| `trip_tags` | `trip-service` | — |
| `notifications` | `notification-service` | — |
| `places_cache` | `search-service` | — |
| `audit_log` | shared (any service writes) | — |

---

## Shared Package Dependencies

```mermaid
graph LR
    subgraph "backend/pkg/"
        CFG["config/"]
        DB["database/"]
        ERR["errors/"]
        GRPC["grpc/"]
        JWT["jwt/"]
        KAFKA_PKG["kafka/"]
        LOG["logger/"]
        MW["middleware/"]
        RESP["response/"]
    end

    subgraph "services"
        GW2["api-gateway"]
        AUTH2["auth-service"]
        TRIP2["trip-service"]
        USER2["user-service"]
        COLLAB2["collaboration-service"]
        NOTIF2["notification-service"]
        SEARCH2["search-service"]
    end

    GW2 --> CFG
    GW2 --> ERR
    GW2 --> JWT
    GW2 --> MW
    GW2 --> RESP
    GW2 --> GRPC

    AUTH2 --> CFG
    AUTH2 --> DB
    AUTH2 --> ERR
    AUTH2 --> JWT
    AUTH2 --> KAFKA_PKG
    AUTH2 --> LOG
    AUTH2 --> MW
    AUTH2 --> RESP

    TRIP2 --> CFG
    TRIP2 --> DB
    TRIP2 --> ERR
    TRIP2 --> KAFKA_PKG
    TRIP2 --> LOG
    TRIP2 --> MW
    TRIP2 --> RESP

    USER2 --> CFG
    USER2 --> DB
    USER2 --> ERR
    USER2 --> GRPC
    USER2 --> LOG
    USER2 --> MW
    USER2 --> RESP

    COLLAB2 --> CFG
    COLLAB2 --> DB
    COLLAB2 --> ERR
    COLLAB2 --> KAFKA_PKG
    COLLAB2 --> LOG
    COLLAB2 --> MW
    COLLAB2 --> RESP

    NOTIF2 --> CFG
    NOTIF2 --> DB
    NOTIF2 --> ERR
    NOTIF2 --> KAFKA_PKG
    NOTIF2 --> LOG
    NOTIF2 --> MW
    NOTIF2 --> RESP

    SEARCH2 --> CFG
    SEARCH2 --> DB
    SEARCH2 --> ERR
    SEARCH2 --> LOG
    SEARCH2 --> MW
    SEARCH2 --> RESP
```

---

## Proto File Dependency Map

| Proto File | Generated gRPC Server | Generated gRPC Client Used By |
|-----------|----------------------|-------------------------------|
| `wanderplan/v1/auth.proto` | `auth-service` | `api-gateway` |
| `wanderplan/v1/trip.proto` | `trip-service` | `user-service` |
| `wanderplan/v1/user.proto` | `user-service` | _(not called yet)_ |

---

## Frontend → Backend API Map

| Frontend Component / Hook | HTTP Method + Path | Backend Service |
|--------------------------|-------------------|-----------------|
| `Login.tsx` | `GET /auth/{provider}/login` | auth-service |
| `Signup.tsx` (callback) | `GET /auth/{provider}/callback` | auth-service |
| `api.ts` (refresh interceptor) | `POST /auth/refresh` | auth-service |
| `Header.tsx` logout | `POST /auth/logout` | auth-service |
| `Header.tsx` profile | `GET /auth/me` | auth-service |
| `Dashboard.tsx` stats | `GET /users/me/stats` | user-service → trip-service |
| `Trips.tsx` | `GET /trips` | trip-service |
| `TripDetail.tsx` load | `GET /trips/:id` | trip-service |
| `TripDetail.tsx` create | `POST /trips` | trip-service |
| `TripDetail.tsx` update | `PATCH /trips/:id` | trip-service |
| `TripDetail.tsx` delete | `DELETE /trips/:id` | trip-service |
| `TripDetail.tsx` add day | `POST /trips/:id/days` | trip-service |
| `TripDetail.tsx` add item | `POST /trips/:id/days/:dayId/items` | trip-service |
| `TripDetail.tsx` reorder | `PATCH /trips/:id/items/reorder` | trip-service |
| `TripDetail.tsx` collaborators | `GET /trips/:id/collaborators` | collaboration-service |
| `TripDetail.tsx` invite | `POST /trips/:id/collaborators` | collaboration-service |
| `Header.tsx` notifications | `GET /notifications` | notification-service |
| `Header.tsx` mark read | `PATCH /notifications/:id/read` | notification-service |
| `Header.tsx` WS | `WS /ws/notifications` | notification-service |
| `CityDetail.tsx` places | `GET /search/places?q=` | search-service |
| `TripDetail.tsx` autocomplete | `GET /search/places?q=&lat=&lng=` | search-service |
| `Dashboard.tsx` trip search | `GET /search/trips?q=` | search-service |

---

## File-Level Dependency Graph (per service)

See individual context documents in `docs/context/` for file-level relationships within each service.

- [api-gateway](context/api-gateway.md)
- [auth-service](context/auth-service.md)
- [trip-service](context/trip-service.md)
- [user-service](context/user-service.md)
- [collaboration-service](context/collaboration-service.md)
- [notification-service](context/notification-service.md)
- [search-service](context/search-service.md)
- [frontend](context/frontend.md)
- [shared packages](context/shared-packages.md)
