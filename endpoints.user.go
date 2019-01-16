package main

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

// InitUserEndpoints creates all API routes on the supplied RouterGroup
func InitUserEndpoints(rg *gin.RouterGroup) {
	user := rg.Group("/user")
	user.Use(RequirePayload())
	{
		user.POST("/login", userLoginHandl)
		user.POST("/signup", userSignupHandl)
	}
	userToken := rg.Group("/user")
	userToken.Use(ExtractToken())
	userToken.Use(ValidateToken(Conf))
	userToken.Use(RequirePayload())
	{
		user.POST("/change-password", userPasswordChangeHandl)
	}
	userGet := rg.Group("/user")
	userGet.Use(ExtractURLToken())
	{
		userGet.GET("/verify-email", userEmailVerifyHandl)
	}
}

func userLoginHandl(c *gin.Context) {
	var client = &http.Client{}
	req, err := http.NewRequest("POST", Conf.URLAuthenticationService+Conf.AuthenticationLogin, c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

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

func userSignupHandl(c *gin.Context) {
	var client = &http.Client{}
	req, err := http.NewRequest("POST", Conf.URLUserManagementService+Conf.AuthenticationSignup, c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

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

func userPasswordChangeHandl(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func userEmailVerifyHandl(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}
