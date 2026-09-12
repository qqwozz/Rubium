package notebooks

import "api/internal/supabase"

type Handler struct {
	client *supabase.Client
}

func NewHandler(client *supabase.Client) *Handler {
	return &Handler{
		client: client,
	}
}