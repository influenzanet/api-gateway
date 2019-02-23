package main

import (
	"context"
	"log"
	"net/http"

	mw "github.com/Influenzanet/api-gateway/middlewares"
	auth_api "github.com/Influenzanet/api/dist/go/auth-service"
	"github.com/gin-gonic/gin"
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
		token.GET("/renew", tokenRenewHandl)
	}
}

func tokenRenewHandl(c *gin.Context) {
	token := c.MustGet("encodedToken").(string)

	req := auth_api.EncodedToken{
		Token: token,
	}
	newToken, err := clients.authService.RenewJWT(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		log.Println(st.Message())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error during token renewal"})
		return
	}
	c.JSON(http.StatusOK, newToken)
}
