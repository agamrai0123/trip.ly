# WanderPlan — Copilot Rate Limit Strategy

GitHub Copilot enforces two overlapping limits: a **session limit** (resets after
a short wait) and a **weekly limit** (resets on a 7-day rolling window).
Agentic workflows — especially anything that touches the microservice layer,
runs migrations, and regenerates protos in one turn — are the fastest way to
exhaust both.

---

## Tier 1 — Warning: approaching the limit

VS Code shows a warning in the Copilot status bar. Act before you hit the wall:

| Action | How |
|---|---|
| Switch to a cheaper model for the current task | Click the model picker in the Chat input → choose `GPT-4o mini` or `Claude Haiku` for simple questions |
| Enable Plan Mode | Type `/plan` in chat — Copilot reasons first, writes less, uses ~30–50% fewer tokens |
| Reduce parallel agent sessions | Close extra Copilot CLI sessions; each parallel tool call multiplies token usage |
| Tighten `.copilotignore` | Stop Copilot indexing `node_modules/`, `dist/`, `bin/`, and generated proto files |

---

## Tier 2 — Session limit hit (wait X minutes)

```
"You've hit your global rate limit. Wait 3h 42m or switch to Auto."
```

**Immediate steps:**

1. **Switch to Auto model** — in the Chat input model picker, select `Auto`.
   Auto uses smaller models for simple requests and still has budget remaining.
2. **Use Plan mode aggressively** — `/plan` before every agentic task so the
   agent writes a minimal diff instead of regenerating whole files.
3. **Break the task into isolated scopes** — open a second VS Code window
   pointed at a single service directory. Each window tracks rate limits
   independently, effectively doubling available capacity.

---

## Tier 3 — Weekly limit hit

```
"Sorry, you've exceeded your weekly rate limit."
```

You can still use Copilot with Auto model selection until premium requests
are also exhausted. To keep working:

1. Switch to **Auto** — free tier of requests is still available.
2. Use `/plan` only — no agentic edits, just planning and advice.
3. Use **inline completions** (ghost text) — these consume far fewer tokens
   than Chat or Agent mode.
4. For large refactors, open only the files actively being edited. Copilot
   won't scan the full workspace for context, cutting payload size.

---

## WanderPlan-specific token-saving tips

| Situation | Tip |
|---|---|
| Adding a new endpoint | Write the proto change manually, run `make proto`, then ask Copilot only to implement the handler — not the full service scaffold |
| Writing migrations | Use the `migrations.instructions.md` rules so Copilot generates the right file in one shot, not iteratively |
| Running tests | Ask Copilot to fix one failing test at a time, not all of them in a single agent turn |
| Reviewing generated code | Use inline chat (`Cmd+I`) instead of the Chat panel — it has a narrower context window and costs less |
| Proto regeneration | The `PostToolUse` hook runs `make proto` automatically — don't ask Copilot to do it again manually |

---

## `.copilotignore` for WanderPlan

Place this at the repo root. Copilot will not index these paths for context,
significantly reducing token usage per request.

```
# Generated code — never needs Copilot context
backend/proto/**/*.pb.go
backend/proto/**/*_grpc.pb.go

# Dependencies & build artefacts
node_modules/
frontend/dist/
bin/
.air/

# Logs & coverage
**/*.log
**/coverage.out
**/coverage/

# Infrastructure config (not useful as code context)
deployments/grafana/dashboards/*.json
```

---

## Upgrade path

If you hit weekly limits regularly, the options in order of cost:

1. **Copilot Pro+** — significantly higher weekly token budget.
2. **Usage-based billing** (GitHub AI Credits, rolling out June 2026) —
   pay per token above the plan cap instead of being hard-blocked.
