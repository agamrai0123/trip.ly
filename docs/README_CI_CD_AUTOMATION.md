# WanderPlan CI/CD & Automation - Complete Setup

> Production-grade automated CI/CD system with error detection, auto-fixing, automatic retry, and auto-merge.

## 🎯 What This Provides

- ✅ **Automated Testing** - Frontend, Backend, Proto, Docker on every push
- ✅ **Error Detection & Auto-Fix** - Detects and fixes linting/formatting errors automatically
- ✅ **Automatic Retry** - Retries on timeout with state preservation
- ✅ **Auto-Merge** - Merges PRs automatically when tests pass
- ✅ **Git Hooks** - Pre-commit and pre-push validation
- ✅ **Local Development** - Docker stack, dev servers, monitoring
- ✅ **Comprehensive Monitoring** - Status dashboard and health checks

---

## 🚀 Quick Start

### 1. One-Time Setup (Choose one)

**Option A: Using Bash/Shell (macOS, Linux, WSL)**
```bash
chmod +x scripts/setup-ci-cd.sh
./scripts/setup-ci-cd.sh
# Follow the interactive menu
```

**Option B: Using Windows Batch**
```bash
scripts\setup-ci-cd.bat
# Follow the interactive menu
```

**Option C: Using Makefile (Any OS)**
```bash
make setup
```

### 2. Start Development
```bash
make docker-up    # Start PostgreSQL, Kafka, Redis, etc.
make dev          # Start frontend + backend dev servers
```

### 3. Make Changes & Push
```bash
git add .
git commit -m "Your commit message"
git push origin main

# ✅ Automatic:
# - Tests run (6-8 minutes)
# - Errors auto-fixed
# - PR auto-merged if successful
# - Deployment artifact created
```

---

## 📚 Documentation

| File | Read Time | Purpose |
|------|-----------|---------|
| [CI_QUICK_REFERENCE.md](CI_QUICK_REFERENCE.md) | 5 min | Quick commands & workflows |
| [SETUP_SUMMARY.md](SETUP_SUMMARY.md) | 15 min | Setup overview & examples |
| [AUTOMATION_SETUP.md](AUTOMATION_SETUP.md) | 30 min | Complete detailed guide |
| [CI_CD_INDEX.md](CI_CD_INDEX.md) | 10 min | Documentation index |

**New to this? Start with [CI_QUICK_REFERENCE.md](CI_QUICK_REFERENCE.md)** ✨

---

## 📋 Common Commands

### Testing & Validation
```bash
make test              # Run all tests
make test-frontend     # Frontend only
make test-backend      # Backend only
make validate-all      # Complete validation
make lint              # Check for errors
make fix-all           # Auto-fix all errors
```

### Development
```bash
make dev              # Start dev environment
make dev-frontend     # Frontend dev server
make dev-backend      # Backend dev server
make docker-up        # Start services
make docker-down      # Stop services
```

### Building
```bash
make build            # Build frontend & backend
make docker-build     # Build Docker images
make ci               # Run CI pipeline
```

### Monitoring
```bash
./scripts/ci-status-dashboard.sh          # View status
./scripts/ci-status-dashboard.sh --watch  # Watch live
./scripts/ci-status-dashboard.sh --health # Health check

# GitHub CLI
gh run list --workflow=ci-cd.yml
gh run view <run-id> --log
```

### Help
```bash
make help             # List all commands
make setup            # Setup everything
make clean            # Clean build artifacts
```

---

## 🔄 Automation Workflow

```
Developer Push
     ↓
GitHub Actions: ci-cd.yml
├─ Frontend: ESLint, TS, Tests
├─ Backend: Vet, Lint, Tests, Coverage
├─ Proto: Validation & Generation
├─ Docker: Build & Security Scan
└─ (All in parallel - 6-8 minutes)
     ↓
  Tests Pass?
   /        \
 YES        NO
  │          │
  │      Error Recovery: error-recovery.yml
  │      ├─ Auto-detect errors
  │      ├─ Auto-fix issues
  │      ├─ Auto-commit
  │      └─ Retry tests
  │          ↓
  │      Tests Pass?
  │       /        \
  │     YES        NO
  │      │          │
  │      └──┐     Create Issue
  │         │     (Manual Review)
  ↓        ↓
Auto-Merge & Deploy Ready
     ↓
🚀 Ready for Production
```

---

## 📁 Files & Workflows

### GitHub Actions Workflows
- **`.github/workflows/ci-cd.yml`** (487 lines) - Main CI/CD pipeline
  - 10 jobs running in parallel
  - Frontend, Backend, Proto, Docker validation
  - Error detection and auto-merge

- **`.github/workflows/error-recovery.yml`** (184 lines) - Error recovery
  - Auto-detect error types
  - Auto-fix formatting and imports
  - Retry failed workflows

### Scripts
- **`scripts/setup-ci-cd.sh`** - Interactive setup for Unix/Linux/macOS
- **`scripts/setup-ci-cd.bat`** - Interactive setup for Windows
- **`scripts/ci-status-dashboard.sh`** - Monitoring dashboard

### Build Automation
- **`Makefile`** (434 lines) - 50+ command shortcuts
- **`.env.example`** - Configuration template (197 variables)

### Documentation
- **`CI_QUICK_REFERENCE.md`** - 5 minute quick start
- **`SETUP_SUMMARY.md`** - Complete setup overview
- **`AUTOMATION_SETUP.md`** - Detailed guide (600+ lines)
- **`CI_CD_INDEX.md`** - Documentation index

