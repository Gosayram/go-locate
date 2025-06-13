# go-locate: Modern File Search Tool

## Overview
`go-locate` (binary: `glocate`) is a modern, fast file search tool designed to replace the outdated `locate` command. It combines the ergonomics of Go with the performance of Rust to provide real-time file system searching without relying on outdated databases.

## Core Problems with Traditional `locate`
1. **Stale Database**: Relies on `updatedb` which may be outdated
2. **Limited Scope**: Filters paths via `updatedb.conf`, missing important directories
3. **Name-only Search**: Cannot search by content, metadata, or advanced criteria
4. **No Live Recursion**: Only searches pre-indexed data

## Architecture

### Hybrid Go + Rust Design
- **Go Frontend**: CLI interface, argument parsing, result formatting
- **Rust Backend**: High-performance file system traversal and searching
- **Integration**: Use CGO bindings or subprocess execution for optimal performance

### Performance Strategy
```
Go CLI → Rust Core → File System
    ↓        ↓           ↓
 Parsing  Traversal   Results
    ↓        ↓           ↓
 Format   Filter     Display
```

## Features

### Basic Mode (Default)
- **Exact Match**: Find files with exact name matching
- **Access-Safe**: Skip directories without read permissions
- **Fast Traversal**: Real-time file system scanning
- **Smart Filtering**: Exclude obviously irrelevant system directories

### Advanced Mode
- **Pattern Matching**: Wildcard support (`*`, `?`, `[]`)
- **Partial Matching**: Substring and fuzzy matching
- **System Directory Control**: Configurable exclusion of kernel-related paths
- **Multi-threaded Search**: Parallel directory traversal

### Extended Functionality
- **File Extension Filtering**: `--ext .go,.rs,.md`
- **Size Filtering**: `--size +100M`, `--size -1K`
- **Modification Time**: `--mtime -7d`, `--mtime +1h`
- **Content Search**: `--content "pattern"` (optional)
- **Metadata Search**: File permissions, ownership, etc.

## CLI Interface Design

### Basic Usage
```bash
glocate filename.txt              # Exact match
glocate "*.go"                    # Pattern match
glocate --advanced pattern        # Advanced mode
```

### Advanced Usage
```bash
glocate --ext go --size +1M --mtime -7d "main"
glocate --exclude /proc,/sys --include /home,/opt "config"
glocate --content "TODO" --ext rs,go
```

## Technical Implementation

### Go Components
- CLI argument parsing (cobra/cli)
- Result formatting and display
- Configuration management
- Cross-platform compatibility layer

### Rust Components
- High-performance file system walker
- Pattern matching engine
- Parallel directory traversal
- Memory-efficient result filtering

### Integration Options
1. **CGO Bindings**: Direct function calls for maximum performance
2. **Subprocess**: Execute Rust binary and parse JSON output
3. **Shared Library**: Dynamic linking to Rust-compiled library

## Performance Goals
- **Speed**: 10x faster than traditional `locate` on cold searches
- **Memory**: Minimal memory footprint, streaming results
- **Scalability**: Handle millions of files efficiently
- **Responsiveness**: Show results as they are found

## Development Phases

### Phase 1: Core Implementation
- Basic file traversal in Rust
- Simple Go CLI wrapper
- Exact name matching
- Basic error handling

### Phase 2: Advanced Features
- Pattern matching
- File metadata filtering
- Performance optimization
- Comprehensive testing

### Phase 3: Extended Functionality
- Content searching
- Advanced CLI options
- Configuration files
- Cross-platform packaging

## Version Management
- Semantic versioning (SemVer)
- Go version tracking via `.go-version`
- Release version in `.release-version`
- Automated release pipeline

## Quality Assurance
- GolangCI linting configuration
- Comprehensive unit tests
- Integration tests
- Performance benchmarks
- Cross-platform compatibility testing

## Future Enhancements
- GUI interface
- Background indexing (optional)
- Network file system support
- Plugin architecture for custom filters
- Integration with popular editors/IDEs 