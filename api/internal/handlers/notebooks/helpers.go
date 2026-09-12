package notebooks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

func countSectionsPages(raw json.RawMessage) (sections, pages int) {
	if len(raw) == 0 {
		return 0, 0
	}

	var nc NotebookContent

	if err := json.Unmarshal(raw, &nc); err != nil {
		return 0, 0
	}

	sections = len(nc.Sections)

	for _, section := range nc.Sections {
		pages += len(section.Pages)
	}

	return sections, pages
}

func (r notebookRow) toResponse(includeContent bool) gin.H {
	sectionsCount, pagesCount := countSectionsPages(r.Content)

	tags := r.Tags
	if tags == nil {
		tags = []string{}
	}

	resp := gin.H{
		"id":             r.ID,
		"title":          r.Title,
		"description":    r.Description,
		"color":          r.Color,
		"tags":           tags,
		"is_public":      r.IsPublic,
		"sections_count": sectionsCount,
		"pages_count":    pagesCount,
		"views_count":    r.ViewsCount,
		"copies_count":   r.CopiesCount,
		"average_rating": r.AverageRating,
		"ratings_count":  r.RatingsCount,
		"created_at":     r.CreatedAt,
		"updated_at":     r.UpdatedAt,
	}

	if includeContent {
		resp["content"] = r.Content
	}

	return resp
}

func (h *NotebooksHandler) getRubiumUserID(
	ctx context.Context,
	authID string,
) (string, error) {
	var users []struct {
		ID string `json:"id"`
	}

	usersEndpoint := fmt.Sprintf(
		"rubium_users?select=id&auth_id=eq.%s",
		authID,
	)

	rawUsers, err := h.client.RawQuery(
		ctx,
		usersEndpoint,
		false,
	)
	if err != nil {
		return "", err
	}

	// Supabase может вернуть один объект вместо массива.
	if len(rawUsers) > 0 && rawUsers[0] == '{' {
		rawUsers = append(
			[]byte("["),
			append(rawUsers, []byte("]")...)...,
		)
	}

	if err := json.Unmarshal(rawUsers, &users); err != nil {
		return "", err
	}

	if len(users) == 0 {
		return "", fmt.Errorf("пользователь не найден")
	}

	return users[0].ID, nil
}

func (h *NotebooksHandler) getOwner(
	ctx context.Context,
	id string,
) (string, error) {
	var rows []notebookRow

	endpoint := fmt.Sprintf(
		"notebooks?select=user_id&id=eq.%s&limit=1",
		id,
	)

	if err := h.client.Query(
		ctx,
		endpoint,
		true,
		&rows,
	); err != nil {
		return "", err
	}

	if len(rows) == 0 {
		return "", fmt.Errorf("не найдено")
	}

	return rows[0].UserID, nil
}