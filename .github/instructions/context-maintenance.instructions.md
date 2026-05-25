---
applyTo: "backend/**/*.go,frontend/src/**/*.ts,frontend/src/**/*.tsx,backend/proto/**/*.proto,migrations/**/*.sql"
---

# Context Document Maintenance Rules

These rules apply whenever you modify any Go, TypeScript, proto, or SQL file in the WanderPlan project.

## The Rule: Keep Context Docs In Sync

Whenever you change a file that affects a documented entity, you MUST update the corresponding context document in `docs/context/` and, if the inter-service relationship graph changes, update `docs/ARCHITECTURE.md`.

---

## Which file maps to which context doc

| Files Changed | Context Doc to Update |
|--------------|----------------------|
| `backend/services/api-gateway/**` | `docs/context/api-gateway.md` |
| `backend/services/auth-service/**` | `docs/context/auth-service.md` |
| `backend/services/trip-service/**` | `docs/context/trip-service.md` |
| `backend/services/user-service/**` | `docs/context/user-service.md` |
| `backend/services/collaboration-service/**` | `docs/context/collaboration-service.md` |
| `backend/services/notification-service/**` | `docs/context/notification-service.md` |
| `backend/services/search-service/**` | `docs/context/search-service.md` |
| `backend/pkg/**` | `docs/context/shared-packages.md` |
| `frontend/src/**` | `docs/context/frontend.md` |
| `backend/proto/**/*.proto` | `docs/ARCHITECTURE.md` (Proto File Dependency Map section) |
| `migrations/**/*.sql` | `docs/ARCHITECTURE.md` (Database Table Ownership section) |

---

## What to update in the context doc

Update the **Last updated** date at the top.

Then update whichever section changed:

- **New endpoint added** → add row to the HTTP Endpoints table
- **Endpoint removed or path changed** → update or remove the row
- **New gRPC RPC added** → add row to gRPC Server / gRPC Clients table
- **New Kafka event** → add row to Kafka Producer or Kafka Consumer table
- **New env var** → add row to Config / Environment Variables table
- **New file added to service** → add line to File Map with its purpose
- **File renamed or deleted** → update File Map
- **New shared package** → add entry to `docs/context/shared-packages.md`
- **New function signature in pkg/** → update the Package Details section
- **New DB table** → update Database Table Ownership in `docs/ARCHITECTURE.md`
- **New inter-service call** → update Inter-Service Communication Map in `docs/ARCHITECTURE.md`

---

## What to update in ARCHITECTURE.md

Update `docs/ARCHITECTURE.md` whenever:
- A new service is added or removed → update System Overview Mermaid diagram
- A new gRPC call between services → add row to Inter-Service Communication Map + update Mermaid
- A new Kafka producer/consumer → add row to Kafka Event Flow table + update Mermaid
- A new proxy route in api-gateway → add row to HTTP Proxy Rules table
- A new database table is added → add row to Database Table Ownership table
- A new proto file is created → add row to Proto File Dependency Map
- A new frontend API call → add row to Frontend → Backend API Map

---

## What NOT to do

- Do NOT skip updating context docs because "the change is small". Small changes accumulate into stale docs.
- Do NOT rewrite context docs from scratch. Edit only the affected sections.
- Do NOT update the context doc BEFORE making the code change. Always update AFTER verifying the code works.

---

## Context doc format rules

- The `Last updated` date at the top uses `YYYY-MM-DD` format.
- Tables must remain aligned with `|` separators.
- File Map entries use `←` to describe purpose, indented 4 spaces for continuation lines.
- All endpoint paths use backtick formatting: `` `GET /trips/:id` ``.
