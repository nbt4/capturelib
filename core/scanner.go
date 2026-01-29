package core

import (
	"github.com/nbt4/capturelib/models"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Scanner scans directories for .c2o files
type Scanner struct {
	db *Database
}

// NewScanner creates a new scanner
func NewScanner(db *Database) *Scanner {
	return &Scanner{db: db}
}

// ScanDirectory scans a directory for .c2o files
func (s *Scanner) ScanDirectory(path string, recursive bool) (int, error) {
	count := 0

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files we can't access
		}

		// Skip directories if not recursive
		if info.IsDir() {
			if !recursive && filePath != path {
				return filepath.SkipDir
			}
			return nil
		}

		// Only process Capture files (.c2o, .c2s, .c2p)
		lowerName := strings.ToLower(info.Name())
		if !strings.HasSuffix(lowerName, ".c2o") && 
		   !strings.HasSuffix(lowerName, ".c2s") && 
		   !strings.HasSuffix(lowerName, ".c2p") {
			return nil
		}

		// Create file model
		file := &models.CaptureFile{
			Filename:   info.Name(),
			Filepath:   filePath,
			Size:       info.Size(),
			ModifiedAt: info.ModTime(),
			CreatedAt:  getCreationTime(info),
		}

		// Add to database
		if err := s.db.AddFile(file); err != nil {
			return err
		}

		count++
		return nil
	})

	return count, err
}

// getCreationTime extracts creation time from FileInfo
// Falls back to ModTime if creation time is not available
func getCreationTime(info os.FileInfo) time.Time {
	// On most systems, we can't easily get creation time
	// So we use ModTime as fallback
	return info.ModTime()
}
