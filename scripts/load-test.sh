#!/bin/bash
# WanderPlan load test runner
# Runs a suite of load tests against the local stack using hey + k6.
# Usage: ./scripts/load-test.sh [--vus N] [--duration Xs] [--k6-only] [--hey-only]
#
# Prerequisites:
#   hey:  go install github.com/rakyll/hey@latest
#   k6:   docker pull grafana/k6  (or install k6 binary)
#
# The script reads AUTH_TOKEN from env, or generates one using gen-test-jwt.
# Make sure the docker-compose stack is running: make local-up

set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
VUS="${VUS:-20}"
DURATION="${DURATION:-30s}"
RUN_K6=true
RUN_HEY=true
RESULTS_DIR="tests/load/results/$(date +%Y%m%d_%H%M%S)"

# Parse flags
for arg in "$@"; do
  case $arg in
    --vus=*)   VUS="${arg#*=}" ;;
    --duration=*) DURATION="${arg#*=}" ;;
    --k6-only) RUN_HEY=false ;;
    --hey-only) RUN_K6=false ;;
  esac
done

mkdir -p "$RESULTS_DIR"

# ── Colour helpers ────────────────────────────────────────────────────────────
GREEN="\033[0;32m"; YELLOW="\033[1;33m"; RED="\033[0;31m"; NC="\033[0m"
info()  { echo -e "${GREEN}[INFO]${NC}  $*"; }
warn()  { echo -e "${YELLOW}[WARN]${NC}  $*"; }
error() { echo -e "${RED}[ERROR]${NC} $*"; }

# ── Check services are up ─────────────────────────────────────────────────────
info "Checking local stack health..."
SERVICES=(
  "api-gateway:$BASE_URL/healthz"
  "auth-service:http://localhost:8081/healthz"
  "trip-service:http://localhost:8082/healthz"
  "user-service:http://localhost:8083/healthz"
  "collaboration-service:http://localhost:8084/healthz"
  "notification-service:http://localhost:8085/healthz"
  "search-service:http://localhost:8086/healthz"
)

ALL_UP=true
for svc in "${SERVICES[@]}"; do
  name="${svc%%:*}"
  url="${svc#*:}"
  if curl -sf "$url" -o /dev/null 2>/dev/null; then
    echo -e "  ${GREEN}✓${NC} $name"
  else
    echo -e "  ${RED}✗${NC} $name (not reachable at $url)"
    ALL_UP=false
  fi
done

if [ "$ALL_UP" = false ]; then
  error "Some services are not running. Start the stack with: make local-up"
  exit 1
fi

# ── Get / generate AUTH_TOKEN ─────────────────────────────────────────────────
if [ -z "${AUTH_TOKEN:-}" ]; then
  info "Generating load-test JWT..."
  if command -v go &>/dev/null; then
    AUTH_TOKEN=$(go run ./tests/load/gen-test-jwt/)
    if [ -z "$AUTH_TOKEN" ]; then
      warn "JWT generation failed; authenticated endpoint tests will be skipped."
      AUTH_TOKEN=""
    else
      info "JWT generated (${#AUTH_TOKEN} bytes)."
    fi
  else
    warn "go not found; set AUTH_TOKEN env var to test authenticated endpoints."
    AUTH_TOKEN=""
  fi
fi

echo ""
info "Running load tests: VUS=$VUS DURATION=$DURATION"
echo "  Results directory: $RESULTS_DIR"
echo ""

# ── hey benchmarks (quick throughput) ────────────────────────────────────────
if [ "$RUN_HEY" = true ]; then
  HEY_BIN="$(command -v hey 2>/dev/null || echo "")"
  if [ -z "$HEY_BIN" ]; then
    HEY_BIN="$(go env GOPATH)/bin/hey"
  fi
  if [ ! -x "$HEY_BIN" ]; then
    warn "hey not found — skipping hey benchmarks. Install with: go install github.com/rakyll/hey@latest"
    RUN_HEY=false
  fi
fi

