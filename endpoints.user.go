package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	mw "github.com/Influenzanet/api-gateway/middlewares"
	user_api "github.com/Influenzanet/api/dist/go/user-management"
)

func connectToUserManagementServer() *grpc.ClientConn {
	conn, err := grpc.Dial(conf.ServiceURLs.UserManagement, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	return conn
}

// InitUserEndpoints creates all API routes on the supplied RouterGroup
func InitUserEndpoints(rg *gin.RouterGroup) {
	user := rg.Group("/user")
	user.Use(mw.RequirePayload())
	{
		user.POST("/login", userLoginHandl)
		user.POST("/signup", userSignupHandl)
	}
	/*
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
		}*/
}

func userLoginHandl(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func userSignupHandl(c *gin.Context) {
	var req user_api.NewUser
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := clients.authService.SignupWithEmail(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		return
	}
	c.JSON(http.StatusCreated, token)
}

func userPasswordChangeHandl(c *gin.Context) {
	// TODO: validate token
	pwReq := &user_api.PasswordChangeMsg{}
	status, err := clients.userManagement.ChangePassword(context.Background(), pwReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println(status)
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func userEmailVerifyHandl(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}
