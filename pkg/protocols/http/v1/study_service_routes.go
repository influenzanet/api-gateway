package v1

import (
	"github.com/gin-gonic/gin"
	mw "github.com/influenzanet/api-gateway/pkg/protocols/http/middlewares"
)

func (h *HttpEndpoints) AddStudyServiceParticipantAPI(rg *gin.RouterGroup) {
	studiesGroup := rg.Group("/studies")
	studiesGroup.Use(mw.ExtractToken())
	studiesGroup.Use(mw.ValidateToken(h.clients.UserManagement))
	{
		studiesGroup.POST("/study/enter", mw.RequirePayload(), h.enterStudyHandl)
		studiesGroup.POST("/study/get-assigned-survey", mw.RequirePayload(), h.getAssignedSurveyHandl)
		studiesGroup.POST("/study/submit-response", mw.RequirePayload(), h.submitSurveyResponseHandl)
		studiesGroup.POST("/study/postpone-survey", mw.RequirePayload(), h.postponeSurveyHandl)
	}
}

func (h *HttpEndpoints) AddStudyServiceAdminAPI(rg *gin.RouterGroup) {
	studiesGroup := rg.Group("/studies")
	studiesGroup.Use(mw.ExtractToken())
	studiesGroup.Use(mw.ValidateToken(h.clients.UserManagement))
	{
		studiesGroup.POST("", mw.RequirePayload(), h.studySystemCreateStudyHandl)
		studiesGroup.POST("/study/save-survey", mw.RequirePayload(), h.saveSurveyToStudyHandl)
		studiesGroup.POST("/remove-survey", mw.RequirePayload(), h.removeSurveyFromStudyHandl)
	}
}
