package v1

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	gjpb "github.com/phev8/gin-protobuf-json-converter"
	"google.golang.org/grpc/status"

	api "github.com/influenzanet/api-gateway/api"
	mw "github.com/influenzanet/api-gateway/middlewares"
	"github.com/influenzanet/api-gateway/utils"
)

func initAuthEndpoints(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	auth.POST("/login-with-email", mw.RequirePayload(), loginWithEmailHandl)
	auth.POST("/signup-with-email", mw.RequirePayload(), signupWithEmailHandl)
	auth.POST("/switch-profile", mw.ExtractToken(), mw.RequirePayload(), switchProfileHandl)
	auth.POST("/renew-token", mw.ExtractToken(), mw.RequirePayload(), tokenRenewHandl)

	// TODO: check if needed:
	// experimental := rg.Group("/exp")
	// experimental.POST("/generate-token", generateTokenHandl)
}

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
	token := c.MustGet("parsedToken").(api.TokenInfos)

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
	token := c.MustGet("parsedToken").(api.TokenInfos)

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
	token := c.MustGet("parsedToken").(api.TokenInfos)

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

func loginWithEmailHandl(c *gin.Context) {
	var req api.LoginWithEmailMsg
	if err := gjpb.JsonToPb(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := clients.AuthService.LoginWithEmail(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	gjpb.SendPBAsJSON(c, http.StatusOK, token)
}

func signupWithEmailHandl(c *gin.Context) {
	var req api.SignupWithEmailMsg
	if err := gjpb.JsonToPb(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := clients.AuthService.SignupWithEmail(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	gjpb.SendPBAsJSON(c, http.StatusOK, token)
}

func switchProfileHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(api.TokenInfos)

	var req api.ProfileRequest
	if err := gjpb.JsonToPb(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = &token
	resp, err := clients.AuthService.SwitchProfile(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	gjpb.SendPBAsJSON(c, http.StatusOK, resp)
}

func tokenRenewHandl(c *gin.Context) {
	var req api.RefreshJWTRequest
	if err := gjpb.JsonToPb(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.AccessToken = c.MustGet("encodedToken").(string)
	token, err := clients.AuthService.RenewJWT(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	gjpb.SendPBAsJSON(c, http.StatusOK, token)
}

func generateTokenHandl(c *gin.Context) {
	var req api.TempTokenInfo
	if err := gjpb.JsonToPb(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := clients.AuthService.GenerateTempToken(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	gjpb.SendPBAsJSON(c, http.StatusOK, token)
}
