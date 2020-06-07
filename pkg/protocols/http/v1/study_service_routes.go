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
		// all surveys accross studies:
		studiesGroup.GET("/all-assigned-surveys", h.getAllAssignedSurveysHandl)
	}

	studyGroup := rg.Group("/study")
	studyGroup.Use(mw.ExtractToken())
	studyGroup.Use(mw.ValidateToken(h.clients.UserManagement))
	{
		studyGroup.GET("/:studyKey/survey-infos", h.getStudySurveyInfosHandl)
		studyGroup.POST("/:studyKey/enter", mw.RequirePayload(), h.enterStudyHandl)
		studyGroup.GET("/:studyKey/survey/:surveyKey", mw.RequirePayload(), h.getAssignedSurveyHandl)
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
		studiesGroup.POST("/:studyKey/get", mw.RequirePayload(), h.getStudyHandl)
		studiesGroup.POST("/:studyKey/get-survey", mw.RequirePayload(), h.getSurveyDefForStudyHandl)
		studiesGroup.POST("/:studyKey/save-survey", mw.RequirePayload(), h.saveSurveyToStudyHandl)
		studiesGroup.POST("/:studyKey/remove-survey", mw.RequirePayload(), h.removeSurveyFromStudyHandl)

		studiesGroup.POST("/:studyKey/save-member", mw.RequirePayload(), h.studySaveMemberHandl)
		studiesGroup.POST("/:studyKey/remove-member", mw.RequirePayload(), h.studyRemoveMemberHandl)
		studiesGroup.POST("/:studyKey/rules", mw.RequirePayload(), h.saveStudyRulesHandl)
		studiesGroup.POST("/:studyKey/status", mw.RequirePayload(), h.saveStudyStatusHandl)
		studiesGroup.POST("/:studyKey/props", mw.RequirePayload(), h.saveStudyPropsHandl)
		studiesGroup.POST("/:studyKey/delete", mw.RequirePayload(), h.deleteStudyHandl)
	}

	responsesGroup := rg.Group("/data/:studyKey")
	responsesGroup.Use(mw.ExtractToken())
	responsesGroup.Use(mw.ValidateToken(h.clients.UserManagement))
	{
		responsesGroup.GET("/statistics", h.getSurveyResponseStatisticsHandl)
		responsesGroup.GET("/responses", h.getSurveyResponsesHandl)
	}
}
