# WanderPlan — Copilot Remote & Mobile Session Guide

Control and monitor long-running Copilot agent sessions (migrations, proto
regeneration, full-service scaffolding) from your phone or any browser while
your machine keeps working.

---

## Prerequisites

1. **GitHub Mobile app** installed on your phone (iOS or Android).
   Log in with the same GitHub account used for Copilot.
2. **Copilot CLI** updated to the latest version:
   ```bash
   npm update -g @githubnext/github-copilot-cli
   # or via the VS Code extension — it auto-updates
   ```
3. VS Code setting already enabled in `.vscode/settings.json`:
   ```json
   "github.copilot.chat.cli.remote.enabled": true
   ```

---

## Starting a remote session

### Option A — From VS Code (recommended for WanderPlan)

1. Open the Agents panel in VS Code (the sidebar icon or `Ctrl+Alt+I`).
2. Start an agent task (e.g. "scaffold the notification-service").
3. In the Chat input, type:
   ```
   /remote on
   ```
4. VS Code displays a **QR code** and a `github.com/copilot/agents` link.
5. Scan the QR code with GitHub Mobile **or** open the link in any browser.

### Option B — From Copilot CLI

```bash
# Start a brand-new remote session
copilot --remote

# Or enable remote on an already-running session
/remote on

# Keep your laptop awake for long runs (migrations, full-service build)
/keep-alive
```

---

## What you can do from your phone

| Action | How |
|---|---|
| Watch the agent in real time | Session stream shows every file write, shell command, and tool call as it happens |
| Send a follow-up instruction | Type in the GitHub Mobile chat input — the agent picks it up immediately |
| Redirect the agent mid-task | Send a new message; the agent pivots without losing context |
| Approve or deny a permission request | Tap the Approve / Deny prompt when the agent asks (e.g. before running `make migrate`) |
| Stop the session | Tap the Stop button in the Agents tab |
| Queue the next task | Type your next prompt while the current one is still running — it queues automatically |

---

## Recommended WanderPlan workflows to run remotely

These are the tasks that take the longest and benefit most from remote monitoring:

```
# Kick off from VS Code, go make coffee, check from your phone
"Build out the notification-service: Kafka consumer for auth-events,
 trip-events, and collab-events; persist to notifications table;
 WebSocket push on /ws/notifications."

"Write all testcontainers-go integration tests for the trip-service
 repository layer. Target 80% coverage."

"Create migration 000003 adding tsvector full-text search index to
 trips.title and trips.destination, with a matching down migration."
```

---

## Session visibility

Sessions are **private by default** — only the GitHub account that started
the session can view or control it.

Find all your active and recent sessions at:
```
github.com/copilot/agents
```
Or tap **Agents** in the GitHub Mobile bottom tab.

---

## Tips for long WanderPlan agent runs

- Always type `/keep-alive` before leaving your desk — prevents macOS/Linux
  sleep from killing the CLI process mid-migration.
- The `Stop` hook in `.vscode/settings.json` runs `make test` at the end of
  each agent turn. Check the test output in the remote session stream before
  approving the next step.
- The `PreToolUse` guard hook will fire on your phone as a permission request
  if the agent tries a blocked command (e.g. `git push origin main`). Deny it
  and send a corrective instruction from mobile.
- If the agent stalls waiting for input, GitHub Mobile will send a push
  notification — you don't need to keep the app open.
