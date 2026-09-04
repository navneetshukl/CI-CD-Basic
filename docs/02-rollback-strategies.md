# Rollback Strategies in CI/CD

This guide teaches you what to do when a deployment fails and how to revert to a previous working version.

---

## 1. What is a Rollback?

A **rollback** is going back to a previous working version of your application.

### Simple Explanation

Imagine you're reading a book:
- You saved your progress at page 50
- You tried reading page 51 but it's confusing
- You go back to page 50 and start again

**Rollback = Going back to a working version**

---

## 2. Why Do We Need Rollback?

### Things That Can Go Wrong

| Problem | What Happens |
|---------|-------------|
| Bug in new code | App crashes |
| Database issues | Data gets corrupted |
| Performance drop | App becomes slow |
| Feature breaks | Users can't use app |
| Security issue | App becomes vulnerable |

### Without Rollback

```
❌ User sees error
❌ Users complain
❌ You scramble to fix
❌ More bugs introduced while fixing
❌ Hours of downtime
```

### With Rollback

```
✅ Instant fix - go back to working version
✅ Users happy
✅ You fix bugs safely in dev
✅ Deploy when ready
```

---

## 3. Types of Rollback Strategies

### Strategy 1: Docker Image Rollback (Recommended)

**How it works:**
1. Previous Docker image is already saved
2. Simply stop current container
3. Start container with previous image
4. App is back to working version

**Benefits:**
- Very fast (1-2 minutes)
- No code changes needed
- Easy to automate
- Can be done manually

### Strategy 2: Blue-Green Deployment

**How it works:**
1. Two identical environments (Blue and Green)
2. One is live, one is standby
3. Switch traffic from one to another
4. To rollback: switch back

**Benefits:**
- Zero downtime rollback
- Instant switch

**Drawbacks:**
- Requires double infrastructure
- More expensive

### Strategy 3: Feature Flag Rollback

**How it works:**
1. New features are behind "flags"
2. To rollback: disable the flag
3. Feature is hidden but code is still there

---

## 4. When to Rollback?

### 🚨 IMMEDIATE Rollback

| Situation | Action |
|-----------|--------|
| App completely down | Rollback NOW |
| Critical bug affecting all users | Rollback NOW |
| Security vulnerability | Rollback NOW |
| Data corruption | Rollback NOW |

### 🤔 MONITOR First

| Situation | Action |
|-----------|--------|
| Minor UI glitch | Monitor, rollback if worsens |
| Performance slightly slower | Monitor for 15 minutes |
| One feature not working | Check logs |

### ❌ DO NOT Rollback

| Situation | Why Not |
|-----------|---------|
| Performance being investigated | Wait for root cause |
| Single user reporting issue | Check their device |
| Scheduled maintenance | Wait for window |

---

## 5. How to Implement Rollback

### Option 1: Manual Rollback (Simple)

**Step 1: SSH into your server**

```bash
ssh deploy@your-server-ip
```

**Step 2: Check running containers**

```bash
docker ps
```

**Step 3: Stop current container**

```bash
docker stop my-go-api
docker rm my-go-api
```

**Step 4: Start previous version**

```bash
# List available images
docker images | grep my-go-api

# Run previous version (use the SHA tag)
docker run -d \
  --name my-go-api \
  -p 8080:8080 \
  --restart unless-stopped \
  navneetshukla/my-go-api:abc1234
```

**Step 5: Verify rollback**

```bash
curl http://localhost:8080/health
```

---


---

## 8. Rollback Best Practices

### 1. Always Tag Your Images

```yaml
# ❌ BAD - Only using 'latest'
docker build -t navneetshukla/my-go-api:latest .

# ✅ GOOD - Use multiple tags
docker build -t navneetshukla/my-go-api:latest \
             -t navneetshukla/my-go-api:${{ github.sha }} .
```

### 2. Keep a Backup Image

```bash
# Before deploying new version
docker tag navneetshukla/my-go-api:latest navneetshukla/my-go-api:backup
docker push navneetshukla/my-go-api:backup
```

### 3. Monitor After Rollback

After rolling back, check:
- ✅ Error rates returning to normal
- ✅ Response times improving
- ✅ No new errors appearing
- ✅ Users able to use the app

### 4. Communicate with Team

```markdown
🚨 INCIDENT: Rollback initiated

What: Production deployment rolled back
When: 2024-09-06 10:30 AM
Why: Critical bug affecting all users
Previous Version: xyz5678
Current Version: abc1234

Next Steps:
1. Bug being fixed in development
2. Will redeploy after QA testing
```

---

## 9. Rollback Decision Tree

```
Deployment fails?
        │
        ├─── No ──→ ✅ Monitor normally
        │
        └─── Yes ──→ Is it critical?
                         │
                         ├─── Yes ──→ 🚨 IMMEDIATE ROLLBACK
                         │
                         └─── No ──→ Is it getting worse?
                                        │
                                        ├─── Yes ──→ 🔄 ROLLBACK
                                        │
                                        └─── No ──→ 📊 MONITOR
```

---

## 10. Common Rollback Commands

### Check Status

```bash
# See running containers
docker ps

# See all images
docker images | grep my-go-api

# Check container logs
docker logs my-go-api --tail 100

# Check app health
curl http://localhost:8080/health
```

### Stop and Remove

```bash
# Stop container
docker stop my-go-api

# Remove container
docker rm my-go-api

# Force remove
docker rm -f my-go-api
```

### Start Previous Version

```bash
# Run specific version
docker run -d \
  --name my-go-api \
  -p 8080:8080 \
  navneetshukla/my-go-api:abc1234
```

---

## 11. Rollback Checklist

### Before Deployment
- [ ] Tag current image as backup
- [ ] Test rollback script
- [ ] Notify team of deployment
- [ ] Have rollback plan ready

### During Incident
- [ ] Assess severity quickly
- [ ] Decide: Rollback or Fix Forward
- [ ] Execute rollback if needed
- [ ] Verify app is working
- [ ] Notify team of status

### After Rollback
- [ ] Monitor for stability
- [ ] Document what went wrong
- [ ] Create fix in development
- [ ] Plan new deployment

---

## 🔗 Related Topics

- [Environment Configurations](./01-environment-configurations.md)
- [Health Checks](./03-health-checks.md)

---

## ❓ Questions to Check Understanding

1. What is a rollback?
2. When should you immediately rollback?
3. What is the fastest way to rollback with Docker?
4. Why should you tag images with commit SHA?
5. What should you do after rolling back?

---

*Last updated: September 2026*
