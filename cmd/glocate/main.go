// Package main provides the command-line interface for go-locate
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Gosayram/go-locate/internal/config"
	"github.com/Gosayram/go-locate/internal/output"
	"github.com/Gosayram/go-locate/internal/search"
	"github.com/Gosayram/go-locate/internal/version"
)

// rootCmd is the root command for glocate
var rootCmd = &cobra.Command{
	Use:   "glocate [pattern]",
	Short: "Modern file search tool",
	Long: `glocate is a modern, fast file search tool that replaces the outdated locate command.
It provides real-time file system searching without relying on outdated databases.`,
	Version: version.GetVersion(),
	Args:    cobra.MaximumNArgs(1),
	RunE:    runSearch,
}

// versionCmd is the command to show the version information
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  "Display detailed version and build information for glocate",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println(version.GetFullVersionInfo())
	},
}

var (
	cfgFile     string
	advanced    bool
	extensions  []string
	size        string
	mtime       string
	content     string
	exclude     []string
	include     []string
	threads     int
	depth       int
	followLinks bool
	format      string
	maxResults  int
	verbose     bool
)

func init() {
	cobra.OnInitialize(initConfig)

	// Add version command
	rootCmd.AddCommand(versionCmd)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.glocate.toml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Search flags
	rootCmd.Flags().BoolVar(&advanced, "advanced", false, "enable advanced search mode with fuzzy matching")
	rootCmd.Flags().StringSliceVar(&extensions, "ext", []string{}, "filter by file extensions (comma-separated)")
	rootCmd.Flags().StringVar(&size, "size", "", "filter by file size (+100M, -1K)")
	rootCmd.Flags().StringVar(&mtime, "mtime", "", "filter by modification time (-7d, +1h)")
	rootCmd.Flags().StringVar(&content, "content", "", "search file content")
	rootCmd.Flags().StringSliceVar(&exclude, "exclude", []string{}, "exclude directories")
	rootCmd.Flags().StringSliceVar(&include, "include", []string{}, "include directories")

	// Performance flags
	rootCmd.Flags().IntVar(&threads, "threads", 0, "number of threads (default: CPU cores)")
	rootCmd.Flags().IntVar(&depth, "depth", 0, "maximum search depth (0 = unlimited)")
	rootCmd.Flags().BoolVar(&followLinks, "follow-symlinks", false, "follow symbolic links")

	// Output flags
	rootCmd.Flags().StringVar(&format, "format", "path", "output format (path, detailed, json)")
	rootCmd.Flags().IntVar(&maxResults, "max-results", config.DefaultMaxResults, "maximum number of results")
}

func initConfig() {
	if cfgFile != "" {
		config.SetConfigFile(cfgFile)
	}

	if err := config.Load(); err != nil {
		if verbose {
			fmt.Fprintf(os.Stderr, "Warning: Could not load config: %v\n", err)
		}
	}
}

func runSearch(_ *cobra.Command, args []string) error {
	var pattern string
	if len(args) > 0 {
		pattern = args[0]
	}

	if pattern == "" {
		return fmt.Errorf("search pattern is required")
	}

	// Create search configuration
	searchConfig := &search.Config{
		Pattern:     pattern,
		Advanced:    advanced,
		Extensions:  extensions,
		Size:        size,
		Mtime:       mtime,
		Content:     content,
		Exclude:     exclude,
		Include:     include,
		Threads:     threads,
		Depth:       depth,
		FollowLinks: followLinks,
		MaxResults:  maxResults,
		Verbose:     verbose,
	}

	// Create searcher
	searcher, err := search.New(searchConfig)
	if err != nil {
		return fmt.Errorf("failed to create searcher: %w", err)
	}

	// Perform search
	results, err := searcher.Search()
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	// Output results
	outputConfig := &output.Config{
		Format:  format,
		Verbose: verbose,
	}

	formatter := output.New(outputConfig)
	return formatter.Print(results)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
