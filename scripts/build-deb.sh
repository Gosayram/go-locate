#!/bin/bash

# DEB Package Builder for glocate
# This script creates a DEB package using standard Debian tools and pre-compiled binaries

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Package information
PACKAGE_NAME="glocate"
VERSION="${VERSION:-$(cat "$PROJECT_ROOT/.release-version" 2>/dev/null || echo "0.1.2")}"
MAINTAINER="abdurakhman.rakhmankulov@gmail.com"
DESCRIPTION="Modern file search tool to replace locate"
HOMEPAGE="https://github.com/Gosayram/go-locate"

# Build information
BUILD_DIR="$PROJECT_ROOT/build-deb"
PACKAGE_DIR="$BUILD_DIR/$PACKAGE_NAME"
DEBIAN_DIR="$PACKAGE_DIR/DEBIAN"

# Architecture detection
ARCH=$(uname -m)
case $ARCH in
    x86_64)
        DEB_ARCH="amd64"
        BINARY_SUFFIX="linux-amd64"
        ;;
    aarch64|arm64)
        DEB_ARCH="arm64"
        BINARY_SUFFIX="linux-arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log() {
    echo -e "${BLUE}[INFO]${NC} $*"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $*"
}

error() {
    echo -e "${RED}[ERROR]${NC} $*" >&2
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $*"
}

