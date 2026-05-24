#!/usr/bin/env bash
# .vscode/hooks/guard-commands.sh
# Called by the PreToolUse hook before Copilot runs any shell command.
# Exit 1 to BLOCK the command; exit 0 to ALLOW it.

COMMAND="$1"

# ── Blocked patterns ────────────────────────────────────────────────────────
BLOCKED_PATTERNS=(
  "DROP DATABASE"
  "DROP SCHEMA"
  "dropdb"
  "rm -rf /"
  "rm -rf ~"
  "git push --force"
  "git push -f"
  "git push origin main"
  "git push origin master"
  "curl.*| bash"
  "wget.*| bash"
  "chmod 777"
  "> /etc/"
  "truncate.*wanderplan"
)

for pattern in "${BLOCKED_PATTERNS[@]}"; do
  if echo "$COMMAND" | grep -qi "$pattern"; then
    echo "🚫 BLOCKED by WanderPlan hook: command matches forbidden pattern \"$pattern\""
    echo "   Command: $COMMAND"
    echo "   If this is intentional, run it manually in the terminal."
    exit 1
  fi
done

# Allow everything else
exit 0
