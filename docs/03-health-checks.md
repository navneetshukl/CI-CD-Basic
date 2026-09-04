# Health Checks & Verification in CI/CD

This guide teaches you how to verify that your application is working correctly after deployment.

---

## 1. What is a Health Check?

A **health check** is a way to verify that your application is running and responding correctly.

### Simple Explanation

Like going to the doctor:
- Doctor checks: "Are you healthy?"
- If yes: You can go back to work
- If no: You need treatment

**Health check = Asking your app "Are you okay?"**

---

## 2. Why Do We Need Health Checks?

### The Problem Without Health Checks

```
Deploy new version
Container starts
App crashes internally
Container still "running" (but not working)
Users see errors
You find out hours later from complaints
```

### The Solution With Health Checks

```
Deploy new version
Container starts
Health check runs automatically
Detects app is broken
Automatically rolls back
Users never see the broken version
```

---

## 3. Types of Health Checks

### Type 1: HTTP Health Check (Most Common)

**How it works:**
1. Send HTTP request to `/health` endpoint
2. If response is 200 OK → App is healthy
3. If response is error → App is unhealthy

**Example:**

```bash
curl http://localhost:8080/health
```

**Expected response (healthy):**

```json
{"status": "healthy"}
```

### Type 2: TCP Health Check

**How it works:**
1. Try to connect to a port
2. If connection succeeds → Service is up
3. If connection fails → Service is down

**Example:**

```bash
nc -z localhost 8080 && echo "Port is open"
```

### Type 3: Process Health Check

**How it works:**
1. Check if the process is running
2. If process exists → App is running
3. If process doesn't exist → App crashed

**Example:**

```bash
pgrep -f "my-go-api" && echo "Running"
```

---

## 4. Implementing Health Check Endpoint

### Basic Health Check (Go + Gin)

```go
package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func main() {
    r := gin.Default()
    
    // Basic health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status": "healthy",
        })
    })
    
    r.Run(":8080")
}
```

### Advanced Health Check (With Database)

```go
package main

import (
    "database/sql"
    "github.com/gin-gonic/gin"
    "net/http"
)

func main() {
    r := gin.Default()
    
    r.GET("/health", func(c *gin.Context) {
        // Check database connection
        db, err := sql.Open("postgres", "your-db-connection")
        if err != nil {
            c.JSON(http.StatusServiceUnavailable, gin.H{
                "status": "unhealthy",
                "error":  "database connection failed",
            })
            return
        }
        defer db.Close()
        
        // Ping database
        if err := db.Ping(); err != nil {
            c.JSON(http.StatusServiceUnavailable, gin.H{
                "status": "unhealthy",
                "error":  "database ping failed",
            })
            return
        }
        
        // All good!
        c.JSON(http.StatusOK, gin.H{
            "status":   "healthy",
            "database": "connected",
            "version":  "1.0.0",
        })
    })
    
    r.Run(":8080")
}
```

### Health Check Response Examples

**Healthy Response (200 OK):**

```json
{
    "status": "healthy",
    "database": "connected",
    "version": "1.0.0"
}
```

**Unhealthy Response (503):**

```json
{
    "status": "unhealthy",
    "error": "database connection failed"
}
```

---

## 5. Docker Health Checks

### Method 1: Dockerfile HEALTHCHECK

```dockerfile
FROM golang:1.25

# ... your app setup ...

HEALTHCHECK --interval=30s --timeout=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1
```

**What each option means:**

| Option | Description |
|--------|-------------|
| `--interval` | How often to check (30s) |
| `--timeout` | How long to wait (5s) |
| `--retries` | Failures before unhealthy (3) |
| `--start-period` | Wait before first check |

### Method 2: Docker Compose Health Check

```yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8080:8080"
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
```

### Method 3: Docker Run with Health Check

```bash
docker run -d \
  --name my-go-api \
  -p 8080:8080 \
  --health-cmd="curl -f http://localhost:8080/health || exit 1" \
  --health-interval=30s \
  --health-timeout=5s \
  --health-retries=3 \
  navneetshukla/my-go-api:latest
```

### Check Health Status

```bash
docker ps
```

**Output shows:**
- `STATUS: Up 5 minutes (healthy)`
- `STATUS: Up 2 minutes (unhealthy)`
- `STATUS: Up 1 minute (health: starting)`

---


## 6. GitHub Actions Health Check

Add this to your deployment workflow:

