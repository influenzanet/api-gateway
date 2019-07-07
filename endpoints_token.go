package main

import (
	"context"
	"log"
	"net/http"

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
	token := rg.Group("/token")
	token.Use(mw.ExtractToken())
	{
		token.POST("/renew", tokenRenewHandl)
	}
}

func tokenRenewHandl(c *gin.Context) {
	token := c.MustGet("token").(string)
	var req api.RefreshJWTRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.AccessToken = token
	newToken, err := clients.authService.RenewJWT(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		log.Println(st.Message())
		c.JSON(grpcStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	c.JSON(http.StatusOK, newToken)
}
