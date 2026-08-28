package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strings"

	"api/internal/supabase"
	"api/internal/validation"

	"github.com/gin-gonic/gin"
)

var uuidRe = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func isValidUUID(s string) bool {
	return uuidRe.MatchString(s)
}

type TasksHandler struct {
	client *supabase.Client
}

func NewTasksHandler(client *supabase.Client) *TasksHandler {
	return &TasksHandler{client: client}
}

// GetTasks — GET /api/v1/tasks
func (h *TasksHandler) GetTasks(c *gin.Context) {
	filters := []string{}

	if subject := c.Query("subject"); subject != "" {
		if val, ok := validation.ValidSubject(subject); ok {
			filters = append(filters, fmt.Sprintf("subject=eq.%s", val))
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный subject"})
			return
		}
	}
	if exam := c.Query("exam"); exam != "" {
		if val, ok := validation.ValidExam(exam); ok {
			filters = append(filters, fmt.Sprintf("exam_type=eq.%s", val))
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный exam"})
			return
		}
	}
	if taskType := c.Query("type"); taskType != "" {
		if val, ok := validation.ValidTaskType(taskType); ok {
			filters = append(filters, fmt.Sprintf("task_type=eq.%s", val))
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный type"})
			return
		}
	}
	if topic := c.Query("topic"); topic != "" {
		if val, ok := validation.SafeString(topic); ok {
			filters = append(filters, fmt.Sprintf("topic=eq.%s", val))
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный topic"})
			return
		}
	}
	if difficulty := c.Query("difficulty"); difficulty != "" {
		if val, ok := validation.ValidDifficulty(difficulty); ok {
			filters = append(filters, fmt.Sprintf("difficulty=lte.%s", val))
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный difficulty"})
			return
		}
	}
	if taskNumber := c.Query("task_number"); taskNumber != "" {
		if val, ok := validation.ValidTaskNumber(taskNumber); ok {
			filters = append(filters, fmt.Sprintf("task_number=eq.%s", val))
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный task_number"})
			return
		}
	}

	limit := 1
	if l := c.Query("limit"); l != "" {
		if parsed, ok := validation.ValidLimit(l); ok {
			limit = parsed
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "limit должен быть от 1 до 100"})
			return
		}
	}

	endpoint := "tasks?select=*"
	if len(filters) > 0 {
		endpoint += "&" + strings.Join(filters, "&")
	}

	// Добавляем limit на уровне БД, а не в памяти
	endpoint += fmt.Sprintf("&limit=%d", limit)

	var tasks []map[string]interface{}
	if err := h.client.Query(endpoint, false, &tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// shuffle только если limit небольшой и данных мало
	if len(tasks) > 0 && len(tasks) <= limit {
		rand.Shuffle(len(tasks), func(i, j int) {
			tasks[i], tasks[j] = tasks[j], tasks[i]
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
		"count": len(tasks),
	})
}

// GetTaskByID — GET /api/v1/tasks/:id
func (h *TasksHandler) GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	if !isValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный UUID"})
		return
	}

	var tasks []map[string]interface{}
	endpoint := fmt.Sprintf("tasks?select=*&id=eq.%s&limit=1", id)
	if err := h.client.Query(endpoint, false, &tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(tasks) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "задание не найдено"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": tasks[0]})
}

type UpdateTaskRequest struct {
	Content    *string `json:"content"`
	Answer     *string `json:"answer"`
	Solution   *string `json:"solution"`
	Subject    *string `json:"subject"`
	ExamType   *string `json:"exam_type"`
	Level      *string `json:"level"`
	Topic      *string `json:"topic"`
	TaskType   *string `json:"task_type"`
	TaskNumber *int    `json:"task_number"`
	Source     *string `json:"source"`
}

// UpdateTask — PUT /api/v1/tasks/:id
func (h *TasksHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	if !isValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный UUID"})
		return
	}

	// Guard: middleware must have set user_id
	if _, exists := c.Get("user_id"); !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "требуется авторизация"})
		return
	}
	// TODO: verify admin/moderator role before allowing edits

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный JSON"})
		return
	}

	updates := map[string]interface{}{}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.Answer != nil {
		updates["answer"] = *req.Answer
	}
	if req.Solution != nil {
		updates["solution"] = *req.Solution
	}
	if req.Subject != nil {
		updates["subject"] = *req.Subject
	}
	if req.ExamType != nil {
		updates["exam_type"] = *req.ExamType
	}
	if req.Level != nil {
		updates["level"] = *req.Level
	}
	if req.Topic != nil {
		updates["topic"] = *req.Topic
	}
	if req.TaskType != nil {
		updates["task_type"] = *req.TaskType
	}
	if req.TaskNumber != nil {
		updates["task_number"] = *req.TaskNumber
	}
	if req.Source != nil {
		updates["source"] = *req.Source
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "нет полей для обновления"})
		return
	}

	endpoint := fmt.Sprintf("tasks?id=eq.%s", id)
	if err := h.client.Patch(endpoint, true, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "задание обновлено"})
}

// DeleteTask — DELETE /api/v1/tasks/:id
func (h *TasksHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	if !isValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный UUID"})
		return
	}

	if _, exists := c.Get("user_id"); !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "требуется авторизация"})
		return
	}
	// TODO: verify admin/moderator role before allowing deletion)

	var tasks []map[string]interface{}
	checkEndpoint := fmt.Sprintf("tasks?select=id&id=eq.%s&limit=1", id)
	if err := h.client.Query(checkEndpoint, false, &tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
		return
	}
	if len(tasks) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "задание не найдено"})
		return
	}

	endpoint := fmt.Sprintf("tasks?id=eq.%s", id)
	if err := h.client.Delete(endpoint, true); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "задание удалено"})
}