# Auto-install dependencies
auto_install_dependencies() {
    log "Checking and installing DEB build dependencies..."

    local missing_deps=()

    if ! command -v fakeroot >/dev/null 2>&1; then
        missing_deps+=("fakeroot")
    fi

    if ! command -v dpkg-deb >/dev/null 2>&1; then
        missing_deps+=("dpkg-dev")
    fi

    if [ ${#missing_deps[@]} -ne 0 ]; then
        warn "Missing dependencies: ${missing_deps[*]}"

        # Try to auto-install if in CI or if user confirms
        if [ "${CI:-false}" = "true" ] || [ "${GITHUB_ACTIONS:-false}" = "true" ]; then
            log "CI environment detected, attempting auto-installation..."
            if command -v apt-get >/dev/null 2>&1; then
                log "Installing dependencies with apt-get..."
                sudo apt-get update -qq && sudo apt-get install -y "${missing_deps[@]}"
            elif command -v dnf >/dev/null 2>&1; then
                log "Installing dependencies with dnf..."
                sudo dnf install -y dpkg-dev fakeroot
            elif command -v yum >/dev/null 2>&1; then
                log "Installing dependencies with yum..."
                sudo yum install -y dpkg-dev fakeroot
            else
                error "No supported package manager found for auto-installation"
                error "Please install manually: ${missing_deps[*]}"
                return 1
            fi
        else
            error "Missing dependencies: ${missing_deps[*]}"
            error "Run the following to install them:"
            error "  make install-deb-tools"
            error "Or manually: sudo apt-get install ${missing_deps[*]}"
            return 1
        fi
    fi

    # Verify installation
    if ! command -v fakeroot >/dev/null 2>&1 || ! command -v dpkg-deb >/dev/null 2>&1; then
        error "Dependencies still missing after installation attempt"
        return 1
    fi

    success "All dependencies are available"
}

# Check dependencies (legacy function, calls auto_install_dependencies)
check_dependencies() {
    auto_install_dependencies
}

# Clean previous build
cleanup() {
    log "Cleaning previous build..."
    rm -rf "$BUILD_DIR"
    rm -f "$PROJECT_ROOT"/*.deb
}

# Create directory structure
create_structure() {
    log "Creating package directory structure..."

    mkdir -p "$DEBIAN_DIR"
    mkdir -p "$PACKAGE_DIR/usr/bin"
    mkdir -p "$PACKAGE_DIR/etc/$PACKAGE_NAME"
    mkdir -p "$PACKAGE_DIR/usr/share/doc/$PACKAGE_NAME"
}

# Copy files
copy_files() {
    log "Copying files for architecture: $DEB_ARCH"

    # Copy binary from pre-built cross-platform binaries
    local binary_path="$PROJECT_ROOT/bin/$PACKAGE_NAME-$BINARY_SUFFIX"
    if [ -f "$binary_path" ]; then
        cp "$binary_path" "$PACKAGE_DIR/usr/bin/$PACKAGE_NAME"
        chmod 755 "$PACKAGE_DIR/usr/bin/$PACKAGE_NAME"
        log "Copied binary from $binary_path"
    else
        error "Pre-compiled binary not found at $binary_path"
        error "Please run 'make build-cross' first to build cross-platform binaries"
        return 1
    fi

    # Copy configuration
    if [ -f "$PROJECT_ROOT/example.glocate.toml" ]; then
        cp "$PROJECT_ROOT/example.glocate.toml" "$PACKAGE_DIR/etc/$PACKAGE_NAME/glocate.toml"
        chmod 644 "$PACKAGE_DIR/etc/$PACKAGE_NAME/glocate.toml"
    fi

    # Copy documentation
    for doc in README.md CHANGELOG.md LICENSE; do
        if [ -f "$PROJECT_ROOT/$doc" ]; then
            cp "$PROJECT_ROOT/$doc" "$PACKAGE_DIR/usr/share/doc/$PACKAGE_NAME/"
        fi
    done
}

# Calculate installed size
calculate_size() {
    du -sk "$PACKAGE_DIR" | cut -f1
}

# Create control file
create_control() {
    log "Creating control file..."

    local installed_size
    installed_size=$(calculate_size)

    cat > "$DEBIAN_DIR/control" << EOF
Package: $PACKAGE_NAME
Version: $VERSION
Section: utils
Priority: optional
Architecture: $DEB_ARCH
Maintainer: $MAINTAINER
Installed-Size: $installed_size
Depends: libc6
Provides: locate
Conflicts: mlocate, findutils-locate
Homepage: $HOMEPAGE
Description: $DESCRIPTION
 glocate is a modern, fast file search tool designed to replace the traditional
 locate command. It provides enhanced search capabilities with better performance
 and more intuitive command-line interface.
 .
 Key features:
  - Fast file indexing and searching
  - Modern command-line interface
  - Configuration file support
  - Backward compatibility with locate
EOF
}

# Create maintainer scripts
create_scripts() {
    log "Creating maintainer scripts..."

    # Post-installation script
    cat > "$DEBIAN_DIR/postinst" << 'EOF'
#!/bin/bash
set -e

case "$1" in
    configure)
        # Create symlink for backward compatibility
        if [ ! -e /usr/bin/locate ]; then
            ln -sf glocate /usr/bin/locate
        fi

        # Update locate database
        echo "Creating initial locate database..."
        /usr/bin/glocate --update-db >/dev/null 2>&1 || true
        ;;
esac

exit 0
EOF
    chmod 755 "$DEBIAN_DIR/postinst"

    # Pre-removal script
    cat > "$DEBIAN_DIR/prerm" << 'EOF'
#!/bin/bash
set -e

case "$1" in
    remove|deconfigure)
        # Remove symlink if it points to glocate
        if [ -L /usr/bin/locate ] && [ "$(readlink /usr/bin/locate)" = "glocate" ]; then
            rm -f /usr/bin/locate
        fi
        ;;
esac

exit 0
EOF
    chmod 755 "$DEBIAN_DIR/prerm"
}

# Create changelog
create_changelog() {
    log "Creating changelog..."

    local changelog_dir="$PACKAGE_DIR/usr/share/doc/$PACKAGE_NAME"

    cat > "$changelog_dir/changelog.Debian" << EOF
$PACKAGE_NAME ($VERSION) unstable; urgency=medium

  * Updated to use pre-compiled binaries following Go packaging best practices
  * Added proper architecture support for amd64 and arm64
  * Improved package metadata and dependencies
  * Added backward compatibility symlink
  * Added Provides: locate and proper Conflicts

 -- $MAINTAINER  $(date -R)

$PACKAGE_NAME (0.1.2-1) unstable; urgency=low

  * Initial package for glocate
  * Modern file search tool to replace locate
  * Configuration file support
  * Backward compatibility with locate command

 -- $MAINTAINER  Mon, 30 Dec 2024 12:00:00 +0000
EOF

    gzip -9 "$changelog_dir/changelog.Debian"
}

# Build package
build_package() {
    log "Building DEB package for $DEB_ARCH..."

    cd "$PROJECT_ROOT"

    # Ensure packages directory exists
    mkdir -p packages

    local deb_filename="packages/${PACKAGE_NAME}_${VERSION}_${DEB_ARCH}.deb"

    if command -v fakeroot >/dev/null 2>&1; then
        fakeroot dpkg-deb --build "$PACKAGE_DIR" "$deb_filename"
    else
        warn "fakeroot not available, building without it (may cause permission issues)"
        dpkg-deb --build "$PACKAGE_DIR" "$deb_filename"
    fi
}

# Verify package
verify_package() {
    log "Verifying package..."

    local deb_file="packages/${PACKAGE_NAME}_${VERSION}_${DEB_ARCH}.deb"
    local full_path="$PROJECT_ROOT/$deb_file"

    if [ -f "$full_path" ]; then
        success "Package created successfully: $deb_file"

        # Show package info
        log "Package information:"
        dpkg-deb --info "$full_path"

        # Show package contents
        log "Package contents:"
        dpkg-deb --contents "$full_path"

        # Run lintian if available
        if command -v lintian >/dev/null 2>&1; then
            log "Running lintian checks..."
            lintian "$full_path" || warn "Lintian found some issues (non-critical)"
        else
            warn "lintian not available, skipping package validation"
        fi
    else
        error "Package creation failed"
        return 1
    fi
}

# Main function
main() {
    log "Starting DEB package build for $PACKAGE_NAME $VERSION ($DEB_ARCH)"

    if ! check_dependencies; then
        exit 1
    fi

    cleanup
    create_structure
    copy_files
    create_control
    create_scripts
    create_changelog
    build_package
    verify_package

    success "DEB package build completed successfully!"
}

# Run main function
main "$@"
