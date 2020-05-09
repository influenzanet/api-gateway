package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ExtractURLToken extracts the url token and saves it in the context
func ExtractURLToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")

		if token == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing token"})
			return
		}

		c.Set("urlToken", token)

		c.Next()
	}
}
