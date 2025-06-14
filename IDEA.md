# go-locate: Modern File Search Tool

## Overview
`go-locate` (binary: `glocate`) is a modern, fast file search tool designed to replace the outdated `locate` command. Built entirely in Go, it provides real-time file system searching without relying on outdated databases, offering excellent performance with a simple, maintainable architecture.

## Core Problems with Traditional `locate`
1. **Stale Database**: Relies on `updatedb` which may be outdated
2. **Limited Scope**: Filters paths via `updatedb.conf`, missing important directories
3. **Name-only Search**: Cannot search by content, metadata, or advanced criteria
4. **No Live Recursion**: Only searches pre-indexed data

## Architecture

### Pure Go Design
- **Single Binary**: Self-contained executable with no external dependencies
- **Modular Structure**: Clean separation of concerns with internal packages
- **High Performance**: Optimized Go implementation with concurrent processing
- **Cross-Platform**: Native builds for Linux, macOS, Windows

### Performance Strategy
```
CLI Interface → Search Engine → File System
     ↓              ↓              ↓
  Parsing      Concurrent      Results
     ↓         Traversal          ↓
  Config        Filtering      Formatting
     ↓              ↓              ↓
  Output        Pattern        Display
              Matching
```

## Current Implementation

### Core Components
- **cmd/glocate**: CLI application with Cobra framework
- **internal/config**: Configuration management with TOML support
- **internal/search**: High-performance search engine with goroutines
- **internal/output**: Multiple output formats (path, detailed, JSON)

### Performance Features
- **Concurrent Search**: Multi-threaded directory traversal
- **Smart Exclusions**: Skip system directories (/proc, /sys, /dev, /tmp)
- **Memory Efficient**: Streaming results with configurable limits
- **Zero Allocations**: Optimized algorithms for hot paths

## Features

### ✅ Implemented Features
- **Exact Match**: Find files with exact name matching
- **Pattern Matching**: Wildcard support (`*`, `?`, `[]`)
- **Fuzzy Search**: Advanced substring matching with `--advanced` flag
- **Extension Filtering**: `--ext go,rs,py`
- **Directory Control**: `--include` and `--exclude` options
- **Output Formats**: Path, detailed, and JSON output
- **Configuration**: TOML configuration file support
- **Threading**: Configurable thread count

### 🚧 Planned Features
- **Size Filtering**: `--size +100M`, `--size -1K`
- **Modification Time**: `--mtime -7d`, `--mtime +1h`
- **Content Search**: `--content "pattern"`
- **Regular Expressions**: Full regex pattern support

## CLI Interface

### Basic Usage
```bash
glocate filename.txt              # Exact match
glocate "*.go"                    # Pattern match
glocate --advanced pattern        # Fuzzy search
```

### Advanced Usage
```bash
glocate --ext go,rs,py "main"                    # Extension filtering
glocate --exclude /proc,/sys --include /home "config"  # Directory control
glocate --format json --max-results 100 "*.md"  # JSON output with limits
glocate --threads 8 --depth 5 "test*"          # Performance tuning
```

## Technical Implementation

### Go Architecture
```
├── cmd/glocate/           # CLI application
│   └── main.go           # Entry point with Cobra CLI
├── internal/
│   ├── config/           # Configuration management
│   │   └── config.go     # TOML config with Viper
│   ├── search/           # Search engine
│   │   ├── search.go     # Core search logic
│   │   └── search_test.go # Unit tests & benchmarks
│   └── output/           # Output formatting
│       └── output.go     # Multiple format support
```

### Performance Optimizations
- **Goroutine Pool**: Concurrent directory traversal
- **Channel-based Communication**: Efficient result streaming
- **Path Filtering**: Early exclusion of irrelevant directories
- **Memory Pooling**: Reuse of common data structures

### Benchmarks (Apple M3 Pro)
```
BenchmarkAdvancedFuzzyMatch/Short-12      191M    6.32 ns/op     0 B/op    0 allocs/op
BenchmarkAdvancedFuzzyMatch/Long-12        70M   16.75 ns/op     0 B/op    0 allocs/op
BenchmarkPatternMatching/Wildcard-12       4M   297.4 ns/op     0 B/op    0 allocs/op
BenchmarkConcurrentSearch-12              1297  940171 ns/op  156606 B/op 1522 allocs/op
```

## Development Workflow

### Quality Assurance
- **Comprehensive Testing**: Unit tests with 55.1% coverage
- **Performance Benchmarks**: 8 different benchmark scenarios
- **Code Quality**: golangci-lint with strict rules
- **Automated Formatting**: goimports integration
- **Static Analysis**: staticcheck for bug detection

### Build System
```bash
make build          # Build binary
make test           # Run tests
make benchmark      # Performance testing
make check-all      # All quality checks
make build-cross    # Cross-platform builds
```

### Development Phases

#### ✅ Phase 1: Core Implementation (Completed)
- [x] Basic file traversal with Go
- [x] CLI interface with Cobra
- [x] Pattern matching with filepath.Match
- [x] Concurrent search implementation
- [x] Multiple output formats

#### ✅ Phase 2: Advanced Features (Completed)
- [x] Fuzzy matching algorithm
- [x] Extension filtering
- [x] Directory inclusion/exclusion
- [x] Configuration file support
- [x] Comprehensive testing and benchmarks

#### 🚧 Phase 3: Extended Functionality (In Progress)
- [ ] Size and time filtering implementation
- [ ] Content search functionality
- [ ] Regular expression support
- [ ] Shell completion scripts
- [ ] Package manager distributions

## Configuration

### Default Configuration (~/.glocate.toml)
```toml
[search]
exclude_dirs = ["/proc", "/sys", "/dev", "/tmp"]
include_dirs = []
max_depth = 20
follow_symlinks = false
default_threads = 0  # 0 = use CPU count

[output]
format = "path"  # or "detailed", "json"
color = true
max_results = 1000
```

## Performance Goals ✅ Achieved

- **Speed**: Real-time search with sub-second response times
- **Memory**: Minimal memory footprint with streaming results
- **Scalability**: Handles thousands of files efficiently
- **Responsiveness**: Concurrent processing with immediate results
- **Zero Dependencies**: Single static binary deployment

## Advantages of Pure Go Architecture

### ✅ Simplicity
- Single language and toolchain
- No CGO complexity
- Simple build process
- Easy debugging and profiling

### ✅ Performance
- Excellent concurrency with goroutines
- Efficient memory management
- Fast compilation
- Native performance without FFI overhead

### ✅ Deployment
- Single static binary
- Cross-platform builds
- No runtime dependencies
- Container-friendly

### ✅ Maintenance
- Unified codebase
- Consistent tooling
- Easier testing
- Better IDE support

## Future Enhancements
- GUI interface with web-based frontend
- Background indexing (optional mode)
- Network file system support
- Plugin architecture for custom filters
- Integration with popular editors/IDEs
- Shell completion for bash/zsh/fish 