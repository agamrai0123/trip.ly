# WanderPlan - Automated CI/CD & Error Recovery Setup

> Complete automation for testing, validation, error detection, fixing, and deployment with automatic retry and resumption.

---

## 🎯 Overview

This automated CI/CD system provides:

1. **Continuous Testing** - Automatic tests on every push/PR
2. **Code Quality Validation** - ESLint, TypeScript, Go Vet, golangci-lint, proto validation
3. **Error Detection** - Identifies issues in logs and test output
4. **Auto-Fix** - Automatically fixes formatting, linting, and import errors
5. **Automatic Retry** - Retries on timeouts or transient failures
6. **Auto-Merge** - Merges successful PRs automatically
7. **Notifications** - Reports status and creates issues for manual review
8. **Resumption** - Continues from where it left off after timeout or reset

---

## 📋 Features

### 1. Automated Testing Pipeline

**On every push to `main` or `develop` branch:**

- ✅ Frontend linting (ESLint)
- ✅ Frontend TypeScript compilation
- ✅ Frontend unit & component tests (Vitest)
- ✅ Backend Go vet checks
- ✅ Backend golangci-lint
- ✅ Backend build (all 7 services)
- ✅ Backend unit tests with coverage
- ✅ Backend integration tests (PostgreSQL, Kafka)
- ✅ Proto file validation (buf)
- ✅ Proto code generation check
- ✅ Docker image build (multi-stage, security scan)

### 2. Error Detection & Recovery

**When tests fail:**

1. Error logs are analyzed automatically
2. Errors are categorized (Frontend, Backend, Proto)
3. Auto-fix is attempted:
   - ESLint auto-fix
   - Go fmt & goimports
   - Proto regeneration
4. Fixed code is committed automatically
5. If auto-fix fails, an issue is created for manual review

### 3. Timeout & Resumption

**If workflow times out:**

1. Workflow automatically triggers a new run
2. Last successful state is preserved
3. On reset, continues from the last incomplete step
4. Full state tracking via artifacts

### 4. Automatic Merging

**On successful tests:**

1. PRs are automatically merged (if configured)
2. Deployment artifacts are created
3. Status is marked as "Ready for Deploy"

---

## 🚀 Quick Start

### Windows Users

```bash
# Run the setup script
scripts\setup-ci-cd.bat

# Or manually run tests
cd backend && go test ./...
cd frontend && npm run test
```

### macOS/Linux Users

```bash
# Make script executable
chmod +x scripts/setup-ci-cd.sh

# Run the setup script
./scripts/setup-ci-cd.sh

# Or manually run tests
cd backend && go test ./...
cd frontend && npm run test
```

---

## 📁 Configuration Files

### GitHub Actions Workflows

**`.github/workflows/ci-cd.yml`** - Main CI/CD Pipeline
- Frontend validation (ESLint, TypeScript, tests)
- Backend validation (go vet, golangci-lint, tests)
- Proto validation (buf)
- Docker build
- Error detection
- Auto-commit & auto-merge

**`.github/workflows/error-recovery.yml`** - Error Recovery
- Monitors failed workflows
- Auto-detects error types
- Applies fixes automatically
- Creates issues for unresolved errors
- Retries on timeout

### Environment Files

Create `.env` file in project root:

```bash
# Backend Configuration
POSTGRES_URL=postgres://user:password@localhost:5432/wanderplan
KAFKA_BROKERS=localhost:9092
JWT_SECRET=your-secret-key
OAUTH_CLIENT_ID=your-client-id
OAUTH_CLIENT_SECRET=your-client-secret

# Frontend Configuration
VITE_API_BASE_URL=http://localhost:8080

# Service Ports
API_GATEWAY_PORT=8080
AUTH_SERVICE_PORT=8081
TRIP_SERVICE_PORT=8082
USER_SERVICE_PORT=8083
COLLABORATION_SERVICE_PORT=8084
NOTIFICATION_SERVICE_PORT=8085
SEARCH_SERVICE_PORT=8086
```

---

## 🔧 Local Development

### Setup Git Hooks

Install pre-commit and pre-push hooks:

```bash
# Windows
scripts\setup-ci-cd.bat
# Select option 2

# macOS/Linux
./scripts/setup-ci-cd.sh
# Select option 2
```

These hooks will:
- Auto-fix linting errors before commit
- Run tests before pushing

### Run Local Validation

