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

### GitHub Actions Example

```yaml
name: Build Packages

on:
  push:
    tags:
      - 'v*'

jobs:
  build-packages:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v3
        with:
          go-version: 1.22

      - name: Install packaging tools
        run: |
          sudo apt-get update
          sudo apt-get install -y rpm dpkg-dev fakeroot lintian

      - name: Build cross-platform binaries
        run: make build-cross

      - name: Build packages
        run: make package-all

      - name: Upload packages
        uses: actions/upload-artifact@v3
        with:
          name: packages
          path: packages/
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
