package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// InitTokenEndpoints creates all API routes on the supplied RouterGroup
func InitTokenEndpoints(rg *gin.RouterGroup) {
	token := rg.Group("/token")
	{
		token.GET("/renew", tokenRenewHandl)
	}
}

func tokenRenewHandl(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}
