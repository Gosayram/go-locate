# Package Building Quick Start Guide

This guide shows how to use the new automated package building system with auto-installation of build tools.

## 🚀 Quick Commands

### For CI/CD Pipelines

```bash
# Build packages optimized for CI (auto-installs tools in CI environments)
make package-ci

# Setup CI environment with all tools
CI=true make package-ci-setup
```

### For Local Development

```bash
# Check your operating system
make detect-os

# Install build tools for your platform
make install-rpm-tools    # For RPM packages
make install-deb-tools    # For DEB packages

# Build specific package types
make package-binaries     # Binary tarballs (works everywhere)
make package-rpm          # RPM packages (auto-installs tools)
make package-deb          # DEB packages (auto-installs tools)
make package-all          # Everything
```

## 🛠️ Supported Operating Systems

### Automatic Tool Installation Support

| OS Family | Package Manager | RPM Tools | DEB Tools | Status |
|-----------|----------------|-----------|-----------|---------|
| **Ubuntu/Debian** | `apt-get` | ✅ Cross-platform | ✅ Native | Full Support |
| **Fedora/RHEL/CentOS** | `dnf`/`yum` | ✅ Native | ✅ Cross-platform | Full Support |
| **openSUSE/SLES** | `zypper` | ✅ Native | ✅ Cross-platform | Full Support |
| **Arch Linux** | `pacman` | ✅ Available | ✅ Available | Full Support |
| **macOS** | `brew` | ✅ Cross-platform | ✅ Cross-platform | Development Only |
| **Other Linux** | Manual | ⚠️ Manual Install | ⚠️ Manual Install | Limited |

## 📦 Package Types

### 1. Binary Tarballs (Recommended for CI)
- **Command**: `make package-binaries`
- **Output**: `glocate-X.Y.Z-linux-{amd64,arm64}.tar.gz`
- **Platforms**: All Linux architectures
- **Use Case**: Distribution, CI/CD, manual installation

### 2. RPM Packages
- **Command**: `make package-rpm`
- **Output**: `glocate-X.Y.Z-1.{x86_64,aarch64}.rpm`
- **Platforms**: RHEL, Fedora, CentOS, openSUSE
- **Use Case**: Enterprise Linux distributions

### 3. DEB Packages
- **Command**: `make package-deb`
- **Output**: `glocate_X.Y.Z_{amd64,arm64}.deb`
- **Platforms**: Debian, Ubuntu, derivatives
- **Use Case**: Debian-based distributions

### 4. Source Tarball
- **Command**: `make package-tarball`
- **Output**: `glocate-X.Y.Z-src.tar.gz`
- **Use Case**: Source distribution, custom builds

## 🔧 CI/CD Integration Examples

### GitHub Actions - Simple

```yaml
name: Build Packages

on: [push, pull_request]

permissions:
  contents: read
  id-token: write  # Required for cosign signing

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Harden Runner
        uses: step-security/harden-runner@002fdce3c6a235733a90a27c80493a3241e56863 # v2.12.1
        with:
          disable-sudo: false
          egress-policy: audit

      - name: Checkout Code
        uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 # v4.2.2

      - name: Get Go Version
        shell: bash
        run: |
          #!/bin/bash
          GOVERSION=$({ [ -f .go-version ] && cat .go-version; })
          echo "GOVERSION=$GOVERSION" >> $GITHUB_ENV

      - name: Setup Go
        uses: actions/setup-go@d35c59abb061a4a6fb18e82ac0862c26744d6ab5 # v5.5.0
        with:
          go-version: ${{ env.GOVERSION }}

      - name: Build packages
        run: make package-ci

      - name: Upload artifacts
        uses: actions/upload-artifact@b4b15b8c7c6ac21ea08fcf65892d2ee8f75cf882 # v4.4.3
        with:
          name: packages
          path: packages/
          retention-days: 7
```

### GitHub Actions - Multi-Platform with Cosign

