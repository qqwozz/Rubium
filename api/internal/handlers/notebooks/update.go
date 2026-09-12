package notebooks

import (
	"fmt"
	"net/http"
	"time"

	"api/internal/validation"

	"github.com/gin-gonic/gin"
)

func (h *NotebooksHandler) UpdateNotebook(c *gin.Context) {
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

	owner, err := h.getOwner(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "тетрадь не найдена"})
		return
	}
	if owner != rubiumUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "нет доступа к тетради"})
		return
	}

	var req UpdateNotebookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный JSON"})
		return
	}

	updates := map[string]interface{}{}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Color != nil {
		updates["color"] = *req.Color
	}
	if req.Tags != nil {
		updates["tags"] = *req.Tags
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "нет полей для обновления"})
		return
	}

	updates["updated_at"] = time.Now().UTC().Format(time.RFC3339)

	endpoint := fmt.Sprintf("notebooks?id=eq.%s", id)
	if err := h.client.Patch(c.Request.Context(), endpoint, true, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "тетрадь обновлена"})
}