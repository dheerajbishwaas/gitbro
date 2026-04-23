#!/bin/bash

echo ""
echo "========================================================"
echo ".                                                      ."
echo ".         GitBro Installation (macOS)                 ."
echo ".                                                      ."
echo "========================================================"
echo ""

INSTALL_DIR="/usr/local/bin"

# Detect architecture
ARCH=$(uname -m)

if [ "$ARCH" = "arm64" ]; then
    SOURCE_FILE="$(dirname "$0")/gitbro-macos-arm64"
    echo "[*] Detected: Apple Silicon (M1/M2/M3)"
elif [ "$ARCH" = "x86_64" ]; then
    SOURCE_FILE="$(dirname "$0")/gitbro-macos-intel"
    echo "[*] Detected: Intel Mac"
else
    echo "ERROR: Unknown architecture: $ARCH"
    exit 1
fi

if [ ! -f "$SOURCE_FILE" ]; then
    echo "ERROR: Binary file not found: $SOURCE_FILE"
    exit 1
fi

# Copy executable (may need sudo)
echo "[*] Copying gitbro to $INSTALL_DIR..."
cp -f "$SOURCE_FILE" "$INSTALL_DIR/gitbro" 2>/dev/null

if [ $? -ne 0 ]; then
    echo "[*] Trying with sudo (you may be prompted for password)..."
    sudo cp -f "$SOURCE_FILE" "$INSTALL_DIR/gitbro"
    
    if [ $? -ne 0 ]; then
        echo ""
        echo "ERROR: Could not copy file to $INSTALL_DIR"
        exit 1
    fi
fi

# Make it executable
chmod +x "$INSTALL_DIR/gitbro" 2>/dev/null || sudo chmod +x "$INSTALL_DIR/gitbro"

echo ""
echo "========================================================"
echo ".         Installation Complete!                      ."
echo "========================================================"
echo ""
echo "GitBro installed to: $INSTALL_DIR/gitbro"
echo ""
echo "Next Steps:"
echo "1. Close and reopen your terminal"
echo "2. Type: gitbro"
echo ""
