#!/bin/bash
set -e

echo "Building GitBro for macOS..."

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$REPO_ROOT"

# Build for Intel Macs.
echo "[1/2] Building macOS executable (Intel/amd64)..."
GOOS=darwin GOARCH=amd64 go build -buildvcs=false -o "$SCRIPT_DIR/gitbro-macos-intel" .

# Build for Apple Silicon (M1/M2/M3).
echo "[2/2] Building macOS executable (Apple Silicon/arm64)..."
GOOS=darwin GOARCH=arm64 go build -buildvcs=false -o "$SCRIPT_DIR/gitbro-macos-arm64" .

echo ""
echo "Build complete!"
echo "Files:"
echo "  - macos/gitbro-macos-intel  (Intel Macs)"
echo "  - macos/gitbro-macos-arm64  (Apple Silicon: M1/M2/M3)"