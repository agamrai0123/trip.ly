#!/usr/bin/env bash
# scripts/run.sh — Start the full WanderPlan local stack.
# Equivalent to: make run / make local-up
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

if [[ ! -f .env ]]; then
  echo "❌  .env not found. Run: cp .env.example .env  and fill in secrets."
  exit 1
fi

echo "▶  Starting WanderPlan local stack..."
docker compose up -d --build

echo "⏳  Waiting for api-gateway to become healthy (max 120s)..."
for i in $(seq 1 24); do
  if curl -sf http://localhost:8080/healthz -o /dev/null 2>/dev/null; then
    echo "✅  Stack is up!"
    break
  fi
  if [[ "$i" == "24" ]]; then
    echo "⚠️   api-gateway not healthy after 120s. Check logs: docker compose logs api-gateway"
  fi
  sleep 5
done

cat <<'EOF'

╔══════════════════════════════════════════════════════╗
║  WanderPlan local stack                              ║
╚══════════════════════════════════════════════════════╝
  Frontend:    http://localhost:5173
  API Gateway: http://localhost:8080
  Auth:        http://localhost:8081
  Trip:        http://localhost:8082
  User:        http://localhost:8083
  Collab:      http://localhost:8084
  Notify:      http://localhost:8085
  Search:      http://localhost:8086
  Prometheus:  http://localhost:9090
  Grafana:     http://localhost:3001  (admin/admin)
EOF
