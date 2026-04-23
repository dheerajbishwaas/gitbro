#!/bin/bash
set -e

echo "Building GitBro for Linux..."

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Build for Linux and place the binary next to the Linux installer.
echo "[1/1] Building Linux executable (amd64)..."
cd "$REPO_ROOT"
GOOS=linux GOARCH=amd64 go build -buildvcs=false -o "$SCRIPT_DIR/gitbro-linux" .

echo "[OK] Linux build successful: linux/gitbro-linux"
echo ""
echo "Build complete!"
echo "File: linux/gitbro-linux"