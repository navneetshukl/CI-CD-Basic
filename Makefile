# ============================================
# MAKEFILE: Shortcut commands for common tasks
# PURPOSE: Type 'make build' instead of long Docker commands
# ============================================

# Variables
IMAGE_NAME := my-go-api
DOCKER_USERNAME := navneetshukla  # CHANGE THIS to your Docker username

# ============================================
# TARGET: Build Docker image
# ============================================
# Usage: make build
build:
	docker build -t $(IMAGE_NAME):latest .

# ============================================
# TARGET: Build with username
# ============================================
# Usage: make build-full
build-full:
	docker build -t $(DOCKER_USERNAME)/$(IMAGE_NAME):latest .

# ============================================
# TARGET: Run container locally
# ============================================
# Usage: make run
run:
	docker run -p 8080:8080 $(IMAGE_NAME):latest

# ============================================
# TARGET: Run in background (detached)
# ============================================
# Usage: make run-detached
run-detached:
	docker run -d -p 8080:8080 --name $(IMAGE_NAME) $(IMAGE_NAME):latest

# ============================================
# TARGET: Stop and remove container
# ============================================
# Usage: make stop
stop:
	docker stop $(IMAGE_NAME) || true
	docker rm $(IMAGE_NAME) || true

# ============================================
# TARGET: Run with docker-compose
# ============================================
# Usage: make up
up:
	docker-compose up --build

# ============================================
# TARGET: Stop docker-compose
# ============================================
# Usage: make down
down:
	docker-compose down

# ============================================
# TARGET: Clean up all Docker resources
# ============================================
# Usage: make clean
clean: stop
	docker rmi $(IMAGE_NAME):latest || true

# ============================================
# TARGET: Show logs
# ============================================
# Usage: make logs
logs:
	docker logs -f $(IMAGE_NAME)

# ============================================
# TARGET: Run tests
# ============================================
# Usage: make test
test:
	go test -v ./...

# ============================================
# TARGET: Run all quality checks
# ============================================
# Usage: make check
check:
	go vet ./...
	gofmt -l .
	./path/to/staticcheck ./...

# ============================================
# TARGET: Build and push to Docker Hub
# ============================================
# Usage: make push
push: build-full
	docker push $(DOCKER_USERNAME)/$(IMAGE_NAME):latest

# ============================================
# TARGET: Pull latest from Docker Hub
# ============================================
# Usage: make pull
pull:
	docker pull $(DOCKER_USERNAME)/$(IMAGE_NAME):latest

# ============================================
# TARGET: Show help
# ============================================
# Usage: make help
help:
	@echo "Available commands:"
	@echo "  make build         - Build Docker image"
	@echo "  make run           - Run container locally"
	@echo "  make run-detached  - Run in background"
	@echo "  make stop          - Stop and remove container"
	@echo "  make up            - Run with docker-compose"
	@echo "  make down          - Stop docker-compose"
	@echo "  make test          - Run tests"
	@echo "  make push          - Build and push to Docker Hub"
	@echo "  make pull          - Pull from Docker Hub"
	@echo "  make clean         - Clean up Docker resources"
	@echo "  make help          - Show this help"

# ============================================
# PHONY: These aren't real files
# ============================================
# This tells make not to look for actual files named 'build', 'run', etc.
.PHONY: build build-full run run-detached stop up down clean logs test check push pull help
