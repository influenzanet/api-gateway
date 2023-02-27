package main

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/coneno/logger"
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
	conf.LogLevel = models.GetLogLevel()
	conf.DebugMode = os.Getenv("DEBUG_MODE") == "true"
	conf.Port = os.Getenv("MANAGEMENT_API_GATEWAY_LISTEN_PORT")
	conf.ServiceURLs.UserManagement = os.Getenv("ADDR_USER_MANAGEMENT_SERVICE")
	conf.ServiceURLs.StudyService = os.Getenv("ADDR_STUDY_SERVICE")
	conf.ServiceURLs.MessagingService = os.Getenv("ADDR_MESSAGING_SERVICE")
	conf.AllowOrigins = strings.Split(os.Getenv("CORS_ALLOW_ORIGINS"), ",")

	conf.MaxMsgSize = models.DefaultGRPCMaxMsgSize
	ms, err := strconv.Atoi(os.Getenv(models.ENV_GRPC_MAX_MSG_SIZE))
	if err != nil {
		logger.Info.Printf("using default max gRPC message size: %d", models.DefaultGRPCMaxMsgSize)
	} else {
		conf.MaxMsgSize = ms
	}

	conf.UseEndpoints.DeleteParticipantData = os.Getenv("USE_DELETE_PARTICIPANT_DATA_ENDPOINT") == "true"

	// SAML configs
	conf.UseEndpoints.LoginWithExternalIDP = os.Getenv("USE_LOGIN_WITH_EXTERNAL_IDP_ENDPOINT") == "true"

	if conf.UseEndpoints.LoginWithExternalIDP {
		conf.SAMLConfig = &models.SAMLConfig{
			IDPUrl:                   os.Getenv("SAML_IDP_URL"),                   // arbitrary name to refer to IDP in the logs
			SPRootUrl:                os.Getenv("SAML_SERVICE_PROVIDER_ROOT_URL"), // url of the management api
			EntityID:                 os.Getenv("SAML_ENTITY_ID"),
			MetaDataURL:              os.Getenv("SAML_IDP_METADATA_URL"),
			SessionCertPath:          os.Getenv("SAML_SESSION_CERT_PATH"),
			SessionKeyPath:           os.Getenv("SAML_SESSION_KEY_PATH"),
			TemplatePathLoginSuccess: os.Getenv("SAML_TEMPLATE_PATH_LOGIN_SUCCESS"),
		}

		if len(conf.SAMLConfig.IDPUrl) > 0 {
			conf.AllowOrigins = append(conf.AllowOrigins, conf.SAMLConfig.IDPUrl)
		}
	}

	// Mutual TLS configs
	conf.TLSPaths = models.TLSPaths{
		ServerCertPath: os.Getenv(models.ENV_MUTUAL_TLS_SERVER_CERT),
		ServerKeyPath:  os.Getenv(models.ENV_MUTUAL_TLS_SERVER_KEY),
		CACertPath:     os.Getenv(models.ENV_MUTUAL_TLS_CA_CERT),
	}
}

func init() {
	grpcClients = &models.APIClients{}

	initConfig()
	if !conf.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}
	logger.SetLevel(conf.LogLevel)
}

func healthCheckHandle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func main() {
	umClient, userManagementClose := gc.ConnectToUserManagement(conf.ServiceURLs.UserManagement, conf.MaxMsgSize)
	defer userManagementClose()
	studyClient, studyServiceClose := gc.ConnectToStudyService(conf.ServiceURLs.StudyService, conf.MaxMsgSize)
	defer studyServiceClose()
	messagingClient, messagingServiceClose := gc.ConnectToMessagingService(conf.ServiceURLs.MessagingService, conf.MaxMsgSize)
	defer messagingServiceClose()

	grpcClients.UserManagement = umClient
	grpcClients.StudyService = studyClient
	grpcClients.MessagingService = messagingClient

	// Start webserver
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		// AllowAllOrigins: true,
		AllowOrigins:     conf.AllowOrigins,
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Content-Length"},
		ExposeHeaders:    []string{"Authorization", "Content-Type", "Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.GET("/", healthCheckHandle)
	v1Root := router.Group("/v1")

	v1APIHandlers := v1.NewHTTPHandler(grpcClients, conf.UseEndpoints, conf.SAMLConfig)
	v1APIHandlers.AddServiceStatusAPI(v1Root)
	v1APIHandlers.AddUserManagementAdminAPI(v1Root)
	v1APIHandlers.AddStudyServiceAdminAPI(v1Root)
	v1APIHandlers.AddMessagingServiceAdminAPI(v1Root)

	logger.Info.Printf("gateway listening on port %s", conf.Port)

	// Create tls config for mutual TLS
	tlsConfig, err := getTLSConfig(conf.TLSPaths)
	if err != nil {
		logger.Error.Fatal(err)
	}

	server := &http.Server{
		Addr:      ":" + conf.Port,
		Handler:   router,
		TLSConfig: tlsConfig,
	}

	err = server.ListenAndServeTLS(conf.TLSPaths.ServerCertPath, conf.TLSPaths.ServerKeyPath)
	if err != nil {
		logger.Error.Fatal(err)
	}
}

func getTLSConfig(paths models.TLSPaths) (*tls.Config, error) {
	serverCert, err := tls.LoadX509KeyPair(paths.ServerCertPath, paths.ServerKeyPath)
	if err != nil {
		return nil, err
	}

	caCert, err := os.ReadFile(paths.CACertPath)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	return &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
	}, nil
}
