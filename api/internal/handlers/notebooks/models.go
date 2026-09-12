package notebooks

import (
	"encoding/json"
)

type NotebookContent struct {
	Sections []Section `json:"sections"`
}

type Section struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Pages []Page `json:"pages"`
}

type Page struct {
	ID      string          `json:"id"`
	Title   string          `json:"title"`
	Content json.RawMessage `json:"content"`
}

type notebookRow struct {
	ID            string          `json:"id"`
	UserID        string          `json:"user_id"`
	Title         string          `json:"title"`
	Description   string          `json:"description"`
	Color         string          `json:"color"`
	Tags          []string        `json:"tags"`
	IsPublic      bool            `json:"is_public"`
	Content       json.RawMessage `json:"content"`
	AverageRating float64         `json:"average_rating"`
	RatingsCount  int             `json:"ratings_count"`
	ViewsCount    int             `json:"views_count"`
	CopiesCount   int             `json:"copies_count"`
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     string          `json:"updated_at"`
}

type CreateNotebookRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description"`
	Color       string   `json:"color"`
	Tags        []string `json:"tags"`
	IsPublic    bool     `json:"is_public"`
}

type UpdateNotebookRequest struct {
	Title       *string          `json:"title"`
	Description *string          `json:"description"`
	Color       *string          `json:"color"`
	Tags        *[]string        `json:"tags"`
	IsPublic    *bool            `json:"is_public"`
	Content     *json.RawMessage `json:"content"`
}