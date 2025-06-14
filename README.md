# go-locate

A modern, fast file search tool that replaces the outdated `locate` command with real-time file system searching.

## Why go-locate?

Traditional `locate` has several limitations:
- ❌ Relies on stale database (`updatedb`)
- ❌ Limited by `updatedb.conf` filters
- ❌ Name-only searching
- ❌ No real-time file system access

**go-locate** solves these problems by:
- ✅ Real-time file system scanning
- ✅ No database dependencies
- ✅ Advanced filtering options
- ✅ High-performance Go implementation
- ✅ Single static binary

## Installation

```bash
# Download from releases
curl -L https://github.com/Gosayram/go-locate/releases/latest/download/glocate-linux-amd64 -o glocate
chmod +x glocate
sudo mv glocate /usr/local/bin/

# Or use the install script
curl -sSL https://raw.githubusercontent.com/Gosayram/go-locate/main/install.sh | bash

# Or build from source
git clone https://github.com/Gosayram/go-locate.git
cd go-locate
make build
sudo make install
```

## Quick Start

```bash
# Basic file search
glocate config.json

# Pattern matching
glocate "*.go"

# Advanced fuzzy search
glocate --advanced "cfg"

# Extension filtering
glocate --ext go,rs,py "main"

# JSON output
glocate --format json "*.md"
```

## Usage

### Basic Commands

```bash
glocate filename                    # Find exact filename
glocate "pattern*"                  # Wildcard search
glocate --advanced substring        # Fuzzy matching
```

### Advanced Options

```bash
# Filter by file extension
glocate --ext go,rs,py "main"

# Filter by file size (planned)
glocate --size +100M               # Files larger than 100MB
glocate --size -1K                 # Files smaller than 1KB

# Filter by modification time (planned)
glocate --mtime -7d                # Modified in last 7 days
glocate --mtime +1h                # Modified more than 1 hour ago

# Exclude/include directories
glocate --exclude /proc,/sys --include /home,/opt "config"

# Search file content (planned)
glocate --content "TODO" --ext go
```

### Performance Options

```bash
glocate --threads 8                # Use 8 threads (default: CPU cores)
glocate --depth 5                  # Limit search depth
glocate --follow-symlinks          # Follow symbolic links
glocate --max-results 1000         # Limit number of results
```

### Output Formats

```bash
glocate --format path "*.go"       # Simple paths (default)
glocate --format detailed "*.go"   # Detailed file info
glocate --format json "*.go"       # JSON output
```

## Configuration

Create `~/.glocate.toml` for default settings:

```toml
[search]
exclude_dirs = ["/proc", "/sys", "/dev", "/tmp"]
include_dirs = ["/home", "/opt", "/usr"]
max_depth = 20
follow_symlinks = false
default_threads = 0  # 0 = use CPU count

[output]
format = "path"  # or "detailed", "json"
color = true
max_results = 1000
```

## Performance

Benchmarks on Apple M3 Pro:

| Operation | Time/op | Memory | Allocations |
|-----------|---------|--------|-------------|
| Fuzzy Match (short) | 6.32 ns | 0 B | 0 |
| Fuzzy Match (long) | 16.75 ns | 0 B | 0 |
| Pattern Matching | 297 ns | 0 B | 0 |
| Directory Exclusion | 1158 ns | 0 B | 0 |
| Extension Filtering | 623 ns | 0 B | 0 |

Real-world comparison:

| Tool | Time | Files Found | Notes |
|------|------|-------------|-------|
| `locate` | 0.05s | 234 | Stale database |
| `find` | 2.3s | 456 | Full scan |
| `glocate` | 0.12s | 456 | Real-time + filtering |

## Development

### Requirements

- Go 1.24.2+
- Make

### Building

```bash
# Install dependencies
make deps

# Build binary
make build

# Run tests with benchmarks
make test
make benchmark

# Run all quality checks
make check-all

# Cross-platform build
make build-cross
```

### Project Structure

```
go-locate/
├── cmd/glocate/           # CLI application
├── internal/
│   ├── config/           # Configuration management
│   ├── search/           # Search engine
│   └── output/           # Output formatting
├── .go-version           # Go version specification
├── .release-version      # Current release version
├── .golangci.yml        # Linting configuration
└── Makefile             # Build automation
```

### Running Benchmarks

```bash
# Run all benchmarks
make benchmark

# Generate benchmark report
make benchmark-report

# Run specific benchmarks
go test -bench=BenchmarkFuzzyMatch ./internal/search/
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests and benchmarks
5. Run `make check-all`
6. Submit a pull request

## License

MIT License - see [LICENSE](LICENSE) for details.

## Roadmap

### ✅ Completed Features
- [x] Basic file search with glob patterns
- [x] Advanced fuzzy matching
- [x] Extension filtering
- [x] Directory inclusion/exclusion
- [x] Multiple output formats (path, detailed, JSON)
- [x] Configurable threading
- [x] TOML configuration support
- [x] Cross-platform builds

### 🚧 Planned Features
- [ ] Size filtering implementation
- [ ] Modification time filtering
- [ ] Content search functionality
- [ ] Regular expression support
- [ ] Shell completion scripts
- [ ] GUI interface
- [ ] Plugin system

## Support

- 📖 [Documentation](docs/)
- 🐛 [Issue Tracker](https://github.com/Gosayram/go-locate/issues)
- 💬 [Discussions](https://github.com/Gosayram/go-locate/discussions)
- 📊 [Benchmarks](benchmark-report.md) 