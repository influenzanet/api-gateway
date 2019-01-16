package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ExtractURLToken extracts the authorization header and saves the bearer token in the context
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
