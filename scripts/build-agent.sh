#!/bin/bash
set -e

# AutoStrike Agent Build Script
# Cross-compiles the agent for multiple platforms

AGENT_DIR="./agent"
OUTPUT_DIR="./dist/agents"
VERSION="${1:-dev}"

# Supported targets
TARGETS=(
    "x86_64-unknown-linux-gnu"
    "x86_64-unknown-linux-musl"
    "x86_64-pc-windows-gnu"
    "x86_64-apple-darwin"
    "aarch64-unknown-linux-gnu"
    "aarch64-apple-darwin"
)

mkdir -p "$OUTPUT_DIR"

echo "=== AutoStrike Agent Builder ==="
echo "Version: $VERSION"
echo ""

cd "$AGENT_DIR"

# Install cross if not present
if ! command -v cross &> /dev/null; then
    echo "Installing cross..."
    cargo install cross
fi

for target in "${TARGETS[@]}"; do
    echo "Building for $target..."

    # Determine output name
    case $target in
        *windows*)
            ext=".exe"
            os="windows"
            ;;
        *darwin*)
            ext=""
            os="macos"
            ;;
        *)
            ext=""
            os="linux"
            ;;
    esac

    # Determine architecture
    case $target in
        x86_64*)
            arch="amd64"
            ;;
        aarch64*)
            arch="arm64"
            ;;
        *)
            arch="unknown"
            ;;
    esac

    output_name="autostrike-agent-${VERSION}-${os}-${arch}${ext}"

    # Build with cross
    if cross build --release --target "$target" 2>/dev/null; then
        cp "target/$target/release/autostrike-agent${ext}" "../$OUTPUT_DIR/$output_name"
        echo "  -> Built: $output_name"
    else
        echo "  -> Skipped: $target (build failed)"
    fi
done

cd ..

echo ""
echo "=== Build Complete ==="
echo "Agents available in: $OUTPUT_DIR"
ls -la "$OUTPUT_DIR"