```bash
# Full validation (all tests)
./scripts/setup-ci-cd.sh  # Option 3

# Frontend only
./scripts/setup-ci-cd.sh  # Option 4

# Backend only
./scripts/setup-ci-cd.sh  # Option 5

# Auto-fix errors
./scripts/setup-ci-cd.sh  # Option 6
```

### Start Docker Stack

```bash
# Start services (PostgreSQL, Kafka, Redis, etc.)
./scripts/setup-ci-cd.sh  # Option 7

# Or manually
docker-compose up -d
```

---

## 📊 Workflow Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                     Push to main/develop                         │
└────────────────────────────┬────────────────────────────────────┘
                             │
                    ┌────────▼────────┐
                    │  Parallel Tests  │
                    └────┬────┬────┬───┘
          ┌─────────────┼─┐  │    │
          │             │ │  │    │
    ┌─────▼──┐   ┌─────▼─┴──▼─┐  │
    │Frontend │   │  Backend   │  │
    │ ESLint  │   │  Go Vet    │  │
    │TypeScript   │  golangci  │  │
    │Tests    │   │  Tests     │  │
    └────┬────┘   └──┬────────┬┘  │
         │           │        │   │
         ├───────────┤        │   │
                     │        │   │
                  ┌──▼────────▼───▼──┐
                  │  Proto & Docker  │
                  │     Builds       │
                  └────────┬─────────┘
                           │
                  ┌────────▼────────┐
                  │  Error Check    │
                  │  & Detection    │
                  └────────┬────────┘
                           │
                    ┌──────▼──────┐
                    │  Passed?    │
                    └──┬──────┬───┘
                      YES    NO
                       │      │
                       │   ┌──▼────────────┐
                       │   │  Auto-Fix     │
                       │   │  & Retry      │
                       │   └──┬──────┬─────┘
                       │      OK   FAIL
                       │      │     │
                       │   ┌──▼──┐┌─▼──────────┐
                       │   │Next │ │Create Issue│
                       │   └──┬──┘│for Review  │
                       │      │   └────────────┘
                       │      │
                  ┌────▼──────▼───┐
                  │  Auto-Merge   │
                  │  Deployment   │
                  │  Ready        │
                  └───────────────┘
```

---

## 🔄 Handling Failures & Timeouts

### If Tests Fail

1. **Automatic fixes** are attempted:
   ```
   ❌ Test fails
   ↓
   🤖 Auto-fix attempted
   ↓
   ✅ Fixed & committed
   ↓
   🔄 Workflow retried
   ```

2. **If auto-fix fails:**
   - Issue is created: `🚨 CI/CD Pipeline Failed - Manual Review Required`
   - Error logs are attached
   - You are assigned to the issue

### If Workflow Times Out

1. **Previous artifacts are preserved**
2. **New workflow run is triggered automatically**
3. **Continues from last successful step**
4. **Maintains state via GitHub artifacts**

### Manual Recovery

If you need to manually intervene:

```bash
# Check workflow status
gh run list --workflow=ci-cd.yml

# View specific run logs
gh run view <run-id>

# View job logs
gh run view <run-id> --log

# Retry failed run
gh run rerun <run-id>

# Rerun with debug
gh run rerun <run-id> --debug
```

---

## ✅ Pre-Commit Checks

When you commit code, these checks run automatically:

```bash
git commit -m "feat: add new feature"

# Automatically:
# 1. ESLint auto-fix frontend files
# 2. Go fmt backend files
# 3. Stage fixed files
# 4. Allow commit to proceed
```

## ✅ Pre-Push Checks

When you push code, these checks run:

```bash
git push origin main

# Automatically:
# 1. Run all frontend tests
# 2. Run all backend tests
# 3. If tests pass, allow push
# 4. If tests fail, block push (can override with --force-push)
```

---

## 📈 Monitoring & Debugging

### View Workflow Status

```bash
# List recent runs
gh run list --workflow=ci-cd.yml

# Watch a specific run
gh run view <run-id> --watch

# Download artifacts
gh run download <run-id> -D artifacts
```

### Check Service Health

```bash
# All services should respond to /healthz
curl http://localhost:8080/healthz
curl http://localhost:8081/healthz
# ... etc for each service
```

### View Logs

```bash
# Frontend logs
npm run build  # or bun run build

# Backend logs (if running locally)
cd backend && go run ./services/api-gateway/cmd/main.go

# Docker logs
docker-compose logs -f api-gateway
docker-compose logs -f auth-service
```

---

## 🛠️ Customization

### Change Merge Strategy

Edit `.github/workflows/ci-cd.yml`:

```yaml
# Change from squash to merge commit
gh pr merge ${{ github.event.pull_request.number }} \
  --auto \
  --merge \  # or --squash or --rebase
  --delete-branch
