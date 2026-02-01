#!/bin/bash
set -e

# AutoStrike Atomic Red Team Importer
# Downloads and imports Atomic Red Team techniques

ATOMIC_REPO="https://github.com/redcanaryco/atomic-red-team.git"
TEMP_DIR="/tmp/atomic-red-team"
OUTPUT_DIR="${1:-./data/techniques}"

echo "=== AutoStrike Atomic Red Team Importer ==="
echo ""

# Clone or update repository
if [ -d "$TEMP_DIR" ]; then
    echo "Updating existing Atomic Red Team repository..."
    cd "$TEMP_DIR"
    git pull
else
    echo "Cloning Atomic Red Team repository..."
    git clone --depth 1 "$ATOMIC_REPO" "$TEMP_DIR"
fi

mkdir -p "$OUTPUT_DIR"

echo ""
echo "Processing techniques..."

# Count techniques
total=$(find "$TEMP_DIR/atomics" -name "*.yaml" -type f | wc -l)
current=0

# Process each technique YAML file
find "$TEMP_DIR/atomics" -name "*.yaml" -type f | while read -r yaml_file; do
    current=$((current + 1))
    technique_id=$(basename "$(dirname "$yaml_file")")

    # Skip non-technique directories
    if [[ ! "$technique_id" =~ ^T[0-9]+ ]]; then
        continue
    fi

    # Copy YAML file
    cp "$yaml_file" "$OUTPUT_DIR/${technique_id}.yaml"

    # Progress indicator
    printf "\r  Processed: %d / %d" "$current" "$total"
done

echo ""
echo ""
echo "=== Import Complete ==="
echo "Techniques imported to: $OUTPUT_DIR"
echo "Total files: $(ls -1 "$OUTPUT_DIR" | wc -l)"