if [ "$RUN_HEY" = true ]; then
  info "=== hey: Health endpoint throughput ==="
  "$HEY_BIN" -n 5000 -c "$VUS" -q 0 "$BASE_URL/healthz" 2>&1 | tee "$RESULTS_DIR/hey_healthz.txt"

  if [ -n "$AUTH_TOKEN" ]; then
    info "=== hey: List trips throughput ==="
    "$HEY_BIN" -n 2000 -c "$VUS" -q 0 \
      -H "Authorization: Bearer $AUTH_TOKEN" \
      "$BASE_URL/api/v1/trips" 2>&1 | tee "$RESULTS_DIR/hey_list_trips.txt"

    info "=== hey: User profile throughput ==="
    "$HEY_BIN" -n 2000 -c "$VUS" -q 0 \
      -H "Authorization: Bearer $AUTH_TOKEN" \
      "$BASE_URL/api/v1/users/me" 2>&1 | tee "$RESULTS_DIR/hey_user_me.txt"
  fi
fi

# ── k6 scenario tests ─────────────────────────────────────────────────────────
if [ "$RUN_K6" = true ]; then
  K6_BIN="$(command -v k6 2>/dev/null || echo "")"
  K6_DOCKER=false

  if [ -z "$K6_BIN" ]; then
    if command -v docker &>/dev/null; then
      K6_DOCKER=true
      info "k6 binary not found — using Docker (grafana/k6)."
    else
      warn "Neither k6 nor docker found — skipping k6 tests."
      RUN_K6=false
    fi
  fi

  if [ "$RUN_K6" = true ]; then
    info "=== k6: Full scenario test ==="
    K6_ARGS="-e BASE_URL=$BASE_URL -e VUS=$VUS -e DURATION=$DURATION"
    if [ -n "$AUTH_TOKEN" ]; then
      K6_ARGS="$K6_ARGS -e AUTH_TOKEN=$AUTH_TOKEN"
    fi

    if [ "$K6_DOCKER" = true ]; then
      # Pull the image if not present (silently)
      docker pull grafana/k6 -q 2>/dev/null || true
      docker run --rm -i \
        --network host \
        -v "$(pwd)/tests/load:/tests/load:ro" \
        $( [ -n "$AUTH_TOKEN" ] && echo "-e AUTH_TOKEN=$AUTH_TOKEN" ) \
        -e BASE_URL="$BASE_URL" \
        -e VUS="$VUS" \
        -e DURATION="$DURATION" \
        grafana/k6 run /tests/load/load_test.js \
        2>&1 | tee "$RESULTS_DIR/k6_results.txt"
    else
      k6 run $K6_ARGS tests/load/load_test.js 2>&1 | tee "$RESULTS_DIR/k6_results.txt"
    fi
  fi
fi

# ── Prometheus metrics summary ────────────────────────────────────────────────
info "=== Pulling Prometheus metrics snapshot ==="
METRICS_ENDPOINTS=(
  "api-gateway:$BASE_URL/metrics"
  "auth-service:http://localhost:8081/metrics"
  "trip-service:http://localhost:8082/metrics"
  "user-service:http://localhost:8083/metrics"
  "collab:http://localhost:8084/metrics"
  "notif:http://localhost:8085/metrics"
  "search:http://localhost:8086/metrics"
)
for m in "${METRICS_ENDPOINTS[@]}"; do
  name="${m%%:*}"
  url="${m#*:}"
  if curl -sf "$url" -o "$RESULTS_DIR/metrics_${name}.txt" 2>/dev/null; then
    req_total=$(grep "^http_requests_total" "$RESULTS_DIR/metrics_${name}.txt" 2>/dev/null | awk '{sum+=$NF} END{print sum}')
    echo -e "  ${GREEN}✓${NC} $name: http_requests_total = ${req_total:-0}"
  else
    echo -e "  ${YELLOW}⚠${NC} $name: could not reach $url"
  fi
done

echo ""
info "Load test complete. Results saved to: $RESULTS_DIR"
echo ""
echo "View Grafana dashboard: http://localhost:3001  (admin/admin)"
echo "View Prometheus:        http://localhost:9090"
