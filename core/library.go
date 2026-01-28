package core

import (
	"fmt"
	"github.com/nbt4/capturelib/models"
	"os"
	"path/filepath"
)

// Library manages the capture file library
type Library struct {
	config  *Config
	db      *Database
	scanner *Scanner
}

// NewLibrary creates a new library manager
func NewLibrary(configPath string) (*Library, error) {
	// Load config
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Setup database path
	dbPath := filepath.Join(filepath.Dir(configPath), "library.db")

	// Open database
	db, err := NewDatabase(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create scanner
	scanner := NewScanner(db)

	lib := &Library{
		config:  config,
		db:      db,
		scanner: scanner,
	}

	// Auto-scan if library path is set
	if config.LibraryPath != "" && config.AutoScan {
		if _, err := os.Stat(config.LibraryPath); err == nil {
			lib.Scan()
		}
	}

	return lib, nil
}

// SetLibraryPath sets the library path and triggers a scan
func (l *Library) SetLibraryPath(path string) error {
	l.config.LibraryPath = path

	// Clear database
	if err := l.db.ClearAll(); err != nil {
		return fmt.Errorf("failed to clear database: %w", err)
	}

	// Scan new path
	_, err := l.Scan()
	return err
}

// Scan scans the library directory
func (l *Library) Scan() (int, error) {
	if l.config.LibraryPath == "" {
		return 0, fmt.Errorf("library path not set")
	}

	return l.scanner.ScanDirectory(l.config.LibraryPath, l.config.ScanSubdirectories)
}

// GetAllFiles returns all files
func (l *Library) GetAllFiles() ([]*models.CaptureFile, error) {
	return l.db.GetAllFiles()
}

// SearchFiles searches for files
func (l *Library) SearchFiles(query string) ([]*models.CaptureFile, error) {
	return l.db.SearchFiles(query)
}

// GetConfig returns the current config
func (l *Library) GetConfig() *Config {
	return l.config
}

// SaveConfig saves the current config
func (l *Library) SaveConfig(path string) error {
	return l.config.Save(path)
}

// GetFileCount returns the total file count
func (l *Library) GetFileCount() (int, error) {
	return l.db.GetFileCount()
}

// Close closes the library
func (l *Library) Close() error {
	return l.db.Close()
}
