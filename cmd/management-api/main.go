package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/influenzanet/api-gateway/pkg/models"
	gc "github.com/influenzanet/api-gateway/pkg/protocols/grpc/clients"
	v1 "github.com/influenzanet/api-gateway/pkg/protocols/http/v1"
)

// Conf holds all static configuration information
var conf models.Config
var grpcClients *models.APIClients

func initConfig() {
	conf.DebugMode = os.Getenv("DEBUG_MODE") == "true"
	conf.Port = os.Getenv("MANAGEMENT_API_GATEWAY_LISTEN_PORT")
	conf.ServiceURLs.UserManagement = os.Getenv("ADDR_USER_MANAGEMENT_SERVICE")
	conf.ServiceURLs.StudyService = os.Getenv("ADDR_STUDY_SERVICE")
	conf.ServiceURLs.MessagingService = os.Getenv("ADDR_MESSAGING_SERVICE")
	conf.AllowOrigins = strings.Split(os.Getenv("CORS_ALLOW_ORIGINS"), ",")
}

func init() {
	grpcClients = &models.APIClients{}

	initConfig()
	if !conf.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}
}

func healthCheckHandle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func main() {
	umClient, close := gc.ConnectToUserManagement(conf.ServiceURLs.UserManagement)
	defer close()
	studyClient, studyServiceClose := gc.ConnectToStudyService(conf.ServiceURLs.StudyService)
	defer studyServiceClose()
	messagingClient, messagingServiceClose := gc.ConnectToMessagingService(conf.ServiceURLs.MessagingService)
	defer messagingServiceClose()

	grpcClients.UserManagement = umClient
	grpcClients.StudyService = studyClient
	grpcClients.MessagingService = messagingClient

	// Start webserver
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		// AllowAllOrigins: true,
		AllowOrigins:     conf.AllowOrigins,
		AllowMethods:     []string{"POST", "GET", "PUT"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Content-Length"},
		ExposeHeaders:    []string{"Authorization", "Content-Type", "Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.GET("/", healthCheckHandle)
	v1Root := router.Group("/v1")

	v1APIHandlers := v1.NewHTTPHandler(grpcClients)
	v1APIHandlers.AddServiceStatusAPI(v1Root)
	v1APIHandlers.AddUserManagementAdminAPI(v1Root)
	v1APIHandlers.AddStudyServiceAdminAPI(v1Root)
	v1APIHandlers.AddMessagingServiceAdminAPI(v1Root)

	log.Printf("gateway listening on port %s", conf.Port)
	log.Fatal(router.Run(":" + conf.Port))
}

// InitExperimentalEndpoints(router.Group(""))
