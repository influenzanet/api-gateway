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

func initUserManagementEndpoints(rg *gin.RouterGroup) {
	userToken := rg.Group("/user")
	userToken.Use(mw.ExtractToken())
	userToken.Use(mw.ValidateToken(clients.AuthService))
	{
		userToken.GET("/", getUserHandl)
		// userToken.GET("/:id", getUserHandl)
		userToken.POST("/change-password", mw.RequirePayload(), userPasswordChangeHandl)
		userToken.POST("/change-account-email", mw.RequirePayload(), changeAccountEmailHandl)
		userToken.POST("/set-language", mw.RequirePayload(), userSetPreferredLanguageHandl)
		userToken.POST("/delete", mw.RequirePayload(), deleteAccountHandl)

		userToken.POST("/profile/save", mw.RequirePayload(), saveProfileHandl)
		userToken.POST("/profile/remove", mw.RequirePayload(), removeProfileHandl)

		//userToken.POST("/contact-preferences", mw.RequirePayload(), todo)
		// userToken.POST("/contact/add-email", mw.RequirePayload(), todo)
		//userToken.POST("/contact/remove-email", mw.RequirePayload(), todo)
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
	if err := gjpb.JsonToPb(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = &token

	resp, err := clients.UserManagement.ChangePassword(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		log.Println(st.Message())
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	gjpb.SendPBAsJSON(c, http.StatusOK, resp)
}

func changeAccountEmailHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(api.TokenInfos)

	var req api.EmailChangeMsg
	if err := gjpb.JsonToPb(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = &token
	resp, err := clients.UserManagement.ChangeAccountIDEmail(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	gjpb.SendPBAsJSON(c, http.StatusOK, resp)
}

func userSetPreferredLanguageHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(api.TokenInfos)

	var req api.LanguageChangeMsg
	if err := gjpb.JsonToPb(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = &token
	resp, err := clients.UserManagement.ChangePreferredLanguage(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	gjpb.SendPBAsJSON(c, http.StatusOK, resp)
}

func getUserHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(api.TokenInfos)

	userID := c.Param("id")

	userRefReq := &api.UserReference{
		Token:  &token,
		UserId: userID,
	}

	resp, err := clients.UserManagement.GetUser(context.Background(), userRefReq)
	if err != nil {
		st := status.Convert(err)
		log.Println(st.Message())
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	gjpb.SendPBAsJSON(c, http.StatusOK, resp)
}

func deleteAccountHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(api.TokenInfos)

	var req api.UserReference
	if err := gjpb.JsonToPb(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = &token
	resp, err := clients.UserManagement.DeleteAccount(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		log.Println(st.Message())
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	gjpb.SendPBAsJSON(c, http.StatusOK, resp)
}

func saveProfileHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(api.TokenInfos)

	var req api.ProfileRequest
	if err := gjpb.JsonToPb(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = &token
	resp, err := clients.UserManagement.SaveProfile(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	gjpb.SendPBAsJSON(c, http.StatusOK, resp)
}

func removeProfileHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(api.TokenInfos)

	var req api.ProfileRequest
	if err := gjpb.JsonToPb(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = &token
	resp, err := clients.UserManagement.RemoveProfile(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	gjpb.SendPBAsJSON(c, http.StatusOK, resp)
}

func userEmailVerifyHandl(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}
