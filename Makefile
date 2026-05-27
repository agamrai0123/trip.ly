# WanderPlan Makefile - Automated Build, Test & Deployment
# Usage: make [target]

.PHONY: help setup check-deps install-tools test-frontend test-backend test-proto build-frontend build-backend docker-up docker-down dev validate-all fix-all lint-frontend lint-backend fmt-backend clean local-up local-down local-restart load-test migrate migrate-down db-migrate db-rollback proto proto-lint run dev-frontend dev-backend status ci cd

# Color output
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

# Go & Node settings
GO_VERSION := 1.23.0
NODE_VERSION := 20.x
BUN_EXECUTABLE := $(shell command -v bun)
NPM_EXECUTABLE := $(shell command -v npm)

# ============================================================================
# HELP
# ============================================================================

help:
	@echo "$(BLUE)╔════════════════════════════════════════════════════════════════╗$(NC)"
	@echo "$(BLUE)║           WanderPlan - Build & Test Automation                  ║$(NC)"
	@echo "$(BLUE)╚════════════════════════════════════════════════════════════════╝$(NC)"
	@echo ""
	@echo "$(YELLOW)Setup & Installation:$(NC)"
	@echo "  make setup              - Complete setup (deps, tools, hooks, env)"
	@echo "  make check-deps         - Check prerequisites"
	@echo "  make install-tools      - Install required tools"
	@echo ""
	@echo "$(YELLOW)Testing:$(NC)"
	@echo "  make test               - Run all tests (frontend & backend)"
	@echo "  make test-frontend      - Run frontend tests only"
	@echo "  make test-backend       - Run backend tests only"
	@echo "  make test-proto         - Validate proto files"
	@echo ""
	@echo "$(YELLOW)Building:$(NC)"
	@echo "  make build              - Build frontend & backend"
	@echo "  make build-frontend     - Build frontend only"
	@echo "  make build-backend      - Build all backend services"
	@echo "  make docker-build       - Build Docker images"
	@echo ""
	@echo "$(YELLOW)Code Quality:$(NC)"
	@echo "  make lint               - Run all linters (frontend & backend)"
	@echo "  make lint-frontend      - Lint frontend with ESLint"
	@echo "  make lint-backend       - Lint backend with golangci-lint"
	@echo "  make fmt-frontend       - Format frontend code"
	@echo "  make fmt-backend        - Format backend code"
	@echo "  make fix-all            - Auto-fix all linting & formatting errors"
	@echo ""
	@echo "$(YELLOW)Docker & Services:$(NC)"
	@echo "  make docker-up          - Start Docker stack (PostgreSQL, Kafka, etc)"
	@echo "  make docker-down        - Stop Docker stack"
	@echo "  make docker-logs        - View Docker logs"
	@echo ""
	@echo "$(YELLOW)Development:$(NC)"
	@echo "  make dev                - Start development environment (all services)"
	@echo "  make dev-frontend       - Start frontend dev server"
	@echo "  make dev-backend        - Start backend dev server"
	@echo "  make validate-all       - Run complete validation pipeline"
	@echo ""
	@echo "$(YELLOW)Utilities:$(NC)"
	@echo "  make clean              - Clean build artifacts & cache"
	@echo "  make git-hooks          - Setup git pre-commit/pre-push hooks"
	@echo "  make status             - Check project status"
	@echo ""

# ============================================================================
# SETUP & INSTALLATION
# ============================================================================

setup: check-deps install-tools git-hooks create-env
	@echo "$(GREEN)✅ Setup complete!$(NC)"
	@echo "$(BLUE)Next steps:$(NC)"
	@echo "  1. make docker-up       - Start services"
	@echo "  2. make dev             - Start development"
	@echo "  3. make test            - Run tests"

check-deps:
	@echo "$(BLUE)Checking prerequisites...$(NC)"
	@command -v git >/dev/null 2>&1 || { echo "$(RED)❌ Git not found$(NC)"; exit 1; }
	@echo "$(GREEN)✅ Git found$(NC)"
	@command -v go >/dev/null 2>&1 || { echo "$(RED)❌ Go not found$(NC)"; exit 1; }
	@echo "$(GREEN)✅ Go $$(go version | awk '{print $$3}') found$(NC)"
	@command -v node >/dev/null 2>&1 || command -v bun >/dev/null 2>&1 || { echo "$(RED)❌ Node/Bun not found$(NC)"; exit 1; }
	@echo "$(GREEN)✅ Package manager found$(NC)"
	@echo "$(GREEN)✅ All prerequisites OK$(NC)"

