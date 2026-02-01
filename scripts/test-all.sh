#!/bin/bash
# Run all tests with coverage for AutoStrike project

set -e

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
BOLD='\033[1m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BOLD}${BLUE}=== AutoStrike Test Suite ===${NC}\n"

# Go Server
echo -e "${BOLD}[1/3] Go Server${NC}"
cd "$ROOT_DIR/server"
go test -cover ./...
echo ""

# React Dashboard
echo -e "${BOLD}[2/3] React Dashboard${NC}"
cd "$ROOT_DIR/dashboard"
npm test -- --run
echo ""

# Rust Agent
echo -e "${BOLD}[3/3] Rust Agent${NC}"
cd "$ROOT_DIR/agent"
cargo test
echo ""

echo -e "${GREEN}${BOLD}All tests passed!${NC}"
