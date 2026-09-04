# 🚨 INTENTIONAL BUGS FOR LEARNING

This file shows you **exact line changes** to trigger specific CI/CD failures.

---

## ⚠️ Important

**Never paste new blocks that redeclare existing code!**

Your main.go already has: package main, type User struct, func main()
Your main_test.go already has: func TestGetUserName(t *testing.T)

For bugs 1-3: we show which line to CHANGE.
For bugs 4-5: safe to append (new unique names).

---

## Bug 1: Code Formatting (gofmt fails)

**File:** main.go

**Find (line ~27):**
```
func main() {
    // Initialize Gin router
    r := gin.Default()
```

**Change to:**
```
func main() {
       // Initialize Gin router   ← extra spaces (bad!)
    r := gin.Default()
```

Or run:
```bash
sed -i "" 's|// Initialize Gin router|       // Initialize Gin router|' main.go
```

**Expected:** `❌ Files need formatting: main.go`
**Fix:** `gofmt -w main.go`


---

## Bug 2: Printf Format Mismatch (go vet fails)

**File:** main.go
**Find (line 123):**
```go
fmt.Printf("User: %v\n", user) // Fixed: changed %d to %v
```

**Change %v to %d:**
```go
fmt.Printf("User: %d\n", user) // BAD: %d for struct
```

**Expected:** `fmt.Printf format %d has arg user of wrong type *main.User`
**Fix:** Change %d back to %v


---

## Bug 3: Failing Test

**File:** main_test.go
**Find (line 166):**
```go
if createdUser["name"] != "Alice" {
    t.Errorf("Expected name 'Alice', got '%s'", createdUser["name"])
}
```

**Change "Alice" to "Bob":**
```go
if createdUser["name"] != "Bob" {  // BAD: actual is "Alice"
    t.Errorf("Expected name 'Bob', got '%s'", createdUser["name"])
}
```

**Expected:** `Expected name 'Bob', got 'Alice' --- FAIL: TestGetUserName`
**Fix:** Change "Bob" back to "Alice"


---

## Bug 4: Unused Function (staticcheck fails)

**File:** main.go
**Add to END of main.go (safe - new unique name):**
```go
// Add to bottom of main.go
func UnusedFunctionForCI() {
    fmt.Println("This function is never called")
    x := 10
    y := 20
    _ = x + y
}
```

**Expected:** `func UnusedFunctionForCI is unused (U1000)`
**Fix:** Delete the function


---

## Bug 5: Syntax Error (build fails)

**File:** main.go
**Add to END of main.go (safe - new unique name):**
```go
// Add to bottom of main.go
func BrokenFunctionForCI() {
    fmt.Println("missing closing brace")
// Missing } here
```

**Expected:** `missing '{' or ')' at end of code block`
**Fix:** Add the missing }


---

## Bug 6: Wrong Go Version

**File:** .github/workflows/00-complete-ci-cd.yml
**Change:**
```yaml
go-version: '1.25'
```
**To:**
```yaml
go-version: '99.99'  # BAD: does not exist
```

**Expected:** `Unable to find Go version '99.99'`
**Fix:** Change back to '1.25'


---

## Bug 7: Docker Image Wrong

**File:** Dockerfile
**Change:**
```dockerfile
FROM golang:1.25-alpine
```
**To:**
```dockerfile
FROM golang:nonexistent-version
```

**Expected:** `failed to solve: golang:nonexistent-version`
**Fix:** Change back to `FROM golang:1.25-alpine`


---

## Bug 8: Missing Secret

**File:** GitHub Settings → Secrets and variables → Actions
**Action:** Delete the DOCKER_USERNAME secret

**Expected:** `Username and password required`
**Fix:** Re-add the secret in GitHub Settings


---

## 🧪 Safe Reproduction Script

### What is `.bak`?
When you run the script, it creates backup files (`.bak`) of your original code:
- `main.go.bak` = copy of your original `main.go`
- `main_test.go.bak` = copy of your original `main_test.go`

These let you restore the original code after testing the bugs.

### Run the Script

```bash
#!/bin/bash
# apply-bugs.sh - Apply bugs with LINE-LEVEL changes

set -e

echo "Backing up..."
cp main.go main.go.bak
cp main_test.go main_test.go.bak

echo "Bug 1: Bad formatting..."
sed -i "" 's|// Initialize Gin router|       // Initialize Gin router|' main.go

echo "Bug 2: Printf wrong format..."
sed -i '' 's/%v\\n/%d\\n/' main.go

echo "Bug 3: Wrong expected name..."
sed -i "" 's|!= "Alice"|!= "Bob"|' main_test.go

echo "Bug 4: Unused function (append)..."
cat >> main.go << 'EOF'
func UnusedFunctionForCI() {
    fmt.Println("This function is never called")
}
EOF

echo "Bug 5: Syntax error (append)..."
cat >> main.go << 'EOF'
func BrokenFunctionForCI() {
    fmt.Println("missing closing brace")
# Missing }
EOF

echo "Done! Push to see failures:"
echo "git add . && git commit -m 'Apply bugs' && git push"
```

**Restore:**
```bash
cp main.go.bak main.go && cp main_test.go.bak main_test.go
```


---

## What Fails in Pipeline

| Bugs Present | Pipeline Stops At |
|--------------|-------------------|
| Bug 5 (syntax) | Step 4: Build |
| Bug 5 fixed, Bug 3 | Step 5: Tests |
| Bugs 3,5 fixed, Bug 1 | Step 6: gofmt |
| Bugs 1,3,5 fixed, Bug 2 | Step 7: go vet |
| Bugs 1,2,3,5 fixed, Bug 4 | Step 8: staticcheck |

---

## Try One Bug at a Time

Best learning: apply one, fix it, push, watch it pass, then try the next.

**Remember:** Each bug should produce its specific error - not "redeclared" errors! 🎓