```yaml
name: Multi-Platform Build with Signing

on:
  workflow_dispatch:
    inputs:
      sign_packages:
        description: 'Sign packages with cosign'
        required: false
        default: true
        type: boolean

permissions:
  contents: read
  id-token: write  # Required for cosign signing

env:
  COSIGN_EXPERIMENTAL: 1

jobs:
  build:
    strategy:
      matrix:
        include:
          - os: ubuntu-latest
            type: deb
            name: "DEB Packages"
          - os: ubuntu-latest
            container: fedora:latest
            type: rpm
            name: "RPM Packages"

    runs-on: ${{ matrix.os }}
    container: ${{ matrix.container }}

    steps:
      - name: Install tools
        run: |
          if command -v dnf >/dev/null; then
            dnf install -y git make golang
          else
            apt-get update && apt-get install -y git make
          fi

      - name: Checkout Code
        uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 # v4.2.2

      - name: Get Go Version
        if: "!matrix.container"
        shell: bash
        run: |
          #!/bin/bash
          GOVERSION=$({ [ -f .go-version ] && cat .go-version; })
          echo "GOVERSION=$GOVERSION" >> $GITHUB_ENV

      - name: Setup Go
        if: "!matrix.container"
        uses: actions/setup-go@d35c59abb061a4a6fb18e82ac0862c26744d6ab5 # v5.5.0
        with:
          go-version: ${{ env.GOVERSION }}

      - name: Build packages
        run: make package-${{ matrix.type }}

      - name: Install Cosign
        if: github.event.inputs.sign_packages != 'false'
        uses: sigstore/cosign-installer@3454372f43399081ed03b604cb2d021dabca52bb # v3.8.2
        with:
          cosign-release: 'v2.4.3'

      - name: Sign packages
        if: github.event.inputs.sign_packages != 'false'
        run: |
          cd packages
          for pkg in *.${{ matrix.type == 'deb' && 'deb' || 'rpm' }}; do
            cosign sign-blob --yes "$pkg" --output-signature "${pkg}.sig"
          done

      - name: Upload packages
        uses: actions/upload-artifact@b4b15b8c7c6ac21ea08fcf65892d2ee8f75cf882 # v4.4.3
        with:
          name: ${{ matrix.name }}
          path: |
            packages/*.${{ matrix.type == 'deb' && 'deb' || 'rpm' }}
            packages/*.sig
          retention-days: 7
```

### GitLab CI

```yaml
stages:
  - build

package-build:
  stage: build
  image: ubuntu:22.04
  before_script:
    - apt-get update && apt-get install -y make git curl
  script:
    - make package-ci
  artifacts:
    paths:
      - packages/
    expire_in: 1 week
```

### Jenkins Pipeline

```groovy
pipeline {
    agent any

    stages {
        stage('Setup') {
            steps {
                sh 'make detect-os'
            }
        }

        stage('Build Packages') {
            steps {
                sh 'make package-ci'
            }
        }

        stage('Archive') {
            steps {
                archiveArtifacts artifacts: 'packages/*', fingerprint: true
            }
        }
    }
}
```

## 🐛 Troubleshooting

### Tool Installation Fails

```bash
# Check OS detection
make detect-os

# Manual installation
sudo apt-get install dpkg-dev fakeroot    # Ubuntu/Debian
sudo dnf install rpm-build rpmdevtools    # Fedora/RHEL
```

### Permission Issues

```bash
# Add user to necessary groups (may require logout/login)
sudo usermod -a -G docker $USER    # For Docker-based builds
```

### Missing Dependencies

```bash
# Install Go if not available
curl -OL https://golang.org/dl/go1.22.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

## 📋 Build Output

After successful build, you'll find packages in the `packages/` directory:

```
packages/
├── glocate-0.1.2-linux-amd64.tar.gz      # AMD64 binary tarball
├── glocate-0.1.2-linux-arm64.tar.gz      # ARM64 binary tarball
├── glocate-0.1.2-src.tar.gz              # Source tarball
├── glocate-0.1.2-1.x86_64.rpm            # AMD64 RPM package
├── glocate-0.1.2-1.aarch64.rpm           # ARM64 RPM package
├── glocate_0.1.2_amd64.deb               # AMD64 DEB package
└── glocate_0.1.2_arm64.deb               # ARM64 DEB package
```

## 🎯 Best Practices

1. **Use `make package-ci` for CI/CD** - it's optimized for automated environments
2. **Test locally first** - run `make package-binaries` to verify builds work
3. **Platform-specific builds** - use containers for consistent RPM/DEB builds
4. **Version management** - ensure `.release-version` file is up to date
5. **Artifact storage** - upload packages to appropriate repositories

## 🔗 Related Documentation

- [Full Packaging Guide](PACKAGING.md) - Complete documentation
- [Development Guide](DEVELOPMENT.md) - Development setup
- [Testing Guide](TESTING.md) - Testing procedures
