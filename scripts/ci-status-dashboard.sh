#!/bin/bash

# ============================================================================
# WanderPlan - Automated CI/CD Status & Monitoring Dashboard
# ============================================================================

set -e

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)"
LOGS_DIR="$PROJECT_ROOT/.ci-logs"
ARTIFACTS_DIR="$PROJECT_ROOT/.ci-artifacts"

# Create log directories if they don't exist
mkdir -p "$LOGS_DIR" "$ARTIFACTS_DIR"

# Color codes
BLUE='\033[0;34m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# ============================================================================
# LOGGING FUNCTIONS
# ============================================================================

log_section() {
  echo -e "\n${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
  echo -e "${BLUE}  $1${NC}"
  echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"
}

log_success() {
  echo -e "${GREEN}✅ $1${NC}"
}

log_error() {
  echo -e "${RED}❌ $1${NC}"
}

log_warning() {
  echo -e "${YELLOW}⚠️  $1${NC}"
}

log_info() {
  echo -e "${BLUE}ℹ️  $1${NC}"
}

# ============================================================================
# STATUS CHECKS
# ============================================================================

check_git_status() {
  log_section "Git Repository Status"

  if [ -z "$(git status --porcelain)" ]; then
    log_success "Working directory is clean"
  else
    log_warning "Uncommitted changes detected:"
    git status --short | head -10
  fi

  echo -e "\n${BLUE}Recent commits:${NC}"
  git log --oneline -5

  echo -e "\n${BLUE}Current branch:${NC}"
  git branch --show-current
}

check_frontend_status() {
  log_section "Frontend Status"

  cd "$PROJECT_ROOT/frontend"

  # Check dependencies
  if [ -d "node_modules" ]; then
    log_success "Dependencies installed"
  else
    log_warning "Dependencies not installed. Run: npm install"
  fi

  # Check recent build
  if [ -d "dist" ]; then
    BUILD_TIME=$(stat -f%Sm -t%Y-%m-%d\ %H:%M:%S "dist" 2>/dev/null || stat -c%y "dist" 2>/dev/null | cut -d' ' -f1-2)
    log_success "Build available (Last: $BUILD_TIME)"
  else
    log_warning "No build found. Run: npm run build"
  fi

  cd "$PROJECT_ROOT"
}

check_backend_status() {
  log_section "Backend Status"

  cd "$PROJECT_ROOT/backend"

  # Check Go version
  GO_VERSION=$(go version | awk '{print $3}')
  log_success "Go version: $GO_VERSION"

  # Check dependencies
  if go mod verify >/dev/null 2>&1; then
    log_success "Go dependencies verified"
  else
    log_warning "Go dependencies need update. Run: go mod tidy"
  fi

  # Check for build artifacts
  if [ -f "api-gateway" ] || [ -d "bin" ]; then
    log_success "Build artifacts found"
  else
    log_warning "No build artifacts. Run: make build-backend"
  fi

  cd "$PROJECT_ROOT"
}

check_services_running() {
  log_section "Running Services"

  if ! command -v docker-compose &>/dev/null; then
    log_warning "Docker Compose not installed"
    return 1
  fi

  RUNNING=$(docker-compose ps --services --filter "status=running" 2>/dev/null | wc -l)
  TOTAL=$(docker-compose ps --services 2>/dev/null | wc -l)

  if [ "$RUNNING" -eq 0 ]; then
    log_warning "No services running. Use: make docker-up"
  else
    log_success "$RUNNING/$TOTAL services running"
    docker-compose ps
  fi
}

check_github_workflows() {
  log_section "GitHub Actions Workflows"

  if ! command -v gh &>/dev/null; then
    log_warning "GitHub CLI not installed"
    return 1
  fi

  log_info "Recent workflow runs:"
  gh run list --limit 5 2>/dev/null || log_warning "Could not fetch workflow runs"

  log_info "\nTo see detailed logs:"
  echo "  gh run view <run-id>"
  echo "  gh run view <run-id> --log"
}

# ============================================================================
# TEST RESULTS
# ============================================================================

check_test_results() {
  log_section "Test Results"

  FRONTEND_COVERAGE=$([ -f "frontend/coverage/coverage-final.json" ] && echo "Available" || echo "Not found")
  BACKEND_COVERAGE=$([ -f "backend/coverage.out" ] && echo "Available" || echo "Not found")

  echo -e "${BLUE}Frontend Coverage:${NC} $FRONTEND_COVERAGE"
  echo -e "${BLUE}Backend Coverage:${NC} $BACKEND_COVERAGE"

  if [ -f "backend/coverage.out" ]; then
    COVERAGE_PCT=$(go tool cover -func=backend/coverage.out 2>/dev/null | tail -1 | awk '{print $NF}' || echo "Unknown")
    echo -e "${BLUE}Backend Coverage %:${NC} $COVERAGE_PCT"
  fi
}

