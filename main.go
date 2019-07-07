package main

import (
	"log"

	"github.com/gin-gonic/gin"

	api "github.com/influenzanet/api-gateway/api"
)

// APIClients holds the service clients to the internal services
type APIClients struct {
	userManagement api.UserManagementApiClient
	authService    api.AuthServiceApiClient
}

var clients = APIClients{}

func main() {
	ReadConfig()

	// Connect to user management service
	userManagementServerConn := connectToUserManagementServer()
	defer userManagementServerConn.Close()
	clients.userManagement = api.NewUserManagementApiClient(userManagementServerConn)

	// Connect to authentication service
	authenticationServerConn := connectToAuthServiceServer()
	defer authenticationServerConn.Close()
	clients.authService = api.NewAuthServiceApiClient(authenticationServerConn)

	// Start webserver
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	v1 := router.Group("/v1")

	InitUserEndpoints(v1)
	InitTokenEndpoints(v1)
	InitSurveyEndpoints(v1)

	log.Fatal(router.Run(":3000"))
}
