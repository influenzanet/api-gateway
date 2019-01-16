package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ExtractToken extracts the authorization header and saves the bearer token in the context
func ExtractToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := c.Request

		var token string
		tokens, ok := req.Header["Authorization"]
		if ok && len(tokens) >= 1 {
			token = tokens[0]
			token = strings.TrimPrefix(token, "Bearer ")
			if len(token) == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "no Authorization token found"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no Authorization token found"})
			c.Abort()
			return
		}

		c.Set("encodedToken", token)

		c.Next()
	}
}
