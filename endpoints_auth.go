package main

import (
	"context"
	"log"
	"net/http"

	"coneno.de/rechnungstool/api-gateway/utils"
	"github.com/gin-gonic/gin"
	api "github.com/influenzanet/api-gateway/api"
	mw "github.com/influenzanet/api-gateway/middlewares"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func connectToAuthServiceServer() *grpc.ClientConn {
	conn, err := grpc.Dial(conf.ServiceURLs.Authentication, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	return conn
}

// InitTokenEndpoints creates all API routes on the supplied RouterGroup
func InitTokenEndpoints(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	auth.POST("/loginWithEmail", mw.RequirePayload(), loginWithEmailHandl)
	auth.POST("/signupWithEmail", mw.RequirePayload(), signupWithEmailHandl)

	token := auth.Group("/token")
	token.Use(mw.ExtractToken())
	{
		token.POST("/renew", tokenRenewHandl)
	}
}

func loginWithEmailHandl(c *gin.Context) {
	var req api.UserCredentials
	if err := gjpb.JsonToPb(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := clients.authService.LoginWithEmail(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	gjpb.SendPBAsJSON(c, http.StatusOK, token)
}

func signupWithEmailHandl(c *gin.Context) {
	var req api.UserCredentials
	if err := gjpb.JsonToPb(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := clients.authService.SignupWithEmail(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	gjpb.SendPBAsJSON(c, http.StatusOK, token)
}

func tokenRenewHandl(c *gin.Context) {
	var req api.RefreshJWTRequest
	if err := gjpb.JsonToPb(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.AccessToken = c.MustGet("encodedToken").(string)
	token, err := clients.authService.RenewJWT(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	gjpb.SendPBAsJSON(c, http.StatusOK, token)
}
