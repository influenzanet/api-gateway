package v1

import (
	"github.com/gin-gonic/gin"
	mw "github.com/influenzanet/api-gateway/pkg/protocols/http/middlewares"
)

func (h *HttpEndpoints) AddStudyServiceParticipantAPI(rg *gin.RouterGroup) {
	studiesGroup := rg.Group("/studies")
	studiesGroup.Use(mw.ExtractToken())
	studiesGroup.Use(mw.ValidateToken(h.clients.UserManagement))
	studiesGroup.Use(mw.CheckAccountConfirmed())
	{
		studiesGroup.GET("/for-user-profiles", h.getStudiesForUserHandl)
		studiesGroup.GET("/active", h.getAllActiveStudiesHandl)
		// all surveys accross studies:
		studiesGroup.GET("/all-assigned-surveys", h.getAllAssignedSurveysHandl)
	}

	studyGroup := rg.Group("/study")
	studyGroup.Use(mw.ExtractToken())
	studyGroup.Use(mw.ValidateToken(h.clients.UserManagement))
	studyGroup.Use(mw.CheckAccountConfirmed())
	{
		studyGroup.GET("/:studyKey/survey-infos", h.getStudySurveyInfosHandl)
		studyGroup.POST("/:studyKey/enter", mw.RequirePayload(), h.enterStudyHandl)
		studyGroup.GET("/:studyKey/survey/:surveyKey", h.getAssignedSurveyHandl)
		studyGroup.POST("/:studyKey/submit-response", mw.RequirePayload(), h.submitSurveyResponseHandl)
		studyGroup.POST("/:studyKey/postpone-survey", mw.RequirePayload(), h.postponeSurveyHandl)
		studyGroup.POST("/:studyKey/leave", mw.RequirePayload(), h.leaveStudyHandl)
	}
}

func (h *HttpEndpoints) AddStudyServiceAdminAPI(rg *gin.RouterGroup) {
	studiesGroup := rg.Group("/studies")
	studiesGroup.Use(mw.ExtractToken())
	studiesGroup.Use(mw.ValidateToken(h.clients.UserManagement))
	{
		studiesGroup.POST("", mw.RequirePayload(), h.studySystemCreateStudyHandl)
		studiesGroup.GET("", h.getAllStudiesHandl)

	}

	studyGroup := rg.Group("/study")
	studyGroup.Use(mw.ExtractToken())
	studyGroup.Use(mw.ValidateToken(h.clients.UserManagement))
	{
		studyGroup.GET("/:studyKey", h.getStudyHandl)
		studyGroup.GET("/:studyKey/survey/:surveyKey", h.getSurveyDefForStudyHandl)
		studyGroup.POST("/:studyKey/save-survey", mw.RequirePayload(), h.saveSurveyToStudyHandl)
		studyGroup.POST("/:studyKey/remove-survey", mw.RequirePayload(), h.removeSurveyFromStudyHandl)

		studyGroup.POST("/:studyKey/save-member", mw.RequirePayload(), h.studySaveMemberHandl)
		studyGroup.POST("/:studyKey/remove-member", mw.RequirePayload(), h.studyRemoveMemberHandl)
		studyGroup.POST("/:studyKey/rules", mw.RequirePayload(), h.saveStudyRulesHandl)
		studyGroup.POST("/:studyKey/status", mw.RequirePayload(), h.saveStudyStatusHandl)
		studyGroup.POST("/:studyKey/props", mw.RequirePayload(), h.saveStudyPropsHandl)
		studyGroup.DELETE("/:studyKey", mw.RequirePayload(), h.deleteStudyHandl)
	}

	responsesGroup := rg.Group("/data/:studyKey")
	responsesGroup.Use(mw.ExtractToken())
	responsesGroup.Use(mw.ValidateToken(h.clients.UserManagement))
	{
		responsesGroup.GET("/statistics", h.getSurveyResponseStatisticsHandl)
		responsesGroup.GET("/responses", h.getSurveyResponsesHandl)
	}
}
