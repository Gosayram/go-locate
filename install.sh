#!/bin/bash

# go-locate installation script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="glocate"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="$HOME"
CONFIG_FILE="$CONFIG_DIR/.glocate.toml"

echo -e "${GREEN}Installing go-locate...${NC}"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed. Please install Go first.${NC}"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.21"

if ! printf '%s\n%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V -C; then
    echo -e "${YELLOW}Warning: Go version $GO_VERSION detected. Recommended version is $REQUIRED_VERSION or higher.${NC}"
fi

# Build the binary
echo "Building $BINARY_NAME..."
if ! make build; then
    echo -e "${RED}Error: Failed to build $BINARY_NAME${NC}"
    exit 1
fi

# Install binary
echo "Installing $BINARY_NAME to $INSTALL_DIR..."
if ! sudo cp "bin/$BINARY_NAME" "$INSTALL_DIR/"; then
    echo -e "${RED}Error: Failed to install $BINARY_NAME to $INSTALL_DIR${NC}"
    echo "You may need to run this script with sudo or check your permissions."
    exit 1
fi

# Make executable
sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"

# Install example config if it doesn't exist
if [ ! -f "$CONFIG_FILE" ]; then
    echo "Installing example configuration to $CONFIG_FILE..."
    cp example.glocate.toml "$CONFIG_FILE"
    echo -e "${GREEN}Example configuration installed. Edit $CONFIG_FILE to customize.${NC}"
else
    echo -e "${YELLOW}Configuration file already exists at $CONFIG_FILE${NC}"
fi

# Verify installation
if command -v "$BINARY_NAME" &> /dev/null; then
    VERSION=$("$BINARY_NAME" --version 2>&1 | head -n1)
    echo -e "${GREEN}✓ Installation successful!${NC}"
    echo -e "${GREEN}✓ $VERSION${NC}"
    echo ""
    echo "Usage examples:"
    echo "  $BINARY_NAME \"*.go\"                    # Find all Go files"
    echo "  $BINARY_NAME --advanced myfile          # Fuzzy search"
    echo "  $BINARY_NAME --ext go,rs main           # Search by extension"
    echo "  $BINARY_NAME --format detailed \"*.md\"   # Detailed output"
    echo ""
    echo "Run '$BINARY_NAME --help' for more options."
else
    echo -e "${RED}Error: Installation verification failed${NC}"
    exit 1
fi 