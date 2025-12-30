#!/bin/bash

# Initiative installer script
# Downloads and installs the latest release of initiative

set -e

REPO="joelzwarrington/initiative"
BINARY_NAME="initiative"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $OS in
  linux)
    OS="linux"
    ;;
  darwin)
    OS="darwin"
    ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

case $ARCH in
  x86_64)
    ARCH="x86_64"
    ;;
  aarch64 | arm64)
    ARCH="arm64"
    ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Get the latest release version
echo "Fetching latest release..."
LATEST_VERSION=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep -o '"tag_name": "[^"]*"' | cut -d'"' -f4)

if [ -z "$LATEST_VERSION" ]; then
    echo "Failed to fetch latest release version"
    exit 1
fi

echo "Latest version: $LATEST_VERSION"

# Build download URL
if [ "$OS" = "darwin" ]; then
    OS_NAME="Darwin"
elif [ "$OS" = "linux" ]; then
    OS_NAME="Linux"
fi

FILENAME="initiative_${OS_NAME}_${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/${FILENAME}"

# Create temp directory
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

echo "Downloading $FILENAME..."
curl -L -o "$TMP_DIR/$FILENAME" "$DOWNLOAD_URL"

# Extract and install
echo "Installing initiative..."
cd "$TMP_DIR"
tar -xzf "$FILENAME"

# Find install directory
INSTALL_DIR="/usr/local/bin"
if [ ! -w "$INSTALL_DIR" ]; then
    echo "Installing to ~/bin (add to PATH if not already there)"
    INSTALL_DIR="$HOME/bin"
    mkdir -p "$INSTALL_DIR"
fi

# Move binary
mv "$BINARY_NAME" "$INSTALL_DIR/"
chmod +x "$INSTALL_DIR/$BINARY_NAME"

echo "✓ Initiative installed successfully to $INSTALL_DIR/$BINARY_NAME"