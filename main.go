package main

import (
	"log"

	"github.com/gin-gonic/gin"

	auth_api "github.com/influenzanet/api/dist/go/auth-service"
	user_api "github.com/influenzanet/api/dist/go/user-management"
)

// APIClients holds the service clients to the internal services
type APIClients struct {
	userManagement user_api.UserManagementApiClient
	authService    auth_api.AuthServiceApiClient
}

var clients = APIClients{}

func main() {
	ReadConfig()

	// Connect to user management service
	userManagementServerConn := connectToUserManagementServer()
	defer userManagementServerConn.Close()
	clients.userManagement = user_api.NewUserManagementApiClient(userManagementServerConn)

	// Connect to authentication service
	authenticationServerConn := connectToAuthServiceServer()
	defer authenticationServerConn.Close()
	clients.authService = auth_api.NewAuthServiceApiClient(authenticationServerConn)

	// Start webserver
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	v1 := router.Group("/v1")

	InitUserEndpoints(v1)
	InitTokenEndpoints(v1)
	InitSurveyEndpoints(v1)

	log.Fatal(router.Run(":3000"))
}
