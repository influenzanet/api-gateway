package main

import (
	"net/http"

	mw "github.com/Influenzanet/middlewares"
	"github.com/Influenzanet/utils"
	"github.com/gin-gonic/gin"
)

// InitUserEndpoints creates all API routes on the supplied RouterGroup
func InitUserEndpoints(rg *gin.RouterGroup) {
	user := rg.Group("/user")
	user.Use(mw.RequirePayload())
	{
		user.POST("/login", userLoginHandl)
		user.POST("/signup", userSignupHandl)
	}
	userToken := rg.Group("/user")
	userToken.Use(mw.ExtractToken())
	userToken.Use(mw.ValidateToken(Conf.ServiceURL.Authentication + "/v1/token/validate"))
	userToken.Use(mw.RequirePayload())
	{
		user.POST("/change-password", userPasswordChangeHandl)
	}
	userGet := rg.Group("/user")
	userGet.Use(mw.ExtractURLToken())
	{
		userGet.GET("/verify-email", userEmailVerifyHandl)
	}
}

func userLoginHandl(c *gin.Context) {
	utils.UntouchedPostForward(Conf.ServiceURL.Authentication+"/v1/user/login", c)
}

func userSignupHandl(c *gin.Context) {
	utils.UntouchedPostForward(Conf.ServiceURL.Authentication+"/v1/user/signup", c)
}

func userPasswordChangeHandl(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func userEmailVerifyHandl(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}
