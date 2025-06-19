# Security Analysis

This document describes the security analysis tools and processes used in the go-locate project.

## Security Scanning

The project uses multiple security analysis tools to ensure code security:

### Static Code Analysis with gosec

[gosec](https://github.com/securego/gosec) performs static security analysis of Go code by scanning the Go AST and SSA code representation. It detects common security issues like hardcoded credentials, insecure TLS configurations, and SQL injection vulnerabilities.

### Vulnerability Detection with govulncheck  

[govulncheck](https://golang.org/x/vuln/cmd/govulncheck) is the official Go vulnerability scanner that checks for known vulnerabilities in Go modules. It analyzes your codebase and dependencies to identify vulnerabilities that actually affect your code, reducing noise by only reporting issues in functions your code actually calls.

### Error Handling Analysis with errcheck

[errcheck](https://github.com/kisielk/errcheck) is a Go static analysis tool that checks for unchecked errors in Go code. It ensures that error return values are properly handled, which is crucial for robust and secure applications. Unhandled errors can lead to unexpected behavior, resource leaks, and security vulnerabilities.

### Code Style and Naming with revive

[revive](https://github.com/mgechev/revive) is a fast, configurable, and extensible linter for Go that enforces coding standards and best practices. It's integrated into our golangci-lint configuration and provides rules for naming conventions, code style, error handling patterns, and more. This helps maintain consistent, readable, and maintainable code across the project.

### Available Make Targets

#### Static Security Analysis (gosec)
- `make security-scan` - Run gosec security scanner with SARIF output
- `make security-scan-json` - Run gosec security scanner with JSON output  
- `make security-scan-html` - Run gosec security scanner with HTML output
- `make security-scan-ci` - Run gosec security scanner for CI (no-fail mode)

#### Vulnerability Detection (govulncheck)
- `make vuln-check` - Run govulncheck vulnerability scanner
- `make vuln-check-json` - Run govulncheck with JSON output
- `make vuln-check-ci` - Run govulncheck for CI (no-fail mode)

#### Error Handling Analysis (errcheck)
- `make errcheck` - Check for unchecked errors in Go code

#### Combined Analysis
- `make check-all` - Run all code quality checks including error checking, security and vulnerability scans

### Configuration

The gosec scanner is configured via the `.gosec.json` file in the project root. This configuration:

- Enables security auditing and nosec directive handling
- Configures specific rules for credential detection (G101), error handling (G104), and other security checks
- Excludes test files from certain security checks where appropriate
- Excludes vendor directories and test data from scanning

The errcheck tool is configured via the `.errcheck_excludes.txt` file which excludes common safe patterns:

- File close operations in defer statements
- Output operations to stdout/stderr
- Printf family functions that rarely fail
- Buffer write operations with predictable behavior

### GitHub Actions Integration

Security and vulnerability scanning is integrated into the CI/CD pipeline via GitHub Actions:

#### Automated Scans
- **gosec**: Static security analysis runs on every push and PR
- **govulncheck**: Vulnerability detection runs on every push and PR
- **errcheck**: Error handling analysis runs on every push and PR
- **revive**: Code style and naming conventions (via golangci-lint)
- **testify**: Comprehensive unit testing framework for robust test coverage
- **Trivy**: Additional vulnerability scanning for comprehensive coverage
- **Nancy**: Sonatype vulnerability checking
- **OpenSSF Scorecard**: Security posture assessment

#### Scheduling and Reports
- Runs on every push to main branch
- Runs on pull requests
- Scheduled weekly scans (Tuesdays at 07:20 UTC)
- Results uploaded to GitHub Security tab
- SARIF and JSON reports stored as artifacts

### Security Rules Coverage

The following gosec rules are actively monitored:

#### Credential Security
- **G101**: Hard coded credentials detection
- **G401**: Weak cryptographic hashes (MD5, SHA1)
- **G402**: Insecure TLS connection settings
- **G404**: Insecure random number sources

#### Code Injection Prevention  
- **G201**: SQL query construction using format strings
- **G202**: SQL query construction using string concatenation
- **G204**: Command execution audit

#### File System Security
- **G301**: Poor file permissions when creating directories
- **G302**: Poor file permissions with chmod
- **G306**: Poor file permissions when writing files

#### Error Handling
- **G104**: Audit unchecked errors
- **G115**: Integer overflow detection

### Suppressing False Positives

For legitimate cases where security warnings are false positives, use the `#nosec` annotation:

```go
// Example: Acceptable use of weak random for non-security purposes
randomValue := rand.Intn(100) // #nosec G404 -- Used for test data generation only
```

### Running Security Scans Locally

#### Quick Start
1. Install all security tools:
   ```bash
   make install-tools
   ```

2. Run all security checks:
   ```bash
   make check-all
   ```

#### Individual Scans

**Static Security Analysis (gosec):**
```bash
make security-scan          # SARIF output
make security-scan-json     # JSON output  
make security-scan-html     # HTML output
```

**Vulnerability Detection (govulncheck):**
```bash
make vuln-check             # Standard output
make vuln-check-json        # JSON output
```

**Error Handling Analysis (errcheck):**
```bash
make errcheck               # Check for unchecked errors
```

#### View Results
```bash
# View gosec results
cat gosec-report.sarif
cat gosec-report.json

# View govulncheck results  
cat vulncheck-report.json
```

### Continuous Security Monitoring

The project maintains continuous security monitoring through:

1. **Automated Scanning**: Every code change triggers security analysis
2. **Dependency Monitoring**: Regular checks for vulnerable dependencies
3. **Security Advisories**: GitHub Security Advisory monitoring
4. **SARIF Integration**: Results integrated with GitHub Security tab

### Security Best Practices

1. **Never commit secrets**: Use environment variables or secure secret management
2. **Validate input**: Always validate and sanitize user input
3. **Use strong cryptography**: Prefer secure algorithms and proper key management
4. **Handle errors**: Never ignore security-related errors
5. **Principle of least privilege**: Use minimal required permissions
6. **Regular updates**: Keep dependencies updated with security patches

### Known Limitations and Issues

#### Outdated CWE Taxonomy

**Current Status**: gosec v2.22.5 uses hardcoded CWE (Common Weakness Enumeration) taxonomy version **4.4** released on March 15, 2021.

**Latest Available**: CWE taxonomy version **4.17** was released on April 3, 2025.

**Impact**: 
- SARIF reports contain outdated CWE classifications
- Newer weakness categories and updated descriptions are not available
- Security tooling integration may reference deprecated CWE entries

**Workaround**: 
- This is a limitation of the gosec tool itself, not our configuration
- gosec hardcodes the CWE version in its source code without configuration options
- Consider upgrading gosec when newer versions with updated CWE taxonomy become available

**Tracking**: Monitor [gosec releases](https://github.com/securego/gosec/releases) for CWE taxonomy updates

### Reporting Security Issues

If you discover a security vulnerability, please report it privately to the maintainers through GitHub Security Advisories rather than opening a public issue. 