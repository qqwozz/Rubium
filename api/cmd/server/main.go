package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api/internal/config"
	"api/internal/handlers"
	"api/internal/middleware"
	"api/internal/supabase"

	"github.com/gin-gonic/gin"
)

func main() {
	// Запускаем тесты перед стартом сервера
	// tests.RunAll()

	// Загружаем конфигурацию (порт, ключи Supabase)
	cfg := config.Load()

	// Инициализируем клиент Supabase (PostgREST + Auth)
	client := supabase.NewClient(cfg.SupabaseURL, cfg.SupabaseAnonKey, cfg.SupabaseServiceKey)

	// Настраиваем Gin в production-режиме (без debug-логов)
	gin.SetMode(gin.ReleaseMode)

	// Создаём роутер
	r := gin.New()
	r.Use(gin.Recovery())         // Recovery от паник
	r.Use(corsMiddleware())       // CORS для фронтенда
	r.Use(bodySizeLimit(5 << 20)) // Лимит тела запроса: 5 MB

	// Инициализируем хендлеры
	tasks := handlers.NewTasksHandler(client)
	check := handlers.NewCheckHandler(client)
	notebooks := handlers.NewNotebooksHandler(client)

	// --- Tasks: public read ---
	r.GET("/api/v1/tasks", tasks.GetTasks)
	r.GET("/api/v1/tasks/:id", tasks.GetTaskByID)
	r.POST("/api/v1/check", check.Check)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// --- Tasks: protected write (auth required) ---
	// TODO: add admin/moderator role check inside handlers when role system is ready
	tasksAdmin := r.Group("/api/v1/tasks")
	tasksAdmin.Use(middleware.RequireAuth(client))
	{
		tasksAdmin.PUT("/:id", tasks.UpdateTask)
		tasksAdmin.DELETE("/:id", tasks.DeleteTask)
	}

	// --- Notebooks: community (без middleware, отдельный роут) ---
	r.GET("/api/v1/notebooks/community", notebooks.GetCommunityNotebooks)

	// --- Notebooks: публичные / опциональная авторизация ---
	nb := r.Group("/api/v1/notebooks")
	nb.Use(middleware.OptionalAuth(client))
	{
		nb.GET("/:id", notebooks.GetNotebookByID)  // Просмотр тетради
		nb.GET("/:id/rating", notebooks.GetRating) // Получить рейтинг
	}

	// --- Notebooks: только авторизованные ---
	nbPrivate := r.Group("/api/v1/notebooks")
	nbPrivate.Use(middleware.RequireAuth(client))
	{
		nbPrivate.GET("", notebooks.GetNotebooks)             // Мои тетради
		nbPrivate.POST("", notebooks.CreateNotebook)          // Создать
		nbPrivate.PUT("/:id", notebooks.UpdateNotebook)       // Обновить
		nbPrivate.DELETE("/:id", notebooks.DeleteNotebook)    // Удалить
		nbPrivate.POST("/:id/copy", notebooks.CopyNotebook)   // Копировать
		nbPrivate.POST("/:id/rate", notebooks.RateNotebook)   // Оценить
		nbPrivate.POST("/:id/view", notebooks.IncrementViews) // Увеличить просмотры
	}

	// HTTP-сервер с таймаутами (prod-ready)
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// Graceful shutdown: ловим SIGINT / SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		fmt.Println("\nShutting down...")

		// Даём активным запросам 10 секунд на завершение
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Force shutdown: %v", err)
		}
		fmt.Println("Server stopped")
	}()

	log.Printf("Server started on http://localhost:%s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// corsMiddleware — разрешает запросы с локальных фронтендов
func corsMiddleware() gin.HandlerFunc {
	origins := map[string]bool{
		"http://localhost:5500": true,
		"http://localhost:5080": true,
		"http://localhost:5081": true,
		"http://localhost:3000": true,
		"http://127.0.0.1:5500": true,
		"http://127.0.0.1:5080": true,
		"http://localhost:5173": true,
		"http://rubium.tech":    true,
		"https://rubium.tech":   true,
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origins[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS, DELETE")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Header("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// bodySizeLimit — ограничивает размер тела запроса (5 MB)
func bodySizeLimit(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}
