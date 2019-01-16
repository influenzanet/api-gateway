package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequirePayload blocks requests that have no payload attached
func RequirePayload() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" { // TODO url encoded?
			c.Next()
			return
		}
		if c.Request.ContentLength == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "payload missing"})
			c.Abort()
			return
		}
		c.Next()
	}
}
