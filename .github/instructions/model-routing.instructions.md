---
applyTo: "**"
---
# Model Selection — Online vs Offline (Ollama)

## TL;DR — Start every request with `/delegate`

`/delegate` runs on **Qwen3.6 35B (local)**. It classifies your task and either:
- Executes it directly on Qwen3 (zero cloud tokens), or
- Emits a copy-paste escalation block naming the exact cloud prompt to use next.

This is the "switch between models within a session" mechanism.

---

## Decision Matrix

| Task type | Model | Prompt |
|---|---|---|
| **Any task (auto-route)** | **Qwen3.6 35B → cloud if needed** | **`/delegate` ← start here** |
| Quick explain / Q&A | Qwen3.6 35B (offline) | `/explain-offline` |
| Single-function refactor | Qwen3.6 35B (offline) | `/refactor-offline` |
| Format / style / lint fix | Qwen3.6 35B (offline) | `/offline-task` |
| Simple boilerplate | Qwen3.6 35B (offline) | `/offline-task` |
| Privacy-sensitive review | Qwen3.6 35B (offline) | `/offline-task` |
| **Build + lint + test + Render health** | **Qwen3.6 35B (offline)** | **`/check` ← run after every cloud turn** |
| **Go mod tidy / bun install / pkg sync** | **Qwen3.6 35B (offline)** | **`/pkg`** |
| **Render.com deploy health polling** | **Qwen3.6 35B (offline)** | **`/wait-render`** |
| Multi-file feature | Claude Sonnet 4.5 (cloud) | `/implement-feature` |
| gRPC / proto design | Claude Sonnet 4.5 (cloud) | `/online-task` |
| Full test suites | Claude Sonnet 4.5 (cloud) | `/online-task` |
| Security audit (OWASP/CVE) | Claude Sonnet 4.5 (cloud) | `/security-review-online` |
| Architecture decisions | Claude Sonnet 4.5 (cloud) | `/online-task` |
| Complex runtime debugging | Claude Sonnet 4.5 (cloud) | `/online-task` |
| Any task needing tool calls | Claude Sonnet 4.5 (cloud) | `/implement-feature` |

---

## How to Invoke

In the Copilot Chat input box, type `/` and select the prompt file:

```
/delegate            — Qwen3.6 35B smart router (START HERE)
/check               — Qwen3.6 35B, per-iteration gate (build+lint+test+Render)
/pkg                 — Qwen3.6 35B, Go mod + frontend package sync
/wait-render         — Qwen3.6 35B, Render.com deploy health poller
/offline-task        — Qwen3.6 35B, ask mode, no tool calls
/online-task         — Claude Sonnet 4.5, agent mode, full tools
/explain-offline     — Qwen3.6 35B, quick Q&A
/refactor-offline    — Qwen3.6 35B, edit mode
/implement-feature   — Claude Sonnet 4.5, full agentic build session
/security-review-online — Claude Sonnet 4.5, OWASP-focused audit
```

---

## How model switching works mid-session

1. User types: `/delegate implement the trip search endpoint`
2. Qwen3.6 35B analyses the request
3. **If simple** → Qwen3 handles it directly (offline, private, free)
4. **If complex** → Qwen3 outputs:
   ```
   ╔═══════════════════════════════════════╗
   ║  🔁  ESCALATING TO CLOUD MODEL        ║
   ╚═══════════════════════════════════════╝
   Reason: requires file writes across trip-service + proto regeneration
   Prompt: /implement-feature
   ...paste block...
   ```
5. User pastes the suggested `/implement-feature ...` line → Claude Sonnet takes over with full tools

---

## Ollama Setup (required for offline prompts)

1. Install Ollama: https://ollama.com/download
2. Active model (already pulled):
   ```bash
   ollama list    # should show qwen3.6:35b
   ```
3. Ollama runs on `http://localhost:11434` — VS Code detects it automatically.
4. In VS Code Copilot Chat, open the **model picker** to confirm `ollama/qwen3.6:35b` appears under **Local Models**.

---

## Changing the Ollama Model

Edit any offline prompt file (e.g. `offline-task.prompt.md`) and change:
```yaml
model: "ollama:qwen3.6:35b"
```
to any model you have pulled, e.g. `"ollama:llama3"`, `"ollama:mistral"`, `"ollama:codellama"`.

## Changing the Cloud Model

Edit `online-task.prompt.md` or `implement-feature.prompt.md` and change:
```yaml
model: claude-sonnet-4-5
```
to any GitHub Copilot cloud model: `gpt-4o`, `o3`, `o4-mini`, `claude-opus-4-5`, etc.
