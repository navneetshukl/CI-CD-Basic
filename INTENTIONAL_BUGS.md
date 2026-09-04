# INTENTIONAL_BUGS.md

> A hands-on playbook to **break** this CI/CD pipeline on purpose.
> Each stage lists a tiny, surgical change you can make to see *exactly* why and where the pipeline fails.
> Roll it back, move to the next stage, learn the failure mode.

---

## ⚠️ How to use this file

1. Pick a stage below.
2. Apply **only** that stage's change.
3. Commit & push to `main` (or run the relevant workflow).
4. Watch the pipeline fail at that step. Read the error.
5. Revert the change (`git checkout -- <file>`) and move to the next stage.

> 💡 Tip: Commit each bug separately so you can clearly see which one broke the pipeline.

---

## 🧱 Stage 1 — Build fails (`01-basic-build.yml` / Build job in `00-complete-ci-cd.yml`)

**Why this fails:** the compiler cannot produce a binary because of a syntax error.

**File:** `main.go`

**Change:** remove a closing brace (or any `{` / `}`) so the file no longer compiles.

```go
// In func getUserByID, delete the opening brace of the function:
// Before:
// func getUserByID(c *gin.Context) {
//   ...
// }

// After (missing `{`):
// func getUserByID(c *gin.Context)
//   ...
// }
```

**Expected error in pipeline:**

```
./main.go:95:6: missing function body for "getUserByID"
# github.com/navneetshukla/test
exit status 1
```

---

## 🧪 Stage 2 — Tests fail (`02-add-tests.yml` / Test job)

**Why this fails:** a unit test asserts the wrong expected value.

**File:** `main_test.go`

**Change:** in `TestGetUserName`, change the expected name to something that won't match the user we create ("Alice").

```go
// Before (line ~166):
if createdUser["name"] != "Alice" {

// After:
if createdUser["name"] != "Bob" {
```

**Expected error in pipeline:**

```
--- FAIL: TestGetUserName (0.00s)
    main_test.go:167: Expected name 'Bob', got 'Alice'
FAIL
exit status 1
FAIL    github.com/navneetshukla/test  0.XXXs
```

---

## 🧹 Stage 3 — Code Quality fails (`03-code-quality.yml`)

This stage runs **gofmt** *and* **go vet**. Pick one (or try both!).

### 3a. `gofmt` failure — bad indentation

**File:** `main.go`

**Change:** inside `healthCheck`, replace the tab indentation with spaces on the `c.JSON` line (gofmt rejects any line that mixes tabs and spaces incorrectly).

```go
func healthCheck(c *gin.Context) {
   c.JSON(http.StatusOK, gin.H{   // ← replace the leading tab with 3 spaces
		"status":  "healthy",
		"message": "Server is running",
	})
}
```

**Expected error:**

```
❌ Files need formatting:
main.go
Run 'gofmt -w .' to fix formatting
```

### 3b. `go vet` failure — printf format mismatch

**File:** `main.go` (line ~123)

**Change:** use a `%d` (integer) verb for a non-integer value.

```go
// Before:
fmt.Printf("User: %v\n", user)

// After:
fmt.Printf("User: %d\n", user)   // %d expects an integer, but `user` is a *User
```

**Expected error:**

```
# github.com/navneetshukla/test
./main.go:123:13: Printf format %d has arg user of wrong type *User
```

---

## 🐳 Stage 4 — Docker build fails (`04-docker.yml`)

**Why this fails:** the Dockerfile tries to copy a file that doesn't exist, so `docker build` errors out before the image is produced.

**File:** `Dockerfile`

**Change:** change the source filename in the `COPY` instruction for the binary (or for the source — pick one).

```dockerfile
# Before (line 75):
COPY --from=builder /app/myapp .

# After:
COPY --from=builder /app/myappp .    # ← file does not exist
```

**Expected error in pipeline:**

```
ERROR: failed to solve: failed to compute cache key:
failed to calculate checksum of ref ...: "/app/myappp": not found
```

> ℹ️ Note: this pipeline uses Docker Buildx, so the error will mention `failed to solve` — not the older `open /var/lib/docker/tmp` style error.

> 🔁 **Bonus Docker failure:** in the same Dockerfile, change the Go version in the `FROM` line to one that doesn't exist, e.g. `FROM golang:9.99-alpine AS builder`. This will fail at the very first `FROM` step.

---

## 🚀 Stage 5 — Deploy fails (`05-deploy.yml` / Deploy job in `00-complete-ci-cd.yml`)

**Why this fails:** the container starts, but the health check pings the wrong port and gets nothing back.

**File:** the deploy script inside `.github/workflows/05-deploy.yml` (or the `deploy` job in `00-complete-ci-cd.yml`).

**Change:** in the `docker run` command, change the container port from `8080` to `9090`, but keep the health check hitting `8080`.

```bash
# Before (in the SSH deploy script):
docker run -d \
  --name ${{ env.IMAGE_NAME }} \
  -p 8080:8080 \
  --restart unless-stopped \
  ${{ env.DOCKER_USERNAME }}/${{ env.IMAGE_NAME }}:latest

# After:
docker run -d \
  --name ${{ env.IMAGE_NAME }} \
  -p 8080:9090 \          # ← host:container — app is on 9090, health check hits 8080
  --restart unless-stopped \
  ${{ env.DOCKER_USERNAME }}/${{ env.IMAGE_NAME }}:latest
```

**Expected error in pipeline:**

```
❌ Application health check failed!
📋 Container logs:
curl: (7) Failed to connect to localhost port 8080 after 5 ms: Connection refused
```

---

## 🛠️ Bonus — Workflow syntax fails (any `.yml` in `.github/workflows/`)

**Why this fails:** a YAML indentation/typo makes GitHub Actions unable to parse the workflow at all.

**File:** any workflow file, e.g. `.github/workflows/01-basic-build.yml`

**Change:** remove the space before `- name:` in any step.

```yaml
    steps:
      - name: Get source code
      - name: Set up Go         # ← delete the leading two spaces on this line
        uses: actions/setup-go@v5
```

**Expected error:**

```
Invalid workflow file: .github/workflows/01-basic-build.yml
Line: 44  Column: 7
Error: (Line: 44, Col: 7): Sequence entries are not allowed here
```

---

## ✅ Reset checklist

After each experiment, undo the change:

```bash
git diff                 # see what you changed
git checkout -- <file>   # revert one file
# or
git restore <file>
```

Then move on to the next stage.

Happy breaking! 💥
