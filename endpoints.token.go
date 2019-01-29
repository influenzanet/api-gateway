package main

import (
	mw "github.com/Influenzanet/middlewares"
	"github.com/Influenzanet/utils"
	"github.com/gin-gonic/gin"
)

// InitTokenEndpoints creates all API routes on the supplied RouterGroup
func InitTokenEndpoints(rg *gin.RouterGroup) {
	token := rg.Group("/token")
	token.Use(mw.ExtractToken())
	{
		token.GET("/renew", tokenRenewHandl)
	}
}

func tokenRenewHandl(c *gin.Context) {
	token := c.MustGet("encodedToken").(string)

	utils.ForwardPostRequestWithAuth(Conf.ServiceURL.Authentication+"/v1/token/renew", "Bearer "+token, c)
}
