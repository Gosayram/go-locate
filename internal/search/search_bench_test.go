package search

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// BenchmarkAdvancedFuzzyMatch tests the performance of fuzzy matching with various patterns
func BenchmarkAdvancedFuzzyMatch(b *testing.B) {
	searcher := &Searcher{}

	tests := []struct {
		name    string
		pattern string
		text    string
	}{
		{"Short", "go", "main.go"},
		{"Medium", "config", "application-config.yaml"},
		{"Long", "search", "very-long-filename-with-search-functionality.go"},
		{"Complex", "glocate", "go-locate-modern-file-search-tool.exe"},
		{"NoMatch", "xyz", "completely-different-filename.txt"},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				searcher.fuzzyMatch(tt.text, tt.pattern)
			}
		})
	}
}

// BenchmarkAdvancedShouldExclude tests directory exclusion with more paths
func BenchmarkAdvancedShouldExclude(b *testing.B) {
	searcher := &Searcher{
		config: &Config{
			Exclude: []string{"/proc", "/sys", "/dev", "/tmp", "node_modules", ".git", "target", "build", "dist"},
		},
	}

	paths := []string{
		"/home/user/documents",
		"/proc/cpuinfo",
		"/sys/devices",
		"/tmp/temp-file",
		"/home/user/project/node_modules/package",
		"/home/user/project/.git/config",
		"/usr/local/bin",
		"/home/user/project/target/debug",
		"/home/user/project/build/output",
		"/home/user/project/dist/bundle.js",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range paths {
			searcher.shouldExclude(path)
		}
	}
}

// BenchmarkAdvancedMatches tests file matching with different configurations
func BenchmarkAdvancedMatches(b *testing.B) {
	searcher := &Searcher{
		config: &Config{
			Pattern:    "*.go",
			Extensions: []string{"go", "rs", "py", "js", "ts"},
		},
	}

	files := []string{
		"main.go",
		"config.rs",
		"script.py",
		"app.js",
		"types.ts",
		"document.txt",
		"image.png",
		"archive.zip",
		"Makefile",
		"README.md",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, file := range files {
			searcher.matches(file, nil)
		}
	}
}

// BenchmarkSearchRoots tests search root determination performance
func BenchmarkSearchRoots(b *testing.B) {
	searcher := &Searcher{
		config: &Config{
			Include: []string{"/home", "/usr/local", "/opt", "/var", "/etc"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		searcher.getSearchRoots()
	}
}

// BenchmarkConcurrentSearch tests concurrent search performance
func BenchmarkConcurrentSearch(b *testing.B) {
	tempDir := b.TempDir()

	// Create files for concurrent testing
	for i := 0; i < 20; i++ {
		dir := filepath.Join(tempDir, fmt.Sprintf("dir%d", i))
		os.MkdirAll(dir, 0755)

		for j := 0; j < 10; j++ {
			file := filepath.Join(dir, fmt.Sprintf("test%d.go", j))
			os.WriteFile(file, []byte("package main\nfunc main() {}"), 0644)
		}
	}

	config := &Config{
		Pattern: "test*.go",
		Include: []string{tempDir},
		Threads: 4,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		searcher, _ := New(config)
		results, _ := searcher.Search()

		// Consume all results
		count := 0
		for range results {
			count++
		}
	}
}

// BenchmarkPatternMatching tests different pattern matching scenarios
func BenchmarkPatternMatching(b *testing.B) {
	patterns := []struct {
		name    string
		pattern string
		files   []string
	}{
		{
			name:    "Wildcard",
			pattern: "*.go",
			files:   []string{"main.go", "config.go", "test.py", "doc.md"},
		},
		{
			name:    "Exact",
			pattern: "main.go",
			files:   []string{"main.go", "main.py", "test.go", "config.go"},
		},
		{
			name:    "Complex",
			pattern: "test_*.go",
			files:   []string{"test_main.go", "test_config.go", "main_test.go", "test.py"},
		},
		{
			name:    "MultipleWildcards",
			pattern: "*_test_*.go",
			files:   []string{"unit_test_main.go", "integration_test_api.go", "test_helper.go", "main.go"},
		},
	}

	for _, tt := range patterns {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, file := range tt.files {
					filepath.Match(tt.pattern, file)
				}
			}
		})
	}
}

// BenchmarkMemoryAllocation tests memory allocation patterns
func BenchmarkMemoryAllocation(b *testing.B) {
	config := &Config{
		Pattern:    "*.go",
		Extensions: []string{"go"},
		MaxResults: 1000,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		searcher, _ := New(config)

		// Simulate creating results
		for j := 0; j < 100; j++ {
			result := &Result{
				Path:    fmt.Sprintf("/path/to/file%d.go", j),
				Size:    1024,
				ModTime: time.Now(),
				IsDir:   false,
				Mode:    "-rw-r--r--",
			}
			_ = result
		}

		_ = searcher
	}
}

// BenchmarkLargeDirectorySearch tests performance with many directories
func BenchmarkLargeDirectorySearch(b *testing.B) {
	tempDir := b.TempDir()

	// Create a large directory structure
	for i := 0; i < 50; i++ {
		dir := filepath.Join(tempDir, fmt.Sprintf("level1_%d", i))
		os.MkdirAll(dir, 0755)

		for j := 0; j < 5; j++ {
			subdir := filepath.Join(dir, fmt.Sprintf("level2_%d", j))
			os.MkdirAll(subdir, 0755)

			for k := 0; k < 3; k++ {
				file := filepath.Join(subdir, fmt.Sprintf("file_%d.go", k))
				os.WriteFile(file, []byte("package main"), 0644)
			}
		}
	}

	config := &Config{
		Pattern:    "*.go",
		Include:    []string{tempDir},
		MaxResults: 1000,
		Threads:    4,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		searcher, _ := New(config)
		results, _ := searcher.Search()

		// Count results to avoid too much overhead
		count := 0
		for range results {
			count++
			if count > 500 { // Limit to avoid too much overhead
				break
			}
		}
	}
}

// BenchmarkExtensionFiltering tests extension-based filtering performance
func BenchmarkExtensionFiltering(b *testing.B) {
	searcher := &Searcher{
		config: &Config{
			Pattern:    "*",
			Extensions: []string{"go", "rs", "py", "js", "ts", "java", "cpp", "c", "h"},
		},
	}

	files := []string{
		"main.go", "config.rs", "script.py", "app.js", "types.ts",
		"Main.java", "program.cpp", "header.h", "source.c",
		"document.txt", "image.png", "data.json", "style.css",
		"README.md", "Makefile", "Dockerfile", "config.yaml",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, file := range files {
			searcher.matches(file, nil)
		}
	}
}
