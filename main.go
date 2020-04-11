package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/golang/protobuf/jsonpb"
	"github.com/influenzanet/api-gateway/api"
	gjpb "github.com/phev8/gin-protobuf-json-converter"
)

// Conf holds all static configuration information
var conf Config
var clients = APIClients{}

func init() {
	initConfig()
	if !conf.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	gjpb.SetMarshaler(jsonpb.Marshaler{
		// EmitDefaults: true,
	})
}

func healthCheckHandle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func main() {
	// Connect to user management service
	userManagementServerConn := connectToUserManagementServer()
	defer userManagementServerConn.Close()
	clients.userManagement = api.NewUserManagementApiClient(userManagementServerConn)

	// Connect to authentication service
	authenticationServerConn := connectToAuthServiceServer()
	defer authenticationServerConn.Close()
	clients.authService = api.NewAuthServiceApiClient(authenticationServerConn)

	// Connect to study service
	studyServiceServerConn := connectToStudyServiceServer()
	defer studyServiceServerConn.Close()
	clients.studyService = api.NewStudyServiceApiClient(studyServiceServerConn)
	// Start webserver
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		// AllowOrigins:     []string{"https://inxp.de", "http://localhost:4200", "https://-1539512783514.firebaseapp.com"},
		AllowMethods:     []string{"POST", "GET", "PUT"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Content-Length"},
		ExposeHeaders:    []string{"Authorization", "Content-Type", "Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.GET("/", healthCheckHandle)
	v1 := router.Group("/v1")

	InitUserEndpoints(v1)
	InitTokenEndpoints(v1)
	InitStudyEndpoints(v1)

	InitExperimentalEndpoints(router.Group(""))

	log.Printf("gateway listening on port %s", conf.Port)
	log.Fatal(router.Run(":" + conf.Port))
}
