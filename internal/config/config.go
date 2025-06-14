// Package config provides configuration management for go-locate
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Search SearchConfig `mapstructure:"search"`
	Output OutputConfig `mapstructure:"output"`
}

// SearchConfig holds search-related configuration
type SearchConfig struct {
	ExcludeDirs    []string `mapstructure:"exclude_dirs"`
	IncludeDirs    []string `mapstructure:"include_dirs"`
	MaxDepth       int      `mapstructure:"max_depth"`
	FollowSymlinks bool     `mapstructure:"follow_symlinks"`
	DefaultThreads int      `mapstructure:"default_threads"`
}

// OutputConfig holds output-related configuration
type OutputConfig struct {
	Format     string `mapstructure:"format"`
	Color      bool   `mapstructure:"color"`
	MaxResults int    `mapstructure:"max_results"`
}

var (
	cfg        *Config
	configFile string
)

// SetConfigFile sets the config file path
func SetConfigFile(file string) {
	configFile = file
}

// Load loads the configuration from file
func Load() error {
	viper.SetConfigName(".glocate")
	viper.SetConfigType("toml")

	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		// Add config search paths
		home, err := os.UserHomeDir()
		if err == nil {
			viper.AddConfigPath(home)
		}
		viper.AddConfigPath(".")
	}

	// Set defaults
	setDefaults()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, use defaults
			return nil
		}
		return fmt.Errorf("error reading config file: %w", err)
	}

	// Unmarshal config
	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("error unmarshaling config: %w", err)
	}

	return nil
}

// Get returns the current configuration
func Get() *Config {
	if cfg == nil {
		cfg = &Config{}
		setDefaults()
	}
	return cfg
}

const (
	// DefaultMaxDepth is the default maximum search depth
	DefaultMaxDepth = 20
	// DefaultMaxResults is the default maximum number of results
	DefaultMaxResults = 1000
)

// setDefaults sets default configuration values
func setDefaults() {
	// Search defaults
	viper.SetDefault("search.exclude_dirs", []string{"/proc", "/sys", "/dev", "/tmp"})
	viper.SetDefault("search.include_dirs", []string{})
	viper.SetDefault("search.max_depth", DefaultMaxDepth)
	viper.SetDefault("search.follow_symlinks", false)
	viper.SetDefault("search.default_threads", 0) // 0 means use CPU count

	// Output defaults
	viper.SetDefault("output.format", "path")
	viper.SetDefault("output.color", true)
	viper.SetDefault("output.max_results", DefaultMaxResults)
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	if configFile != "" {
		return configFile
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ".glocate.toml"
	}

	return filepath.Join(home, ".glocate.toml")
}
