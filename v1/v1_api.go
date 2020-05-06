package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/influenzanet/api-gateway/models"
)

var clients = &models.APIClients{}

func InitAPI(clientRef *models.APIClients, rootGroup *gin.RouterGroup) {
	clients = clientRef
	initStatusEndpoints(rootGroup)
	initStudyEndpoints(rootGroup)
	initUserManagementEndpoints(rootGroup)
	initAuthEndpoints(rootGroup)
}
