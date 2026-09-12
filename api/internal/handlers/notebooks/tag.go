package notebooks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"api/internal/validation"

	"github.com/gin-gonic/gin"
)

func (h *NotebooksHandler) GetNotebooksByTag(c *gin.Context) {
	tag := c.Query("tag")

	if !validation.IsSafeTag(tag) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "невалидный tag",
		})
		return
	}

	encodedTag := url.QueryEscape(tag)

	limitVal := c.DefaultQuery("limit", "50")

	limit, ok := validation.ValidLimit(limitVal)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "limit должен быть от 1 до 100",
		})
		return
	}

	selectFields := "id,user_id,title,description,color,tags,is_public,average_rating,ratings_count,views_count,copies_count,created_at,updated_at"

	filters := []string{
		"select=" + selectFields,
		"is_public=eq.true",
		fmt.Sprintf("tags=cs.{%s}", encodedTag),
		fmt.Sprintf("limit=%d", limit),
	}

	endpoint := "notebooks?" + strings.Join(filters, "&")

	rawNotebooks, err := h.client.RawQuery(
		c.Request.Context(),
		endpoint,
		true,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "не удалось получить тетради",
		})
		return
	}

	// Supabase может вернуть один объект вместо массива.
	if len(rawNotebooks) > 0 && rawNotebooks[0] == '{' {
		rawNotebooks = append(
			[]byte("["),
			append(rawNotebooks, []byte("]")...)...,
		)
	}

	var rows []notebookRow

	if err := json.Unmarshal(rawNotebooks, &rows); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "не удалось обработать ответ от базы данных",
		})
		return
	}

	notebooks := make([]gin.H, 0, len(rows))

	for _, row := range rows {
		notebooks = append(
			notebooks,
			row.toResponse(false),
		)
	}

	c.JSON(http.StatusOK, gin.H{
		"notebooks": notebooks,
	})
}