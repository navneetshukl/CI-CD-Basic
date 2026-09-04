# CI/CD Learning Guide 🚀

A beginner-friendly guide to understanding and implementing CI/CD pipelines from scratch.

---

## 📚 Table of Contents

1. [What is CI/CD?](#what-is-cicd)
2. [Key Concepts](#-key-concepts)
3. [Pipeline Flow](#-pipeline-flow)
4. [Step 1: Basic Build Pipeline](#-step-1-basic-build-pipeline)
5. [Line-by-Line Explanation](#-line-by-line-explanation)
6. [How GitHub Actions Processes This](#-how-github-actions-processes-this)
7. [How to Use](#-how-to-use)
8. [Common Workflow File Names](#-common-workflow-file-names)
9. [Workflow Syntax Reference](#-workflow-syntax-reference)
10. [Complete Learning Path](#-complete-learning-path)
11. [Step 2: Adding Automated Tests](#-step-2-adding-automated-tests)
12. [Step 3: Code Quality Checks](#-step-3-code-quality-checks)
13. [Step 4: Docker - Containerize Your App](#-step-4-docker-containerize-your-app)
14. [Complete CI/CD Pipeline (Production)](#-complete-cicd-pipeline-production)
15. [Step 5: Deployment Pipeline](#-step-5-deployment-pipeline)
16. [What is Jenkins? (Bonus)](#-what-is-jenkins-bonus)
17. [📚 Additional Resources](#-additional-resources)
    - [Environment Configurations](./docs/01-environment-configurations.md)
    - [Rollback Strategies](./docs/02-rollback-strategies.md)
    - [Health Checks & Verification](./docs/03-health-checks.md)

---

## What is CI/CD?

### CI = Continuous Integration
- Developers merge code changes into a shared repository frequently (multiple times per day)
- Automated builds and tests run on every change
- Catches bugs early before they become problems

### CD = Continuous Delivery (or Deployment)
- **Continuous Delivery**: Code changes are automatically prepared for release to production
- **Continuous Deployment**: Changes go live automatically without manual approval

---

## 🔑 Key Concepts

| Concept | What it means |
|---------|---------------|
| **Pipeline** | A series of automated steps that code goes through (build → test → deploy) |
| **Build** | Converting source code into a runnable artifact |
| **Test** | Running automated checks to verify the code works |
| **Artifact** | The deployable output (Docker image, zip file, etc.) |
| **Environment** | Where the app runs (dev, staging, production) |
| **Job** | A unit of work that runs on a runner (can have multiple jobs in parallel) |
| **Step** | Individual tasks within a job (run sequentially) |
| **Runner** | A server/machine that executes your pipeline jobs |
| **Trigger** | An event that starts your pipeline (push, pull request, schedule, etc.) |

---

## 🔄 Pipeline Flow

```
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│  Commit  │───▶│  Build   │───▶│   Test   │───▶│ Quality  │───▶│  Deploy  │
│   Code   │    │          │    │          │    │  Check   │    │          │
└──────────┘    └──────────┘    └──────────┘    └──────────┘    └──────────┘
```

---

## 🔰 Step 1: Basic Build Pipeline

### File Location
`.github/workflows/01-basic-build.yml`

### The Code
```yaml
name: Basic Go Build

on:
  push:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
      - name: Get source code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.25'

      - name: Download dependencies
        run: go mod download

      - name: Build the Go app
        run: go build -v ./...
```

---

## 📖 Line-by-Line Explanation

### 1. Pipeline Name
```yaml
name: Basic Go Build
```
👉 The name of your pipeline (shows up in GitHub Actions UI)

---

### 2. Trigger Configuration
```yaml
on:
  push:
    branches: [ main ]
```
👉 **WHEN to run**: This pipeline triggers when code is pushed to the `main` branch

Common trigger options:
- `push` - when code is pushed
- `pull_request` - when a PR is created/updated
- `schedule` - on a cron schedule (e.g., daily)
- `workflow_dispatch` - manually triggered

---

### 3. Job Definition
```yaml
jobs:
  build:
    runs-on: ubuntu-latest
```
👉 **WHAT to run**: A job called "build" that runs on a fresh Ubuntu Linux machine

Common `runs-on` options:
- `ubuntu-latest` - Ubuntu Linux (free)
- `windows-latest` - Windows Server
- `macos-latest` - macOS

---

### 4. Step 1 - Checkout Code
```yaml
    steps:
      - name: Get source code
        uses: actions/checkout@v4
```
👉 **STEP 1**: Download your code from GitHub onto the runner machine

`actions/checkout@v4` is a pre-built action that clones your repository.

---

### 5. Step 2 - Setup Go
```yaml
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.25'
```
👉 **STEP 2**: Install Go version 1.25 on the runner

The `with` section passes parameters to the action.

---

### 6. Step 3 - Download Dependencies
```yaml
      - name: Download dependencies
        run: go mod download
```
👉 **STEP 3**: Fetch all project dependencies (Gin, etc.)

Required before building so the compiler knows about all external packages.

---

### 7. Step 4 - Build
```yaml
      - name: Build the Go app
        run: go build -v ./...
```
👉 **STEP 4**: Compile your Go code

`./...` means "all packages in current directory and subdirectories"
`-v` flag shows verbose output (what's being built)

> **Note:** We only BUILD (compile) here, not RUN. Running a web server would make the workflow hang forever!

---

## 🔄 How GitHub Actions Processes This

```
You push code to GitHub
         ↓
GitHub detects the trigger (push to main)
         ↓
GitHub spins up a fresh Ubuntu machine (runner)
         ↓
Step 1: Downloads your code via actions/checkout
         ↓
Step 2: Installs Go 1.25 via actions/setup-go
         ↓
Step 3: Downloads dependencies via go mod download
         ↓
Step 4: Compiles code via go build -v ./...
         ↓
Reports ✅ Success or ❌ Failure in GitHub Actions tab
```

---

## 🚀 How to Use

### Prerequisites
1. A GitHub account
2. A GitHub repository with your Go code
3. Git installed locally

### Setup Steps

1. **Create a new repository on GitHub** (or use existing one)

2. **Add this file to your repository:**
   ```
   .github/workflows/01-basic-build.yml
   ```

3. **Push your code to the main branch:**
   ```bash
   git add .
   git commit -m "Add basic CI pipeline"
   git push origin main
   ```

4. **View the pipeline in GitHub:**
   - Go to your repository on GitHub
   - Click the **"Actions"** tab
   - You should see your pipeline running!

5. **Click on the workflow run** to see:
   - Each step's output
   - Build logs
   - Success/failure status

---

## 🔧 Common Workflow File Names

You can name your workflow files anything, but common conventions:
- `ci.yml` - for continuous integration
- `build.yml` - for build pipelines
- `01-*.yml`, `02-*.yml` - for numbered steps in learning

---

## 📖 Workflow Syntax Reference

### Jobs
```yaml
jobs:
  job-name:           # Must be unique within the file
    runs-on: ubuntu-latest
    needs: other-job  # Wait for another job to complete first
```

### Steps
```yaml
steps:
  - name: Step name     # Display name
    uses: action@version # Use a pre-built action
    run: command        # Run a shell command
    with:               # Parameters for actions
      key: value
```

### Environment Variables
```yaml
env:
  NODE_ENV: production
  API_URL: https://api.example.com

steps:
  - name: Use env var
    run: echo $NODE_ENV
```

### Conditional Execution
```yaml
- name: Deploy
  if: github.ref == 'refs/heads/main'
  run: ./deploy.sh
```

---

## 🎯 Complete Learning Path

Here's what we've built step by step:

### ✅ Level 1: Basic Build Pipeline
- Simple CI pipeline that builds the application
- File: `01-basic-build.yml`

### ✅ Level 2: Add Tests
- Automated unit tests that run on every push
- File: `02-add-tests.yml`

### ✅ Level 3: Code Quality
- gofmt (code formatting)
- go vet (bug detection)
- staticcheck (advanced linting)
- Test coverage reports
- File: `03-code-quality.yml`

### ✅ Level 4: Docker
- Containerize the application
- Multi-stage Dockerfile for small images
- Push to Docker Hub
- File: `04-docker.yml`
- Helper files: `Dockerfile`, `docker-compose.yml`, `Makefile`

### ✅ Complete CI/CD Pipeline
- All steps in one workflow (Build → Test → Quality → Deploy)
- Uses `needs:` to ensure sequential execution
- File: `00-complete-ci-cd.yml`

### 📌 Level 5: Deployment
- Deploy to a server automatically via SSH
- Pull Docker image and run container
- Health check verification after deployment
- Files: `05-deploy.yml` (manual) and `00-complete-ci-cd.yml` (automatic)
---

## 🧪 Step 2: Adding Automated Tests

### Why Tests Matter
- **Catch bugs early** - Find issues before they reach production
- **Confidence** - Know your changes don't break existing code
- **Documentation** - Tests show how your code should be used
- **Safety net** - Refactor code without fear of breaking things

### What We Added
1. **Gin web framework** - For handling HTTP requests
2. **REST API endpoints** - Health check + User CRUD routes
3. **Unit tests** - 7 tests covering all endpoints
4. **Test pipeline** - Automatically runs tests on every push

### Project Structure (after Step 2)
```
test/
└── ci-cd/                          ← Everything CI/CD related lives here
    ├── main.go                     ← Gin server with API endpoints
    ├── main_test.go                ← Tests for the API
    ├── go.mod                      ← Go module file
    ├── go.sum                      ← Dependency checksums
    ├── README.md                   ← This learning guide
    ├── .github/
    │   └── workflows/
    │       ├── 01-basic-build.yml
    │       └── 02-add-tests.yml
    └── scripts/                    ← Helper scripts (for future steps)
```

### The API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/api/users` | Get all users |
| POST | `/api/users` | Create a new user |
| GET | `/api/users/:id` | Get user by ID |

### Testing Locally

Run tests on your machine:
```bash
go test -v ./...
```

Expected output:
```
=== RUN   TestHealthCheck
--- PASS: TestHealthCheck (0.00s)
=== RUN   TestCreateUser
--- PASS: TestCreateUser (0.00s)
...
PASS
ok      github.com/navneetshukla/test      0.787s
```

### The New Pipeline

File: `.github/workflows/02-add-tests.yml`

```yaml
name: Go CI with Tests

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest

    steps:
      - name: Get source code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.25'

      - name: Download dependencies
        run: go mod download

      - name: Run tests
        run: go test -v ./...

      - name: Build the app
        run: go build -v ./...
```

### What's New Compared to Step 1?

| Addition | Why |
|----------|-----|
| `on: pull_request` | Tests run on PRs too, not just main branch |
| `go mod download` | Download dependencies before testing |
| `go test -v ./...` | Run all tests with verbose output |
| `go build -v ./...` | Verify the app still compiles |

### The Flow

```
Developer pushes code
         ↓
GitHub triggers pipeline
         ↓
Checkout code
         ↓
Install Go 1.25
         ↓
Download dependencies (gin, etc.)
         ↓
Run all tests (7 tests)
         ↓
If tests pass → Build the app
         ↓
If anything fails → Report error, block deployment
```

### Testing in GitHub

1. Push your code to GitHub
2. Go to the **Actions** tab
3. Click on the workflow run
4. Expand the **"Run tests"** step
5. You'll see output like:
   ```
   === RUN   TestHealthCheck
   --- PASS: TestHealthCheck (0.00s)
   === RUN   TestCreateUser
   --- PASS: TestCreateUser (0.00s)
   ...
   PASS
   ```

### Test Coverage (Optional Enhancement)

To check how much of your code is tested, add this step:
```yaml
- name: Test with coverage
  run: go test -v -cover ./...
```

Output will show:
```
PASS    coverage: 68.3% of statements
```

### Benefits You Get Now

✅ **Automatic testing** - Tests run on every push  
✅ **PR validation** - Can't merge code that breaks tests  
✅ **Fast feedback** - Know within minutes if something broke  
✅ **Confidence** - Deploy with confidence knowing tests pass  

---

## 🔍 Step 3: Code Quality Checks

### Why Code Quality Matters
- **Catch bugs automatically** - Before code review
- **Consistent style** - Everyone writes code the same way
- **Industry best practices** - Learn from tools used by professionals

### The Four Quality Tools

| Tool | What it does | Example Error |
|------|-------------|---------------|
| **gofmt** | Checks code formatting | `main.go needs formatting` |
| **go vet** | Finds suspicious code | `Printf format mismatch` |
| **staticcheck** | Advanced linting | `Unused function` |
| **Coverage** | % of code tested | `68.3% coverage` |

### What Happens When Checks Fail

| Check | Fix Command |
|-------|------------|
| **gofmt** | `gofmt -w .` |
| **go vet** | Fix the code issue |
| **staticcheck** | Remove unused code |
| **Test** | Fix the test |

### Benefits
✅ Automatic bug detection  
✅ Consistent code style  
✅ Higher test coverage  
✅ Faster code reviews  

---

## 🐳 Step 4: Docker - Containerize Your App

### What is Docker?
Packages your application with everything it needs to run, so it works everywhere!

### Files We Created

| File | Purpose |
|-------|---------|
| `Dockerfile` | Build instructions for container |
| `.dockerignore` | Exclude files from image |
| `docker-compose.yml` | Run locally with one command |
| `Makefile` | Shortcut commands |
| `04-docker.yml` | Pipeline to push to Docker Hub |

### Quick Commands

```bash
# Build locally
docker build -t my-go-api:latest .

# Run locally
docker run -p 8080:8080 my-go-api:latest

# Or use Make
make build
make run
```

### Before Docker Pipeline Works
1. Create Docker Hub account: https://hub.docker.com/
2. Create access token
3. Add `DOCKER_USERNAME` and `DOCKER_TOKEN` to GitHub Secrets

### Benefits
✅ Consistent environments everywhere  
✅ Easy deployment  
✅ Isolated dependencies  
✅ Version control  

---

## 🎯 Complete CI/CD Pipeline (Production)

### Problem with Separate Workflows
All workflows run **independently** - Docker might deploy even if tests fail!

### Solution: `needs:` Keyword
Each job waits for the previous one to pass:

```yaml
build:
  # First job - no dependencies
  
test:
  needs: build    # Wait for build
  
quality:
  needs: test     # Wait for tests
  
docker:
  needs: quality  # Wait for quality

deploy:
  needs: docker   # Wait for Docker build to complete
```

### Flow
```
Build → If fail, STOP!
  ↓
Tests → If fail, STOP!
  ↓
Quality → If fail, STOP!
  ↓
Docker → Build & push image
  ↓
Deploy → SSH to server & run container
```

### Files Created

| File | Purpose |
|------|---------|
| `00-complete-ci-cd.yml` | Complete pipeline with all steps |
| `05-deploy.yml` | Standalone manual deployment |

### When to Use Which

| Workflow | Use For |
|----------|---------|
| `00-complete-ci-cd.yml` | Production deployment (automatic) |
| `05-deploy.yml` | Manual hotfix deployments |
| `01-04.yml` | Learning and development |

### Benefits
✅ No broken deployments
✅ Clear failure points
✅ Sequential execution
✅ Industry standard
✅ Automatic health checks
  needs: quality  # Wait for quality

---

## 🚀 Step 5: Deployment Pipeline

### Overview

Now that your app builds, tests pass, quality checks succeed, and Docker image is created — it's time to **deploy it to a real server**!

This guide covers:
- SSH-based deployment to a remote server
- Docker pull + run on production host
- Environment-specific configurations
- Secrets management
- Health checks after deployment
- Rollback strategies

---

### Prerequisites

Before we start, you need:

1. **A remote server** with Docker installed
2. **SSH access** to that server
3. **GitHub Secrets** configured in your repository

---

### Deployment Architecture

```
GitHub Actions → SSH → Remote Server → Docker Container
     ↓                ↓            ↓
  Push code    Execute commands  Run container
```

The workflow will:
1. Build Docker image (already done in Step 4)
2. Push image to Docker Hub (already done in Step 4)
3. SSH into remote server
4. Pull latest image from Docker Hub
5. Run container with new image
6. Verify health check passes

---

## 🐳 What is Jenkins? (Bonus)

**Jenkins** = Open-source automation server (like GitHub Actions, but self-hosted).

| Aspect | Jenkins | GitHub Actions |
|--------|---------|----------------|
| **Where** | Your server | GitHub's cloud |
| **Cost** | Free (pay for server) | Free limits |
| **Setup** | Complex | Ready to use |
| **Maintenance** | You do it | GitHub does it |

**When to use Jenkins**: Enterprise needing full control  
**When to use GitHub Actions**: Quick setup, less maintenance  

---

## 📚 Additional Resources

These detailed guides cover advanced deployment topics:

| Document | Description |
|----------|-------------|
| [Environment Configurations](./docs/01-environment-configurations.md) | Managing Dev/Staging/Production environments, GitHub Secrets setup |
| [Rollback Strategies](./docs/02-rollback-strategies.md) | What to do when deployment fails, rollback scripts & workflows |
| [Health Checks & Verification](./docs/03-health-checks.md) | Verifying deployments, Docker health checks, auto-rollback |

---

*Last updated: September 2026*