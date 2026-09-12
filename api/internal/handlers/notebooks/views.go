package notebooks

import (
	"fmt"
	"net/http"

	"api/internal/validation"

	"github.com/gin-gonic/gin"
)

func (h *NotebooksHandler) IncrementViews(c *gin.Context) {
	id := c.Param("id")
	if !validation.IsValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный UUID"})
		return
	}

	// Если пользователь авторизован — проверяем, не его ли это тетрадь
	if authID, authed := c.Get("user_id"); authed {
		rubiumUserID, err := h.getRubiumUserID(c.Request.Context(), authID.(string))
		if err == nil {
			var rows []notebookRow
			endpoint := fmt.Sprintf("notebooks?select=user_id&id=eq.%s&limit=1", id)
			if err := h.client.Query(c.Request.Context(), endpoint, true, &rows); err == nil && len(rows) > 0 {
				if rows[0].UserID == rubiumUserID {
					c.JSON(http.StatusOK, gin.H{"message": "своя тетрадь"})
					return
				}
			}
		}
	}

	if err := h.client.RPC(c.Request.Context(), "increment_notebook_views", true, map[string]interface{}{
		"notebook_id": id,
	}, nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "просмотр засчитан"})
}