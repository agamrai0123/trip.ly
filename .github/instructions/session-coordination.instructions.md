---
applyTo: "SESSIONS.md,WORKLOG.md"
---

# Session Coordination Rules

These rules prevent concurrent Copilot sessions from working on the same files and causing merge conflicts or duplicate work.

## The Three-Step Protocol

### Step 1 — CLAIM before touching any code

At the very start of a session, before reading any source files or planning changes:

1. Read `SESSIONS.md` — check the **Active Claims** table.
2. Identify which entity/files your task touches (use the entity ownership map in `SESSIONS.md`).
3. **If your entity is already claimed:** do NOT start. Pick a different unclaimed task from `WORKLOG.md`, or wait.
4. **If your entity is free:** add your claim row immediately:

   ```
   | <4-char-id>  | <entity>  | <one-line task>  | <YYYY-MM-DD HH:MM> |
   ```
   
   Use the first 4 characters of a random value (e.g., timestamp digits `0525`) as your session ID.

5. Commit `SESSIONS.md` with **only** this change first:
   ```
   git commit -m "session: claim <entity> — <task>"
   git push origin main
   ```
   This commit is the lock. If a concurrent session pushed first, you will see a conflict on pull — that means they won the claim; you must back off.

### Step 2 — WORK exclusively on claimed files

- Only modify files within your claimed entity's scope.
- If a change becomes necessary in a **different** entity's files, add a second claim row for that entity (check it's free first) before touching those files.
- The `shared-packages` entity (`backend/pkg/**`) is special: claiming it blocks all other sessions that depend on it. Keep `shared-packages` changes small and release quickly.

### Step 3 — RELEASE when done

Before calling `task_complete`:

1. Remove your row from the **Active Claims** table in `SESSIONS.md`.
2. Update the affected `docs/context/<entity>.md` (`todo:` → `done:`).
3. Update `WORKLOG.md` status table.
4. Commit everything together:
   ```
   git commit -m "session: release <entity> — <task> done"
   git push origin main
   ```

---

## Stale claim policy

A claim is **stale** if:
- Its `Claimed at` timestamp is more than **4 hours** old, AND
- No commit touching the claimed entity has been made in the last 2 hours (check `git log --since="2 hours ago"`).

If you find a stale claim, you may remove it and replace it with your own claim. Leave a note in the commit message:
```
git commit -m "session: evict stale claim for <entity> by <id>, claim <entity> — <task>"
```

---

## Overlap detection table

Two sessions overlap (and cannot run concurrently) if they touch the same entity **or** any of these cross-cutting combinations:

| If session A claims... | Then session B CANNOT claim... |
|---|---|
| `shared-packages` | any service (`api-gateway`, `auth-service`, etc.) |
| `proto` | any service that imports the changed proto |
| `migrations` | any service that uses the affected tables |
| `frontend` | `api-gateway` (if changing CORS or proxy routes) |
| any service | `infra` (if Dockerfile or render.yaml for that service changes) |

---

## Quick reference

```
# Start of session
git pull origin main
cat SESSIONS.md          # check active claims
# add your row, then:
git add SESSIONS.md && git commit -m "session: claim <entity> — <task>" && git push

# End of session
# remove your row, update context doc + WORKLOG, then:
git add SESSIONS.md WORKLOG.md docs/context/<entity>.md
git commit -m "session: release <entity> — <task> done"
git push origin main
```
