# Packaging Documentation for go-locate

This document describes how to build RPM and DEB packages for the go-locate project using modern Go packaging best practices.

## Overview

The go-locate project supports creating distribution packages using pre-compiled binaries:
- **Binary tarballs** for multiple architectures (amd64, arm64)
- **RPM packages** for Red Hat, Fedora, CentOS, and other RPM-based systems
- **DEB packages** for Debian, Ubuntu, and other DEB-based systems
- **Source tarballs** for manual compilation and distribution

## Modern Go Packaging Approach

Following Go packaging best practices, this project:
- **Uses pre-compiled binaries** instead of compiling during package build
- **Supports multiple architectures** (amd64, arm64) with proper detection
- **Follows distribution guidelines** for Go applications
- **Maintains static binary benefits** with no runtime dependencies
- **Provides proper package metadata** and backward compatibility

## Prerequisites

### For RPM Building

On Red Hat/Fedora/CentOS systems:
```bash
sudo dnf install rpm-build rpmdevtools
# or on older systems:
sudo yum install rpm-build rpmdevtools
```

On Debian/Ubuntu systems:
```bash
sudo apt-get install rpm rpmbuild
```

### For DEB Building

On Debian/Ubuntu systems:
```bash
sudo apt-get install dpkg-dev fakeroot lintian
```

### For Cross-Platform Building

On macOS (for development/testing):
```bash
# RPM building is not supported on macOS
# DEB building requires Linux tools and will fail gracefully
```

## Building Packages

### Build All Packages

To build all package types:
```bash
make package-all
```

This will create:
- Binary tarballs for amd64 and arm64 architectures
- Source tarball (`glocate-X.Y.Z-src.tar.gz`)
- RPM packages (`.rpm` files)
- DEB packages (`.deb` files)

### Build Individual Package Types

#### Binary Tarballs (Recommended for Distribution)
```bash
make package-binaries
```

Creates:
- `glocate-X.Y.Z-linux-amd64.tar.gz`
- `glocate-X.Y.Z-linux-arm64.tar.gz`

#### RPM Package
```bash
make package-rpm
```

#### DEB Package
```bash
make package-deb
```

#### Source Tarball
```bash
make package-tarball
```

## Package Contents

### Binary Files
- `/usr/bin/glocate` - Main executable (pre-compiled for target architecture)
- `/usr/bin/locate` - Symlink for backward compatibility (created post-install)

### Configuration Files
- `/etc/glocate/glocate.toml` - Default configuration file

### Documentation
- `/usr/share/doc/glocate/README.md` - Project documentation
- `/usr/share/doc/glocate/CHANGELOG.md` - Version history
- `/usr/share/doc/glocate/LICENSE` - Software license
- `/usr/share/doc/glocate/glocate.toml.example` - Configuration example

## Installation

### Binary Tarball Installation
```bash
# Extract and install manually
tar -xzf glocate-X.Y.Z-linux-amd64.tar.gz
sudo cp glocate-X.Y.Z-linux-amd64/glocate /usr/local/bin/
sudo ln -sf /usr/local/bin/glocate /usr/local/bin/locate
```

### RPM Installation
```bash
# Install the package
sudo rpm -ivh glocate-X.Y.Z-1.x86_64.rpm

# Or using dnf/yum
sudo dnf install ./glocate-X.Y.Z-1.x86_64.rpm
```

### DEB Installation
```bash
# Install the package
sudo dpkg -i glocate_X.Y.Z_amd64.deb

# Install dependencies if needed
sudo apt-get install -f
```

## Removal

### RPM Removal
```bash
sudo rpm -e glocate
```

### DEB Removal
```bash
sudo dpkg -r glocate
```

## Package Metadata

### RPM Package Information
- **Name**: glocate
- **Summary**: Modern file search tool to replace locate
- **License**: MIT
- **Group**: Applications/System
- **Architecture**: x86_64, aarch64 (multi-arch support)
- **Requires**: No special dependencies (static binary)
- **Provides**: locate (for compatibility)
- **Obsoletes**: mlocate, findutils-locate

### DEB Package Information
- **Package**: glocate
- **Section**: utils
- **Priority**: optional
- **Architecture**: amd64, arm64 (automatic detection)
- **Depends**: libc6 (minimal dependencies)
- **Provides**: locate
- **Conflicts**: mlocate, findutils-locate

## Architecture Support

### Supported Architectures
- **amd64/x86_64**: Intel/AMD 64-bit systems
- **arm64/aarch64**: ARM 64-bit systems (Apple Silicon, ARM servers)

### Architecture Detection
- **RPM**: Supports both x86_64 and aarch64 in spec file
- **DEB**: Automatically detects architecture during build
- **Binary tarballs**: Separate archives for each architecture

## Troubleshooting

### Common Issues

#### RPM Build Fails
```bash
# Check if rpm-build is installed
rpm -qa | grep rpm-build

# Setup RPM build environment
rpmdev-setuptree

# Ensure cross-platform binaries are built
make build-cross
```

