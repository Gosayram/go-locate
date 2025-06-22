# Matrix Testing Guide

This document describes the matrix testing strategy for go-locate, which ensures compatibility across multiple Go versions and operating systems.

## Overview

Matrix testing runs the complete test suite across different combinations of:
- **Go versions**: 1.22, 1.23, 1.24.4, 1.24
- **Operating systems**: Ubuntu, macOS, Windows
- **Architectures**: amd64, arm64 (limited combinations)
- **Experimental versions**: Go 1.25rc1 (optional)

## Features

### 🔄 Automated Triggers
- **Push to main/develop**: Runs matrix tests automatically
- **Pull requests**: Tests all combinations for comprehensive validation
- **Scheduled runs**: Daily at 02:00 UTC to catch compatibility issues early
- **Manual dispatch**: On-demand testing with configurable options

### ⚙️ Configuration Options

#### Skip Failures Mode
- **Default**: `true` - Non-blocking testing continues even if some combinations fail
- **Purpose**: Prevents experimental or edge-case failures from blocking development
- **Control**: Can be disabled for critical release testing

#### Experimental Testing
- **Default**: `false` - Stable versions only
- **When enabled**: Includes Go 1.25rc1 and other pre-release versions
- **Use case**: Early compatibility testing with upcoming Go releases

### 🎯 Test Matrix

#### Standard Matrix (12 combinations)
```
Go Versions × Operating Systems:
├── Go 1.22  × [Ubuntu, macOS]           (2 combinations)
├── Go 1.23  × [Ubuntu, macOS, Windows]  (3 combinations)
├── Go 1.24.4× [Ubuntu, macOS, Windows]  (3 combinations)
└── Go 1.24  × [Ubuntu, macOS, Windows]  (3 combinations)
Plus Go 1.24.4 × Ubuntu × ARM64          (1 combination)
Total: 12 combinations
```

#### With Experimental (13 combinations)
```
Standard Matrix + Go 1.25rc1 × Ubuntu = 13 combinations
```

#### Excluded Combinations
- **Windows + Go 1.22**: Compatibility issues
- **macOS + Go 1.22**: Platform-specific problems

## Test Execution

### Each Matrix Job Runs:
1. **Environment Setup**
   - Go version installation and verification
   - Dependency caching and download
   - Environment variable display

2. **Code Quality Checks**
   - `go vet` static analysis
   - Module verification

3. **Testing Suite**
   - Standard tests with timeout (10 minutes)
   - Race detection tests (except Windows)
   - Coverage testing with threshold (50%)

4. **Integration Testing**
   - Binary compilation for target platform
   - Binary execution verification
   - CLI functionality tests
   - Search functionality validation

5. **Artifacts Collection**
   - Coverage reports
   - Test binaries
   - Benchmark results (for stable versions)

### Coverage Requirements
- **Minimum threshold**: 50%
- **Current coverage**: ~55.1% (search package)
- **Experimental versions**: Threshold skipped to prevent blocking

## Local Matrix Testing

### Quick Setup
```bash
# Check current Go version compatibility
make test-go-versions

# Run tests with current Go version
make test-multi-go

# View matrix configuration
make matrix-info
```

### Multi-Version Testing
```bash
# Install additional Go versions
go install golang.org/dl/go1.22@latest && go1.22 download
go install golang.org/dl/go1.23@latest && go1.23 download

# Run local matrix tests
make matrix-test-local
```

## CI/CD Integration

### Matrix Test Results
- **Parallel execution**: All combinations run simultaneously
- **Failure isolation**: One failing combination doesn't stop others
- **Artifacts storage**: 30-day retention for debugging
- **Coverage upload**: Only stable versions upload to Codecov

### Status Reporting
- **Matrix summary job**: Aggregates results from all combinations
- **GitHub Actions UI**: Clear visualization of matrix results
- **Notifications**: Warnings for failures, success confirmations

## Manual Triggering

### Via GitHub UI
1. Go to **Actions** → **Matrix Testing**
2. Click **Run workflow**
3. Configure options:
   - **Skip failures**: `true`/`false`
   - **Test experimental**: `true`/`false`
4. Click **Run workflow**

### Via GitHub CLI
```bash
# Trigger with default settings
gh workflow run matrix-test.yml

# Trigger with experimental testing enabled
gh workflow run matrix-test.yml -f test_experimental=true -f skip_failures=false
```

## Performance Considerations

### Resource Usage
- **Parallel jobs**: Up to 13 simultaneous runners
- **Duration**: ~5-10 minutes per combination
- **Total time**: ~10-15 minutes (parallel execution)
- **Cost**: GitHub Actions minutes consumed

### Optimization Strategies
- **Caching**: Go modules and build cache across jobs
- **Selective runs**: Experimental versions only when requested
- **Artifact limits**: 30-day retention to manage storage
- **Windows optimization**: Skip race detection for performance

## Troubleshooting

### Common Issues

#### Go Version Compatibility
```bash
# Check minimum version requirement
make test-go-versions

# Expected output for compatible version:
# ✅ Go version 1.24.4 meets minimum requirement
```

#### Local Matrix Testing
```bash
# If go1.22/go1.23 commands not found:
go install golang.org/dl/go1.22@latest
go1.22 download
```

#### CI Failures
1. **Check matrix summary**: Review aggregated results
2. **Examine individual jobs**: Identify specific failure patterns
3. **Review artifacts**: Download coverage and binary artifacts
4. **Local reproduction**: Use same Go version locally

### Matrix Job Status Interpretation
- **✅ Success**: All tests passed for this combination
- **❌ Failure**: Tests failed (may be acceptable if `skip_failures=true`)
- **⚠️ Skipped**: Job excluded from matrix
- **🔄 In Progress**: Currently running

## Benefits

### For Development
- **Early detection**: Compatibility issues caught before release
- **Confidence**: Code works across supported environments
- **Documentation**: Clear support matrix for users

### For Users
- **Reliability**: Tested across their likely environments
- **Transparency**: Known supported versions and platforms
- **Quality assurance**: Multiple validation layers

### For Maintenance
- **Automated**: No manual intervention required
- **Scalable**: Easy to add new Go versions or platforms
- **Historical data**: Trend analysis of compatibility

## Configuration Constants

```makefile
MATRIX_MIN_GO_VERSION := 1.22
MATRIX_STABLE_GO_VERSION := 1.24.4
MATRIX_LATEST_GO_VERSION := 1.24
MATRIX_TEST_TIMEOUT := 10m
MATRIX_COVERAGE_THRESHOLD := 50
```

## Future Enhancements

### Planned Improvements
- [ ] **Database compatibility**: Add database backend testing
- [ ] **Extended architectures**: ARM32, RISC-V support
- [ ] **Container testing**: Docker image validation
- [ ] **Performance benchmarks**: Cross-version performance comparison
- [ ] **Integration matrix**: Test with different dependency versions

### Version Strategy
- **Add new versions**: When Go releases new versions
- **Deprecate old versions**: Following Go's support policy
- **Experimental inclusion**: Beta/RC versions for early testing

## Support

For matrix testing issues:
1. **Check workflow logs**: GitHub Actions provides detailed output
2. **Review documentation**: This guide and inline comments
3. **Local reproduction**: Use `make matrix-test-local`
4. **Issue reporting**: Include matrix job URLs and logs
