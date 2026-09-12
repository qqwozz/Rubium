package notebooks

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *NotebooksHandler) CreateNotebook(c *gin.Context) {
	var req CreateNotebookRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "нужен title",
		})
		return
	}

	authID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "требуется авторизация",
		})
		return
	}

	rubiumUserID, err := h.getRubiumUserID(
		c.Request.Context(),
		authID.(string),
	)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	payload := map[string]interface{}{
		"user_id": rubiumUserID,
		"title":   req.Title,
		"content": NotebookContent{
			Sections: []Section{},
		},
	}

	if req.Description != "" {
		payload["description"] = req.Description
	}

	if req.Color != "" {
		payload["color"] = req.Color
	}

	if req.Tags != nil {
		payload["tags"] = req.Tags
	}

	payload["is_public"] = req.IsPublic

	var rows []notebookRow

	if err := h.client.Post(
		c.Request.Context(),
		"notebooks",
		true,
		payload,
		&rows,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(rows) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "не удалось создать тетрадь",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"notebook": rows[0].toResponse(false),
	})
}