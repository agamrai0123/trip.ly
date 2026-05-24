# WanderPlan - Git Branching Strategy

## 📋 Overview

This document outlines the Git branching strategy for the WanderPlan project. It follows a **modified GitFlow** pattern with three main branches and automated testing at each merge point.

---

## 🌳 Branch Structure

```
production (stable, production-ready code)
    ↑
    ├─ develop (integration branch for features)
    ↑
    ├─ feature/* (new features)
    ├─ bugfix/* (bug fixes)
    └─ hotfix/* (emergency production fixes)
```

### **Four Main Branches**

| Branch | Purpose | Protection | Auto-Test |
|--------|---------|-----------|-----------|
| **main** | Official releases & production | ✅ Yes | ✅ Yes |
| **production** | Production-ready deployments | ✅ Yes | ✅ Yes |
| **develop** | Integration of features | ✅ Yes | ✅ Yes |
| **feature/\*** | Individual feature development | ❌ No | ✅ Yes |

---

## 🔄 Development Workflow

### **1. Starting a New Feature**

```bash
# Update develop branch
git checkout develop
git pull origin develop

# Create feature branch
git checkout -b feature/user-authentication
# or: git checkout -b feature/trip-filtering
# or: git checkout -b bugfix/payment-bug
```

**Branch Naming Conventions:**
- Features: `feature/descriptive-name`
- Bug fixes: `bugfix/bug-description`
- Hotfixes: `hotfix/critical-issue`
- Releases: `release/v1.2.0`

### **2. Development Phase**

```bash
# Work on your feature
# Make commits regularly with descriptive messages

git add .
git commit -m "feat: add user authentication with OAuth2"

# Push to remote
git push origin feature/user-authentication
```

**Commit Message Convention:**
```
feat: add new feature
fix: fix a bug
docs: update documentation
chore: update dependencies
test: add tests
refactor: code restructuring
perf: performance improvements
```

### **3. Ready for Code Review - Create Pull Request**

When your feature is complete:

```bash
# Push your branch
git push origin feature/user-authentication

# Create Pull Request on GitHub
# - Base: develop
# - Compare: feature/user-authentication
# - Title: Clear feature description
# - Description: What changed and why
```

**PR Checklist:**
- ✅ Code passes all automated tests
- ✅ No conflicts with develop branch
- ✅ Tests are included for new features
- ✅ Documentation is updated
- ✅ No hardcoded secrets or credentials

### **4. Automated Testing (PR Merge Gate)**

When a PR is created:

1. **Frontend Tests Run**
   - ESLint validation
   - TypeScript type checking
   - Unit/component tests
   - Coverage upload

2. **Backend Tests Run**
   - Go vet validation
   - Golangci-lint linting
   - Build all services
   - Unit tests with coverage

3. **Proto Validation**
   - Buf lint
   - Generation verification

4. **Merge Gate Decision**
   - ✅ If all pass → PR can be merged
   - ❌ If any fail → PR blocked until fixed

### **5. Code Review & Approval**

- Team members review the code
- Automated tests must pass
- At least one approval required
- Address feedback comments

### **6. Merge to Develop**

```bash
# GitHub UI: Click "Merge pull request"
# OR via CLI:
gh pr merge <PR_NUMBER> --squash

# Branch is automatically deleted
```

---

## 📊 Branch Progression

### **Feature → Develop Merge**

```
feature/new-feature (your work)
         ↓
   [AUTOMATED TESTS]
   - Frontend tests
   - Backend tests
   - Proto validation
         ↓
   [MERGE GATE CHECK]
   - All tests must pass
   - Review approved
         ↓
      develop (integration)
```

### **Develop → Production Merge**

When develop has accumulated several features ready for release:

```bash
# Create PR: develop → production
git checkout -b release/v1.2.0
git push origin release/v1.2.0

# On GitHub:
# 1. Create PR from develop to production
# 2. All automated tests run again
# 3. If pass, merge to production
```

**What Triggers Release:**
- Marketing milestone
- Feature set complete
- Bug fixes accumulated
- Performance improvements ready
- Security patches needed

### **Production → Main Merge**

When production is stable and ready for official release:

```bash
# Create PR: production → main (with version tag)
# This is typically done by release manager
```

---

## 🔥 Hotfix Workflow (Emergency Fixes)

For critical production bugs:

```bash
# Create hotfix branch from production
git checkout -b hotfix/payment-processing production
git push origin hotfix/payment-processing

# Fix the bug and push
git add .
git commit -m "fix: resolve payment processing timeout"
git push origin hotfix/payment-processing

# Create PR: hotfix/... → production
# After merge: hotfix/... → develop (keep develop in sync)
```

---

## 🛡️ Branch Protection Rules

### **Develop Branch Protection**
- ✅ Require pull request reviews (minimum 1)
- ✅ Require all checks to pass
- ✅ Require branches to be up to date
- ✅ Require code owners approval
- ✅ Dismiss stale reviews when new commits pushed
- ✅ Require status checks to pass

### **Production Branch Protection**
- ✅ Require pull request reviews (minimum 2)
- ✅ Require all checks to pass
- ✅ Require branches to be up to date
- ✅ Require code owners approval
- ✅ Dismiss stale reviews when new commits pushed
- ✅ Require status checks to pass
- ✅ Require approved reviews from code owners

