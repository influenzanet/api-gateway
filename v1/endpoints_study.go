package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/influenzanet/api-gateway/api"
	mw "github.com/influenzanet/api-gateway/middlewares"
	"github.com/influenzanet/api-gateway/utils"
	gjpb "github.com/phev8/gin-protobuf-json-converter"
	"google.golang.org/grpc/status"
)

// initStudyEndpoints creates all API routes on the supplied RouterGroup
func initStudyEndpoints(rg *gin.RouterGroup) {
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

func studySystemCreateStudyHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api.TokenInfos)

	var req api.NewStudyRequest
	if err := gjpb.JsonToPb(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token

	resp, err := clients.StudyService.CreateNewStudy(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	gjpb.SendPBAsJSON(c, http.StatusOK, resp)
}

func saveSurveyToStudyHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api.TokenInfos)

	var req api.AddSurveyReq
	if err := gjpb.JsonToPb(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token

	resp, err := clients.StudyService.SaveSurveyToStudy(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	gjpb.SendPBAsJSON(c, http.StatusOK, resp)
}

func removeSurveyFromStudyHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api.TokenInfos)

	var req api.SurveyReferenceRequest
	if err := gjpb.JsonToPb(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token

	resp, err := clients.StudyService.RemoveSurveyFromStudy(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	gjpb.SendPBAsJSON(c, http.StatusOK, resp)
}

func getAssignedSurveyHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api.TokenInfos)

	var req api.SurveyReferenceRequest
	if err := gjpb.JsonToPb(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	resp, err := clients.StudyService.GetAssignedSurvey(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	gjpb.SendPBAsJSON(c, http.StatusOK, resp)
}

func submitSurveyResponseHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api.TokenInfos)

	var req api.SubmitResponseReq
	if err := gjpb.JsonToPb(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	resp, err := clients.StudyService.SubmitResponse(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	gjpb.SendPBAsJSON(c, http.StatusOK, resp)
}
