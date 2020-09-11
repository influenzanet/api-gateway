package v1

import (
	"github.com/gin-gonic/gin"
	mw "github.com/influenzanet/api-gateway/pkg/protocols/http/middlewares"
)

func (h *HttpEndpoints) AddMessagingServiceAdminAPI(rg *gin.RouterGroup) {
	messagingGroup := rg.Group("/messaging")
	messagingGroup.Use(mw.ExtractToken())
	messagingGroup.Use(mw.ValidateToken(h.clients.UserManagement))
	messagingGroup.Use(mw.CheckAccountConfirmed())
	{
		messagingGroup.GET("/email-templates", h.getEmailTemplatesHandl)
		messagingGroup.POST("/email-templates", mw.RequirePayload(), h.saveEmailTemplateHandl)
		messagingGroup.POST("/email-templates/delete", mw.RequirePayload(), h.deleteEmailTemplateHandl)

		messagingGroup.GET("/auto-messages", h.getAutoMessagesHandl)
		messagingGroup.POST("/auto-messages", mw.RequirePayload(), h.saveAutoMessageHandl)
		messagingGroup.DELETE("/auto-message/:id", h.deleteAutoMessageHandl)

		messagingGroup.POST("/send-message/all-users", mw.RequirePayload(), h.sendMessageToAllUsersHandl)
		messagingGroup.POST("/send-message/study-participants", mw.RequirePayload(), h.sendMessageToStudyParticipantsHandl)
	}
}