### **Main Branch Protection**
- ✅ Require pull request reviews (minimum 2)
- ✅ Require all checks to pass
- ✅ Require branches to be up to date
- ✅ Require code owners approval
- ✅ Dismiss stale reviews when new commits pushed
- ✅ Require status checks to pass
- ✅ Require approved reviews from code owners
- ✅ Allow force pushes: NO
- ✅ Allow deletions: NO

---

## 📈 Continuous Integration Pipeline

### **On Every Push (Pre-push Hook)**
```
Local Validation:
  1. Frontend linting & type check
  2. Backend build & vet
  3. Tests on modified files
```

### **On Every PR Creation/Update**
```
GitHub Actions - branch-merge.yml:
  1. Validate branch structure
  2. Run all frontend tests
  3. Run all backend tests
  4. Validate proto files
  5. Merge gate decision
  6. Auto-merge if successful
```

### **On Every Push to Main Branch**
```
GitHub Actions - ci-cd.yml:
  1. Full test suite
  2. Docker build & security scan
  3. Deployment readiness check
  4. Notification
```

---

## 📝 Commit Message Standards

Follow this format for consistency:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only
- `style`: Changes that don't affect code meaning
- `refactor`: Code change that neither fixes bug nor adds feature
- `perf`: Code change that improves performance
- `test`: Adding or updating tests
- `chore`: Changes to build process, dependencies, etc.

**Examples:**
```
feat(auth): implement OAuth2 authentication
fix(payments): resolve timeout on large transactions
docs(README): update installation instructions
test(users): add unit tests for user service
chore(deps): update Go dependencies to latest versions
```

---

## 🔄 Merging Strategy

### **Squash Merge (Default)**
Used for most feature branches to keep history clean:
```bash
gh pr merge <PR_NUMBER> --squash
```

### **Create a Merge Commit (For Release Branches)**
Used for develop → production to preserve history:
```bash
gh pr merge <PR_NUMBER> --create-merge-commit
```

### **Rebase & Merge (For Small Fixes)**
Used for small, single-commit changes:
```bash
gh pr merge <PR_NUMBER> --rebase
```

---

## 📊 Status Check Examples

### **✅ Passing Checks**
```
Status: All checks have passed

✅ Frontend - ESLint & TypeScript Check
✅ Frontend - Unit & Component Tests
✅ Backend - Go Vet & Linting
✅ Backend - Build All Services
✅ Backend - Tests
✅ Proto - Validation & Generation Check
✅ Merge Gate - All Tests Must Pass
```

### **❌ Failing Checks**
```
Status: Some checks did not complete successfully

✅ Frontend - ESLint & TypeScript Check
❌ Frontend - Unit & Component Tests
   └─ Test failure in example.test.ts:12
✅ Backend - Go Vet & Linting
⏳ Backend - Build All Services (Running)
⏸️ Other checks blocked by failures
```

---

## 🚀 Ready for Production

### **Checklist Before Production Release**

- [ ] All tests passing on production branch
- [ ] All commits reviewed and approved
- [ ] Version number updated
- [ ] CHANGELOG updated
- [ ] Release notes written
- [ ] Dependencies up to date
- [ ] Security scan passed
- [ ] Performance benchmarks acceptable
- [ ] Documentation updated
- [ ] Database migrations tested
- [ ] Rollback plan documented

---

## 📞 Common Scenarios

### **Scenario 1: Merge Feature to Develop**
```bash
# 1. Create PR on GitHub (feature/X → develop)
# 2. Tests run automatically
# 3. After review approval
# 4. Click "Merge pull request"
# 5. Delete branch
```

### **Scenario 2: Fix Failed Tests on Feature Branch**
```bash
# Tests failed on PR
git checkout feature/my-feature

# Fix the issue
git add .
git commit -m "fix: resolve test failure"
git push origin feature/my-feature

# Tests run again automatically
```

### **Scenario 3: Update Feature from Develop**
```bash
# Develop changed since you started your feature
git checkout feature/my-feature
git fetch origin
git merge origin/develop

# Resolve any conflicts
git add .
git commit -m "chore: merge develop updates"
git push origin feature/my-feature
```

### **Scenario 4: Emergency Production Fix**
```bash
git checkout production
git pull origin production
git checkout -b hotfix/critical-bug

# Fix the bug
git add .
git commit -m "fix: critical payment bug"
git push origin hotfix/critical-bug

# Create 2 PRs:
# 1. hotfix/critical-bug → production
# 2. hotfix/critical-bug → develop
```

---

## 🎓 Best Practices

1. **Keep branches short-lived** (1-2 weeks maximum)
2. **Commit frequently** with clear messages
3. **Test locally** before pushing
4. **Rebase with develop** before creating PR
5. **Keep PR focused** on one feature/fix
6. **Link to issues** in PR description
7. **Request reviews** from code owners
8. **Address feedback** promptly
9. **Squash unnecessary commits** before merge
10. **Delete branches** after merge

---

## 🔗 Related Documentation

- [CI/CD Pipeline Documentation](./CI_CD_INDEX.md)
- [GitHub Actions Workflows](./.github/workflows/)
- [Development Setup Guide](./SETUP_SUMMARY.md)
- [Code Quality Standards](./README.md)

---

## 📧 Questions?

For questions about the branching strategy:
1. Check this document
2. Review recent PR examples
3. Ask your team lead or DevOps engineer

---

**Last Updated:** May 24, 2024  
**Strategy Version:** 1.0  
**Status:** 🟢 Active
