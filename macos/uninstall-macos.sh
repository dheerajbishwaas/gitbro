#!/bin/bash

echo ""
echo "========================================================"
echo ".                                                      ."
echo ".         GitBro Uninstaller (macOS)                  ."
echo ".                                                      ."
echo "========================================================"
echo ""

INSTALL_DIR="/usr/local/bin"

if [ ! -f "$INSTALL_DIR/gitbro" ]; then
    echo "ERROR: GitBro is not installed"
    exit 1
fi

read -p "Are you sure? (Y/N): " confirm

if [ "$confirm" != "Y" ] && [ "$confirm" != "y" ]; then
    echo ""
    echo "Uninstallation cancelled."
    exit 0
fi

echo ""
echo "[*] Removing gitbro..."
rm -f "$INSTALL_DIR/gitbro" 2>/dev/null || sudo rm -f "$INSTALL_DIR/gitbro"

if [ $? -eq 0 ]; then
    echo ""
    echo "========================================================"
    echo ".      Uninstallation Complete!                        ."
    echo "========================================================"
    echo ""
    echo "GitBro has been successfully removed."
    echo ""
else
    echo ""
    echo "ERROR: Could not remove GitBro"
    exit 1
fi
