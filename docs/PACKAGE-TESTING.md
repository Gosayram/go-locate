# Package Testing Guide

This document describes how to test the automated package building system for go-locate.

## 🎯 Overview

The package testing system builds, validates, and signs packages automatically:

- **Binary Tarballs**: Cross-platform distribution packages
- **DEB Packages**: Debian/Ubuntu native packages  
- **RPM Packages**: Red Hat/Fedora/CentOS native packages
- **Cosign Signatures**: Security verification for all packages

## 🚀 Quick Testing

### Local Testing

```bash
# Test OS detection
make detect-os

# Test binary package creation
make package-binaries

# Test platform-specific packages
make package-deb    # On DEB-based systems
make package-rpm    # On RPM-based systems

# Test CI-style build
CI=true make package-ci
```

### GitHub Actions Testing

1. **Manual Trigger**: Go to Actions → "Package Build and Test" → "Run workflow"
2. **Choose Build Type**:
   - `test`: Build binaries + one platform package (fastest)
   - `full`: Build all package types (comprehensive)
   - `binaries-only`: Just binary tarballs (minimal)
3. **Enable Signing**: Check "Sign packages with cosign" for security testing

## 📦 Package Validation

### Automated Tests

The CI system automatically validates:

```bash
# Binary tarball tests
tar -tf package.tar.gz                    # List contents
./extracted/glocate --version            # Test execution
./extracted/glocate --help               # Test help

# DEB package tests  
dpkg-deb --info package.deb              # Package metadata
dpkg-deb --contents package.deb          # File listing
dpkg-deb --field package.deb Version     # Version check

# RPM package tests
rpm -qip package.rpm                     # Package info
rpm -qlp package.rpm                     # File listing  
rpm -qp --queryformat '%{VERSION}' package.rpm  # Version check
```

### Manual Validation

```bash
# Download artifacts from GitHub Actions
cd packages/

# Verify checksums
sha256sum -c checksums.txt

# Verify cosign signatures (if signed)
cosign verify-blob --signature package.tar.gz.sig package.tar.gz

# Test package installation
sudo dpkg -i glocate_*.deb              # DEB installation
sudo rpm -i glocate-*.rpm               # RPM installation
```

## 🔐 Security Testing

### Cosign Signature Verification

```bash
# Install cosign
curl -O -L "https://github.com/sigstore/cosign/releases/latest/download/cosign-linux-amd64"
sudo mv cosign-linux-amd64 /usr/local/bin/cosign
sudo chmod +x /usr/local/bin/cosign

# Verify package signatures
cosign verify-blob --signature package.tar.gz.sig package.tar.gz
cosign verify-blob --signature package.deb.sig package.deb  
cosign verify-blob --signature package.rpm.sig package.rpm
```

### Package Integrity

```bash
# Check package integrity
dpkg-deb --fsys-tarfile package.deb | tar -tv  # DEB contents
rpm2cpio package.rpm | cpio -tv               # RPM contents

# Verify no malicious content
find extracted/ -type f -executable -exec file {} \;
strings extracted/glocate | grep -i -E "(wget|curl|http|ftp)"
```

## 🧪 Test Scenarios

### Scenario 1: Basic Package Build

```yaml
# Trigger: Push to main branch
# Expected: Binary tarballs created and tested
# Validation:
#   - glocate-X.Y.Z-linux-amd64.tar.gz exists
#   - glocate-X.Y.Z-linux-arm64.tar.gz exists
#   - Both execute --version successfully
```

### Scenario 2: Multi-Platform Build

```yaml
# Trigger: Manual workflow with "full" build type
# Expected: All package types created
# Validation:
#   - Binary tarballs (2 files)
#   - DEB packages (2 architectures)  
#   - RPM packages (2 architectures)
#   - All packages signed with cosign
```

### Scenario 3: Cross-Platform Compatibility

```yaml
# Trigger: Matrix test on different OS containers
# Expected: Packages build on Fedora, Ubuntu, Debian, CentOS
# Validation:
#   - Each OS produces valid packages
#   - Package metadata is correct
#   - No OS-specific dependencies
```

## 🔍 Troubleshooting

### Common Issues

#### Build Tool Installation Fails

```bash
# Check OS detection
make detect-os

# Manual tool installation
sudo apt-get install dpkg-dev fakeroot    # Ubuntu/Debian
sudo dnf install rpm-build rpmdevtools    # Fedora/RHEL

# Verify installation
command -v dpkg-deb && echo "DEB tools OK"
command -v rpmbuild && echo "RPM tools OK"
```

#### Package Validation Fails

```bash
# Check package format
file packages/*.deb packages/*.rpm packages/*.tar.gz

# Verify package contents
dpkg-deb --contents packages/*.deb | grep glocate
rpm -qlp packages/*.rpm | grep glocate
tar -tzf packages/*.tar.gz | grep glocate
```

#### Cosign Signing Fails

```bash
# Check cosign installation
cosign version

# Verify OIDC token (in CI)
echo $ACTIONS_ID_TOKEN_REQUEST_TOKEN

# Test manual signing
cosign sign-blob --yes testfile.txt --output-signature testfile.sig
cosign verify-blob --signature testfile.sig testfile.txt
```

### Debug Commands

```bash
# Verbose package building
make package-deb VERBOSE=1

# Check build environment
make detect-os
go env
rpm --showrc | grep _topdir
dpkg-architecture

# Validate package dependencies
dpkg-deb --field package.deb Depends
rpm -qp --requires package.rpm
```

## 📊 Test Results

### Success Criteria

- ✅ All package types build without errors
- ✅ Packages contain correct files and permissions
- ✅ Binary executes and shows correct version
- ✅ Package metadata is accurate
- ✅ Cosign signatures verify successfully
- ✅ No security warnings or vulnerabilities

### Performance Benchmarks

| Build Type | Time | Artifacts | Size |
|------------|------|-----------|------|
| Binary Only | ~2 min | 2 files | ~4MB |
| DEB Build | ~3 min | 2 files | ~4MB |
| RPM Build | ~4 min | 2 files | ~4MB |
| Full Build | ~6 min | 6 files | ~12MB |

## 🔗 Related Documentation

- [Packaging Quick Start](PACKAGING-QUICKSTART.md) - Getting started guide
- [Full Packaging Guide](PACKAGING.md) - Complete documentation  
- [Security Guide](SECURITY.md) - Security best practices
