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
	}
}

/*

	user := rg.Group("/user")
	user.Use(mw.ExtractToken())
	user.Use(mw.ValidateToken(h.clients.UserManagement))

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
