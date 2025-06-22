package search

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	config := &Config{
		Pattern:    "test",
		MaxResults: 100,
	}

	searcher, err := New(config)
	require.NoError(t, err, "Expected no error creating searcher")
	require.NotNil(t, searcher, "Expected searcher to be created")
	assert.Equal(t, "test", searcher.config.Pattern, "Expected pattern to match")
}

func TestNewWithNilConfig(t *testing.T) {
	_, err := New(nil)
	assert.Error(t, err, "Expected error for nil config")
}

func TestNewWithEmptyPattern(t *testing.T) {
	config := &Config{
		Pattern: "",
	}

	_, err := New(config)
	assert.Error(t, err, "Expected error for empty pattern")
}

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
		{"sys filesystem", "/sys/devices", true},
		{"dev filesystem", "/dev/null", true},
		{"tmp directory", "/tmp/test", true},
		{"node_modules", "/home/user/node_modules/package", true},
		{"git directory", "/home/user/.git/config", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := searcher.shouldExclude(tt.path)
			assert.Equal(t, tt.expected, result, "shouldExclude result should match expected")
		})
	}
}

func TestFuzzyMatch(t *testing.T) {
	searcher := &Searcher{}

	tests := []struct {
		name     string
		text     string
		pattern  string
		expected bool
	}{
		{"hello world match", "hello_world.txt", "hw", true},
		{"main go match", "main.go", "mgo", true},
		{"no match", "test.rs", "xyz", false},
		{"config match", "config.yaml", "cfg", true},
		{"readme no match", "README.md", "rm", false}, // R-E-A-D-M-E doesn't contain 'r' followed by 'm'
		{"empty text", "", "test", false},
		{"empty pattern", "test", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := searcher.fuzzyMatch(tt.text, tt.pattern)
			assert.Equal(t, tt.expected, result, "fuzzyMatch result should match expected")
		})
	}
}

func TestMatches(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "glocate_test")
	require.NoError(t, err, "Failed to create temp dir")

	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	// Create test files
	testFiles := []string{
		"test.go",
		"main.go",
		"config.yaml",
		"README.md",
	}

	for _, filename := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		file, err := os.Create(filePath)
		require.NoError(t, err, "Failed to create test file %s", filename)
		require.NoError(t, file.Close(), "Failed to close test file %s", filename)
	}

	searcher := &Searcher{
		config: &Config{
			Pattern:    "*.go",
			Extensions: []string{"go"},
		},
	}

	// Test pattern matching
	testGoPath := filepath.Join(tempDir, "test.go")
	info, err := os.Stat(testGoPath)
	require.NoError(t, err, "Failed to stat test file")

	assert.True(t, searcher.matches(testGoPath, info), "Expected test.go to match *.go pattern")

	// Test with different pattern
	searcher.config.Pattern = "*.yaml"
	searcher.config.Extensions = []string{"yaml"}

	testYamlPath := filepath.Join(tempDir, "config.yaml")
	info, err = os.Stat(testYamlPath)
	require.NoError(t, err, "Failed to stat yaml file")

	assert.True(t, searcher.matches(testYamlPath, info), "Expected config.yaml to match *.yaml pattern")
}

func TestGetSearchRoots(t *testing.T) {
	// Test with include directories
	searcher := &Searcher{
		config: &Config{
			Include: []string{"/home", "/opt"},
		},
	}

	roots := searcher.getSearchRoots()
	assert.Len(t, roots, 2, "Expected 2 search roots")
	assert.Equal(t, []string{"/home", "/opt"}, roots, "Expected specific search roots")

	// Test without include directories (should use defaults)
	searcher.config.Include = []string{}
	roots = searcher.getSearchRoots()
	assert.NotEmpty(t, roots, "Expected default search roots")
}
