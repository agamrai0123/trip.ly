# WanderPlan - AI-Powered Travel Itinerary Planning Platform

![Status](https://img.shields.io/badge/status-active-brightgreen)
![Go](https://img.shields.io/badge/Go-1.22+-blue)
![React](https://img.shields.io/badge/React-18+-blue)
![License](https://img.shields.io/badge/license-MIT-green)

**WanderPlan** is an intelligent travel itinerary planning platform combining React frontend with seven Go microservices. AI-assisted trip planning with real-time collaboration and smart recommendations.

---

## 🚀 Quick Start

### Prerequisites
- **Go 1.22+** - Backend runtime
- **Node.js v24+** & **Bun 1.3+** - Frontend tooling
- **Docker 29.0+** - Containerization
- **PostgreSQL 16+** - Database

### Setup

```bash
# 1. Clone repository
git clone https://github.com/agamrai0123/trip.ly.git
cd trip.ly

# 2. Install dependencies
cd frontend && bun install && cd ..
cd backend && go mod download && cd ..

# 3. Run pre-commit hooks setup
git config core.hooksPath .git/hooks
chmod +x .git/hooks/pre-commit .git/hooks/pre-push

# 4. Start development environment
docker-compose up -d
make dev
```

---

## 📚 Documentation

### Core Guides
| Document | Purpose |
|----------|---------|
| [Branching Strategy](docs/BRANCHING_STRATEGY.md) | Git workflow & branch management |
| [CI/CD Setup](docs/CI_CD_IMPLEMENTATION_COMPLETE.md) | Pipeline overview & status |
| [Setup Summary](docs/SETUP_SUMMARY.md) | Initial environment setup |

### Detailed References
| Document | Purpose |
|----------|---------|
| [CI/CD Fixes & Branching](docs/CI_CD_FIXES_AND_BRANCHING.md) | All pipeline fixes explained |
| [CI Quick Reference](docs/CI_QUICK_REFERENCE.md) | Quick command reference |
| [Automation Setup](docs/AUTOMATION_SETUP.md) | Previous automation setup |
| [CI/CD Index](docs/CI_CD_INDEX.md) | Pipeline components index |
| [Automation Complete](docs/AUTOMATION_SESSION_COMPLETE.md) | Session completion summary |

---

## 🏗️ Architecture

### Backend Services (Go)
```
api-gateway:8080        → Entry point
auth-service:8081       → Authentication (OAuth2, JWT)
trip-service:8082       → Trip management
user-service:8083       → User profiles
collaboration-service:8084 → Real-time collaboration
notification-service:8085  → Notifications & WebSocket
search-service:8086     → Place & destination search
```

### Communication
- **gRPC** - Synchronous inter-service calls (proto files in `backend/proto/`)
- **Kafka** - Async events (Sarama client)
- **HTTP** - REST API (Gin framework)

### Frontend (React 18 + TypeScript)
```
React 18              → UI framework
Vite                  → Build tool
Tailwind CSS          → Styling
React Router v6       → Routing
TanStack Query v5     → Server state
React Hook Form       → Form handling
shadcn/ui             → Components
```

---

## 🔄 Development Workflow

### 1. Create Feature Branch
```bash
git checkout develop
git checkout -b feature/my-feature
```

### 2. Develop & Commit
```bash
# Pre-commit hook validates automatically
git add .
git commit -m "feat: describe feature"
```

### 3. Push to GitHub
```bash
# Pre-push hook tests before allowing push
git push origin feature/my-feature
```

### 4. Create Pull Request
- Create PR: `feature/my-feature` → `develop`
- All tests run automatically (branch-merge.yml)
- Auto-merge on success

### 5. Release to Production
```bash
# develop → production (staging)
# production → main (production)
```

**Full details:** See [Branching Strategy](docs/BRANCHING_STRATEGY.md)

---

## ✅ CI/CD Pipeline

### Automated Testing at Every Level

#### **Local Commits** ← Pre-commit hook
- ESLint validation
- Go fmt check
- Code formatting

#### **Local Pushes** ← Pre-push hook  
- Frontend tests (Vitest)
- Backend tests
- Proto validation

#### **GitHub PRs** ← branch-merge.yml
- Comprehensive frontend tests
- Backend build & tests (all 7 services)
- Proto lint & generation
- Merge gate (all must pass)
- Auto-merge on success

#### **GitHub Pushes** ← ci-cd.yml
- Full pipeline validation
- Docker build & security scan
- Integration tests
- Error recovery

**Status:** ✅ All workflows operational

---

## 📊 Branches

| Branch | Purpose | Protection |
|--------|---------|-----------|
| **main** | Production releases | 2 approvals required |
| **production** | Staging environment | 1 approval required |
| **develop** | Feature integration | Auto-merge on tests |
| **feature/\*** | New features | No protection |
| **bugfix/\*** | Bug fixes | No protection |

---

## 🛠️ Common Commands

### Frontend
```bash
cd frontend

# Development
bun dev              # Start dev server
bun build            # Production build
bun test             # Run tests

# Linting
bun run lint         # Check ESLint
bun run lint:fix     # Fix ESLint errors
```

### Backend
```bash
cd backend

# Development
make dev             # All services with hot-reload
make dev-service SERVICE=api-gateway  # Single service

# Testing
make test            # All tests
make test-service SERVICE=trip-service # Single service

# Building
make build           # Build all services
make docker-build    # Docker images

# Database
make migrate         # Run migrations
make migrate-down    # Rollback migrations
```

### Root Commands
```bash
make fmt             # Format all code
make lint            # Lint all code
make test            # Test everything
make docker-compose  # Start full stack
```

---

## 🔐 Authentication

- **OAuth 2.0 PKCE** - Google & GitHub
- **JWT Tokens** - 15-minute expiry, RS256
- **Refresh Tokens** - Opaque, stored in httpOnly cookies
- **Validation** - api-gateway delegates to auth-service

---

## 📊 Observability

### Prometheus Metrics
- HTTP request count & latency
- Error rates
- Database query duration
- Goroutine count
- Kafka consumer lag

### Grafana Dashboards
- Service RPS
- P50/P95/P99 latency
- Error rates
- Database connections
- Consumer lag

**Access:** http://localhost:3000 (Grafana)

---

## 🗄️ Database

**Schema Migrations:**
```bash
make migrate         # Apply all
make migrate-down    # Rollback one
```

**Files:** `migrations/` directory with numbered SQL files

**Migration Tool:** golang-migrate

---

## 📦 Dependencies

### Go Modules
```
Gin Framework        → HTTP routing
pgx                  → PostgreSQL driver
Sarama               → Kafka client
GORM                 → ORM
protobuf/gRPC        → RPC framework
Prometheus           → Metrics
Zerolog              → Structured logging
```

### NPM Packages
```
React 18             → UI framework
Vite                 → Build tool
TanStack Query       → Server state
React Hook Form      → Forms
Tailwind CSS         → Styling
shadcn/ui            → Components
date-fns             → Date handling
recharts             → Charts
@dnd-kit             → Drag & drop
```

---

## 🧪 Testing

### Frontend
- **Framework:** Vitest + @testing-library/react
- **Coverage:** Minimum 80%
- **Command:** `cd frontend && bun test`

### Backend
- **Framework:** testify + table-driven tests
- **Integration:** TestContainers (PostgreSQL)
- **Coverage:** Minimum 80%
- **Command:** `cd backend && make test`

---

## 🔒 Security

- **Code Scanning:** Trivy (Docker images)
- **Dependency Scanning:** CodeCov
- **HTTPS:** TLS certificates in service certs/
- **Secrets:** Environment variables (never hardcoded)
- **SQL:** Parameterized queries only

---

## 📈 Project Statistics

| Metric | Value |
|--------|-------|
| **Services** | 7 microservices |
| **Frontend** | React 18 + TypeScript |
| **Backend** | Go 1.22+ |
| **Workflows** | 3 GitHub Actions |
| **Tests** | Comprehensive suite |
| **Branches** | main, production, develop |

---

## 🤝 Contributing

1. **Create feature branch** from `develop`
2. **Make changes** with clear commits
3. **Push & create PR** to `develop`
4. **Wait for tests** to pass (automatic)
5. **Address feedback** if any
6. **Auto-merge** or manual merge

**Guidelines:**
- Follow conventional commits (feat:, fix:, docs:)
- Keep commits focused
- Add tests for new features
- Update documentation
- Ensure all tests pass locally

---

## 📞 Support

### Documentation
- [Setup Guide](docs/SETUP_SUMMARY.md)
- [Branching Strategy](docs/BRANCHING_STRATEGY.md)
- [CI/CD Overview](docs/CI_CD_IMPLEMENTATION_COMPLETE.md)

### Issues & PRs
- GitHub Issues for bugs & features
- Pull Requests for code review
- Discussions for questions

---

## 📄 License

MIT License - see LICENSE file

---

## 👥 Team

**Project:** WanderPlan - AI Travel Planning  
**Current Status:** Active Development  
**Latest Update:** May 24, 2026

---

## 🗺️ Roadmap

- [ ] Advanced itinerary AI recommendations
- [ ] Real-time collaboration improvements
- [ ] Mobile app support
- [ ] Offline mode
- [ ] Advanced search filters
- [ ] Budget tracking
- [ ] Social sharing features

---

**🚀 Ready to contribute? Start with [Branching Strategy](docs/BRANCHING_STRATEGY.md)**

