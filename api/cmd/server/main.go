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
	"api/internal/handlers/notebooks"
	"api/internal/middleware"
	"api/internal/supabase"

	"github.com/gin-gonic/gin"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Инициализируем клиент Supabase
	client := supabase.NewClient(
		cfg.SupabaseURL,
		cfg.SupabaseAnonKey,
		cfg.SupabaseServiceKey,
	)

	// Production mode
	gin.SetMode(gin.ReleaseMode)

	// Создаём роутер
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(corsMiddleware())
	r.Use(bodySizeLimit(5 << 20)) // 5 MB

	// =========================
	// Handlers
	// =========================

	tasksHandler := handlers.NewTasksHandler(client)
	checkHandler := handlers.NewCheckHandler(client)
	notebooksHandler := notebooks.NewNotebooksHandler(client)

	// =========================
	// Health
	// =========================

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// =========================
	// Tasks
	// =========================

	// Public read
	r.GET(
		"/api/v1/tasks",
		tasksHandler.GetTasks,
	)

	r.GET(
		"/api/v1/tasks/:id",
		tasksHandler.GetTaskByID,
	)

	// Check task
	r.POST(
		"/api/v1/check",
		checkHandler.Check,
	)

	// Protected write
	tasksAdmin := r.Group("/api/v1/tasks")
	tasksAdmin.Use(middleware.RequireAuth(client))

	{
		tasksAdmin.PUT(
			"/:id",
			tasksHandler.UpdateTask,
		)

		tasksAdmin.DELETE(
			"/:id",
			tasksHandler.DeleteTask,
		)
	}

	// =========================
	// Notebooks
	// =========================

	// Community — публичный endpoint
	r.GET(
		"/api/v1/notebooks/community",
		notebooksHandler.GetCommunityNotebooks,
	)

	// Public / Optional Auth
	nb := r.Group("/api/v1/notebooks")
	nb.Use(middleware.OptionalAuth(client))

	{
		// Получить одну тетрадь
		nb.GET(
			"/:id",
			notebooksHandler.GetNotebookByID,
		)

		// Получить рейтинг
		nb.GET(
			"/:id/rating",
			notebooksHandler.GetRating,
		)
	}

	// Protected notebooks
	nbPrivate := r.Group("/api/v1/notebooks")
	nbPrivate.Use(middleware.RequireAuth(client))

	{
		// Мои тетради
		nbPrivate.GET(
			"",
			notebooksHandler.GetNotebooks,
		)

		// Создать
		nbPrivate.POST(
			"",
			notebooksHandler.CreateNotebook,
		)

		// Обновить
		nbPrivate.PUT(
			"/:id",
			notebooksHandler.UpdateNotebook,
		)

		// Удалить
		nbPrivate.DELETE(
			"/:id",
			notebooksHandler.DeleteNotebook,
		)

		// Копировать
		nbPrivate.POST(
			"/:id/copy",
			notebooksHandler.CopyNotebook,
		)

		// Оценить
		nbPrivate.POST(
			"/:id/rate",
			notebooksHandler.RateNotebook,
		)

		// Увеличить просмотры
		nbPrivate.POST(
			"/:id/view",
			notebooksHandler.IncrementViews,
		)
	}

	// =========================
	// Internal API
	// =========================

	internal := r.Group("/internal/v1")
	internal.Use(
		middleware.RequireInternalKey(cfg.InternalAPIKey),
	)

	{
		internal.GET(
			"/notebooks/by-tag",
			notebooksHandler.GetNotebooksByTag,
		)
	}

	// =========================
	// HTTP Server
	// =========================

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,

		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// =========================
	// Graceful Shutdown
	// =========================

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	go func() {
		<-quit

		fmt.Println("\nShutting down...")

		ctx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf(
				"Force shutdown: %v",
				err,
			)
		}

		fmt.Println("Server stopped")
	}()

	// =========================
	// Start Server
	// =========================

	log.Printf(
		"Server started on http://localhost:%s",
		cfg.Port,
	)

	if err := srv.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// =========================
// CORS
// =========================

func corsMiddleware() gin.HandlerFunc {
	origins := map[string]bool{
		"http://localhost:5500": true,
		"http://localhost:5080": true,
		"http://localhost:5081": true,
		"http://localhost:3000": true,
		"http://localhost:5173": true,

		"http://127.0.0.1:5500": true,
		"http://127.0.0.1:5080": true,

		"http://rubium.tech":  true,
		"https://rubium.tech": true,
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if origins[origin] {
			c.Header(
				"Access-Control-Allow-Origin",
				origin,
			)

			c.Header(
				"Access-Control-Allow-Methods",
				"GET, POST, PUT, OPTIONS, DELETE",
			)

			c.Header(
				"Access-Control-Allow-Headers",
				"Content-Type, Authorization",
			)

			c.Header(
				"Access-Control-Max-Age",
				"86400",
			)
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(
				http.StatusNoContent,
			)
			return
		}

		c.Next()
	}
}

// =========================
// Body Size Limit
// =========================

func bodySizeLimit(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(
			c.Writer,
			c.Request.Body,
			maxBytes,
		)

		c.Next()
	}
}