#!/bin/bash

echo ""
echo "========================================================"
echo ".                                                      ."
echo ".         GitBro Installation (Linux)                 ."
echo ".                                                      ."
echo "========================================================"
echo ""

INSTALL_DIR="$HOME/.local/bin"

# Create directory if it doesn't exist
mkdir -p "$INSTALL_DIR"

# Copy executable
echo "[*] Copying gitbro to $INSTALL_DIR..."
cp -f "$(dirname "$0")/gitbro-linux" "$INSTALL_DIR/gitbro"

if [ $? -ne 0 ]; then
    echo ""
    echo "ERROR: Could not copy file"
    echo "Try: sudo bash $0"
    exit 1
fi

# Make it executable
chmod +x "$INSTALL_DIR/gitbro"

# Check if ~/.local/bin is in PATH
if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
    echo ""
    echo "⚠ WARNING: $HOME/.local/bin is not in your PATH"
    echo "Add this line to your ~/.bashrc or ~/.zshrc:"
    echo ""
    echo "export PATH=\"\$HOME/.local/bin:\$PATH\""
    echo ""
    echo "Then run: source ~/.bashrc"
    echo ""
else
    echo "[OK] $HOME/.local/bin is already in PATH"
fi

echo ""
echo "========================================================"
echo ".         Installation Complete!                      ."
echo "========================================================"
echo ""
echo "GitBro installed to: $INSTALL_DIR/gitbro"
echo ""
echo "Next Steps:"
echo "1. If you added PATH, run: source ~/.bashrc"
echo "2. Open a new terminal"
echo "3. Type: gitbro"
echo ""
