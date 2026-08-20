package mocks

import (
	"lets-go-edwards/internal/models"
	"time"
)

var mockSnippet = models.Snippet{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

// DB wrapper now doesn't really link to any DB
type SnippetModel struct{}

// Fake insert entry to DB and return a fake ID of entry
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	return 2, nil
}

// Get fake Snippet
func (m *SnippetModel) Get(id int) (models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return models.Snippet{}, models.ErrNoRecord
	}
}

// Get latest Snippets (freshly clones from Fake Snippet)
func (m *SnippetModel) Latest() ([]models.Snippet, error) {
	return []models.Snippet{mockSnippet}, nil
}
