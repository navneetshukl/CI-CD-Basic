# Environment Configurations in CI/CD

This guide teaches you everything about managing environments.

---

## 1. What is an Environment?

An **environment** is a place where your application runs.

### Real-World Example

Like a chef opening a restaurant:
- **Kitchen (Dev)** = Try new recipes, make mistakes
- **Test Kitchen (Staging)** = Friends test food before opening
- **Restaurant (Production)** = Real customers eat

**Never serve experimental food directly to customers!**

---

## 2. The Three Main Environments

### Development (Dev)

| Aspect | Details |
|--------|---------|
| Who uses it | Developers |
| Branch | `develop` |
| Risk | LOW - Can break |
| Log Level | `debug` |
| Data | Fake data |

### Staging (QA)

| Aspect | Details |
|--------|---------|
| Who uses it | QA Team |
| Branch | `staging` |
| Risk | MEDIUM |
| Log Level | `info` |
| Data | Copy of production |

### Production (Prod)

| Aspect | Details |
|--------|---------|
| Who uses it | End Users |
| Branch | `main` |
| Risk | HIGH |
| Log Level | `error` |
| Data | Real user data |

---

## 3. What is a GitHub Secret?

A **GitHub Secret** is like a locked safe for sensitive data.

### What to Store
- Passwords
- API keys
- Server IP addresses
- Database credentials
- SSH keys

### Why Not in Code?
- Code is shared/public
- Secrets would be exposed
- Anyone could access your servers!

---

## 4. How to Create Secrets

**Step 1:** Go to GitHub Repository → **Settings**

**Step 2:** Click **Secrets and variables → Actions**

**Step 3:** Click **"New repository secret"**

**Step 4:** Add name and value, click **Add secret**

### Required Secrets

| Secret Name | Example Value |
|-------------|---------------|
| `SERVER_HOST_DEV` | `192.168.1.10` |
| `SERVER_HOST_STAGING` | `192.168.1.20` |
| `SERVER_HOST_PROD` | `192.168.1.30` |
| `SERVER_USER` | `deploy` |
| `SSH_PRIVATE_KEY` | Your SSH key |
| `DATABASE_URL_DEV` | `postgres://localhost:5432/dev` |
| `DATABASE_URL_PROD` | `postgres://prod:5432/prod` |

---

## 5. Branch-Based Deployment

```
develop  ──→  Development Server
staging  ──→  Staging Server
main     ──→  Production Server
```

### What Happens

**Push to `develop`:**
1. CI/CD runs tests
2. Auto-deploys to Dev server

**Push to `staging`:**
1. CI/CD runs tests
2. Auto-deploys to Staging server
3. QA team tests

**Push to `main`:**
1. CI/CD runs tests
2. Manual approval required
3. Deploys to Production

---

## 6. Complete Workflow Example

```yaml
name: Deploy by Branch

on:
  push:
    branches:
      - develop
      - staging
      - main

jobs:
  deploy:
    runs-on: ubuntu-latest
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Determine environment
        id: env
        run: |
          if [[ "${{ github.ref }}" == "refs/heads/develop" ]]; then
            echo "ENV=development" >> $GITHUB_OUTPUT
            echo "HOST=${{ secrets.SERVER_HOST_DEV }}" >> $GITHUB_OUTPUT
          elif [[ "${{ github.ref }}" == "refs/heads/staging" ]]; then
            echo "ENV=staging" >> $GITHUB_OUTPUT
            echo "HOST=${{ secrets.SERVER_HOST_STAGING }}" >> $GITHUB_OUTPUT
          else
            echo "ENV=production" >> $GITHUB_OUTPUT
            echo "HOST=${{ secrets.SERVER_HOST_PROD }}" >> $GITHUB_OUTPUT
          fi
      
      - name: Deploy via SSH
        uses: appleboy/ssh-action@v1.0.3
        with:
          host: ${{ steps.env.outputs.HOST }}
          username: ${{ secrets.SERVER_USER }}
          key: ${{ secrets.SSH_PRIVATE_KEY }}
          script: |
            echo "Deploying to ${{ steps.env.outputs.ENV }}"
            docker pull navneetshukla/my-go-api:${{ github.sha }}
            docker stop my-go-api || true
            docker rm my-go-api || true
            docker run -d \
              --name my-go-api \
              -p 8080:8080 \
              -e ENVIRONMENT=${{ steps.env.outputs.ENV }} \
              navneetshukla/my-go-api:${{ github.sha }}
```

> ⚠️ **Note:** The `ENVIRONMENT` variable is passed to the container as an example of environment-aware deployment. To make your app actually use it, update `main.go` to read it with `os.Getenv("ENVIRONMENT")` and act on it (e.g. log levels, feature flags, DB selection). The sample app in this project has a hardcoded port — see the Dockerfile (`EXPOSE 8080` and `CMD ["./myapp"]`).

---

## 7. Best Practices

### 1. Use Secrets, Not Hardcoded Values

```yaml
# ❌ WRONG
password="my-secret"

# ✅ CORRECT
password="${{ secrets.DATABASE_PASSWORD }}"
```

### 2. Different Secrets per Environment

```yaml
# ❌ WRONG - One password for all
DATABASE_URL: ${{ secrets.DATABASE_URL }}

# ✅ CORRECT - Separate passwords
DATABASE_URL_DEV: ${{ secrets.DATABASE_URL_DEV }}
DATABASE_URL_PROD: ${{ secrets.DATABASE_URL_PROD }}
```

### 3. Protect Main Branch

In GitHub: **Settings → Branches → Add rule**
- ✅ Require pull request reviews
- ✅ Require status checks to pass

### 4. Log Levels

| Env | Level | What Logs |
|-----|-------|-----------|
| Dev | `debug` | Everything |
| Staging | `info` | Important events |
| Prod | `error` | Only errors |

### 5. Naming Convention

```
SERVER_HOST_DEV      ← Development
SERVER_HOST_STAGING  ← Staging
SERVER_HOST_PROD     ← Production
```

---

## 8. Common Mistakes

### Mistake 1: Production Credentials in Dev

```yaml
# ❌ WRONG
DATABASE_URL="postgres://prod/prod_db"

# ✅ CORRECT
DATABASE_URL="${{ secrets.DATABASE_URL_DEV }}"
```

### Mistake 2: Committing Secrets

```bash
# .gitignore
.env
*.pem
*.key
```

### Mistake 3: One Server for All

```yaml
# ❌ RISKY
SERVER_HOST: "192.168.1.100"

# ✅ SAFE
SERVER_HOST_DEV: "192.168.1.10"
SERVER_HOST_PROD: "192.168.1.30"
```

---

## 📋 Checklist

- [ ] Create GitHub Secrets for each environment
- [ ] Set up branch protection for `main`
- [ ] Add `.env` to .gitignore
- [ ] Use separate servers for prod
- [ ] Test in dev/staging before prod

---

## 🔗 Related Topics

- [Rollback Strategies](./02-rollback-strategies.md)
- [Health Checks](./03-health-checks.md)

---

*Last updated: September 2026*
