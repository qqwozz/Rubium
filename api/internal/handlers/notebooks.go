package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"api/internal/supabase"
	"api/internal/validation"

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
	ID            string          `json:"id"`
	UserID        string          `json:"user_id"`
	Title         string          `json:"title"`
	Description   string          `json:"description"`
	Color         string          `json:"color"`
	Tags          []string        `json:"tags"`
	IsPublic      bool            `json:"is_public"`
	Content       json.RawMessage `json:"content"`
	AverageRating float64         `json:"average_rating"`
	RatingsCount  int             `json:"ratings_count"`
	ViewsCount    int             `json:"views_count"`
	CopiesCount   int             `json:"copies_count"`
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     string          `json:"updated_at"`
}

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
		"sections_count": sc,
		"pages_count":    pc,
		"views_count":    r.ViewsCount,
		"copies_count":   r.CopiesCount,
		"average_rating": r.AverageRating,
		"created_at":     r.CreatedAt,
		"updated_at":     r.UpdatedAt,
	}
	if includeContent {
		resp["content"] = r.Content
	}
	return resp
}

func (h *NotebooksHandler) getRubiumUserID(ctx context.Context, authID string) (string, error) {
	var users []struct {
		ID string `json:"id"`
	}
	usersEndpoint := fmt.Sprintf("rubium_users?select=id&auth_id=eq.%s", authID)
	rawUsers, err := h.client.RawQuery(ctx, usersEndpoint, false)
	if err != nil {
		return "", err
	}
	if len(rawUsers) > 0 && rawUsers[0] == '{' {
		rawUsers = append([]byte("["), append(rawUsers, []byte("]")...)...)
	}
	if err := json.Unmarshal(rawUsers, &users); err != nil {
		return "", err
	}
	if len(users) == 0 {
		return "", fmt.Errorf("пользователь не найден")
	}
	return users[0].ID, nil
}

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

	selectFields := "id,user_id,title,description,color,tags,is_public,average_rating,views_count,copies_count,created_at,updated_at,content"

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
	if !isValidUUID(id) {
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

type CreateNotebookRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description"`
	Color       string   `json:"color"`
	Tags        []string `json:"tags"`
	IsPublic    bool     `json:"is_public"`
}

func (h *NotebooksHandler) CreateNotebook(c *gin.Context) {
	var req CreateNotebookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "нужен title"})
		return
	}

	authID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "требуется авторизация"})
		return
	}

	rubiumUserID, err := h.getRubiumUserID(c.Request.Context(), authID.(string))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	payload := map[string]interface{}{
		"user_id": rubiumUserID,
		"title":   req.Title,
		"content": NotebookContent{Sections: []Section{}},
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
	if err := h.client.Post(c.Request.Context(), "notebooks", true, payload, &rows); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(rows) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось создать тетрадь"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"notebook": rows[0].toResponse(false)})
}

type UpdateNotebookRequest struct {
	Title       *string          `json:"title"`
	Description *string          `json:"description"`
	Color       *string          `json:"color"`
	Tags        *[]string        `json:"tags"`
	IsPublic    *bool            `json:"is_public"`
	Content     *json.RawMessage `json:"content"`
}

func (h *NotebooksHandler) UpdateNotebook(c *gin.Context) {
	id := c.Param("id")
	if !isValidUUID(id) {
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

func (h *NotebooksHandler) DeleteNotebook(c *gin.Context) {
	id := c.Param("id")
	if !isValidUUID(id) {
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

	endpoint := fmt.Sprintf("notebooks?id=eq.%s", id)
	if err := h.client.Delete(c.Request.Context(), endpoint, true); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "тетрадь удалена"})
}

func (h *NotebooksHandler) CopyNotebook(c *gin.Context) {
	id := c.Param("id")
	if !isValidUUID(id) {
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

type RateRequest struct {
	Rating int `json:"rating" binding:"required"`
}

func (h *NotebooksHandler) RateNotebook(c *gin.Context) {
	id := c.Param("id")
	if !isValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный UUID"})
		return
	}

	authID := c.MustGet("user_id").(string)

	rubiumUserID, err := h.getRubiumUserID(c.Request.Context(), authID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	var req RateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "нужен rating"})
		return
	}
	if req.Rating < 1 || req.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rating должен быть от 1 до 5"})
		return
	}

	// Проверяем, не оценивал ли уже пользователь
	var userRows []struct {
		RatedNotebooks map[string]int `json:"rated_notebooks"`
	}
	userEndpoint := fmt.Sprintf("rubium_users?select=rated_notebooks&id=eq.%s&limit=1", rubiumUserID)
	rawUser, err := h.client.RawQuery(c.Request.Context(), userEndpoint, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(rawUser) > 0 && rawUser[0] == '{' {
		rawUser = append([]byte("["), append(rawUser, []byte("]")...)...)
	}
	if err := json.Unmarshal(rawUser, &userRows); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(userRows) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "пользователь не найден"})
		return
	}

	ratedNotebooks := userRows[0].RatedNotebooks
	if ratedNotebooks == nil {
		ratedNotebooks = make(map[string]int)
	}

	if _, exists := ratedNotebooks[id]; exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "вы уже оценивали эту тетрадь"})
		return
	}

	var rows []notebookRow
	endpoint := fmt.Sprintf("notebooks?select=id,user_id,is_public,average_rating,ratings_count&id=eq.%s&limit=1", id)
	if err := h.client.Query(c.Request.Context(), endpoint, true, &rows); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(rows) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "тетрадь не найдена"})
		return
	}
	nb := rows[0]

	if !nb.IsPublic {
		c.JSON(http.StatusForbidden, gin.H{"error": "оценивать можно только публичные тетради"})
		return
	}
	if nb.UserID == rubiumUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "нельзя оценивать свою тетрадь"})
		return
	}

	newCount := nb.RatingsCount + 1
	newAvg := (nb.AverageRating*float64(nb.RatingsCount) + float64(req.Rating)) / float64(newCount)

	patchEndpoint := fmt.Sprintf("notebooks?id=eq.%s", id)
	if err := h.client.Patch(c.Request.Context(), patchEndpoint, true, map[string]interface{}{
		"average_rating": newAvg,
		"ratings_count":  newCount,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Сохраняем оценку пользователя
	ratedNotebooks[id] = req.Rating
	userPatchEndpoint := fmt.Sprintf("rubium_users?id=eq.%s", rubiumUserID)
	if err := h.client.Patch(c.Request.Context(), userPatchEndpoint, false, map[string]interface{}{
		"rated_notebooks": ratedNotebooks,
	}); err != nil {
		// Не критично — логируем, но не ломаем ответ
		fmt.Printf("⚠ не удалось сохранить оценку пользователя %s: %v\n", rubiumUserID, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"average_rating": newAvg,
		"ratings_count":  newCount,
	})
}

func (h *NotebooksHandler) GetRating(c *gin.Context) {
	id := c.Param("id")
	if !isValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный UUID"})
		return
	}

	var rows []notebookRow
	endpoint := fmt.Sprintf("notebooks?select=id,user_id,is_public,average_rating&id=eq.%s&limit=1", id)
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
		c.JSON(http.StatusOK, gin.H{
			"average_rating": nb.AverageRating,
		})
		return
	}

	uid, authed := c.Get("user_id")
	if !authed {
		c.JSON(http.StatusForbidden, gin.H{"error": "нет доступа"})
		return
	}

	rubiumUserID, err := h.getRubiumUserID(c.Request.Context(), uid.(string))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	if rubiumUserID != nb.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "нет доступа"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"average_rating": nb.AverageRating,
	})
}

