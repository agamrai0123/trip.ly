# 🔧 GitHub Branch Protection & Approver Configuration

**Status:** ACTION REQUIRED IN GITHUB UI  
**Setup Date:** May 24, 2026

---

## 📋 What Needs to Be Done in GitHub

### **Step 1: Configure Develop Branch**

Navigate to: **Repository Settings → Branches → Branch Protection Rules**

**Add Rule for: `develop`**

```
✓ Require pull request reviews before merging
  └─ Required number of reviewals: 1

✓ Require status checks to pass before merging
  └─ Status checks that must pass:
     • frontend-lint
     • frontend-test
     • backend-lint
     • backend-build
     • proto-check

✓ Require branches to be up to date before merging
  ✓ Yes

✓ Dismiss stale pull request approvals when new commits are pushed
  ✓ Yes

✓ Require code owner reviews
  ✓ Yes (if you set up CODEOWNERS)

✓ Restrict who can push to matching branches
  ✓ Allow specified actors to push
     • agamrai0123 (you)
```

---

### **Step 2: Configure Production Branch**

**Add Rule for: `production`**

```
✓ Require pull request reviews before merging
  └─ Required number of reviewals: 1

✓ Require status checks to pass before merging
  └─ Status checks that must pass:
     • All from develop PLUS
     • integration-tests
     • docker-build
     • security-scan

✓ Require branches to be up to date before merging
  ✓ Yes

✓ Dismiss stale pull request approvals when new commits are pushed
  ✓ Yes

✓ Require code owner reviews
  ✓ Yes

✓ Require signed commits
  ✓ No (optional)

✓ Restrict who can push to matching branches
  ✓ Allow specified actors to push
     • agamrai0123 (you)
```

---

### **Step 3: Configure Main Branch**

**Add Rule for: `main`**

```
✓ Require pull request reviews before merging
  └─ Required number of reviewals: 2  ← TWO approvals for production

✓ Require status checks to pass before merging
  └─ Status checks that must pass:
     • All previous checks

✓ Require branches to be up to date before merging
  ✓ Yes

✓ Dismiss stale pull request approvals when new commits are pushed
  ✓ Yes

✓ Require code owner reviews
  ✓ Yes

✓ Allow force pushes
  ✓ Allow specified actors: (leave empty - no force pushes to main)

✓ Allow deletions
  ✓ Checked: No (prevent accidental deletion)

✓ Require signed commits
  ✓ No (optional)
```

---

## 👤 Set You as Sole Approver

### **Option 1: Via GitHub CODEOWNERS**

Create file: `.github/CODEOWNERS`

```
# All code requires your approval
* @agamrai0123

# Specific approvers for different areas
/frontend/ @agamrai0123
/backend/ @agamrai0123
/migrations/ @agamrai0123
/deployments/ @agamrai0123
```

**Then in branch protection:**
- Enable "Require code owner reviews"
- This makes you required approver

### **Option 2: Via Team/User Settings**

**Settings → Branches → Protection Rules**

For each rule:
- Under "Require pull request reviews"
- Restrict to user: `@agamrai0123`

---

## 🔄 Approval Workflow After Setup

```
Feature Branch Created
      ↓
Developer pushes: feature/name
      ↓
Tests run (pre-push hook)
      ↓
PR created: feature/name → develop
      ↓
GitHub runs branch-merge.yml tests
      ↓
PR shows "1 approval required"
      ↓
You review & approve PR
      ↓
GitHub shows "Ready to auto-merge"
      ↓
Auto-merge activated (if configured)
      ↓
PR merged automatically
      ↓
Feature in develop ✓
      ↓
Later: Release PR: develop → production
      ↓
Same approval process
      ↓
Then: Release PR: production → main
      ↓
Same approval process
      ↓
Released to production ✓
```

---

## 🚀 Quick Configuration Steps

1. **Go to GitHub repository**
2. **Click Settings tab**
3. **Select Branches (left menu)**
4. **Click "Add rule"**
5. **Fill in configuration above**
6. **Repeat for each branch** (develop, production, main)

**Est. Time:** 5-10 minutes

---

## ✅ Verification

After configuration:

```bash
# Check if rules are applied
curl https://api.github.com/repos/agamrai0123/trip.ly/branches/develop/protection

# Should show:
{
  "required_pull_request_reviews": {
    "required_approving_review_count": 1,
    "dismiss_stale_reviews": true,
    "require_code_owner_reviews": true
  },
  "required_status_checks": {
    "enforcement_level": "everyone",
    "contexts": [...]
  }
}
```

---

## 📝 Next Steps

1. ✅ Read this configuration guide
2. ⏳ Go to GitHub and apply branch protection rules
3. ✅ Set CODEOWNERS to make you required approver
4. ✅ Test with sample feature branch (automated below)
5. ✅ Start using workflow for real features

---

**Status:** Ready for you to configure in GitHub UI  
**Time Estimate:** 5-10 minutes in GitHub

See `.github/BRANCHING_STRATEGY.md` for complete workflow documentation.
