#!/bin/bash
# Generate coverage reports for AutoStrike project

set -e

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
COVERAGE_DIR="$ROOT_DIR/coverage"
BOLD='\033[1m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

mkdir -p "$COVERAGE_DIR"

echo -e "${BOLD}${BLUE}=== AutoStrike Coverage Report ===${NC}\n"

# Go Server
echo -e "${BOLD}[1/3] Go Server${NC}"
cd "$ROOT_DIR/server"
go test -coverprofile="$COVERAGE_DIR/go.out" ./...
go tool cover -func="$COVERAGE_DIR/go.out" | tail -1
echo ""

# React Dashboard
echo -e "${BOLD}[2/3] React Dashboard${NC}"
cd "$ROOT_DIR/dashboard"
npm test -- --coverage --run 2>&1 | grep -A 10 "Coverage summary" || npm test -- --coverage --run
echo ""

# Rust Agent (requires cargo-tarpaulin: cargo install cargo-tarpaulin)
echo -e "${BOLD}[3/3] Rust Agent${NC}"
cd "$ROOT_DIR/agent"
if command -v cargo-tarpaulin &> /dev/null; then
    cargo tarpaulin --out Html --output-dir "$COVERAGE_DIR" 2>&1 | tail -5
else
    echo "cargo-tarpaulin not installed. Install with: cargo install cargo-tarpaulin"
    cargo test
fi
echo ""

echo -e "${GREEN}${BOLD}Coverage reports saved in: $COVERAGE_DIR/${NC}"
