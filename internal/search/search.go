// Package search provides file search functionality with various filtering options
package search

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Constants
const (
	// Default buffer size for results channel
	defaultResultsBufferSize = 100
)

// Config holds search configuration
type Config struct {
	Pattern     string
	Advanced    bool
	Extensions  []string
	Size        string
	Mtime       string
	Content     string
	Exclude     []string
	Include     []string
	Threads     int
	Depth       int
	FollowLinks bool
	MaxResults  int
	Verbose     bool
}

// Result represents a search result
type Result struct {
	Path    string    `json:"path"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mod_time"`
	IsDir   bool      `json:"is_dir"`
	Mode    string    `json:"mode"`
}

// Searcher performs file searches
type Searcher struct {
	config  *Config
	results chan *Result
	done    chan struct{}
	wg      sync.WaitGroup
}

// New creates a new searcher instance
func New(config *Config) (*Searcher, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	if config.Pattern == "" {
		return nil, fmt.Errorf("search pattern cannot be empty")
	}

	// Set default thread count
	if config.Threads <= 0 {
		config.Threads = runtime.NumCPU()
	}

	return &Searcher{
		config:  config,
		results: make(chan *Result, defaultResultsBufferSize),
		done:    make(chan struct{}),
	}, nil
}

// Search performs the file search
func (s *Searcher) Search() ([]*Result, error) {
	var results []*Result
	var mu sync.Mutex

	// Start result collector
	go func() {
		for result := range s.results {
			mu.Lock()
			if len(results) < s.config.MaxResults {
				results = append(results, result)
			}
			mu.Unlock()
		}
	}()

	// Determine search roots
	searchRoots := s.getSearchRoots()

	// Start search workers
	for _, root := range searchRoots {
		s.wg.Add(1)
		go s.searchWorker(root)
	}

	// Wait for all workers to complete
	s.wg.Wait()
	close(s.results)

	return results, nil
}

// getSearchRoots returns the directories to search
func (s *Searcher) getSearchRoots() []string {
	if len(s.config.Include) > 0 {
		return s.config.Include
	}

	// Default search roots
	roots := []string{"/"}

	// On macOS, also search common user directories
	if runtime.GOOS == "darwin" {
		if home, err := os.UserHomeDir(); err == nil {
			roots = append(roots, home)
		}
		roots = append(roots, "/Applications", "/usr/local")
	}

	return roots
}

// searchWorker performs search in a specific directory tree
func (s *Searcher) searchWorker(root string) {
	defer s.wg.Done()

	// Add support for parallel search
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Skip directories we can't access
			if s.config.Verbose {
				fmt.Fprintf(os.Stderr, "Warning: Cannot access %s: %v\n", path, err)
			}
			return nil
		}

		// Check if we should exclude this path
		if s.shouldExclude(path) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Check depth limit
		if s.config.Depth > 0 {
			depth := strings.Count(strings.TrimPrefix(path, root), string(os.PathSeparator))
			if depth > s.config.Depth {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// Check if this matches our search criteria
		if s.matches(path, info) {
			result := &Result{
				Path:    path,
				Size:    info.Size(),
				ModTime: info.ModTime(),
				IsDir:   info.IsDir(),
				Mode:    info.Mode().String(),
			}

			select {
			case s.results <- result:
			case <-s.done:
				return fmt.Errorf("search canceled")
			}
		}

		return nil
	})

	if err != nil && s.config.Verbose {
		fmt.Fprintf(os.Stderr, "Warning: Error walking %s: %v\n", root, err)
	}
}

// shouldExclude checks if a path should be excluded
func (s *Searcher) shouldExclude(path string) bool {
	for _, exclude := range s.config.Exclude {
		if strings.Contains(path, exclude) {
			return true
		}
	}

	// Default exclusions for system directories
	systemDirs := []string{"/proc", "/sys", "/dev", "/tmp"}
	for _, sysDir := range systemDirs {
		if strings.HasPrefix(path, sysDir) {
			return true
		}
	}

	return false
}

// matches checks if a file matches the search criteria
func (s *Searcher) matches(path string, _ os.FileInfo) bool {
	filename := filepath.Base(path)

	// Pattern matching
	if s.config.Advanced {
		// Fuzzy matching - check if pattern characters appear in order
		if !s.fuzzyMatch(strings.ToLower(filename), strings.ToLower(s.config.Pattern)) {
			return false
		}
	} else {
		// Exact or wildcard matching
		matched, err := filepath.Match(s.config.Pattern, filename)
		if err != nil || !matched {
			// Try case-insensitive match
			matched, _ = filepath.Match(strings.ToLower(s.config.Pattern), strings.ToLower(filename))
			if !matched {
				return false
			}
		}
	}

	// Extension filtering
	if len(s.config.Extensions) > 0 {
		ext := strings.TrimPrefix(filepath.Ext(filename), ".")
		found := false
		for _, allowedExt := range s.config.Extensions {
			if strings.EqualFold(ext, allowedExt) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Size filtering (basic implementation)
	// TODO: Implement size filtering
	// For now, we accept all files regardless of size
	_ = s.config.Size

	// Modification time filtering (basic implementation)
	// TODO: Implement mtime filtering
	// For now, we accept all files regardless of modification time
	_ = s.config.Mtime

	return true
}

// fuzzyMatch performs fuzzy string matching
func (s *Searcher) fuzzyMatch(text, pattern string) bool {
	if pattern == "" {
		return true
	}
	if text == "" {
		return false
	}

	textIdx := 0
	for _, patternChar := range pattern {
		found := false
		for textIdx < len(text) {
			if rune(text[textIdx]) == patternChar {
				found = true
				textIdx++
				break
			}
			textIdx++
		}
		if !found {
			return false
		}
	}
	return true
}
