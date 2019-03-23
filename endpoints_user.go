package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	mw "github.com/influenzanet/api-gateway/middlewares"
	infl_api "github.com/influenzanet/api/dist/go"
	user_api "github.com/influenzanet/api/dist/go/user-management"
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
		user.POST("/loginWithEmail", userLoginHandl)
		user.POST("/signupWithEmail", userSignupHandl)

	}
	userToken := rg.Group("/user")
	userToken.Use(mw.ExtractToken())
	userToken.Use(mw.ValidateToken(clients.authService))
	{
		userToken.GET("/", getUserHandl)
		userToken.GET("/:id", getUserHandl)

		userTokenWithPayload := userToken.Group("/")
		userTokenWithPayload.Use(mw.RequirePayload())
		{
			userTokenWithPayload.POST("/changePassword", userPasswordChangeHandl)
			// userTokenWithPayload.PUT("/profile", updateProfileHandl)
			userTokenWithPayload.POST("/subprofile", addSubProfileHandl)
			userTokenWithPayload.PUT("/subprofile", updateSubProfileHandl)
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

func userLoginHandl(c *gin.Context) {
	var req infl_api.UserCredentials
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := clients.authService.LoginWithEmail(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		return
	}
	c.JSON(http.StatusOK, token)
}

func userSignupHandl(c *gin.Context) {
	var req infl_api.UserCredentials
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
	parsedToken := c.MustGet("parsedToken").(infl_api.ParsedToken)

	type PwChangeReq struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	var req PwChangeReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pwReq := &user_api.PasswordChangeMsg{
		Auth:        &parsedToken,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}

	res, err := clients.userManagement.ChangePassword(context.Background(), pwReq)
	if err != nil {
		st := status.Convert(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func getUserHandl(c *gin.Context) {
	parsedToken := c.MustGet("parsedToken").(infl_api.ParsedToken)

	userID := c.Param("id")
	if userID == "" {
		userID = parsedToken.UserId
	}

	userRefReq := &user_api.UserReference{
		Auth:   &parsedToken,
		UserId: userID,
	}

	res, err := clients.userManagement.GetUser(context.Background(), userRefReq)
	if err != nil {
		st := status.Convert(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		return
	}
	c.JSON(http.StatusOK, res)
}

/* TODO: remove
func updateProfileHandl(c *gin.Context) {
	parsedToken := c.MustGet("parsedToken").(infl_api.ParsedToken)

	var profile user_api.Profile
	if err := c.BindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req := &user_api.ProfileRequest{
		Auth:    &parsedToken,
		Profile: &profile,
	}

	res, err := clients.userManagement.UpdateProfile(context.Background(), req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		return
	}
	c.JSON(http.StatusOK, res)
}*/

func parseSubProfileRequest(c *gin.Context) *user_api.SubProfileRequest {
	parsedToken := c.MustGet("parsedToken").(infl_api.ParsedToken)

	var subProfile user_api.SubProfile
	if err := c.BindJSON(&subProfile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return &user_api.SubProfileRequest{}
	}

	return &user_api.SubProfileRequest{
		Auth:       &parsedToken,
		SubProfile: &subProfile,
	}
}

func addSubProfileHandl(c *gin.Context) {
	req := parseSubProfileRequest(c)
	res, err := clients.userManagement.AddSubprofile(context.Background(), req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func updateSubProfileHandl(c *gin.Context) {
	req := parseSubProfileRequest(c)
	res, err := clients.userManagement.EditSubprofile(context.Background(), req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func userEmailVerifyHandl(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}