# ============================================================================
# LOGS ANALYSIS
# ============================================================================

analyze_logs() {
  log_section "Recent Error Logs"

  if [ ! -d "$LOGS_DIR" ]; then
    log_info "No logs found"
    return
  fi

  # Find recent errors
  ERRORS=$(find "$LOGS_DIR" -type f -name "*.log" -mtime -1 -exec grep -l "ERROR\|FATAL" {} \; 2>/dev/null | wc -l)

  if [ "$ERRORS" -gt 0 ]; then
    log_warning "Found $ERRORS log files with errors"
    find "$LOGS_DIR" -type f -name "*.log" -mtime -1 | head -5
  else
    log_success "No recent errors detected"
  fi
}

# ============================================================================
# CONFIGURATION CHECK
# ============================================================================

check_configuration() {
  log_section "Configuration Status"

  # Check .env
  if [ -f "$PROJECT_ROOT/.env" ]; then
    ENV_VARS=$(grep -v "^#" "$PROJECT_ROOT/.env" | grep -v "^$" | wc -l)
    log_success ".env file configured ($ENV_VARS variables)"
  else
    log_warning ".env file not found. Run: cp .env.example .env"
  fi

  # Check workflows
  if [ -f "$PROJECT_ROOT/.github/workflows/ci-cd.yml" ]; then
    log_success "CI/CD workflow configured"
  else
    log_error "CI/CD workflow not found"
  fi

  if [ -f "$PROJECT_ROOT/.github/workflows/error-recovery.yml" ]; then
    log_success "Error recovery workflow configured"
  else
    log_error "Error recovery workflow not found"
  fi

  # Check docker-compose
  if [ -f "$PROJECT_ROOT/docker-compose.yml" ]; then
    SERVICES=$(grep "^  [a-z].*:$" "$PROJECT_ROOT/docker-compose.yml" | wc -l)
    log_success "Docker Compose configured ($SERVICES services)"
  else
    log_warning "docker-compose.yml not found"
  fi
}

# ============================================================================
# PERFORMANCE METRICS
# ============================================================================

show_performance_metrics() {
  log_section "Performance Metrics"

  # Git repo size
  REPO_SIZE=$(du -sh "$PROJECT_ROOT" 2>/dev/null | cut -f1)
  echo -e "${BLUE}Repository size:${NC} $REPO_SIZE"

  # Frontend size
  if [ -d "$PROJECT_ROOT/frontend/dist" ]; then
    FRONTEND_SIZE=$(du -sh "$PROJECT_ROOT/frontend/dist" 2>/dev/null | cut -f1)
    echo -e "${BLUE}Frontend build size:${NC} $FRONTEND_SIZE"
  fi

  # Number of commits
  COMMIT_COUNT=$(git rev-list --count HEAD 2>/dev/null || echo "Unknown")
  echo -e "${BLUE}Total commits:${NC} $COMMIT_COUNT"

  # Current branch info
  BRANCH=$(git rev-parse --abbrev-ref HEAD)
  COMMITS_AHEAD=$(git rev-list --count origin/$BRANCH..$BRANCH 2>/dev/null || echo "0")
  echo -e "${BLUE}Current branch:${NC} $BRANCH ($([ "$COMMITS_AHEAD" -eq 0 ] && echo "up to date" || echo "$COMMITS_AHEAD commits ahead"))"
}

# ============================================================================
# RECOMMENDATIONS
# ============================================================================

