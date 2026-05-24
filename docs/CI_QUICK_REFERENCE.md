# WanderPlan CI/CD - Quick Reference

## 🚀 Quick Start (5 minutes)

### 1. One-Time Setup
```bash
# Windows
scripts\setup-ci-cd.bat

# macOS/Linux  
chmod +x scripts/setup-ci-cd.sh
./scripts/setup-ci-cd.sh

# OR with Makefile (any OS)
make setup
```

### 2. Start Development
```bash
make docker-up    # Start services
make dev          # Start frontend + backend in dev mode
```

### 3. Make Changes & Push
```bash
git add .
git commit -m "Your changes"
git push origin main

# ✅ GitHub Actions automatically:
# 1. Tests your code
# 2. Fixes any errors
# 3. Merges if successful
# 4. Creates deployment artifact
```

---

## 📋 Common Commands

### Testing
```bash
make test              # Run all tests
make test-frontend     # Frontend only
make test-backend      # Backend only
make validate-all      # Full validation
```

### Building
```bash
make build             # Build frontend & backend
make docker-build      # Build Docker images
make docker-up         # Start services
make docker-down       # Stop services
```

### Code Quality
```bash
make lint              # Check for errors
make fix-all           # Auto-fix all errors
make fmt-frontend      # Format frontend
make fmt-backend       # Format backend
```

### Monitoring
```bash
./scripts/ci-status-dashboard.sh          # View status
./scripts/ci-status-dashboard.sh --watch  # Watch live
./scripts/ci-status-dashboard.sh --health # Health check

# GitHub CLI
gh run list --workflow=ci-cd.yml
gh run view <run-id>
```

---

## 🔧 Daily Workflow

### Morning Setup
```bash
git pull origin main
make docker-up      # Start services
make dev            # Start dev servers
```

### During Work
```bash
# Edit files, test locally
make test
make lint
make fix-all        # Auto-fix issues

# Commit changes
git add .
git commit -m "feat: add feature"
git push origin main
```

### Monitor CI/CD
```bash
# Option 1: GitHub CLI
gh run view --watch

# Option 2: Dashboard
./scripts/ci-status-dashboard.sh

# Option 3: GitHub Website
https://github.com/agamrai0123/trip.ly/actions
```

### End of Day
```bash
make docker-down    # Stop services
git push origin main # Ensure everything is pushed
```

---

## ✅ Automation Features at a Glance

| Feature | What It Does | When It Runs |
|---------|-------------|-------------|
| **Auto Test** | Runs all tests on every push | On every `git push` |
| **Auto Fix** | Fixes lint & format errors | When tests fail |
| **Auto Retry** | Retries after fixing | After auto-fix completes |
| **Auto Merge** | Merges PR if tests pass | When tests pass |
| **Auto Resume** | Continues after timeout | On workflow timeout |
| **Git Hooks** | Validates before commit | On `git commit` & `git push` |

---

## 🚨 Troubleshooting

### Tests failing locally?
```bash
make fix-all        # Auto-fix issues
make validate-all   # Re-run all tests
```

### Docker not working?
```bash
docker-compose down -v
make docker-up
docker-compose logs -f
```

### Need to debug CI/CD?
```bash
gh run list --workflow=ci-cd.yml
gh run view <run-id> --log
./scripts/ci-status-dashboard.sh --report
```

### Need to manually retry?
```bash
gh run rerun <run-id>
```

---

## 📁 Key Files

| File | Purpose |
|------|---------|
| `.github/workflows/ci-cd.yml` | Main automation workflow |
| `.github/workflows/error-recovery.yml` | Error detection & recovery |
| `Makefile` | Quick command shortcuts |
| `AUTOMATION_SETUP.md` | Complete documentation |
| `SETUP_SUMMARY.md` | Detailed setup guide |
| `scripts/setup-ci-cd.sh` | Unix setup script |
| `scripts/setup-ci-cd.bat` | Windows setup script |
| `scripts/ci-status-dashboard.sh` | Monitoring dashboard |

---

## 🎯 What Gets Automated

### Before Commit
```bash
Git Hook → Auto-fix linting → Stage files → Allow commit
```

### On Push
```bash
Push → Tests → Errors? → Auto-fix → Retry → Merge → Deploy Ready
```

### On Timeout
```bash
Timeout → Preserve state → New run → Resume → Complete
```

---

## 💡 Pro Tips

1. **Always use `make`** for consistency
2. **Check dashboard** after pushing
3. **Let auto-fix work** - don't force push
4. **Run local tests** before pushing (saves time)
5. **Review auto-fix commits** in PRs

---

## 📞 Need Help?

```bash
# View documentation
cat AUTOMATION_SETUP.md          # Detailed guide
cat SETUP_SUMMARY.md            # Setup overview
cat CI_QUICK_REFERENCE.md       # This file

# Check status
./scripts/ci-status-dashboard.sh --health

# View recent runs
gh run list --workflow=ci-cd.yml

# View specific error
gh run view <run-id> --log
```

---

## 🎉 You're All Set!

Everything is automated. Just:
1. **Write code**
2. **Push to GitHub**
3. **Let automation handle the rest**

**That's it!** ✨

---

*For complete documentation, see `AUTOMATION_SETUP.md`*
