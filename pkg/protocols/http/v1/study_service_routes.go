package v1

/*
// initStudyEndpoints creates all API routes on the supplied RouterGroup
func (h *HttpEndpoints) initStudyEndpoints(rg *gin.RouterGroup) {

		studySystem := rg.Group("/study-system")
		{
			studySystemWithAuth := studySystem.Group("")
			studySystemWithAuth.Use(mw.ExtractToken())
			studySystemWithAuth.Use(mw.ValidateToken(clients.AuthService))
			{
				studySystemWithAuth.POST("/create-study", mw.RequirePayload(), studySystemCreateStudyHandl)

				studyRoutes := studySystemWithAuth.Group("/study")
				{
					studyRoutes.POST("/save-survey", mw.RequirePayload(), saveSurveyToStudyHandl)
					studyRoutes.POST("/remove-survey", mw.RequirePayload(), removeSurveyFromStudyHandl)
					studyRoutes.POST("/get-assigned-survey", mw.RequirePayload(), getAssignedSurveyHandl)
					studyRoutes.POST("/submit-response", mw.RequirePayload(), submitSurveyResponseHandl)
				}
			}
		}
}
*/
