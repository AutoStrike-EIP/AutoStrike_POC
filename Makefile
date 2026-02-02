# AutoStrike - Breach and Attack Simulation Platform
# Main Makefile

.PHONY: all build clean test dev docker help
.PHONY: server-build server-dev server-test
.PHONY: agent-build agent-test
.PHONY: dashboard-build dashboard-test
.PHONY: certs docker-build docker-up docker-down
.PHONY: run stop logs agent

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
	@echo "$(GREEN)Quick Start:$(RESET)"
	@echo "  $(CYAN)make run$(RESET)    - Build and start everything"
	@echo "  $(CYAN)make agent$(RESET)  - Connect an agent"
	@echo "  $(CYAN)make stop$(RESET)   - Stop all services"
	@echo ""
	@echo "$(GREEN)Available targets:$(RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(CYAN)%-20s$(RESET) %s\n", $$1, $$2}'

all: build ## Build all components

# =============================================================================
# Quick Start - Single command to run everything
# =============================================================================

run: stop server-build-quick dashboard-build-quick ## Build and run the server (serves API + Dashboard)
	@echo ""
	@echo "$(GREEN)════════════════════════════════════════$(RESET)"
	@echo "$(GREEN)  Starting AutoStrike...$(RESET)"
	@echo "$(GREEN)════════════════════════════════════════$(RESET)"
	@echo ""
	@# Setup directories
	@mkdir -p server/data server/certs
	@# Generate certs if missing
	@if [ ! -f server/certs/server.crt ]; then \
		echo "$(YELLOW)Generating TLS certificates...$(RESET)"; \
		openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
			-keyout server/certs/server.key \
			-out server/certs/server.crt \
			-subj "/CN=localhost" 2>/dev/null; \
		cp server/certs/server.crt server/certs/ca.crt; \
	fi
	@# Start server (serves both API and Dashboard)
	@echo "$(YELLOW)Starting server...$(RESET)"
	@cd server && ./autostrike-server > /tmp/autostrike-server.log 2>&1 &
	@sleep 2
	@curl -s http://localhost:8443/health > /dev/null && echo "$(GREEN)✓ Server running$(RESET)" || echo "$(YELLOW)Server starting...$(RESET)"
	@echo ""
	@echo "$(GREEN)════════════════════════════════════════$(RESET)"
	@echo "$(GREEN)  AutoStrike is running!$(RESET)"
	@echo "$(GREEN)════════════════════════════════════════$(RESET)"
	@echo ""
	@echo "  $(CYAN)http://localhost:8443$(RESET)"
	@echo ""
	@echo "  Routes:"
	@echo "    /           Dashboard"
	@echo "    /api/v1/*   REST API"
	@echo "    /ws/*       WebSocket"
	@echo "    /health     Health check"
	@echo ""
	@echo "  Commands:"
	@echo "    $(CYAN)make agent$(RESET)  - Connect an agent"
	@echo "    $(CYAN)make stop$(RESET)   - Stop all services"
	@echo "    $(CYAN)make logs$(RESET)   - View logs"
	@echo ""

dev: run ## Alias for 'make run'

# =============================================================================
# Build
# =============================================================================

build: server-build agent-build dashboard-build ## Build all components
	@echo "$(GREEN)All components built successfully!$(RESET)"

server-build: ## Build the Go server
	@echo "$(YELLOW)Building server...$(RESET)"
	cd server && CGO_ENABLED=1 go build -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)" -o autostrike-server ./cmd/autostrike

server-build-quick: ## Build server (quick, for run target)
	@cd server && go build -o autostrike-server ./cmd/autostrike 2>/dev/null || \
		(echo "$(YELLOW)Building server...$(RESET)" && go build -o autostrike-server ./cmd/autostrike)

agent-build: ## Build the Rust agent
	@echo "$(YELLOW)Building agent...$(RESET)"
	cd agent && cargo build --release

agent-build-quick: ## Build agent (quick)
	@cd agent && cargo build --release 2>/dev/null

dashboard-build: ## Build the React dashboard
	@echo "$(YELLOW)Building dashboard...$(RESET)"
	cd dashboard && npm run build

dashboard-build-quick: ## Build dashboard if not exists
	@if [ ! -f dashboard/dist/index.html ]; then \
		echo "$(YELLOW)Building dashboard...$(RESET)"; \
		cd dashboard && npm install --silent && npm run build; \
	fi

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
	cd dashboard && npm run test -- --run

lint: ## Run linters
	@echo "$(YELLOW)Linting...$(RESET)"
	cd server && go vet ./...
	cd agent && cargo clippy
	cd dashboard && npm run lint

# =============================================================================
# Agent
# =============================================================================

agent: agent-build-quick ## Build and run the agent
	@echo "$(YELLOW)Starting agent...$(RESET)"
	@cd agent && ./target/release/autostrike-agent \
		--server http://localhost:8443 \
		--paw agent-$$(hostname)-$$$$

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

install: ## Install all dependencies
	@echo "$(YELLOW)Installing dependencies...$(RESET)"
	cd server && go mod download
	cd agent && cargo fetch
	cd dashboard && npm install

clean: ## Clean build artifacts
	@echo "$(YELLOW)Cleaning...$(RESET)"
	rm -rf dist/
	rm -rf server/autostrike-server
	rm -rf agent/target/
	rm -rf dashboard/dist/
	rm -rf dashboard/node_modules/

setup: install certs ## Initial project setup
	@echo "$(GREEN)Project setup complete!$(RESET)"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Run 'make run' to start AutoStrike"
	@echo "  2. Open http://localhost:8443"

stop: ## Stop all running services
	@echo "$(YELLOW)Stopping services...$(RESET)"
	-@pkill -f "autostrike-server" 2>/dev/null
	-@pkill -f "autostrike-agent" 2>/dev/null
	@echo "$(GREEN)✓ Services stopped$(RESET)"

logs: ## Show server logs
	@echo "$(CYAN)=== Server Logs ===$(RESET)"
	@tail -30 /tmp/autostrike-server.log 2>/dev/null || echo "No server logs"

dist: ## Create distribution package
	mkdir -p dist
	$(MAKE) build
	cp server/autostrike-server dist/
	cp agent/target/release/autostrike-agent dist/
	cp -r dashboard/dist dist/dashboard
	tar -czvf autostrike-$(VERSION).tar.gz dist/
