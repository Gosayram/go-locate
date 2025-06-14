// Package output provides result formatting and display functionality
package output

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Gosayram/go-locate/internal/search"
)

// Config holds output configuration
type Config struct {
	Format  string
	Verbose bool
	Color   bool
}

// Formatter handles result output formatting
type Formatter struct {
	config *Config
}

// New creates a new formatter instance
func New(config *Config) *Formatter {
	if config == nil {
		config = &Config{
			Format: "path",
			Color:  true,
		}
	}
	return &Formatter{config: config}
}

// Print outputs the search results
func (f *Formatter) Print(results []*search.Result) error {
	if len(results) == 0 {
		if f.config.Verbose {
			fmt.Println("No results found")
		}
		return nil
	}

	switch f.config.Format {
	case "json":
		return f.printJSON(results)
	case "detailed":
		return f.printDetailed(results)
	default:
		return f.printPath(results)
	}
}

// printPath prints only the file paths
func (f *Formatter) printPath(results []*search.Result) error {
	for _, result := range results {
		if f.config.Color && result.IsDir {
			fmt.Printf("\033[34m%s\033[0m\n", result.Path) // Blue for directories
		} else {
			fmt.Println(result.Path)
		}
	}
	return nil
}

// printDetailed prints detailed information about each result
func (f *Formatter) printDetailed(results []*search.Result) error {
	for _, result := range results {
		var typeStr string
		if result.IsDir {
			typeStr = "DIR"
		} else {
			typeStr = "FILE"
		}

		sizeStr := f.formatSize(result.Size)
		timeStr := result.ModTime.Format("2006-01-02 15:04:05")

		if f.config.Color {
			if result.IsDir {
				fmt.Printf("\033[34m%-4s\033[0m %8s %s \033[34m%s\033[0m\n",
					typeStr, sizeStr, timeStr, result.Path)
			} else {
				fmt.Printf("%-4s %8s %s %s\n",
					typeStr, sizeStr, timeStr, result.Path)
			}
		} else {
			fmt.Printf("%-4s %8s %s %s\n",
				typeStr, sizeStr, timeStr, result.Path)
		}
	}
	return nil
}

// printJSON prints results in JSON format
func (f *Formatter) printJSON(results []*search.Result) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	output := map[string]interface{}{
		"results":   results,
		"count":     len(results),
		"timestamp": time.Now().Format(time.RFC3339),
	}

	return encoder.Encode(output)
}

// formatSize formats file size in human-readable format
func (f *Formatter) formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%dB", size)
	}

	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"K", "M", "G", "T", "P", "E"}
	return fmt.Sprintf("%.1f%s", float64(size)/float64(div), units[exp])
}

// PrintSummary prints a summary of the search operation
func (f *Formatter) PrintSummary(results []*search.Result, duration time.Duration) {
	if !f.config.Verbose {
		return
	}

	fileCount := 0
	dirCount := 0

	for _, result := range results {
		if result.IsDir {
			dirCount++
		} else {
			fileCount++
		}
	}

	fmt.Fprintf(os.Stderr, "\nSearch completed in %v\n", duration)
	fmt.Fprintf(os.Stderr, "Found %d files and %d directories\n", fileCount, dirCount)
}

// PrintError prints an error message
func (f *Formatter) PrintError(err error) {
	if f.config.Color {
		fmt.Fprintf(os.Stderr, "\033[31mError:\033[0m %v\n", err)
	} else {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

// PrintWarning prints a warning message
func (f *Formatter) PrintWarning(msg string) {
	if !f.config.Verbose {
		return
	}

	if f.config.Color {
		fmt.Fprintf(os.Stderr, "\033[33mWarning:\033[0m %s\n", msg)
	} else {
		fmt.Fprintf(os.Stderr, "Warning: %s\n", msg)
	}
}

// GetResultStats returns statistics about the results
func (f *Formatter) GetResultStats(results []*search.Result) map[string]interface{} {
	stats := make(map[string]interface{})

	fileCount := 0
	dirCount := 0
	totalSize := int64(0)
	extensions := make(map[string]int)

	for _, result := range results {
		if result.IsDir {
			dirCount++
		} else {
			fileCount++
			totalSize += result.Size

			ext := strings.ToLower(filepath.Ext(result.Path))
			if ext != "" {
				extensions[ext]++
			}
		}
	}

	stats["total_results"] = len(results)
	stats["file_count"] = fileCount
	stats["directory_count"] = dirCount
	stats["total_size"] = totalSize
	stats["total_size_formatted"] = f.formatSize(totalSize)
	stats["extensions"] = extensions

	return stats
}
