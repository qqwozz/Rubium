package notebooks

import (
	"api/internal/supabase"
)

type NotebooksHandler struct {
	client *supabase.Client
}

func NewNotebooksHandler(client *supabase.Client) *NotebooksHandler {
	return &NotebooksHandler{
		client: client,
	}
}