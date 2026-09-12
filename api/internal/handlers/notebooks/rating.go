package notebooks

import (
	"encoding/json"
	"fmt"
	"net/http"

	"api/internal/validation"

	"github.com/gin-gonic/gin"
)

type RateRequest struct {
	Rating int `json:"rating" binding:"required"`
}

// RateNotebook — POST /api/v1/notebooks/:id/rate
func (h *NotebooksHandler) RateNotebook(c *gin.Context) {
	id := c.Param("id")

	if !validation.IsValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "невалидный UUID",
		})
		return
	}

	authID := c.MustGet("user_id").(string)

	rubiumUserID, err := h.getRubiumUserID(
		c.Request.Context(),
		authID,
	)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	var req RateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "нужен rating",
		})
		return
	}

	if req.Rating < 1 || req.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "rating должен быть от 1 до 5",
		})
		return
	}

	// Получаем оценки, которые пользователь уже оставлял.
	var userRows []struct {
		RatedNotebooks map[string]int `json:"rated_notebooks"`
	}

	userEndpoint := fmt.Sprintf(
		"rubium_users?select=rated_notebooks&id=eq.%s&limit=1",
		rubiumUserID,
	)

	rawUser, err := h.client.RawQuery(
		c.Request.Context(),
		userEndpoint,
		false,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Supabase может вернуть один объект вместо массива.
	if len(rawUser) > 0 && rawUser[0] == '{' {
		rawUser = append(
			[]byte("["),
			append(rawUser, []byte("]")...)...,
		)
	}

	if err := json.Unmarshal(rawUser, &userRows); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(userRows) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "пользователь не найден",
		})
		return
	}

	ratedNotebooks := userRows[0].RatedNotebooks

	if ratedNotebooks == nil {
		ratedNotebooks = make(map[string]int)
	}

	// Проверяем, не оценивал ли пользователь эту тетрадь.
	if _, exists := ratedNotebooks[id]; exists {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "вы уже оценивали эту тетрадь",
		})
		return
	}

	// Получаем текущий рейтинг тетради.
	var rows []notebookRow

	endpoint := fmt.Sprintf(
		"notebooks?select=id,user_id,is_public,average_rating,ratings_count&id=eq.%s&limit=1",
		id,
	)

	if err := h.client.Query(
		c.Request.Context(),
		endpoint,
		true,
		&rows,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(rows) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "тетрадь не найдена",
		})
		return
	}

	nb := rows[0]

	// Оценивать можно только публичные тетради.
	if !nb.IsPublic {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "оценивать можно только публичные тетради",
		})
		return
	}

	// Нельзя оценивать собственную тетрадь.
	if nb.UserID == rubiumUserID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "нельзя оценивать свою тетрадь",
		})
		return
	}

	// Пересчитываем рейтинг.
	newCount := nb.RatingsCount + 1

	newAvg := (
		nb.AverageRating*float64(nb.RatingsCount) +
			float64(req.Rating)) / float64(newCount)

	// Обновляем рейтинг тетради.
	patchEndpoint := fmt.Sprintf(
		"notebooks?id=eq.%s",
		id,
	)

	if err := h.client.Patch(
		c.Request.Context(),
		patchEndpoint,
		true,
		map[string]interface{}{
			"average_rating": newAvg,
			"ratings_count":  newCount,
		},
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Запоминаем оценку пользователя.
	ratedNotebooks[id] = req.Rating

	userPatchEndpoint := fmt.Sprintf(
		"rubium_users?id=eq.%s",
		rubiumUserID,
	)

	if err := h.client.Patch(
		c.Request.Context(),
		userPatchEndpoint,
		false,
		map[string]interface{}{
			"rated_notebooks": ratedNotebooks,
		},
	); err != nil {
		// Рейтинг тетради уже сохранён,
		// поэтому не ломаем успешный ответ.
		fmt.Printf(
			"⚠ не удалось сохранить оценку пользователя %s: %v\n",
			rubiumUserID,
			err,
		)
	}

	c.JSON(http.StatusOK, gin.H{
		"average_rating": newAvg,
		"ratings_count":  newCount,
	})
}

// GetRating — GET /api/v1/notebooks/:id/rating
func (h *NotebooksHandler) GetRating(c *gin.Context) {
	id := c.Param("id")

	if !validation.IsValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "невалидный UUID",
		})
		return
	}

	var rows []notebookRow

	endpoint := fmt.Sprintf(
		"notebooks?select=id,user_id,is_public,average_rating,ratings_count&id=eq.%s&limit=1",
		id,
	)

	if err := h.client.Query(
		c.Request.Context(),
		endpoint,
		true,
		&rows,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(rows) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "тетрадь не найдена",
		})
		return
	}

	nb := rows[0]

	// Для публичной тетради рейтинг доступен всем.
	if nb.IsPublic {
		c.JSON(http.StatusOK, gin.H{
			"average_rating": nb.AverageRating,
			"ratings_count":  nb.RatingsCount,
		})
		return
	}

	// Для приватной тетради нужен авторизованный пользователь.
	uid, authed := c.Get("user_id")
	if !authed {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "нет доступа",
		})
		return
	}

	authID, ok := uid.(string)
	if !ok || authID == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "некорректный пользователь",
		})
		return
	}

	rubiumUserID, err := h.getRubiumUserID(
		c.Request.Context(),
		authID,
	)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Приватную тетрадь может смотреть только её владелец.
	if rubiumUserID != nb.UserID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "нет доступа",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"average_rating": nb.AverageRating,
		"ratings_count":  nb.RatingsCount,
	})
}