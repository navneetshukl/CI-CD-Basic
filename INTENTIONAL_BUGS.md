# 🚨 INTENTIONAL BUGS FOR LEARNING

This file documents the bugs we intentionally introduced to demonstrate what CI/CD pipeline failures look like in GitHub Actions.

---

## What We Broke

### Bug 1: Code Formatting Issue
**File:** `main.go` (multiple locations)
**Issue:** Code needs to be formatted with `gofmt`

```bash
# Check locally:
gofmt -l .
# Output: main.go  ← This file needs formatting
```

**Expected GitHub Actions Error:**
```
❌ Files need formatting:
main.go

Run 'gofmt -w .' to fix formatting
```

---

### Bug 2: Printf Format Mismatch (go vet)
**File:** `main.go` line 128
**Issue:** Using `%d` (integer format) for a struct

```go
// WRONG:
fmt.Printf("User: %d\n", user)

// RIGHT:
fmt.Printf("User: %v\n", user)
// or
fmt.Printf("User: %+v\n", user)
```

**Expected GitHub Actions Error:**
```
./main.go:128:20: fmt.Printf format %d has arg user of wrong type *github.com/navneetshukla/test.User
```

---

### Bug 3: Staticcheck Warning
**File:** `main.go` line 128
**Issue:** Same as above, staticcheck catches it too

**Expected GitHub Actions Error:**
```
main.go:128:13: Printf format %d has arg #1 of wrong type *github.com/navneetshukla/test.User (SA5009)
```

---

### Bug 4: Failing Test
**File:** `main_test.go` line 166
**Issue:** Test expects "Bob" but gets "Alice"

```go
// WRONG:
if createdUser["name"] != "Bob" {
    t.Errorf("Expected name 'Bob', got '%s'", createdUser["name"])
}

// RIGHT:
if createdUser["name"] != "Alice" {
    t.Errorf("Expected name 'Alice', got '%s'", createdUser["name"])
}
```

**Expected GitHub Actions Error:**
```
main_test.go:167: Expected name 'Bob', got 'Alice'
--- FAIL: TestGetUserName (0.00s)
```

---

### Bug 5: Unused Function
**File:** `main.go` lines 26-29
**Issue:** `UnusedFunction()` is defined but never called

```go
// UnusedFunction is not used anywhere
func UnusedFunction() {
    fmt.Println("This function is never called")
}
```

**Note:** Staticcheck may or may not catch this depending on the version.

---

## What Will Happen in GitHub Actions

When you push this code:

```
Checkout code
      ↓
Set up Go
      ↓
Download dependencies
      ↓
Check code formatting (gofmt)
      ↓
❌ FAIL: main.go needs formatting
      ↓
[Pipeline stops here - remaining steps don't run]
```

OR if formatting passes:

```
...
      ↓
Run go vet (bug detection)
      ↓
❌ FAIL: Printf format %d has wrong type
      ↓
[Pipeline stops here]
```

OR if vet passes:

```
...
      ↓
Run staticcheck
      ↓
❌ FAIL: Printf format %d has wrong type (SA5009)
      ↓
[Pipeline stops here]
```

OR if staticcheck passes:

```
...
      ↓
Run tests with coverage
      ↓
❌ FAIL: TestGetUserName - Expected 'Bob', got 'Alice'
      ↓
[Pipeline stops here]
```

---

## How to See the Failures

1. Push the code to GitHub:
   ```bash
   cd ci-cd
   git add .
   git commit -m "Add intentional bugs for learning"
   git push origin main
   ```

2. Go to your GitHub repository → **Actions** tab

3. Click on the failing workflow run

4. You'll see something like:
   ```
   ❌ Code Quality Checks / Quality Checks
   
   Run go vet (bug detection)
   
   Error:
   ./main.go:128:20: fmt.Printf format %d has arg user of wrong type
   ```

5. Click on individual steps to see detailed logs

---

## How to Fix Each Issue

### Fix 1: Format Code
```bash
gofmt -w main.go main_test.go
```

### Fix 2: Fix Printf
```go
// Change this:
fmt.Printf("User: %d\n", user)

// To this:
fmt.Printf("User: %v\n", user)
```

### Fix 3: Fix Test
```go
// Change this:
if createdUser["name"] != "Bob" {

// To this:
if createdUser["name"] != "Alice" {
```

### Fix 4: Remove Unused Function
Delete these lines from main.go:
```go
// UnusedFunction is not used anywhere
func UnusedFunction() {
    fmt.Println("This function is never called")
}
```

---

## Learning Goals

After seeing these failures, you should understand:

1. ✅ How GitHub Actions displays errors
2. ✅ What each quality check tool catches
3. ✅ How to read error messages
4. ✅ How to fix each type of issue
5. ✅ Why these checks exist (they catch real bugs!)

---

## The Complete Fix

If you want to fix everything at once:

```bash
# 1. Format code
gofmt -w main.go main_test.go

# 2. Fix the Printf
sed -i '' 's/fmt.Printf("User: %d\\n", user)/fmt.Printf("User: %v\\n", user)/g' main.go

# 3. Fix the test
sed -i '' 's/createdUser\["name"\] != "Bob"/createdUser["name"] != "Alice"/g' main_test.go

# 4. Remove unused function
# (manually delete lines 26-29 in main.go)

# 5. Verify everything passes
go vet ./...
go test ./...
gofmt -l .
```

---

**Remember:** Breaking things on purpose is one of the best ways to learn! 🎓
