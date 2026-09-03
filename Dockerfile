# ============================================
# DOCKERFILE: Build configuration for your Go app
# PURPOSE: Create a container image of your application
# ============================================

# ============================================
# STAGE 1: Build Stage
# ============================================
# We use a multi-stage build for smaller final image
# This stage compiles the code

# Base image: Go 1.25 on Alpine Linux (small, secure)
# Alpine is a minimal Linux distribution - great for containers!
FROM golang:1.25-alpine AS builder

# ============================================
# STEP 1: Set working directory
# ============================================
# Everything after this happens inside /app folder
WORKDIR /app

# ============================================
# STEP 2: Copy go mod files first (for caching)
# ============================================
# Docker caches layers. By copying go.mod first,
# we can reuse the dependency download layer even if code changes
COPY go.mod go.sum ./

# ============================================
# STEP 3: Download dependencies
# ============================================
# This layer is cached unless go.mod/go.sum change
RUN go mod download

# ============================================
# STEP 4: Copy source code
# ============================================
# Now copy the actual code
COPY main.go ./

# ============================================
# STEP 5: Build the application
# ============================================
# CGO_ENABLED=0 = build without C bindings (pure Go)
# -o myapp = output binary named "myapp"
# . = current directory (where main.go is)
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp .

# ============================================
# STAGE 2: Production Stage
# ============================================
# This creates the final, small production image

# Base image: Alpine Linux (only ~5MB!)
# No Go installed - just the runtime essentials
FROM alpine:latest

# ============================================
# STEP 6: Install CA certificates
# ============================================
# Some apps need to verify HTTPS connections
# ca-certificates contains root CA certificates
RUN apk --no-cache add ca-certificates

# ============================================
# STEP 7: Create app directory
# ============================================
WORKDIR /root/

# ============================================
# STEP 8: Copy binary from build stage
# ============================================
# --from=builder = copy from the "builder" stage we created earlier
# Only copies the compiled binary, not the source code!
COPY --from=builder /app/myapp .

# ============================================
# STEP 9: Expose port
# ============================================
# Document which port the app listens on
# This doesn't actually publish the port - just documents it
EXPOSE 8080

# ============================================
# STEP 10: Run the application
# ============================================
# CMD is the command that runs when container starts
# ./myapp = run our compiled binary
CMD ["./myapp"]
