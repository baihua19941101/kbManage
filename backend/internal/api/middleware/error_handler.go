package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}

		lastErr := c.Errors.Last()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      lastErr.Error(),
			"request_id": c.GetString(RequestIDKey),
		})
	}
}
