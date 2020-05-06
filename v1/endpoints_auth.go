package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/influenzanet/api-gateway/api"
	mw "github.com/influenzanet/api-gateway/middlewares"
	"github.com/influenzanet/api-gateway/utils"
	gjpb "github.com/phev8/gin-protobuf-json-converter"
	"google.golang.org/grpc/status"
)

func initAuthEndpoints(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	auth.POST("/login-with-email", mw.RequirePayload(), loginWithEmailHandl)
	auth.POST("/signup-with-email", mw.RequirePayload(), signupWithEmailHandl)
	auth.POST("/switch-profile", mw.ExtractToken(), mw.ValidateToken(clients.AuthService), mw.RequirePayload(), switchProfileHandl)
	auth.POST("/renew-token", mw.ExtractToken(), mw.RequirePayload(), tokenRenewHandl)
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