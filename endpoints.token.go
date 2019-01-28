package main

import (
	"io/ioutil"
	"net/http"
	"strings"

	middlewares "github.com/Influenzanet/middlewares"
	"github.com/gin-gonic/gin"
)

// InitTokenEndpoints creates all API routes on the supplied RouterGroup
func InitTokenEndpoints(rg *gin.RouterGroup) {
	token := rg.Group("/token")
	token.Use(middlewares.ExtractToken())
	{
		token.GET("/renew", tokenRenewHandl)
	}
}

func tokenRenewHandl(c *gin.Context) {
	token := c.MustGet("encodedToken").(string)

	req, err := http.NewRequest("Get", Conf.URLAuthenticationService+Conf.AuthenticationRenew, nil)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	req.Header.Add("Authorization", strings.Join([]string{"Bearer", token}, " "))

	res, err := client.Do(req)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	rawBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(res.StatusCode, rawBody)
}
