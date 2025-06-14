package search

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	config := &Config{
		Pattern:    "test",
		MaxResults: 100,
	}

	searcher, err := New(config)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if searcher == nil {
		t.Fatal("Expected searcher to be created")
	}

	if searcher.config.Pattern != "test" {
		t.Errorf("Expected pattern 'test', got '%s'", searcher.config.Pattern)
	}
}

func TestNewWithNilConfig(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("Expected error for nil config")
	}
}

func TestNewWithEmptyPattern(t *testing.T) {
	config := &Config{
		Pattern: "",
	}

	_, err := New(config)
	if err == nil {
		t.Fatal("Expected error for empty pattern")
	}
}

func TestShouldExclude(t *testing.T) {
	searcher := &Searcher{
		config: &Config{
			Exclude: []string{"node_modules", ".git"},
		},
	}

	tests := []struct {
		path     string
		expected bool
	}{
		{"/home/user/project/file.go", false},
		{"/proc/cpuinfo", true},
		{"/sys/devices", true},
		{"/dev/null", true},
		{"/tmp/test", true},
		{"/home/user/node_modules/package", true},
		{"/home/user/.git/config", true},
	}

	for _, test := range tests {
		result := searcher.shouldExclude(test.path)
		if result != test.expected {
			t.Errorf("shouldExclude(%s) = %v, expected %v", test.path, result, test.expected)
		}
	}
}

func TestFuzzyMatch(t *testing.T) {
	searcher := &Searcher{}

	tests := []struct {
		text     string
		pattern  string
		expected bool
	}{
		{"hello_world.txt", "hw", true},
		{"main.go", "mgo", true},
		{"test.rs", "xyz", false},
		{"config.yaml", "cfg", true},
		{"README.md", "rm", false}, // R-E-A-D-M-E doesn't contain 'r' followed by 'm'
		{"", "test", false},
		{"test", "", true},
	}

	for _, test := range tests {
		result := searcher.fuzzyMatch(test.text, test.pattern)
		if result != test.expected {
			t.Errorf("fuzzyMatch(%s, %s) = %v, expected %v", test.text, test.pattern, result, test.expected)
		}
	}
}

func TestMatches(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "glocate_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

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
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
		file.Close()
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
	if err != nil {
		t.Fatalf("Failed to stat test file: %v", err)
	}

	if !searcher.matches(testGoPath, info) {
		t.Error("Expected test.go to match *.go pattern")
	}

	// Test with different pattern
	searcher.config.Pattern = "*.yaml"
	searcher.config.Extensions = []string{"yaml"}

	testYamlPath := filepath.Join(tempDir, "config.yaml")
	info, err = os.Stat(testYamlPath)
	if err != nil {
		t.Fatalf("Failed to stat yaml file: %v", err)
	}

	if !searcher.matches(testYamlPath, info) {
		t.Error("Expected config.yaml to match *.yaml pattern")
	}
}

func TestGetSearchRoots(t *testing.T) {
	// Test with include directories
	searcher := &Searcher{
		config: &Config{
			Include: []string{"/home", "/opt"},
		},
	}

	roots := searcher.getSearchRoots()
	if len(roots) != 2 {
		t.Errorf("Expected 2 search roots, got %d", len(roots))
	}

	if roots[0] != "/home" || roots[1] != "/opt" {
		t.Errorf("Expected [/home, /opt], got %v", roots)
	}

	// Test without include directories (should use defaults)
	searcher.config.Include = []string{}
	roots = searcher.getSearchRoots()
	if len(roots) == 0 {
		t.Error("Expected default search roots")
	}
}
