# CI/CD Learning Guide 🚀

A beginner-friendly guide to understanding and implementing CI/CD pipelines from scratch.

---

## 📚 Table of Contents

1. [What is CI/CD?](#what-is-cicd)
2. [Key Concepts](#-key-concepts)
3. [Pipeline Flow](#-pipeline-flow)
4. [Step 1: Basic Build Pipeline](#-step-1-basic-build-pipeline)
5. [How to Use](#-how-to-use)
6. [Next Steps](#-next-steps)

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
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│  Commit  │───▶│  Build   │───▶│   Test   │───▶│  Deploy  │
│   Code   │    │          │    │          │    │          │
└──────────┘    └──────────┘    └──────────┘    └──────────┘
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

      - name: Build the Go app
        run: go build -v ./...

      - name: Run the app
        run: go run main.go
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

### 6. Step 3 - Build
```yaml
      - name: Build the Go app
        run: go build -v ./...
```
👉 **STEP 3**: Compile your Go code

`./...` means "all packages in current directory and subdirectories"
`-v` flag shows verbose output (what's being built)

---

### 7. Step 4 - Run
```yaml
      - name: Run the app
        run: go run main.go
```

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
Step 3: Runs "go build -v ./..."
         ↓
Step 4: Runs "go run main.go"
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

## 🎯 Next Steps

Now that you understand the basics, here are the natural progressions:

### Level 2: Add Tests
- Add automated unit tests to catch bugs
- See `02-add-tests.yml` (coming soon)

### Level 3: Code Quality
- Add linting (go vet, golint)
- Add code formatting checks
- See `03-code-quality.yml` (coming soon)

### Level 4: Docker
- Build a container image
- Push to Docker Hub
- See `04-docker-build.yml` (coming soon)

### Level 5: Deployment
- Deploy to a server automatically
- Multiple environments (dev, staging, prod)
- See `05-deploy.yml` (coming soon)

---

## 🧪 Step 2: Adding Automated Tests

### Why Tests Matter
- **Catch bugs early** - Find issues before they reach production
- **Confidence** - Know your changes don't break existing code
- **Documentation** - Tests show how your code should be used
- **Safety net** - Refactor code without fear of breaking things

### What We Added
1. **Gin web framework** - For handling HTTP requests
2. **REST API endpoints** - GET and POST routes
3. **Unit tests** - 6 tests covering all endpoints
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
Run all tests (6 tests)
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
PASS    coverage: 85.7% of statements
```

### Benefits You Get Now

✅ **Automatic testing** - Tests run on every push  
✅ **PR validation** - Can't merge code that breaks tests  
✅ **Fast feedback** - Know within minutes if something broke  
✅ **Confidence** - Deploy with confidence knowing tests pass  

---

*Last updated: September 2026*

---

## 🔗 Useful Resources

| Resource | Link |
|----------|------|
| GitHub Actions Docs | https://docs.github.com/en/actions |
| Go GitHub Action | https://github.com/actions/setup-go |
| Checkout Action | https://github.com/actions/checkout |

---

## 💡 Tips for Beginners

1. **Start simple** - Get a basic pipeline working first, then add complexity
2. **Check the Actions tab** - It's your best friend for debugging
3. **Use the marketplace** - Thousands of pre-built actions available
4. **Read the logs** - Pipeline failures usually have clear error messages
5. **Use `run: echo "debug"`** - Add debug steps to see what's happening

---

*Last updated: September 2026*
👉 **STEP 4**: Execute the program to verify it works