#### DEB Build Fails
```bash
# Check if required tools are installed
dpkg --version
fakeroot --version

# Install missing dependencies
sudo apt-get install dpkg-dev fakeroot

# Ensure cross-platform binaries are built
make build-cross
```

#### Missing Pre-compiled Binaries
```bash
# Build cross-platform binaries first
make build-cross

# Verify binaries exist
ls -la bin/glocate-*

# Expected files:
# bin/glocate-linux-amd64
# bin/glocate-linux-arm64
# bin/glocate-darwin-amd64
# bin/glocate-darwin-arm64
```

#### Permission Issues
```bash
# Ensure proper permissions
chmod +x scripts/build-deb.sh

# Use fakeroot for DEB building
fakeroot dpkg-deb --build package_dir
```

### Package Validation

#### RPM Validation
```bash
# Check package info
rpm -qip glocate-X.Y.Z-1.x86_64.rpm

# List package contents
rpm -qlp glocate-X.Y.Z-1.x86_64.rpm

# Verify package
rpm -K glocate-X.Y.Z-1.x86_64.rpm
```

#### DEB Validation
```bash
# Check package info
dpkg-deb --info glocate_X.Y.Z_amd64.deb

# List package contents
dpkg-deb --contents glocate_X.Y.Z_amd64.deb

# Run lintian checks
lintian glocate_X.Y.Z_amd64.deb
```

## CI/CD Integration

### Automated Tool Installation

The packaging system now supports automatic installation of build tools:

```bash
# Auto-detect OS and install appropriate tools
make install-rpm-tools    # Install RPM build tools
make install-deb-tools    # Install DEB build tools
make detect-os            # Show detected OS information

# CI-optimized packaging (auto-installs tools in CI environments)
make package-ci           # Build binary tarballs + source (CI-friendly)
make package-ci-setup     # Setup CI environment with tools
```

### Supported Operating Systems

#### RPM Building Support:
- **Fedora/RHEL/CentOS/Rocky/AlmaLinux**: `dnf install` or `yum install`
- **Ubuntu/Debian**: `apt-get install rpm` (cross-platform)
- **openSUSE/SLES**: `zypper install`
- **Arch/Manjaro**: `pacman -S rpm-tools`

#### DEB Building Support:
- **Ubuntu/Debian**: `apt-get install dpkg-dev fakeroot lintian` (native)
- **Fedora/RHEL/CentOS**: `dnf install dpkg-dev fakeroot` (cross-platform)
- **openSUSE/SLES**: `zypper install dpkg fakeroot`
- **Arch/Manjaro**: `pacman -S dpkg fakeroot`

### GitHub Actions Example

#### Complete Multi-Platform Packaging

```yaml
name: Build and Package

on:
  push:
    tags:
      - 'v*'
  pull_request:
    branches: [ main, develop ]

jobs:
  # Build binary tarballs (works on any OS)
  build-binaries:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.22

      - name: Build cross-platform binaries and tarballs
        run: make package-ci

      - name: Upload binary artifacts
        uses: actions/upload-artifact@v3
        with:
          name: binary-packages
          path: packages/*.tar.gz

  # Build RPM packages on RPM-based system
  build-rpm:
    runs-on: ubuntu-latest
    container: fedora:latest
    steps:
      - uses: actions/checkout@v4

      - name: Install Go
        run: dnf install -y golang make git

      - name: Build RPM packages
        run: make package-rpm

      - name: Upload RPM artifacts
        uses: actions/upload-artifact@v3
        with:
          name: rpm-packages
          path: packages/*.rpm

  # Build DEB packages on DEB-based system
  build-deb:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.22

      - name: Build DEB packages
        run: make package-deb

      - name: Upload DEB artifacts
        uses: actions/upload-artifact@v3
        with:
          name: deb-packages
          path: packages/*.deb

  # Create release with all packages
  release:
    if: startsWith(github.ref, 'refs/tags/v')
    needs: [build-binaries, build-rpm, build-deb]
    runs-on: ubuntu-latest
    steps:
      - name: Download all artifacts
        uses: actions/download-artifact@v3

      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          files: |
            binary-packages/*
            rpm-packages/*
            deb-packages/*
          generate_release_notes: true
```

#### CI-Optimized Workflow (Single Job)

```yaml
name: Quick Package Build

on:
  push:
    branches: [ main ]

jobs:
  package:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.22

      - name: Build packages with auto-tool installation
        run: make package-ci

      - name: Show package information
        run: |
          echo "Built packages:"
          ls -la packages/

      - name: Upload packages
        uses: actions/upload-artifact@v3
        with:
          name: packages
          path: packages/
```

### Docker-based Building

For consistent cross-platform building:

