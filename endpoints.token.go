package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
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
	/* token := rg.Group("/token")
	token.Use(mw.ExtractToken())
	{
		token.GET("/renew", tokenRenewHandl)
	}
	*/
}

func tokenRenewHandl(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
	// token := c.MustGet("encodedToken").(string)

	// utils.ForwardPostRequestWithAuth(Conf.ServiceURL.Authentication+"/v1/token/renew", "Bearer "+token, c)
}
