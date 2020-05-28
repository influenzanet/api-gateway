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
		studiesGroup.GET("/for-user-profiles", h.getStudiesForUserHandl)
		studiesGroup.GET("/active", h.getAllActiveStudiesHandl)
		studiesGroup.POST("/study/get-survey-infos", mw.RequirePayload(), h.getStudySurveyInfosHandl)

		studiesGroup.POST("/study/enter", mw.RequirePayload(), h.enterStudyHandl)
		studiesGroup.POST("/study/get-assigned-survey", mw.RequirePayload(), h.getAssignedSurveyHandl)
		studiesGroup.POST("/study/submit-response", mw.RequirePayload(), h.submitSurveyResponseHandl)
		studiesGroup.POST("/study/postpone-survey", mw.RequirePayload(), h.postponeSurveyHandl)
		studiesGroup.POST("/study/leave", mw.RequirePayload(), h.leaveStudyHandl)

		// all surveys accross studies:
		studiesGroup.GET("/all-assigned-surveys", h.getAllAssignedSurveysHandl)
	}
}

func (h *HttpEndpoints) AddStudyServiceAdminAPI(rg *gin.RouterGroup) {
	studiesGroup := rg.Group("/studies")
	studiesGroup.Use(mw.ExtractToken())
	studiesGroup.Use(mw.ValidateToken(h.clients.UserManagement))
	{
		studiesGroup.POST("", mw.RequirePayload(), h.studySystemCreateStudyHandl)
		studiesGroup.GET("", h.getAllStudiesHandl)
		studiesGroup.POST("/study/get", mw.RequirePayload(), h.getStudyHandl)
		studiesGroup.POST("/study/get-survey", mw.RequirePayload(), h.getSurveyDefForStudyHandl)
		studiesGroup.POST("/study/save-survey", mw.RequirePayload(), h.saveSurveyToStudyHandl)
		studiesGroup.POST("/study/remove-survey", mw.RequirePayload(), h.removeSurveyFromStudyHandl)
	}
}
