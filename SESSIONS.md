# WanderPlan — Active Session Registry

> **Protocol:** Before starting any task, a session MUST claim it here. Before ending, release it.
> This file is the single source of truth for what work is in-flight across concurrent sessions.
> Coordinate by committing claims first; if two sessions claim the same task, the later commit wins and the earlier session must back off.

## Active Claims

<!-- Sessions write rows here when they start a task. Remove the row when done. -->

| Session ID | Entity / File scope | Task description | Claimed at |
|---|---|---|---|
| _(none)_ | — | — | — |

---

## How to use this file

### On session start
1. Read this file to see all active claims.
2. **Do not start any task already claimed** by another session (same entity or overlapping files).
3. Pick your task from `WORKLOG.md`. Add a row to **Active Claims** above:
   ```
   | <short-id>  | <entity, e.g. "frontend"> | <task, e.g. "wire Dashboard page to API"> | <YYYY-MM-DD HH:MM> |
   ```
4. Commit `SESSIONS.md` immediately with message `session: claim <entity> — <task>`.
5. Begin work.

### During work
- If you discover your task depends on a claimed entity, **stop** and coordinate:
  - Either wait for the other session to finish (check git log).
  - Or split the work into non-overlapping files and update both claims.

### On session end
1. Remove your row from **Active Claims**.
2. Update `WORKLOG.md` and the affected context doc(s) with what was completed.
3. Commit `SESSIONS.md` + `WORKLOG.md` + context docs together:
   ```
   git commit -m "session: release <entity> — <task> done"
   ```

---

## Entity ownership map (for determining overlap)

| Entity | Files it owns |
|---|---|
| `api-gateway` | `backend/services/api-gateway/**` |
| `auth-service` | `backend/services/auth-service/**` |
| `trip-service` | `backend/services/trip-service/**` |
| `user-service` | `backend/services/user-service/**` |
| `collaboration-service` | `backend/services/collaboration-service/**` |
| `notification-service` | `backend/services/notification-service/**` |
| `search-service` | `backend/services/search-service/**` |
| `frontend` | `frontend/src/**` |
| `shared-packages` | `backend/pkg/**` |
| `proto` | `backend/proto/**` |
| `migrations` | `migrations/**` |
| `infra` | `docker-compose.yml`, `render.yaml`, `deployments/**`, `Makefile` |
| `docs` | `docs/**`, `WORKLOG.md`, `README.md` |

> Two sessions overlap if they touch the **same entity** or if one touches `shared-packages` / `proto` / `migrations` while the other touches any service that depends on the changed package.
