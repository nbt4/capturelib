package models

import "time"

// CaptureFile represents a .c2o stage design file
type CaptureFile struct {
	ID         int64     `json:"id"`
	Filename   string    `json:"filename"`
	Filepath   string    `json:"filepath"`
	Size       int64     `json:"size"`
	ModifiedAt time.Time `json:"modified_at"`
	CreatedAt  time.Time `json:"created_at"`
	IndexedAt  time.Time `json:"indexed_at"`
}

// Tag represents a category/tag for organizing files
type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