show_recommendations() {
  log_section "Recommendations"

  RECOMMENDATIONS=()

  # Check if tests are passing
  if [ ! -f "backend/coverage.out" ]; then
    RECOMMENDATIONS+=("Run backend tests: make test-backend")
  fi

  # Check if frontend is built
  if [ ! -d "frontend/dist" ]; then
    RECOMMENDATIONS+=("Build frontend: make build-frontend")
  fi

  # Check if Docker stack is running
  if ! docker-compose ps 2>/dev/null | grep -q "running" ; then
    RECOMMENDATIONS+=("Start Docker stack: make docker-up")
  fi

  # Check for uncommitted changes
  if [ -n "$(git status --porcelain)" ]; then
    RECOMMENDATIONS+=("Commit pending changes: git add . && git commit -m 'message'")
  fi

  # Check for unmerged PRs
  if command -v gh &>/dev/null; then
    PR_COUNT=$(gh pr list --search "is:open" --json number 2>/dev/null | grep -c number || echo "0")
    if [ "$PR_COUNT" -gt 0 ]; then
      RECOMMENDATIONS+=("Review open PRs: gh pr list")
    fi
  fi

  if [ ${#RECOMMENDATIONS[@]} -eq 0 ]; then
    log_success "Everything looks good! No recommendations."
  else
    for i in "${!RECOMMENDATIONS[@]}"; do
      echo "  $((i + 1)). ${RECOMMENDATIONS[$i]}"
    done
  fi
}

# ============================================================================
# GENERATE REPORT
# ============================================================================

generate_report() {
  local REPORT_FILE="$ARTIFACTS_DIR/ci-status-$(date +%Y%m%d-%H%M%S).txt"

  {
    echo "════════════════════════════════════════════════════════════"
    echo "   WanderPlan CI/CD Status Report"
    echo "   Generated: $(date)"
    echo "════════════════════════════════════════════════════════════"
    echo ""

    echo "GIT REPOSITORY STATUS"
    echo "────────────────────────────────────────────────────────────"
    git log --oneline -1
    echo "Branch: $(git rev-parse --abbrev-ref HEAD)"
    echo "Status: $([ -z "$(git status --porcelain)" ] && echo "Clean" || echo "Dirty")"
    echo ""

    echo "SERVICES STATUS"
    echo "────────────────────────────────────────────────────────────"
    if command -v docker-compose &>/dev/null; then
      docker-compose ps 2>/dev/null || echo "Docker not available"
    fi
    echo ""

    echo "BUILD STATUS"
    echo "────────────────────────────────────────────────────────────"
    echo "Frontend: $([ -d "frontend/dist" ] && echo "✅ Built" || echo "❌ Not built")"
    echo "Backend: $([ -f "backend/coverage.out" ] && echo "✅ Tested" || echo "❌ Not tested")"
    echo ""

  } | tee "$REPORT_FILE"

  log_success "Report saved to: $REPORT_FILE"
}

# ============================================================================
# HEALTH CHECK
# ============================================================================

health_check() {
  log_section "System Health Check"

  HEALTH_SCORE=0
  MAX_SCORE=10

  # Git clean
  if [ -z "$(git status --porcelain)" ]; then
    ((HEALTH_SCORE++))
    log_success "Git repository is clean"
  else
    log_warning "Git repository has uncommitted changes"
  fi

  # Frontend built
  if [ -d "frontend/dist" ]; then
    ((HEALTH_SCORE++))
    log_success "Frontend is built"
  fi

  # Backend tested
  if [ -f "backend/coverage.out" ]; then
    ((HEALTH_SCORE++))
    log_success "Backend is tested"
  fi

  # Dependencies
  if [ -d "frontend/node_modules" ] && go mod verify >/dev/null 2>&1; then
    ((HEALTH_SCORE++))
    log_success "Dependencies are installed"
  fi

  # Docker available
  if command -v docker &>/dev/null; then
    ((HEALTH_SCORE++))
    log_success "Docker is available"
  fi

  # Docker Compose running
  if docker-compose ps 2>/dev/null | grep -q "running"; then
    ((HEALTH_SCORE++))
    log_success "Services are running"
  fi

  # Workflows configured
  if [ -f ".github/workflows/ci-cd.yml" ]; then
    ((HEALTH_SCORE++))
    log_success "CI/CD pipeline configured"
  fi

  # Environment configured
  if [ -f ".env" ]; then
    ((HEALTH_SCORE++))
    log_success "Environment configured"
  fi

  echo ""
  echo -e "${BLUE}Overall Health Score: $HEALTH_SCORE / $MAX_SCORE${NC}"

  if [ "$HEALTH_SCORE" -ge 7 ]; then
    log_success "System is healthy! 💪"
  elif [ "$HEALTH_SCORE" -ge 5 ]; then
    log_warning "System is mostly healthy. Some configuration needed."
  else
    log_error "System needs attention. Several issues detected."
  fi
}

# ============================================================================
# MAIN DASHBOARD
# ============================================================================

show_dashboard() {
  clear

  cat << "EOF"
╔════════════════════════════════════════════════════════════════╗
║        WanderPlan - CI/CD Status & Monitoring Dashboard        ║
╚════════════════════════════════════════════════════════════════╝
EOF

  # Run all checks
  cd "$PROJECT_ROOT"
  
  check_git_status
  check_frontend_status
  check_backend_status
  check_configuration
  check_services_running
  check_test_results
  show_performance_metrics
  health_check
  show_recommendations

  log_section "Next Steps"
  echo "Run one of the following commands:"
  echo ""
  echo "  make validate-all    - Run complete validation"
  echo "  make docker-up       - Start service stack"
  echo "  make dev             - Start development environment"
  echo "  make ci              - Run CI pipeline"
  echo ""
}

# ============================================================================
# MAIN
# ============================================================================

main() {
  if [ "$1" == "--report" ]; then
    generate_report
  elif [ "$1" == "--health" ]; then
    health_check
  elif [ "$1" == "--watch" ]; then
    while true; do
      show_dashboard
      sleep 30
      clear
    done
  else
    show_dashboard
  fi
}

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  main "$@"
fi
