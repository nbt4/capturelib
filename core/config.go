package core

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	LibraryPath         string `json:"library_path"`
	Theme               string `json:"theme"` // "dark" or "light"
	AutoScan            bool   `json:"auto_scan"`
	ScanSubdirectories  bool   `json:"scan_subdirectories"`
	WindowWidth         int    `json:"window_width"`
	WindowHeight        int    `json:"window_height"`
}

// DefaultConfig returns a new config with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		LibraryPath:        "",
		Theme:              "dark",
		AutoScan:           true,
		ScanSubdirectories: true,
		WindowWidth:        1200,
		WindowHeight:       800,
	}
}

// LoadConfig loads config from file or returns default
func LoadConfig(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Save saves the config to file
func (c *Config) Save(path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