```yaml
name: Docker Package Build

on:
  workflow_dispatch:

jobs:
  build-packages:
    strategy:
      matrix:
        os:
          - { name: "ubuntu", image: "ubuntu:22.04", pkg: "deb" }
          - { name: "fedora", image: "fedora:latest", pkg: "rpm" }
          - { name: "debian", image: "debian:bookworm", pkg: "deb" }
          - { name: "centos", image: "quay.io/centos/centos:stream9", pkg: "rpm" }

    runs-on: ubuntu-latest
    container: ${{ matrix.os.image }}

    steps:
      - name: Install Git and basic tools
        run: |
          if command -v apt-get >/dev/null 2>&1; then
            apt-get update && apt-get install -y git make curl
          elif command -v dnf >/dev/null 2>&1; then
            dnf install -y git make curl golang
          elif command -v yum >/dev/null 2>&1; then
            yum install -y git make curl golang
          fi

      - uses: actions/checkout@v4

      - name: Set up Go (if not installed)
        if: matrix.os.name == 'ubuntu' || matrix.os.name == 'debian'
        uses: actions/setup-go@v4
        with:
          go-version: 1.22

      - name: Build package
        run: make package-${{ matrix.os.pkg }}

      - name: Upload artifacts
        uses: actions/upload-artifact@v3
        with:
          name: ${{ matrix.os.name }}-packages
          path: packages/
```

### Local Development

For local development and testing:

```bash
# Check what OS you're on
make detect-os

# Install tools for your platform
make install-rpm-tools   # If you want to build RPM
make install-deb-tools   # If you want to build DEB

# Build packages
make package-binaries    # Always works (creates tarballs)
make package-rpm         # Works on any OS with RPM tools
make package-deb         # Works on any OS with DEB tools
make package-all         # Builds everything

# CI-style build (auto-installs tools)
CI=true make package-ci  # Simulates CI environment
```

## Best Practices

### Go Packaging Guidelines
1. **Use pre-compiled binaries** instead of building from source during packaging
2. **Support multiple architectures** with proper detection
3. **Minimize dependencies** by using static binaries
4. **Follow distribution conventions** for file placement and metadata
5. **Provide backward compatibility** through symlinks and package provides

### Distribution Guidelines
- **RPM**: Follow Red Hat packaging guidelines for Go applications
- **DEB**: Follow Debian Go packaging policy
- **File placement**: Use standard FHS locations
- **Dependencies**: Declare minimal runtime dependencies only

## Development Notes

### Package Structure
The packaging system follows these principles:
- Uses pre-compiled binaries following Go best practices
- Maintains version consistency across all package types
- Includes proper dependency declarations
- Provides backward compatibility through symlinks
- Follows distribution-specific packaging guidelines
- Supports multiple architectures automatically

### Version Management
- Version is read from `.release-version` file
- Git commit hash is included in package metadata
- Build timestamp is recorded for traceability
- Architecture is automatically detected and included

### Cross-Platform Considerations
- **RPM building**: Requires `rpmbuild` and related tools
- **DEB building**: Requires `dpkg-dev` and `fakeroot`
- **macOS development**: Shows warnings but doesn't fail
- **Windows**: Not supported for package building
- **Architecture**: Automatically detects and builds for target architecture

## Contributing

When modifying the packaging system:
1. Test on both RPM and DEB-based systems
2. Verify package installation and removal
3. Test on multiple architectures if possible
4. Check backward compatibility with `locate` command
5. Update this documentation for any changes
6. Test in CI environment before merging
7. Follow Go packaging best practices

## Support

For packaging-related issues:
1. Check the troubleshooting section above
2. Verify your system has required tools installed
3. Ensure cross-platform binaries are built (`make build-cross`)
4. Review the build logs for specific error messages
5. Check architecture compatibility
6. Open an issue with your system information and error details

# Package Distribution

This document describes how to build and distribute packages for the glocate project.

## Automated Release Packaging

**Starting from version 0.1.3, all packages are automatically built and published with every GitHub release.**

### Release Process

When a new tag is pushed (e.g., `v0.1.3`), the GitHub Actions release workflow automatically:

1. **Builds cross-platform binaries** for Linux, macOS, and Windows
2. **Creates Linux packages**:
   - DEB packages for Debian/Ubuntu systems (amd64/arm64)
   - RPM packages for Red Hat/Fedora/CentOS systems (x86_64/aarch64)
3. **Signs all packages** with Cosign using private key
4. **Generates checksums** (SHA256/SHA512) for all packages
5. **Creates GitHub release** with comprehensive installation instructions
6. **Publishes Docker images** to GitHub Container Registry

### Package Availability

All packages are available on the [Releases page](https://github.com/Gosayram/go-locate/releases):

- **Binary tarballs**: `glocate-X.Y.Z-OS-ARCH.tar.gz`
- **DEB packages**: `glocate_X.Y.Z_ARCH.deb`
- **RPM packages**: `glocate-X.Y.Z-1.ARCH.rpm`
- **Signatures**: `*.sig` files for Cosign verification
- **Checksums**: `*.sha256` and `*.sha512` files

### Manual Release Creation

To create a release manually:

```bash
# Create and push a new tag
git tag v0.1.3
git push origin v0.1.3

# Or use GitHub CLI
gh release create v0.1.3 --generate-notes
```

The release workflow will automatically trigger and build all packages.

## Manual Package Building

For development and testing purposes, you can still build packages manually:

### Prerequisites
