# Testing Guide for go-locate

This document provides guidelines for writing and maintaining tests in the go-locate project.

## Testing Framework

The project uses **[testify](https://github.com/stretchr/testify)** v1.10.0 as the primary testing framework, which provides:
- Easy assertions with better error messages
- Test setup and teardown utilities  
- Structured test organization with sub-tests
- Mocking capabilities (when needed)

## Quick Start

### Basic Test Structure

```go
package search

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestFunctionName(t *testing.T) {
    // Setup
    input := "test input"
    expected := "expected output"

    // Execute
    result := FunctionToTest(input)

    // Assert
    assert.Equal(t, expected, result, "Function should return expected output")
}
```

### Table-Driven Tests

```go
func TestShouldExclude(t *testing.T) {
    searcher := &Searcher{
        config: &Config{
            Exclude: []string{"node_modules", ".git"},
        },
    }

    tests := []struct {
        name     string
        path     string
        expected bool
    }{
        {"regular file", "/home/user/project/file.go", false},
        {"proc filesystem", "/proc/cpuinfo", true},
        {"node_modules", "/home/user/node_modules/package", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := searcher.shouldExclude(tt.path)
            assert.Equal(t, tt.expected, result, "shouldExclude result should match expected")
        })
    }
}
```

## Assert vs Require

### Use `assert` for non-critical checks
- Test continues even if assertion fails
- Good for multiple validations in one test

```go
assert.Equal(t, expected, actual, "Values should match")
assert.True(t, condition, "Condition should be true")
assert.NotEmpty(t, slice, "Slice should not be empty")
```

### Use `require` for critical prerequisites
- Test stops immediately if assertion fails
- Good for setup validation and preventing panics

```go
require.NoError(t, err, "Setup must succeed")
require.NotNil(t, object, "Object must be created")
require.FileExists(t, filepath, "Test file must exist")
```

## Common Patterns

### File System Tests

```go
func TestFileOperation(t *testing.T) {
    // Create temporary directory
    tempDir, err := os.MkdirTemp("", "test_prefix")
    require.NoError(t, err, "Failed to create temp dir")

    // Clean up after test
    defer func() {
        if err := os.RemoveAll(tempDir); err != nil {
            t.Logf("Failed to remove temp dir: %v", err)
        }
    }()

    // Create test file
    testFile := filepath.Join(tempDir, "test.txt")
    err = os.WriteFile(testFile, []byte("test content"), 0644)
    require.NoError(t, err, "Failed to create test file")

    // Test your function
    result := ProcessFile(testFile)
    assert.True(t, result, "File should be processed successfully")
}
```

### Error Handling Tests

```go
func TestErrorCases(t *testing.T) {
    tests := []struct {
        name          string
        input         string
        expectedError string
    }{
        {"empty input", "", "input cannot be empty"},
        {"invalid format", "invalid", "invalid format"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := FunctionThatCanFail(tt.input)
            require.Error(t, err, "Should return error")
            assert.Contains(t, err.Error(), tt.expectedError, "Error message should match")
        })
    }
}
```

## Running Tests

### Local Development

```bash
# Run all tests with testify
make test

# Run tests with race detection
make test-race

# Run tests with coverage
make test-coverage

# Run only specific package
go test -v ./internal/search/

# Run specific test
go test -v -run TestFuzzyMatch ./internal/search/
```

### Continuous Integration

Tests run automatically in CI with:
- Race detection enabled
- Coverage reporting  
- Parallel execution for performance

## Best Practices

### 1. Test Naming
- Use descriptive names that explain the scenario
- Include expected behavior in test name
- Use sub-tests for multiple scenarios

```go
func TestSearcher_FuzzyMatch_ReturnsTrueForValidPatterns(t *testing.T) {
    // Test implementation
}
```

### 2. Test Organization
- One test file per source file (`search.go` → `search_test.go`)
- Group related tests in sub-tests
- Keep tests focused and independent

### 3. Constants for Test Data
Following .cursorrules, use named constants instead of magic values:

```go
const (
    TestPattern = "*.go"
    TestMaxResults = 100
    TestTimeout = 5 * time.Second
    ExpectedFileCount = 3
)

func TestSearch(t *testing.T) {
    config := &Config{
        Pattern:    TestPattern,
        MaxResults: TestMaxResults,
    }
    // Use constants instead of magic numbers
}
```

### 4. Helpful Assertions
- Always include descriptive messages
- Use specific assertions when available
- Check both positive and negative cases

```go
// Good - specific and descriptive
assert.Len(t, results, ExpectedFileCount, "Should find exactly 3 Go files")
assert.Contains(t, result.Path, "main.go", "Results should include main.go")

// Avoid - too generic
assert.True(t, len(results) == 3)
```

## Current Test Coverage

As of the latest run:
- **Overall**: 55.1% statement coverage
- **internal/search**: 55.1% with comprehensive fuzzy matching tests
- **internal/version**: 100% with full build info validation

## Future Enhancements

When the project grows, consider:
- **Mock objects** for external dependencies using `testify/mock`
- **Test suites** for complex integration scenarios using `testify/suite`
- **HTTP testing utilities** if web interfaces are added

## Contributing Test Code

When adding new tests:
1. Follow the existing patterns in the codebase
2. Use testify assertions and require statements
3. Include both positive and negative test cases
4. Add table-driven tests for multiple scenarios
5. Keep tests fast and focused
6. Use named constants instead of magic numbers
7. Include descriptive error messages

For questions or examples, refer to existing tests in:
- `internal/search/search_test.go` - Core functionality tests
- `internal/version/version_test.go` - Build information tests
- `internal/search/search_bench_test.go` - Performance benchmarks
