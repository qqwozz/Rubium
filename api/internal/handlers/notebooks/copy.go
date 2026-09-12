package notebooks

import (
	"fmt"
	"net/http"

	"api/internal/validation"

	"github.com/gin-gonic/gin"
)

func (h *NotebooksHandler) CopyNotebook(c *gin.Context) {
	id := c.Param("id")
	if !validation.IsValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный UUID"})
		return
	}

	authID := c.MustGet("user_id").(string)

	rubiumUserID, err := h.getRubiumUserID(c.Request.Context(), authID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	var rows []notebookRow
	endpoint := fmt.Sprintf("notebooks?select=*&id=eq.%s&limit=1", id)
	if err := h.client.Query(c.Request.Context(), endpoint, true, &rows); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(rows) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "тетрадь не найдена"})
		return
	}
	original := rows[0]

	if !original.IsPublic {
		c.JSON(http.StatusForbidden, gin.H{"error": "можно копировать только публичные тетради"})
		return
	}

	payload := map[string]interface{}{
		"user_id":      rubiumUserID,
		"title":        original.Title + " (копия)",
		"description":  original.Description,
		"color":        original.Color,
		"tags":         original.Tags,
		"is_public":    false,
		"content":      original.Content,
		"views_count":  0,
		"copies_count": 0,
	}

	var newRows []notebookRow
	if err := h.client.Post(c.Request.Context(), "notebooks", true, payload, &newRows); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(newRows) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось создать копию"})
		return
	}

	if err := h.client.RPC(c.Request.Context(), "increment_notebook_copies", true, map[string]interface{}{
		"notebook_id": id,
	}, nil); err != nil {
		_ = err
	}

	c.JSON(http.StatusCreated, gin.H{"notebook": newRows[0].toResponse(false)})
}