```

### Disable Auto-Merge

Comment out the `auto-merge` job in `.github/workflows/ci-cd.yml`:

```yaml
# auto-merge:
#   name: Auto-Merge on Success
#   ...
```

### Adjust Timeout Duration

In `.github/workflows/ci-cd.yml`:

```yaml
timeout-minutes: 60  # Change from default
```

### Add Custom Tests

```yaml
# Add to ci-cd.yml
- name: Run custom tests
  run: |
    # Your custom test command
    ./scripts/custom-tests.sh
```

---

## 🚀 Deployment

Once tests pass, the system is ready for deployment:

1. **Artifacts created** with build info
2. **Status marked** "Ready for Deploy"
3. **Integration tests passed** (PostgreSQL, Kafka)
4. **Security scan completed** (Trivy)
5. **Coverage reports** uploaded

Deploy using your standard deployment process:

```bash
# Example: Deploy to staging
gh deployment create production \
  --ref main \
  --environment-url https://staging.wanderplan.com \
  --auto-merge

# Or use your CI/CD platform's deployment tools
```

---

## 📚 Useful Commands

### Frontend

```bash
cd frontend

# Run tests
npm run test          # or bun run test
npm run test:watch   # or bun run test:watch

# Lint
npm run lint         # Check
npm run lint --fix   # Fix

# Build
npm run build        # Production
npm run build:dev    # Development

# Preview
npm run preview
```

### Backend

```bash
cd backend

# Build
go build ./...

# Test
go test ./...           # All tests
go test -race ./...     # With race detector
go test -cover ./...    # With coverage

# Lint
go vet ./...
golangci-lint run ./...

# Format
go fmt ./...
goimports -w .

# Run service
go run ./services/api-gateway/cmd/main.go
```

### Proto

```bash
cd backend

# Validate
buf lint

# Generate
buf generate

# Check for changes
buf diff --against HEAD
```

---

## ⚙️ Architecture

The automation uses:

- **GitHub Actions** - CI/CD orchestration
- **Docker Compose** - Local service stack
- **Go tools** - `go vet`, `golangci-lint`, `goimports`
- **Node tools** - `npm`, `bun`, `eslint`, `vitest`
- **Proto tools** - `buf`, `protoc`
- **Scanning** - `trivy` for security

---

## 🎓 Best Practices

1. **Commit small changes** - Easier to debug if tests fail
2. **Write descriptive commit messages** - Helps with auto-fix attribution
3. **Run local tests before pushing** - Saves CI/CD cycles
4. **Check error logs early** - Faster resolution
5. **Use git hooks** - Catch issues before pushing
6. **Review auto-fix commits** - Ensure they're correct

---

## 🆘 Troubleshooting

### Workflow never completes

```bash
# Check for hanging processes
gh run list --workflow=ci-cd.yml

# Retry the workflow
gh run rerun <run-id>

# Or cancel and start fresh
gh run cancel <run-id>
```

### Auto-merge not working

1. Enable in repo settings: **Settings → General → Allow auto-merge**
2. Configure default merge method
3. Ensure branch protection rules allow auto-merge

### Tests pass locally but fail in CI

1. Check environment variables in workflow
2. Check secret configuration
3. Compare Go/Node versions
4. Run in Docker locally to match CI environment

### Git hooks not firing

```bash
# Check hooks are executable
ls -la .git/hooks/

# Make executable
chmod +x .git/hooks/pre-commit
chmod +x .git/hooks/pre-push

# Test manually
./.git/hooks/pre-commit
```

---

## 📞 Support

For issues or questions:

1. Check **error logs** in GitHub Actions
2. Run **local validation** to reproduce
3. Create an **issue** with error details
4. Review **recent commits** for context

---

## 📄 License

This automation configuration is part of the WanderPlan project.

---

## 🎉 Summary

You now have:

✅ Automated testing on every push
✅ Automatic error detection & fixing  
✅ Automatic retry on timeout  
✅ Automatic merge on success  
✅ Git hooks for local validation  
✅ Docker stack for local development  
✅ Comprehensive logging & monitoring  
✅ Security scanning  
✅ Deployment readiness checks  

**Start using it:**

```bash
# Windows
scripts\setup-ci-cd.bat

# macOS/Linux
./scripts/setup-ci-cd.sh
```

Happy building! 🚀
