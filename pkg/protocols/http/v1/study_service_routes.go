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
		studiesGroup.DELETE("/confidential-data", h.removeConfidentialDataHandl) // ?profiles=todo1,todo2 (optional)
		if h.useEndpoints.DeleteParticipantData {
			studiesGroup.DELETE("/participant-data", h.deleteParticipantDataHandl)
		}
	}

	tempStudyParticipantGroup := rg.Group("/temp-participant")
	{
		tempStudyParticipantGroup.GET("/register", h.registerTempParticipant)           // ?instance=todo&study=todo
		tempStudyParticipantGroup.GET("/surveys", h.getTempParticipantSurveys)          // ?instance=todo&study=todo&pid=todo
		tempStudyParticipantGroup.GET("/survey", h.getSurveyDefForTempParticipantHandl) // ?instance=todo&study=todo&pid=todo&survey=todo
		tempStudyParticipantGroup.POST("/submit-response", mw.RequirePayload(), h.submitSurveyResponseForTempParticipantHandl)
	}

	rg.GET("/reports",
		mw.ExtractToken(),
		mw.ValidateToken(h.clients.UserManagement),
		mw.CheckAccountConfirmed(),
		h.getReportsForParticipant) // ?studies=todo1,todo2&profileIds=todo1,todo2&from=time1&until=time2&reportKey=todo3&ignoreReports=todo1,todo2

	studyGroup := rg.Group("/study")
	studyGroup.Use(mw.ExtractToken())
	studyGroup.Use(mw.ValidateToken(h.clients.UserManagement))
	{

		studyGroup.GET("/:studyKey/survey-infos", h.getStudySurveyInfosHandl)
		studyGroup.POST("/:studyKey/enter", mw.RequirePayload(), h.enterStudyHandl)
		studyGroup.GET("/:studyKey/survey/:surveyKey", h.getAssignedSurveyHandl)
		studyGroup.POST("/:studyKey/submit-response", mw.RequirePayload(), h.submitSurveyResponseHandl)
		studyGroup.POST("/:studyKey/leave", mw.RequirePayload(), h.leaveStudyHandl)
		studyGroup.POST("/:studyKey/file-upload", mw.RequirePayload(), h.uploadParticipantFileReq)
		studyGroup.POST("/:studyKey/delete-files", mw.RequirePayload(), h.deleteParticipantFilesReq)
		studyGroup.POST("/:studyKey/assume-temp-participant", mw.RequirePayload(), h.convertTempToActiveParticipant)
	}
}

func (h *HttpEndpoints) AddStudyServiceAdminAPI(rg *gin.RouterGroup) {
	studiesGroup := rg.Group("/studies")
	studiesGroup.Use(mw.ExtractToken())
	studiesGroup.Use(mw.ValidateToken(h.clients.UserManagement))
	studiesGroup.Use(mw.CheckAccountConfirmed())
	{
		studiesGroup.POST("", mw.RequirePayload(), h.studySystemCreateStudyHandl)
		studiesGroup.GET("", h.getAllStudiesHandl)

	}

	studyGroup := rg.Group("/study")
	studyGroup.Use(mw.ExtractToken())
	studyGroup.Use(mw.ValidateToken(h.clients.UserManagement))
	studyGroup.Use(mw.CheckAccountConfirmed())
	{
		studyGroup.GET("/:studyKey", h.getStudyHandl)

		studyGroup.GET("/:studyKey/surveys", h.getStudySurveyInfosHandl)
		studyGroup.GET("/:studyKey/survey-keys", h.getStudySurveySurveyKeysHandl)
		studyGroup.POST("/:studyKey/surveys", mw.RequirePayload(), h.saveSurveyToStudyHandl)
		studyGroup.GET("/:studyKey/survey/:surveyKey/versions", h.getSurveyVersionInfosHandl)
		studyGroup.GET("/:studyKey/survey/:surveyKey", h.getSurveyDefForStudyHandl)
		studyGroup.GET("/:studyKey/survey/:surveyKey/:versionID", h.getSurveyDefForStudyHandl)
		studyGroup.DELETE("/:studyKey/survey/:surveyKey", h.unpublishSurveyHandl)
		studyGroup.DELETE("/:studyKey/survey/:surveyKey/:versionID", h.removeSurveyVersionHandl)

		studyGroup.POST("/:studyKey/save-member", mw.RequirePayload(), h.studySaveMemberHandl)
		studyGroup.POST("/:studyKey/remove-member", mw.RequirePayload(), h.studyRemoveMemberHandl)
		studyGroup.POST("/:studyKey/rules", mw.RequirePayload(), h.saveStudyRulesHandl)
		studyGroup.POST("/:studyKey/run-rules", mw.RequirePayload(), h.runCustomStudyRulesHandl)
		studyGroup.POST("/:studyKey/run-rules-for-single-participant", mw.RequirePayload(), h.runCustomStudyRulesForSingleParticipantHandl)
		studyGroup.POST("/:studyKey/status", mw.RequirePayload(), h.saveStudyStatusHandl)
		studyGroup.POST("/:studyKey/props", mw.RequirePayload(), h.saveStudyPropsHandl)
		studyGroup.GET("/:studyKey/notification-subscriptions", h.getStudyNotificationSubscriptionsHandl)
		studyGroup.POST("/:studyKey/notification-subscriptions", mw.RequirePayload(), h.saveStudyNotificationSubscriptionsHandl)
		studyGroup.DELETE("/:studyKey", mw.RequirePayload(), h.deleteStudyHandl)
	}

	responsesGroup := rg.Group("/data/:studyKey")
	responsesGroup.Use(mw.ExtractToken())
	responsesGroup.Use(mw.ValidateToken(h.clients.UserManagement))
	responsesGroup.Use(mw.CheckAccountConfirmed())
	{
		responsesGroup.GET("/statistics", h.getSurveyResponseStatisticsHandl)
		responsesGroup.GET("/participants/all", h.getParticipantStatesForStudy) // ?&status=active(opt)
		responsesGroup.GET("/file-infos", h.getFileInfosForStudy)               // ?fileType=todo&from=time1&until=time2&participant=todo
		responsesGroup.GET("/file", h.getParticipantFile)                       // ?id=todo
		responsesGroup.POST("/delete-files", mw.RequirePayload(), h.deleteParticipantFilesReq)
		responsesGroup.GET("/reports", h.getReportsForStudy) // ?reportKey=todo&from=time1&until=time2&participant=todo
		responsesGroup.GET("/responses", h.getSurveyResponsesHandl)
		responsesGroup.POST("/fetch-confidential-responses", h.getConfidentialResponsesHandl)
		responsesGroup.GET("/participant", h.getParticipantStateByID)
		responsesGroup.GET("/participants", h.getParticipantStatesWithPagination)
		responsesGroup.GET("/rules", h.getCurrentStudyRules)

		surveyResponsesGroup := responsesGroup.Group("/survey/:surveyKey")
		{
			surveyResponsesGroup.GET("/response", h.getResponseWideFormatCSV)
			surveyResponsesGroup.GET("/response/long-format", h.getResponseLongFormatCSV)
			surveyResponsesGroup.GET("/response/json", h.getResponseFlatJSON)
			surveyResponsesGroup.GET("/survey-info", h.getSurveyInfoPreview)
			surveyResponsesGroup.GET("/survey-info/csv", h.getSurveyInfoPreviewCSV)
		}
	}
}
