package v1

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"

	api "github.com/influenzanet/api-gateway/api"
	mw "github.com/influenzanet/api-gateway/middlewares"
	"github.com/influenzanet/api-gateway/utils"
)

func initUserManagementEndpoints(rg *gin.RouterGroup) {
	userToken := rg.Group("/user")
	userToken.Use(mw.ExtractToken())
	userToken.Use(mw.ValidateToken(clients.AuthService))
	{
		userToken.GET("/", getUserHandl)
		// userToken.GET("/:id", getUserHandl)

		userTokenWithPayload := userToken.Group("/")
		userTokenWithPayload.Use(mw.RequirePayload())
		{
			userTokenWithPayload.POST("/changePassword", userPasswordChangeHandl)
			userTokenWithPayload.DELETE("/", deleteAccountHandl)

		}
	}
	/*
		userGet := rg.Group("/user")
		userGet.Use(mw.ExtractURLToken())
		{
			userGet.GET("/verify-email", userEmailVerifyHandl)
		}
	*/
}

func userPasswordChangeHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(api.TokenInfos)

	var req api.PasswordChangeMsg
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = &token

	res, err := clients.UserManagement.ChangePassword(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		log.Println(st.Message())
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func getUserHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(api.TokenInfos)

	userID := c.Param("id")

	userRefReq := &api.UserReference{
		Token:  &token,
		UserId: userID,
	}

	res, err := clients.UserManagement.GetUser(context.Background(), userRefReq)
	if err != nil {
		st := status.Convert(err)
		log.Println(st.Message())
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func deleteAccountHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(api.TokenInfos)

	var req api.UserReference
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = &token
	res, err := clients.UserManagement.DeleteAccount(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		log.Println(st.Message())
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func userEmailVerifyHandl(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}
