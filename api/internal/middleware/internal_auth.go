package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireInternalKey(expectedKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		got := c.GetHeader("X-Internal-Key")

		if got == "" ||
			subtle.ConstantTimeCompare([]byte(got), []byte(expectedKey)) != 1 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}