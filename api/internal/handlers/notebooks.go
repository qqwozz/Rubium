package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"api/internal/supabase"

	"github.com/gin-gonic/gin"
)

type NotebooksHandler struct {
	client *supabase.Client
}

func NewNotebooksHandler(client *supabase.Client) *NotebooksHandler {
	return &NotebooksHandler{client: client}
}

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
	ID          string          `json:"id"`
	UserID      string          `json:"user_id"`
	Title       string          `json:"title"`
	Color       string          `json:"color"`
	Tags        []string        `json:"tags"`
	IsPublic    bool            `json:"is_public"`
	IsVerified  bool            `json:"is_verified"`
	Content     json.RawMessage `json:"content"`
	ViewsCount  int             `json:"views_count"`
	CopiesCount int             `json:"copies_count"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
}

//! helpers
func countSectionsPages(raw json.RawMessage) (sections, pages int) {
	if len(raw) == 0 {
		return 0, 0
	}
	var nc NotebookContent
	if err := json.Unmarshal(raw, &nc); err != nil {
		return 0, 0
	}
	sections = len(nc.Sections)
	for _, s := range nc.Sections {
		pages += len(s.Pages)
	}
	return sections, pages
}

func (r notebookRow) toResponse(includeContent bool) gin.H {
	sc, pc := countSectionsPages(r.Content)

	// nil tags → [] в JSON
	tags := r.Tags
	if tags == nil {
		tags = []string{}
	}

	resp := gin.H{
		"id":             r.ID,
		"title":          r.Title,
		"color":          r.Color,
		"tags":           tags,
		"is_public":      r.IsPublic,
		"is_verified":    r.IsVerified,
		"sections_count": sc,
		"pages_count":    pc,
		"views_count":    r.ViewsCount,
		"copies_count":   r.CopiesCount,
		"created_at":     r.CreatedAt,
		"updated_at":     r.UpdatedAt,
	}
	if includeContent {
		resp["content"] = r.Content
	}
	return resp
}

//! GET /api/v1/notebooks
// GetNotebooks — список тетрадей **текущего** пользователя.
// Требует авторизации. is_public фильтрует внутри списка пользователя.
func (h *NotebooksHandler) GetNotebooks(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "требуется авторизация"})
		return
	}
	userIDStr := userID.(string)

	filters := []string{
		"select=id,title,color,tags,is_public,is_verified,views_count,copies_count,created_at,updated_at,content",
		fmt.Sprintf("user_id=eq.%s", userIDStr),
	}

	if isPublic := c.Query("is_public"); isPublic != "" {
		filters = append(filters, fmt.Sprintf("is_public=eq.%s", isPublic))
	}

	endpoint := "notebooks?" + strings.Join(filters, "&")

	var rows []notebookRow
	if err := h.client.Query(endpoint, true, &rows); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	notebooks := make([]gin.H, 0, len(rows))
	for _, r := range rows {
		notebooks = append(notebooks, r.toResponse(false))
	}

	c.JSON(http.StatusOK, gin.H{"notebooks": notebooks})
}

//! GET /api/v1/notebooks/:id
func (h *NotebooksHandler) GetNotebookByID(c *gin.Context) {
	id := c.Param("id")
	if !isValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный UUID"})
		return
	}

	var rows []notebookRow
	endpoint := fmt.Sprintf("notebooks?select=*&id=eq.%s&limit=1", id)
	if err := h.client.Query(endpoint, true, &rows); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(rows) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "тетрадь не найдена"})
		return
	}

	nb := rows[0]
	uid, authed := c.Get("user_id")
	if !nb.IsPublic && (!authed || uid.(string) != nb.UserID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "нет доступа к тетради"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notebook": nb.toResponse(true)})
}

//! POST /api/v1/notebooks
type CreateNotebookRequest struct {
	Title    string   `json:"title" binding:"required"`
	Color    string   `json:"color"`
	Tags     []string `json:"tags"`
	IsPublic bool     `json:"is_public"`
}

func (h *NotebooksHandler) CreateNotebook(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	var req CreateNotebookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "нужен title"})
		return
	}

	// Не шлём пустые значения, чтобы сработали DEFAULT в БД
	payload := map[string]interface{}{
		"user_id": userID,
		"title":   req.Title,
		"content": NotebookContent{Sections: []Section{}},
	}
	if req.Color != "" {
		payload["color"] = req.Color
	}
	if req.Tags != nil {
		payload["tags"] = req.Tags
	}
	payload["is_public"] = req.IsPublic

	var rows []notebookRow
	if err := h.client.Post("notebooks", true, payload, &rows); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(rows) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось создать тетрадь"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"notebook": rows[0].toResponse(false)})
}

//! PUT /api/v1/notebooks/:id
type UpdateNotebookRequest struct {
	Title    *string          `json:"title"`
	Color    *string          `json:"color"`
	Tags     *[]string        `json:"tags"`
	IsPublic *bool            `json:"is_public"`
	Content  *json.RawMessage `json:"content"`
}

func (h *NotebooksHandler) UpdateNotebook(c *gin.Context) {
	id := c.Param("id")
	if !isValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный UUID"})
		return
	}
	userID := c.MustGet("user_id").(string)

	owner, err := h.getOwner(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "тетрадь не найдена"})
		return
	}
	if owner != userID {
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
	if err := h.client.Patch(endpoint, true, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "тетрадь обновлена"})
}

//! DELETE /api/v1/notebooks/:id 
func (h *NotebooksHandler) DeleteNotebook(c *gin.Context) {
	id := c.Param("id")
	if !isValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный UUID"})
		return
	}
	userID := c.MustGet("user_id").(string)

	owner, err := h.getOwner(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "тетрадь не найдена"})
		return
	}
	if owner != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "нет доступа к тетради"})
		return
	}

	endpoint := fmt.Sprintf("notebooks?id=eq.%s", id)
	if err := h.client.Delete(endpoint, true); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "тетрадь удалена"})
}

//! helpers
func (h *NotebooksHandler) getOwner(id string) (string, error) {
	var rows []notebookRow
	endpoint := fmt.Sprintf("notebooks?select=user_id&id=eq.%s&limit=1", id)
	if err := h.client.Query(endpoint, true, &rows); err != nil {
		return "", err
	}
	if len(rows) == 0 {
		return "", fmt.Errorf("не найдено")
	}
	return rows[0].UserID, nil
}