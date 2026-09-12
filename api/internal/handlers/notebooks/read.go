package notebooks

import (
	"api/internal/validation"
	"fmt"
	"net/http"
	"strings"

	"encoding/json"

	"github.com/gin-gonic/gin"
)


func (h *NotebooksHandler) GetNotebooks(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "требуется авторизация"})
		return
	}
	authID := userID.(string)

	rubiumUserID, err := h.getRubiumUserID(c.Request.Context(), authID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	selectFields := "id,user_id,title,description,color,tags,is_public,average_rating,ratings_count,views_count,copies_count,created_at,updated_at,content"

	filters := []string{
		"select=" + selectFields,
		fmt.Sprintf("user_id=eq.%s", rubiumUserID),
	}

	if isPublic := c.Query("is_public"); isPublic != "" {
		filters = append(filters, fmt.Sprintf("is_public=eq.%s", isPublic))
	}

	endpoint := "notebooks?" + strings.Join(filters, "&")

	rawNotebooks, err := h.client.RawQuery(c.Request.Context(), endpoint, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(rawNotebooks) > 0 && rawNotebooks[0] == '{' {
		rawNotebooks = append([]byte("["), append(rawNotebooks, []byte("]")...)...)
	}

	var rows []notebookRow
	if err := json.Unmarshal(rawNotebooks, &rows); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	notebooks := make([]gin.H, 0, len(rows))
	for _, r := range rows {
		notebooks = append(notebooks, r.toResponse(false))
	}

	c.JSON(http.StatusOK, gin.H{"notebooks": notebooks})
}

func (h *NotebooksHandler) GetNotebookByID(c *gin.Context) {
	id := c.Param("id")
	if !validation.IsValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный UUID"})
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

	nb := rows[0]

	if nb.IsPublic {
		c.JSON(http.StatusOK, gin.H{"notebook": nb.toResponse(true)})
		return
	}

	uid, authed := c.Get("user_id")
	if !authed {
		c.JSON(http.StatusForbidden, gin.H{"error": "нет доступа к тетради"})
		return
	}

	rubiumUserID, err := h.getRubiumUserID(c.Request.Context(), uid.(string))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	if rubiumUserID != nb.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "нет доступа к тетради"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notebook": nb.toResponse(true)})
}