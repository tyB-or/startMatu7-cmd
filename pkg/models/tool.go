package models

import "time"

// OfflineTool 表示离线工具
type OfflineTool struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Path        string    `json:"path"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	Command     string    `json:"command"`
	URL         string    `json:"url"`
	Icon        string    `json:"icon"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	LastUsedAt  time.Time `json:"last_used_at"`
	UsageCount  int       `json:"usage_count"`
	NoteFile    string    `json:"note_file"`
	KeyPath     string    `json:"key_path"`
}

// WebTool 表示网页工具
type WebTool struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	Icon        string    `json:"icon"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Tags        []string  `json:"tags"`
	UsageCount  int       `json:"usage_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	LastUsedAt  time.Time `json:"last_used_at"`
	NoteFile    string    `json:"note_file"`
}

// Note 表示笔记
type Note struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	Source    string    `json:"source"`
	Tool      string    `json:"tool"`
	Tags      []string  `json:"tags"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
