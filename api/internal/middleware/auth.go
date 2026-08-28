package middleware

import (
	"net/http"
	"strings"

	"api/internal/supabase"

	"github.com/gin-gonic/gin"
)

func RequireAuth(client *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "требуется авторизация"})
			return
		}

		userID, err := client.AuthUser(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "невалидный токен"})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

func OptionalAuth(client *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token != "" {
			if userID, err := client.AuthUser(c.Request.Context(), token); err == nil {
				c.Set("user_id", userID)
			}
		}
		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	h := c.GetHeader("Authorization")
	if !strings.HasPrefix(h, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(h, "Bearer ")
}