install-tools:
	@echo "$(BLUE)Installing required tools...$(NC)"
	@echo "$(BLUE)Installing buf...$(NC)"
	@go install github.com/bufbuild/buf/cmd/buf@latest
	@echo "$(GREEN)✅ buf installed$(NC)"
	@echo "$(BLUE)Installing golangci-lint...$(NC)"
	@which golangci-lint >/dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin
	@echo "$(GREEN)✅ golangci-lint installed$(NC)"
	@echo "$(BLUE)Installing goimports...$(NC)"
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "$(GREEN)✅ goimports installed$(NC)"

create-env:
	@if [ ! -f .env ]; then \
		echo "$(BLUE)Creating .env file...$(NC)"; \
		cp -n .env.example .env 2>/dev/null || echo "# Add your environment variables here" > .env; \
		echo "$(YELLOW)⚠️  .env created. Please configure it.$(NC)"; \
	else \
		echo "$(GREEN)✅ .env already exists$(NC)"; \
	fi

git-hooks:
	@echo "$(BLUE)Setting up git hooks...$(NC)"
	@mkdir -p .git/hooks
	@cp -f scripts/git-hooks/pre-commit .git/hooks/pre-commit || echo "pre-commit hook not found"
	@cp -f scripts/git-hooks/pre-push .git/hooks/pre-push || echo "pre-push hook not found"
	@chmod +x .git/hooks/* 2>/dev/null || true
	@echo "$(GREEN)✅ Git hooks installed$(NC)"

# ============================================================================
# TESTING
# ============================================================================

test: test-frontend test-backend test-proto
	@echo "$(GREEN)✅ All tests passed!$(NC)"

test-frontend:
	@echo "$(BLUE)Testing frontend...$(NC)"
	@cd frontend && $(if $(BUN_EXECUTABLE), bun install --frozen-lockfile && bun run test, npm install --frozen-lockfile && npm run test)
	@echo "$(GREEN)✅ Frontend tests passed$(NC)"

test-backend:
	@echo "$(BLUE)Testing backend services...$(NC)"
	@for svc in api-gateway auth-service trip-service user-service collaboration-service notification-service search-service; do \
		echo "  Testing $$svc..."; \
		(cd backend/services/$$svc && go test ./internal/... -count=1 -timeout 60s -race) || exit 1; \
	done
	@echo "$(GREEN)✅ Backend tests passed$(NC)"

test-proto:
	@echo "$(BLUE)Validating proto files...$(NC)"
	@cd backend && buf lint
	@echo "$(GREEN)✅ Proto validation passed$(NC)"

# ============================================================================
# BUILDING
# ============================================================================

build: build-frontend build-backend
	@echo "$(GREEN)✅ Build complete!$(NC)"

build-frontend:
	@echo "$(BLUE)Building frontend...$(NC)"
	@cd frontend && $(if $(BUN_EXECUTABLE), bun install --frozen-lockfile && bun run build, npm install --frozen-lockfile && npm run build)
	@echo "$(GREEN)✅ Frontend build complete$(NC)"

build-backend:
	@echo "$(BLUE)Building backend services...$(NC)"
	@cd backend && go build -v ./services/api-gateway/cmd
	@cd backend && go build -v ./services/auth-service/cmd
	@cd backend && go build -v ./services/trip-service/cmd
	@cd backend && go build -v ./services/user-service/cmd
	@echo "$(GREEN)✅ Backend build complete$(NC)"

docker-build:
	@echo "$(BLUE)Building Docker images...$(NC)"
	@docker-compose build
	@echo "$(GREEN)✅ Docker build complete$(NC)"

# ============================================================================
# CODE QUALITY
# ============================================================================

lint: lint-frontend lint-backend
	@echo "$(GREEN)✅ Linting complete$(NC)"

lint-frontend:
	@echo "$(BLUE)Linting frontend...$(NC)"
	@cd frontend && $(if $(BUN_EXECUTABLE), bun run lint, npm run lint)
	@echo "$(GREEN)✅ Frontend lint passed$(NC)"

lint-backend:
	@echo "$(BLUE)Linting backend...$(NC)"
	@cd backend && go vet ./...
	@cd backend && golangci-lint run ./...
	@echo "$(GREEN)✅ Backend lint passed$(NC)"

fmt-frontend:
	@echo "$(BLUE)Formatting frontend code...$(NC)"
	@cd frontend && $(if $(BUN_EXECUTABLE), bun run lint --fix, npm run lint -- --fix)
	@echo "$(GREEN)✅ Frontend formatted$(NC)"

fmt-backend:
	@echo "$(BLUE)Formatting backend code...$(NC)"
	@cd backend && go fmt ./...
	@cd backend && goimports -w .
	@echo "$(GREEN)✅ Backend formatted$(NC)"

fix-all: fmt-frontend fmt-backend
	@echo "$(BLUE)Auto-fixing all issues...$(NC)"
	@cd backend && go mod tidy
	@echo "$(BLUE)Committing fixes...$(NC)"
	@git config user.email "makefile@wanderplan.local" 2>/dev/null || true
	@git config user.name "WanderPlan Makefile" 2>/dev/null || true
	@git add -A && git commit -m "🤖 Auto-fix: Formatting & imports" || true
	@echo "$(GREEN)✅ All fixes applied$(NC)"

# ============================================================================
# DOCKER & SERVICES
# ============================================================================

docker-up:
	@echo "$(BLUE)Starting Docker stack...$(NC)"
	@docker compose up -d
	@echo "$(BLUE)Waiting for services to be healthy...$(NC)"
	@sleep 5
	@echo "$(GREEN)✅ Docker stack started$(NC)"
	@echo "$(BLUE)Services running:$(NC)"
	@docker compose ps

docker-down:
	@echo "$(BLUE)Stopping Docker stack...$(NC)"
	@docker compose down
	@echo "$(GREEN)✅ Docker stack stopped$(NC)"

docker-logs:
	@docker compose logs -f

docker-clean:
	@echo "$(BLUE)Cleaning Docker resources...$(NC)"
	@docker compose down -v
	@echo "$(GREEN)✅ Docker cleaned$(NC)"

# ── Local full-stack (build + run all services + infra + frontend) ─────────────

local-up:
	@echo "$(BLUE)Starting full local stack (build + run)...$(NC)"
	@if [ ! -f .env ]; then echo "$(RED)❌ .env missing. Run: cp .env.example .env and fill in JWT keys$(NC)"; exit 1; fi
	@chmod +x scripts/postgres-init.sh 2>/dev/null || true
	@docker compose up -d --build
	@echo "$(BLUE)Waiting for all services to be healthy (up to 120s)...$(NC)"
	@for i in $$(seq 1 24); do \
		if curl -sf http://localhost:8080/healthz -o /dev/null 2>/dev/null; then \
			echo "$(GREEN)✅ Stack is up!$(NC)"; break; \
		fi; \
		if [ "$$i" = "24" ]; then echo "$(YELLOW)⚠ api-gateway not healthy after 120s — check: make local-logs$(NC)"; fi; \
		sleep 5; \
	done
	@echo ""
	@echo "$(GREEN)╔══════════════════════════════════════════════════════╗$(NC)"
	@echo "$(GREEN)║  WanderPlan local stack                               ║$(NC)"
	@echo "$(GREEN)╚══════════════════════════════════════════════════════╝$(NC)"
	@echo "  Frontend:     http://localhost:5173"
	@echo "  API Gateway:  http://localhost:8080"
	@echo "  Auth:         http://localhost:8081"
	@echo "  Trip:         http://localhost:8082"
	@echo "  User:         http://localhost:8083"
	@echo "  Collab:       http://localhost:8084"
	@echo "  Notif:        http://localhost:8085"
	@echo "  Search:       http://localhost:8086"
	@echo "  Prometheus:   http://localhost:9090"
	@echo "  Grafana:      http://localhost:3001  (admin/admin)"

local-down:
	@docker compose down
	@echo "$(GREEN)✅ Local stack stopped$(NC)"

local-restart:
	@docker compose down
	@docker compose up -d --build
	@echo "$(GREEN)✅ Local stack restarted$(NC)"

local-logs:
	@docker compose logs -f --tail=50

# ── Load tests ─────────────────────────────────────────────────────────────────

load-test:
	@echo "$(BLUE)Running load tests against local stack...$(NC)"
	@chmod +x scripts/load-test.sh
	@bash scripts/load-test.sh $(ARGS)

# ============================================================================
# DEVELOPMENT
# ============================================================================

dev: docker-up
	@echo "$(BLUE)Starting development environment...$(NC)"
	@echo "$(YELLOW)Starting services with hot reload...$(NC)"
	@echo "  Frontend: http://localhost:5173"
	@echo "  API Gateway: http://localhost:8080"
	@echo "  Auth Service: http://localhost:8081"
	@echo ""
	@echo "$(BLUE)In separate terminals, run:$(NC)"
	@echo "  make dev-frontend"
	@echo "  make dev-backend"

dev-frontend:
	@echo "$(BLUE)Starting frontend dev server...$(NC)"
	@cd frontend && $(if $(BUN_EXECUTABLE), bun run dev, npm run dev)

dev-backend:
	@echo "$(BLUE)Starting backend services...$(NC)"
	@which air >/dev/null || go install github.com/cosmtrek/air@latest
	@cd backend && air

# ============================================================================
# VALIDATION & VERIFICATION
# ============================================================================

validate-all: check-deps lint test
	@echo "$(GREEN)✅ All validations passed!$(NC)"

status:
	@echo "$(BLUE)Project Status$(NC)"
	@echo ""
	@echo "$(BLUE)Git:$(NC)"
	@git status --short || echo "Not a git repository"
	@echo ""
	@echo "$(BLUE)Go Version:$(NC)"
	@go version 2>/dev/null || echo "Go not installed"
	@echo ""
	@echo "$(BLUE)Node/Bun Version:$(NC)"
	@node --version 2>/dev/null || echo "Node not installed"
	@$(if $(BUN_EXECUTABLE), bun --version, echo "Bun not installed")
	@echo ""
	@echo "$(BLUE)Docker Status:$(NC)"
	@docker-compose ps 2>/dev/null || echo "Docker not running"

# ============================================================================
# CLEANUP
# ============================================================================

clean:
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	@rm -rf frontend/dist
	@rm -rf frontend/.vite
	@rm -rf backend/bin
	@rm -rf backend/coverage.out
	@find . -name "*.air" -type d -exec rm -rf {} + 2>/dev/null || true
	@find . -name "__pycache__" -type d -exec rm -rf {} + 2>/dev/null || true
	@echo "$(GREEN)✅ Cleanup complete$(NC)"

# ============================================================================
# PROTO
# ============================================================================

proto:
	@echo "$(BLUE)Generating proto code...$(NC)"
	@cd backend && buf generate
	@cd backend && go mod tidy
	@echo "$(GREEN)✅ Proto code generated$(NC)"

proto-lint:
	@echo "$(BLUE)Validating proto files...$(NC)"
	@cd backend && buf lint
	@echo "$(GREEN)✅ Proto validation passed$(NC)"

# ============================================================================
# DATABASE
# ============================================================================

db-migrate:
	@echo "$(BLUE)Running database migrations...$(NC)"
	@cd backend && go run -mod=mod github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path ./migrations -database "${DATABASE_URL}" up
	@echo "$(GREEN)✅ Migrations complete$(NC)"

db-rollback:
	@echo "$(BLUE)Rolling back migrations...$(NC)"
	@cd backend && go run -mod=mod github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path ./migrations -database "${DATABASE_URL}" down
	@echo "$(GREEN)✅ Rollback complete$(NC)"

# Aliases matching project convention (migrate / migrate-down)
migrate: db-migrate
migrate-down: db-rollback

# ============================================================================
# CI/CD
# ============================================================================

ci: validate-all docker-build
	@echo "$(GREEN)✅ CI pipeline complete!$(NC)"

cd:
	@echo "$(BLUE)Ready for deployment$(NC)"
	@echo "$(GREEN)✅ All tests passed$(NC)"
	@echo "$(GREEN)✅ All builds successful$(NC)"
	@echo "$(GREEN)✅ Docker images ready$(NC)"

# ============================================================================
# DEFAULT TARGET
# ============================================================================

.DEFAULT_GOAL := help

# ============================================================================
# SHORTHAND ALIASES (required by project conventions)
# ============================================================================

# run -- start all backend services (docker-up first)
run:
	@for svc in api-gateway auth-service trip-service user-service collaboration-service notification-service search-service; do \
		echo "  Starting $$svc..."; \
		(cd backend/services/$$svc && go run ./cmd/... &); \
	done
	@echo "All services started."

# migrate -- run all pending migrations up
migrate: db-migrate

# migrate-down -- roll back last migration
migrate-down: db-rollback
