package core

import (
	"database/sql"
	"github.com/nbt4/capturelib/models"
	_ "modernc.org/sqlite"
	"time"
)

type Database struct {
	db *sql.DB
}

// NewDatabase creates a new database connection
func NewDatabase(path string) (*Database, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	// Initialize schema
	if err := initSchema(db); err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func initSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		filename TEXT NOT NULL,
		filepath TEXT UNIQUE NOT NULL,
		size INTEGER,
		modified_at DATETIME,
		created_at DATETIME,
		indexed_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE INDEX IF NOT EXISTS idx_filename ON files(filename);
	CREATE INDEX IF NOT EXISTS idx_filepath ON files(filepath);
	`

	_, err := db.Exec(schema)
	return err
}

// AddFile adds a file to the database
func (d *Database) AddFile(file *models.CaptureFile) error {
	query := `
		INSERT INTO files (filename, filepath, size, modified_at, created_at, indexed_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(filepath) DO UPDATE SET
			filename = excluded.filename,
			size = excluded.size,
			modified_at = excluded.modified_at,
			indexed_at = excluded.indexed_at
	`

	_, err := d.db.Exec(query,
		file.Filename,
		file.Filepath,
		file.Size,
		file.ModifiedAt,
		file.CreatedAt,
		time.Now(),
	)

	return err
}

// GetAllFiles returns all files from the database
func (d *Database) GetAllFiles() ([]*models.CaptureFile, error) {
	query := `SELECT id, filename, filepath, size, modified_at, created_at, indexed_at FROM files ORDER BY filename`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*models.CaptureFile
	for rows.Next() {
		file := &models.CaptureFile{}
		err := rows.Scan(
			&file.ID,
			&file.Filename,
			&file.Filepath,
			&file.Size,
			&file.ModifiedAt,
			&file.CreatedAt,
			&file.IndexedAt,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

// SearchFiles searches for files by filename
func (d *Database) SearchFiles(query string) ([]*models.CaptureFile, error) {
	sqlQuery := `
		SELECT id, filename, filepath, size, modified_at, created_at, indexed_at 
		FROM files 
		WHERE filename LIKE ? 
		ORDER BY filename
	`

	rows, err := d.db.Query(sqlQuery, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*models.CaptureFile
	for rows.Next() {
		file := &models.CaptureFile{}
		err := rows.Scan(
			&file.ID,
			&file.Filename,
			&file.Filepath,
			&file.Size,
			&file.ModifiedAt,
			&file.CreatedAt,
			&file.IndexedAt,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

// DeleteFile removes a file from the database
func (d *Database) DeleteFile(filepath string) error {
	_, err := d.db.Exec("DELETE FROM files WHERE filepath = ?", filepath)
	return err
}

// GetFileCount returns the total number of files
func (d *Database) GetFileCount() (int, error) {
	var count int
	err := d.db.QueryRow("SELECT COUNT(*) FROM files").Scan(&count)
	return count, err
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// ClearAll removes all files from the database
func (d *Database) ClearAll() error {
	_, err := d.db.Exec("DELETE FROM files")
	return err
}