```yaml
- name: Verify deployment
  uses: appleboy/ssh-action@v1.0.3
  with:
    host: ${{ secrets.SERVER_HOST_PROD }}
    username: ${{ secrets.SERVER_USER }}
    key: ${{ secrets.SSH_PRIVATE_KEY }}
    script: |
      echo "Running health checks..."
      sleep 10
      
      if docker ps | grep -q my-go-api; then
        echo "Container is running"
      else
        echo "Container not running!"
        exit 1
      fi
      
      if curl -f http://localhost:8080/health; then
        echo "Health check passed"
      else
        echo "Health check failed!"
        docker logs my-go-api --tail 20
        exit 1
      fi
```

## 7. Auto-Rollback on Health Check Failure

```yaml
- name: Health Check with Rollback
  uses: appleboy/ssh-action@v1.0.3
  with:
    host: ${{ secrets.SERVER_HOST_PROD }}
    username: ${{ secrets.SERVER_USER }}
    key: ${{ secrets.SSH_PRIVATE_KEY }}
    script: |
      sleep 15
      
      if curl -f http://localhost:8080/health; then
        echo "Health check PASSED"
      else
        echo "FAILED - Rolling back..."
        docker stop my-go-api || true
        docker rm my-go-api || true
        docker run -d --name my-go-api -p 8080:8080 navneetshukla/my-go-api:backup
        exit 1
      fi
```

## 8. Health Check Best Practices

### 1. Check All Critical Dependencies

```go
r.GET("/health", func(c *gin.Context) {
    status := gin.H{"status": "healthy"}
    allHealthy := true
    
    if err := db.Ping(); err != nil {
        status["database"] = "unhealthy"
        allHealthy = false
    } else {
        status["database"] = "healthy"
    }
    
    if !allHealthy {
        c.JSON(503, status)
        return
    }
    
    c.JSON(200, status)
})
```

### 2. Set Appropriate Timeouts

| Check | Timeout | Why |
|-------|---------|-----|
| Health check | 5s | App should respond fast |
| Database | 3s | DB should be quick |
| External API | 10s | Network can be slow |

### 3. Don't Check Too Often

| Environment | Interval | Reason |
|-------------|----------|--------|
| Development | 60s | Less critical |
| Staging | 30s | Normal monitoring |
| Production | 15-30s | Need to know fast |

### 4. Return Useful Information

```json
{
    "status": "healthy",
    "version": "1.2.3",
    "uptime": "2 days 5 hours",
    "database": "connected"
}
```

## 9. Common Health Check Commands

### Check Endpoint

```bash
# Basic health check
curl http://localhost:8080/health

# With verbose output
curl -v http://localhost:8080/health

# Check HTTP status only
curl -o /dev/null -s -w "%{http_code}" http://localhost:8080/health
```

### Check Container

```bash
# Check if running
docker ps | grep my-go-api

# Check health status
docker inspect --format='{{.State.Health.Status}}' my-go-api

# Check container logs
docker logs my-go-api --tail 50
```

## 10. Troubleshooting Health Checks

### Problem 1: Health Check Always Fails

**Check:**
```bash
curl http://localhost:8080/health
docker ps
docker logs my-go-api
```

**Common Causes:**
- Wrong port in health check
- App not listening on expected interface
- Firewall blocking requests

### Problem 2: Container Marked Unhealthy

**Check:**
```bash
docker inspect my-go-api | grep -A 10 Health
```

**Common Causes:**
- App takes too long to start
- Database connection failing
- Memory limit too low

### Problem 3: Intermittent Failures

**Check:**
```bash
docker stats my-go-api
top
```

**Common Causes:**
- High CPU/memory usage
- Database connection pool exhausted

## 11. Health Check Checklist

### Application Code
- [ ] `/health` endpoint implemented
- [ ] Returns 200 when healthy
- [ ] Returns 503 when unhealthy
- [ ] Checks critical dependencies

### Docker Configuration
- [ ] HEALTHCHECK in Dockerfile
- [ ] Appropriate interval (30s)
- [ ] Appropriate timeout (5s)
- [ ] Appropriate retries (3)

### Deployment Workflow
- [ ] Health check after deploy
- [ ] Wait time for startup
- [ ] Auto-rollback on failure
- [ ] Log health check results

## Related Topics

- [Environment Configurations](./01-environment-configurations.md)
- [Rollback Strategies](./02-rollback-strategies.md)

## Questions to Check Understanding

1. What is a health check?
2. Why do we need health checks?
3. How often should you run health checks in production?
4. What should you do if health check fails after deployment?

*Last updated: September 2026*
