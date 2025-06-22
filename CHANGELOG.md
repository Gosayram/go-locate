# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.3] - 2025-01-22

### Added
- **Release Automation**: Integrated RPM and DEB package building into GitHub release workflow
  - Automated package creation for every release tag
  - Multi-architecture support (amd64/arm64) for Linux packages
  - Cosign signing for all packages with checksums
  - Professional release notes with installation instructions for packages
  - Seamless integration with existing binary and Docker releases

### Fixed
- **Package Build Dependencies**: Fixed `package-deb` target to properly depend on `build-cross` instead of just `build`
  - DEB packages now correctly use pre-compiled cross-platform binaries
  - Resolved issue where DEB build script expected architecture-specific binaries
  - Added proper error handling for missing cross-platform binaries
  - Updated build-deb.sh to place packages in correct `packages/` directory

- **RPM Package Building**: Enhanced RPM build system with proper platform detection
  - Added intelligent platform detection for RPM tools availability
  - Clear error messages when attempting RPM builds on unsupported platforms (e.g., macOS)
  - Helpful guidance for using Docker or Linux systems for RPM packaging
  - Maintained functionality for CI/CD environments with proper Linux containers

- **Package Organization**: Improved package file placement and directory structure
  - All packages now correctly placed in `packages/` directory
  - Consistent naming conventions across all package types
  - Proper cleanup and organization of build artifacts

### Enhanced
- **Makefile Package Targets**: Improved reliability and user experience
  - `package-deb` now automatically builds required cross-platform binaries
  - `package-rpm` provides clear feedback when tools are unavailable
  - `package-all` works reliably on all platforms with appropriate warnings
  - Better error messages and guidance for platform-specific limitations

### Added
- **Automated Package Building System**: Complete overhaul of package building with auto-installation of build tools
  - Auto-detection of operating system (Linux distributions, macOS, FreeBSD)
  - Automatic installation of RPM build tools (`rpm-build`, `rpmdevtools`) on supported systems
  - Automatic installation of DEB build tools (`dpkg-dev`, `fakeroot`, `lintian`) on supported systems
  - Support for cross-platform package building (RPM on DEB systems, DEB on RPM systems)
  - CI/CD optimized package building with `make package-ci` command
  - Comprehensive package validation and testing
  - Cosign integration for package signing and verification
  - Multi-architecture support (AMD64, ARM64) for all package types

- **Enhanced CI/CD Integration**:
  - Professional GitHub Actions workflow with pinned action versions
  - Security hardening with step-security/harden-runner
  - Automatic package testing and validation
  - Cosign keyless signing for all packages
  - Comprehensive package reports with checksums and verification
  - Support for manual workflow triggers with configurable options

- **New Makefile Targets**:
  - `detect-os`: Detect operating system for package building
  - `install-rpm-tools`: Install RPM build tools with OS auto-detection
  - `install-deb-tools`: Install DEB build tools with OS auto-detection
  - `package-ci`: CI-optimized package building (auto-installs tools in CI)
  - `package-ci-setup`: Setup CI/CD packaging environment

- **Supported Operating Systems**:
  - **Ubuntu/Debian**: Native DEB building, cross-platform RPM support
  - **Fedora/RHEL/CentOS/Rocky/AlmaLinux**: Native RPM building, cross-platform DEB support
  - **openSUSE/SLES**: Full RPM and DEB support via zypper
  - **Arch Linux/Manjaro**: Full RPM and DEB support via pacman
  - **macOS**: Development support via Homebrew (cross-platform only)

### Enhanced
- **Package Building Documentation**:
  - New [Package Testing Guide](docs/PACKAGE-TESTING.md) with comprehensive testing procedures
  - Updated [Packaging Quick Start](docs/PACKAGING-QUICKSTART.md) with modern examples
  - Enhanced [Packaging Guide](docs/PACKAGING.md) with CI/CD integration examples
  - Professional GitHub Actions examples with security best practices

- **DEB Build Script**: Enhanced with automatic dependency installation and better error handling
- **Package Validation**: Comprehensive testing of package contents, metadata, and signatures
- **Security**: All packages signed with Cosign for verification and supply chain security

### Technical Details
- **Package Formats**: Binary tarballs (`.tar.gz`), DEB packages (`.deb`), RPM packages (`.rpm`)
- **Architectures**: AMD64 (x86_64) and ARM64 (aarch64) support
- **Signing**: Cosign keyless signing with GitHub OIDC for all package types
- **Validation**: Automated package structure, metadata, and signature verification
- **Performance**: Optimized build times with Go module caching and parallel builds

## [0.1.1] - 2025-01-18

### Added
- **Professional Version Management System**
  - Created dedicated `internal/version` package with comprehensive build information
  - Added build-time variables: Version, Commit, Date, BuiltBy, BuildNumber
  - Implemented `GetVersion()`, `GetFullVersionInfo()`, and `Get()` functions
  - Added `BuildInfo` struct with `String()` and `Short()` methods
  - Proper commit hash truncation and date formatting
  - Integration with CLI commands (`glocate version` and `glocate --version`)

- **Enhanced Development Infrastructure**
  - Completely rewritten Makefile with 100+ targets and modern build system
  - Advanced versioning system using COMMIT, DATE, BUILT_BY variables
  - Updated to golangci-lint v2 with auto-installation
  - Comprehensive LDFLAGS for version information injection
  - New sections: Configuration, Test Data, Documentation, CI/CD Support
  - Enhanced build targets with proper cross-compilation support
  - Improved help documentation with categorized targets
  - Integration tests and benchmark improvements

- **Docker and CI/CD Support**
  - Updated Dockerfile with build arguments for version information
  - Proper version injection during Docker builds
  - CI/CD pipeline support in Makefile
  - Cross-platform build support (Linux, macOS, Windows)

- **Testing and Quality Assurance**
  - Created comprehensive test suite for version package (100% coverage)
  - Fixed integration tests with proper command-line syntax
  - Added benchmark targets and performance reporting
  - Test data management with safe copying mechanisms
  - Coverage reporting and race detection tests

### Fixed
- **Integration Tests Issues**
  - Fixed incorrect command-line flag usage (`--pattern` → positional argument)
  - Corrected flag names (`--fuzzy` → `--advanced`, `--output` → `--format`)
  - Resolved command syntax errors (multiple arguments issue)
  - Updated test commands to use proper search scope limiting

- **Build System Improvements**
  - Fixed LDFLAGS to use correct import paths for version package
  - Resolved linter issues with unused parameters in version command
  - Updated import statements to use new version package structure
  - Fixed Docker build process with proper version injection

- **Code Quality Issues**
  - Eliminated magic numbers throughout codebase (following zero tolerance policy)
  - Fixed cyclomatic complexity warnings
  - Resolved linter warnings and staticcheck issues
  - Improved code organization and dependency management
