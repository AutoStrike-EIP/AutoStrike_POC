# AutoStrike - Breach and Attack Simulation Platform
# Main Makefile

.PHONY: all build clean test dev docker help
.PHONY: server-build server-dev server-test
.PHONY: agent-build agent-test
.PHONY: dashboard-build dashboard-dev dashboard-test
.PHONY: certs docker-build docker-up docker-down

# Variables
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Colors
CYAN := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RESET := \033[0m

help: ## Show this help
	@echo "$(CYAN)AutoStrike$(RESET) - Breach and Attack Simulation Platform"
	@echo ""
	@echo "$(GREEN)Available targets:$(RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(CYAN)%-20s$(RESET) %s\n", $$1, $$2}'

all: build ## Build all components

# =============================================================================
# Development
# =============================================================================

dev: ## Start development environment
	@echo "$(YELLOW)Starting development environment...$(RESET)"
	@$(MAKE) -j3 server-dev dashboard-dev

install: ## Install all dependencies
	@echo "$(YELLOW)Installing dependencies...$(RESET)"
	cd server && go mod download
	cd agent && cargo fetch
	cd dashboard && npm install

# =============================================================================
# Build
# =============================================================================

build: server-build agent-build dashboard-build ## Build all components
	@echo "$(GREEN)All components built successfully!$(RESET)"

server-build: ## Build the Go server
	@echo "$(YELLOW)Building server...$(RESET)"
	cd server && CGO_ENABLED=1 go build -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)" -o ../dist/autostrike ./cmd/autostrike

agent-build: ## Build the Rust agent
	@echo "$(YELLOW)Building agent...$(RESET)"
	cd agent && cargo build --release
	cp agent/target/release/autostrike-agent dist/

dashboard-build: ## Build the React dashboard
	@echo "$(YELLOW)Building dashboard...$(RESET)"
	cd dashboard && npm run build
	cp -r dashboard/dist dist/dashboard

# =============================================================================
# Development Servers
# =============================================================================

server-dev: ## Run server in development mode
	cd server && go run ./cmd/autostrike

dashboard-dev: ## Run dashboard in development mode
	cd dashboard && npm run dev

# =============================================================================
# Testing
# =============================================================================

test: server-test agent-test dashboard-test ## Run all tests

server-test: ## Run server tests
	@echo "$(YELLOW)Testing server...$(RESET)"
	cd server && go test -v ./...

agent-test: ## Run agent tests
	@echo "$(YELLOW)Testing agent...$(RESET)"
	cd agent && cargo test

dashboard-test: ## Run dashboard tests
	@echo "$(YELLOW)Testing dashboard...$(RESET)"
	cd dashboard && npm run test

lint: ## Run linters
	@echo "$(YELLOW)Linting...$(RESET)"
	cd server && go vet ./...
	cd agent && cargo clippy
	cd dashboard && npm run lint

# =============================================================================
# Docker
# =============================================================================

docker-build: ## Build Docker images
	@echo "$(YELLOW)Building Docker images...$(RESET)"
	docker-compose build

docker-up: ## Start Docker containers
	@echo "$(YELLOW)Starting containers...$(RESET)"
	docker-compose up -d

docker-down: ## Stop Docker containers
	@echo "$(YELLOW)Stopping containers...$(RESET)"
	docker-compose down

docker-logs: ## Show Docker logs
	docker-compose logs -f

# =============================================================================
# Certificates
# =============================================================================

certs: ## Generate TLS certificates
	@echo "$(YELLOW)Generating certificates...$(RESET)"
	chmod +x scripts/generate-certs.sh
	./scripts/generate-certs.sh ./certs

# =============================================================================
# Utilities
# =============================================================================

clean: ## Clean build artifacts
	@echo "$(YELLOW)Cleaning...$(RESET)"
	rm -rf dist/
	rm -rf server/autostrike
	rm -rf agent/target/
	rm -rf dashboard/dist/
	rm -rf dashboard/node_modules/

setup: install certs ## Initial project setup
	@echo "$(GREEN)Project setup complete!$(RESET)"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Configure server/config/config.yaml"
	@echo "  2. Run 'make dev' to start development"
	@echo "  3. Run 'make docker-up' for production"

dist: ## Create distribution package
	mkdir -p dist
	$(MAKE) build
	./scripts/build-agent.sh $(VERSION)
	tar -czvf autostrike-$(VERSION).tar.gz dist/

import-techniques: ## Import Atomic Red Team techniques
	chmod +x scripts/import-atomic.sh
	./scripts/import-atomic.sh ./data/techniques
