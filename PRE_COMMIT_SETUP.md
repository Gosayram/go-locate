# Pre-commit Setup for go-locate

This document describes the setup and usage of pre-commit hooks for the go-locate project.

## Installing pre-commit

### macOS (with Homebrew)
```bash
brew install pre-commit
```

### Python pip
```bash
pip install pre-commit
```

### Verify installation
```bash
pre-commit --version
```

## Project Setup

### 1. Install hooks
```bash
# Install pre-commit hooks in git
pre-commit install

# Install hook for commit message validation
pre-commit install --hook-type commit-msg
```

### 2. Initial run
```bash
# Run hooks on all files (recommended for initial setup)
pre-commit run --all-files
```

## Configured Hooks

### Basic file checks
- **trailing-whitespace**: Removes trailing whitespace
- **end-of-file-fixer**: Ensures files end with a newline
- **check-yaml**: Validates YAML file syntax
- **check-toml**: Validates TOML file syntax
- **check-json**: Validates JSON file syntax
- **check-added-large-files**: Prevents committing large files (>1MB)
- **check-merge-conflict**: Checks for merge conflict markers
- **mixed-line-ending**: Ensures consistent line endings (LF)

### Go-specific hooks
- **go-fmt-repo**: Code formatting with gofmt
- **go-imports-repo**: Import management with goimports
- **go-mod-tidy-repo**: Clean up go.mod files
- **go-vet-repo-mod**: Static analysis with go vet
- **go-build-repo-mod**: Compilation check
- **go-test-repo-mod**: Run tests with code coverage
- **golangci-lint-repo-mod**: Run golangci-lint with existing configuration

### Security
- **detect-secrets**: Detect secrets in code

### Commit Messages
- **custom-commit-msg**: Validate commit message format (supports both custom [TYPE] - description and conventional formats)

## Usage

### Automatic execution
Pre-commit hooks run automatically on every `git commit`.

### Manual execution
```bash
# Run all hooks
pre-commit run --all-files

# Run specific hook
pre-commit run go-fmt-repo

# Run hooks only on changed files
pre-commit run
```

### Skipping hooks
```bash
# Skip all hooks for one commit
git commit --no-verify -m "commit message"

# Skip specific hooks
SKIP=go-test-repo-mod git commit -m "commit message"

# Skip multiple hooks
SKIP=go-test-repo-mod,golangci-lint-repo-mod git commit -m "commit message"
```

## Commit Message Format

The project supports two commit message formats:

### Format 1: Custom Format (Recommended)

```
[TYPE] - description

[optional body]

[optional footer(s)]
```

**Allowed Types:**
- `ADD` - Adding new features or files
- `CI` - Continuous Integration changes
- `FEATURE` - New feature implementation
- `BUGFIX` - Bug fixes
- `FIX` - General fixes
- `INIT` - Initial project setup
- `DOCS` - Documentation changes
- `TEST` - Adding or modifying tests
- `REFACTOR` - Code refactoring
- `STYLE` - Code style changes
- `CHORE` - Maintenance tasks

**Examples:**
```bash
git commit -m "[ADD] - new search functionality"
git commit -m "[FIX] - resolve configuration parsing error"
git commit -m "[CI] - update GitHub Actions workflow"
git commit -m "[DOCS] - update installation instructions"
git commit -m "[TEST] - add unit tests for search module"
```

### Format 2: Conventional Commits

```
type: description
type(scope): description

[optional body]

[optional footer(s)]
```

**Allowed Types:**
- `feat` - New features
- `fix` - Bug fixes
- `docs` - Documentation changes
- `style` - Code style changes (formatting, etc.)
- `refactor` - Code refactoring
- `test` - Adding or modifying tests
- `chore` - Maintenance tasks
- `ci` - CI/CD changes
- `build` - Build system changes
- `perf` - Performance improvements
- `revert` - Reverting changes

**Examples:**
```bash
git commit -m "feat: add new search functionality"
git commit -m "fix: resolve configuration parsing error"
git commit -m "chore: update version to 0.1.1"
git commit -m "feat(search): add fuzzy matching support"
git commit -m "fix(config): handle missing configuration file"
```

## Updating Hooks

```bash
# Update hooks to latest versions
pre-commit autoupdate

# Clear pre-commit cache (if issues occur)
pre-commit clean
```

## IDE Configuration

### VS Code
Add to `.vscode/settings.json`:
```json
{
    "go.formatTool": "goimports",
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
        "source.organizeImports": true
    }
}
```

### GoLand/IntelliJ IDEA
1. Configure goimports as the default formatter
2. Enable auto-formatting on save
3. Configure golangci-lint as an external tool

## Troubleshooting

### Hook installation issues
```bash
# Reinstall hooks
pre-commit uninstall
pre-commit install --install-hooks
```

### golangci-lint issues
```bash
# Ensure golangci-lint is installed
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Check configuration
golangci-lint config -h
```

### Go module issues
```bash
# Clean modules and reinstall dependencies
go clean -modcache
go mod download
go mod tidy
```

## Additional Resources

- [Pre-commit documentation](https://pre-commit.com/)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [GolangCI-Lint documentation](https://golangci-lint.run/)
- [TekWizely/pre-commit-golang](https://github.com/TekWizely/pre-commit-golang)
