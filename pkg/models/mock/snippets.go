package mock

import (
	"time"

	"github.com/Onizukachi/snippetbox/pkg/models"
)

type SnippetModel struct{}

var mockSnippet = &models.Snippet{
	ID:      1,
	Title:   "Test title",
	Content: "Test content",
	Created: time.Now(),
	Expires: time.Now(),
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	return 2, nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}