---

## ✅ What Gets Tested

### Frontend
- ✅ ESLint linting
- ✅ TypeScript compilation
- ✅ Vitest unit tests
- ✅ Component tests
- ✅ Coverage reporting

### Backend
- ✅ Go vet checks
- ✅ golangci-lint
- ✅ Unit tests (race detector)
- ✅ Integration tests (PostgreSQL, Kafka)
- ✅ Coverage reporting (80% minimum)

### Proto
- ✅ Buf linting
- ✅ Code generation verification
- ✅ Dependency checking

### Docker
- ✅ Multi-stage builds
- ✅ Security scanning (Trivy)
- ✅ Image validation

---

## 🎯 Key Features

### 1. Automatic Testing
- Tests run on every push (no manual trigger)
- Parallel execution (6-8 minutes total)
- Coverage tracking
- Integration tests with real services

### 2. Error Detection & Auto-Fix
- Auto-detect error types
- Auto-fix linting errors
- Auto-fix formatting issues
- Auto-fix import statements
- Auto-commit fixes
- Auto-retry tests

### 3. Timeout Handling
- Preserve workflow state
- Automatic retry
- Resume from checkpoint
- Full history tracking

### 4. Automatic Merging
- Merge on success (if configured)
- Squash commits
- Delete branch
- Create deployment artifact

### 5. Git Hooks
- Pre-commit: Auto-fix before commit
- Pre-push: Validate before push
- Prevents broken code from reaching remote

### 6. Local Development
- Docker stack with all services
- Hot-reload dev servers
- Local validation
- Health checks

---

## 📊 Testing Times

| Component | Time |
|-----------|------|
| Frontend Lint & TS | 2-3 min |
| Frontend Tests | 1-2 min |
| Backend Build & Tests | 3-4 min |
| Integration Tests | 2-3 min |
| Proto Validation | 1 min |
| Docker Build | 2-3 min |
| **Total (Parallel)** | **6-8 min** |

With auto-fix retry: **10-15 min** (if fixes needed)

---

## 🆘 Troubleshooting

### Tests failing locally?
```bash
make fix-all
make validate-all
```

### Docker not working?
```bash
docker-compose down -v
make docker-up
docker-compose logs -f
```

### Need to debug CI?
```bash
gh run list --workflow=ci-cd.yml
gh run view <run-id> --log
```

### Want to understand workflows?
```bash
cat AUTOMATION_SETUP.md
cat CI_CD_INDEX.md
```

See **AUTOMATION_SETUP.md** for complete troubleshooting guide.

---

## 📖 Learning Resources

1. **First Time?** → Read [CI_QUICK_REFERENCE.md](CI_QUICK_REFERENCE.md)
2. **Want Details?** → Read [SETUP_SUMMARY.md](SETUP_SUMMARY.md)
3. **Need Full Guide?** → Read [AUTOMATION_SETUP.md](AUTOMATION_SETUP.md)
4. **Can't Find Info?** → Check [CI_CD_INDEX.md](CI_CD_INDEX.md)

---

## 💡 Pro Tips

1. **Use `make` commands** - Consistent across platforms
2. **Run `make test` locally** - Before pushing (saves CI time)
3. **Watch the dashboard** - Monitor progress in real-time
4. **Let auto-fix work** - Don't force push after failures
5. **Check .env** - Configure required variables

---

## 🔐 Security

- ✅ Code scanning (Trivy)
- ✅ Parameterized database queries
- ✅ Secret management via GitHub Secrets
- ✅ HTTPS ready
- ✅ JWT authentication
- ✅ OAuth 2.0 support

---

## 📈 Monitoring

### Real-Time Dashboard
```bash
./scripts/ci-status-dashboard.sh --watch
```

### Health Check
```bash
./scripts/ci-status-dashboard.sh --health
```

### Generate Report
```bash
./scripts/ci-status-dashboard.sh --report
```

---

## 🚀 Deployment Ready

Once CI/CD passes:
1. ✅ All tests passed
2. ✅ Coverage requirements met
3. ✅ Security scan passed
4. ✅ Build artifacts created
5. ✅ Deployment readiness marked

**Ready to deploy to staging/production!**

---

## 📞 Support

- **Documentation**: See files listed in [CI_CD_INDEX.md](CI_CD_INDEX.md)
- **Health Check**: `./scripts/ci-status-dashboard.sh --health`
- **GitHub CLI**: `gh run view <run-id>`
- **Local Testing**: `make validate-all`

---

## 🎓 Best Practices

✅ Commit small, logical changes  
✅ Write clear commit messages  
✅ Run tests locally before pushing  
✅ Review auto-fix commits  
✅ Monitor workflow status  
✅ Use git hooks effectively  
✅ Keep .env.example updated  
✅ Document your changes  

---

## 📝 Next Steps

1. Run setup: `make setup` (or use script)
2. Start development: `make dev`
3. Read documentation: See [CI_QUICK_REFERENCE.md](CI_QUICK_REFERENCE.md)
4. Make changes and push
5. Watch automation work! 🎉

---

## 📄 License

Part of the WanderPlan project - AI-powered travel itinerary planning platform.

---

**Happy Coding! 🚀** 

*No more manual testing, fixing, and merging — automation handles it all.*
