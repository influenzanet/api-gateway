package v1

import (
	"github.com/gin-gonic/gin"
	mw "github.com/influenzanet/api-gateway/pkg/protocols/http/middlewares"
)

func (h *HttpEndpoints) AddMessagingServiceAdminAPI(rg *gin.RouterGroup) {
	messagingGroup := rg.Group("/messaging")
	messagingGroup.Use(mw.ExtractToken())
	messagingGroup.Use(mw.ValidateToken(h.clients.UserManagement))
	{
		messagingGroup.GET("/email-templates", h.getEmailTemplatesHandl)
		messagingGroup.POST("/email-templates/save", mw.RequirePayload(), h.saveEmailTemplateHandl)
		messagingGroup.POST("/email-templates/delete", mw.RequirePayload(), h.deleteEmailTemplateHandl)

		messagingGroup.POST("/send-message/all-users", mw.RequirePayload(), h.sendMessageToAllUsersHandl)
	}
}