func (h *NotebooksHandler) IncrementViews(c *gin.Context) {
	id := c.Param("id")
	if !isValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный UUID"})
		return
	}

	if err := h.client.RPC(c.Request.Context(), "increment_notebook_views", true, map[string]interface{}{
		"notebook_id": id,
	}, nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "просмотр засчитан"})
}

func (h *NotebooksHandler) GetCommunityNotebooks(c *gin.Context) {
	sortVal := c.DefaultQuery("sort", "rating")
	if _, ok := validation.ValidSort(sortVal); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный sort"})
		return
	}

	search := c.Query("search")
	if search != "" {
		if _, ok := validation.SafeSearchString(search); !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный поисковый запрос"})
			return
		}
	}

	limitVal := c.DefaultQuery("limit", "50")
	limit, ok := validation.ValidLimit(limitVal)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit должен быть от 1 до 100"})
		return
	}

	selectFields := "id,user_id,title,description,color,tags,is_public,average_rating,views_count,copies_count,created_at,updated_at"

	filters := []string{
		"select=" + selectFields,
		"is_public=eq.true",
	}

	if search != "" {
		safeSearch, _ := validation.SafeSearchString(search)
		filters = append(filters, fmt.Sprintf("or=(title.ilike.*%s*,description.ilike.*%s*)", safeSearch, safeSearch))
	}

	switch sortVal {
	case "newest":
		filters = append(filters, "order=created_at.desc")
	case "popular":
		filters = append(filters, "order=views_count.desc")
	default:
		filters = append(filters, "order=average_rating.desc")
	}

	filters = append(filters, fmt.Sprintf("limit=%d", limit))

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

	userIDs := make([]string, 0, len(rows))
	for _, r := range rows {
		userIDs = append(userIDs, r.UserID)
	}

	authors := make(map[string]gin.H)
	if len(userIDs) > 0 {
		usersEndpoint := fmt.Sprintf("rubium_users?select=id,first_name,last_name,email,avatar_url&id=in.(%s)", strings.Join(userIDs, ","))
		rawUsers, err := h.client.RawQuery(c.Request.Context(), usersEndpoint, false)
		if err == nil {
			if len(rawUsers) > 0 && rawUsers[0] == '{' {
				rawUsers = append([]byte("["), append(rawUsers, []byte("]")...)...)
			}
			var users []struct {
				ID        string `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Email     string `json:"email"`
				AvatarURL string `json:"avatar_url"`
			}
			if json.Unmarshal(rawUsers, &users) == nil {
				for _, u := range users {
					authors[u.ID] = gin.H{
						"id":         u.ID,
						"first_name": u.FirstName,
						"last_name":  u.LastName,
						"email":      u.Email,
						"avatar_url": u.AvatarURL,
					}
				}
			}
		}
	}

	notebooks := make([]gin.H, 0, len(rows))
	for _, r := range rows {
		resp := r.toResponse(false)
		if author, ok := authors[r.UserID]; ok {
			resp["author"] = author
		} else {
			resp["author"] = gin.H{
				"id":         "",
				"first_name": "Автор",
				"email":      "",
				"avatar_url": "",
			}
		}
		notebooks = append(notebooks, resp)
	}

	c.JSON(http.StatusOK, gin.H{"notebooks": notebooks})
}

func (h *NotebooksHandler) getOwner(ctx context.Context, id string) (string, error) {
	var rows []notebookRow
	endpoint := fmt.Sprintf("notebooks?select=user_id&id=eq.%s&limit=1", id)
	if err := h.client.Query(ctx, endpoint, true, &rows); err != nil {
		return "", err
	}
	if len(rows) == 0 {
		return "", fmt.Errorf("не найдено")
	}
	return rows[0].UserID, nil
}
