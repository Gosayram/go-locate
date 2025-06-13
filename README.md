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
- ✅ High-performance Rust backend

## Installation

```bash
# Download from releases for Linux (pls, choose via OS requirements)
curl -L https://github.com/Gosayram/go-locate/releases/latest/download/glocate.linux -o glocate
chmod +x glocate
sudo mv glocate /usr/local/bin/

# Or build from source
git clone https://github.com/yourusername/go-locate.git
cd go-locate
make build
```

## Quick Start

```bash
# Basic file search
glocate config.json

# Pattern matching
glocate "*.go"

# Advanced filtering
glocate --ext go --size +1M --mtime -7d "main"

# Search with content filtering
glocate --content "TODO" --ext rs,go
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

# Filter by file size
glocate --size +100M               # Files larger than 100MB
glocate --size -1K                 # Files smaller than 1KB

# Filter by modification time
glocate --mtime -7d                # Modified in last 7 days
glocate --mtime +1h                # Modified more than 1 hour ago

# Exclude/include directories
glocate --exclude /proc,/sys --include /home,/opt "config"

# Search file content
glocate --content "TODO" --ext go
```

### Performance Options

```bash
glocate --threads 8                # Use 8 threads (default: CPU cores)
glocate --depth 5                  # Limit search depth
glocate --follow-symlinks          # Follow symbolic links
```

## Configuration

Create `~/.glocate.toml` for default settings:

```toml
[search]
exclude_dirs = ["/proc", "/sys", "/dev"]
include_dirs = ["/home", "/opt", "/usr"]
max_depth = 20
follow_symlinks = false

[output]
format = "path"  # or "detailed", "json"
color = true
max_results = 1000
```

## Performance

Benchmarks on a typical development machine (not approved statistics):

| Tool | Time | Files Found |
|------|------|-------------|
| `locate` | 0.05s | 234 (stale) |
| `find` | 2.3s | 456 |
| `glocate` | 0.12s | 456 |

## Development

### Requirements

- Go 1.21+
- Rust 1.70+
- Make

### Building

```bash
# Install dependencies
make deps

# Build binary
make build

# Run tests
make test

# Run linter
make lint
```

### Project Structure

```
go-locate/
├── cmd/glocate/           # Go CLI application
├── internal/              # Go internal packages
├── rust-core/             # Rust search engine
├── .go-version            # Go version specification
├── .release-version       # Current release version
├── .golangci.yml         # Linting configuration
└── Makefile              # Build automation
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run `make lint` and `make test`
6. Submit a pull request

## License

MIT License - see [LICENSE](LICENSE) for details.

## Roadmap

- [x] Basic file search
- [ ] Pattern matching
- [ ] Advanced filtering
- [ ] Content search
- [ ] GUI interface
- [ ] Plugin system

## Support

- 📖 [Documentation](docs/)
- 🐛 [Issue Tracker](https://github.com/Gosayram/go-locate/issues)
- 💬 [Discussions](https://github.com/Gosayram/go-locate/discussions) 