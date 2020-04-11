package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/influenzanet/api-gateway/api"
	mw "github.com/influenzanet/api-gateway/middlewares"
	"github.com/influenzanet/api-gateway/utils"
	gjpb "github.com/phev8/gin-protobuf-json-converter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func connectToStudyServiceServer() *grpc.ClientConn {
	conn, err := grpc.Dial(conf.ServiceURLs.StudyService, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	return conn
}

// InitStudyEndpoints creates all API routes on the supplied RouterGroup
func InitStudyEndpoints(rg *gin.RouterGroup) {
	studySystem := rg.Group("/study-system")
	{
		studySystemWithAuth := studySystem.Group("")
		studySystemWithAuth.Use(mw.ExtractToken())
		studySystemWithAuth.Use(mw.ValidateToken(clients.authService))
		{
			studySystemWithAuth.POST("/create-study", mw.RequirePayload(), studySystemCreateStudyHandl)

			studyRoutes := studySystemWithAuth.Group("/study")
			{
				studyRoutes.POST("/save-survey", mw.RequirePayload(), saveSurveyToStudyHandl)
				studyRoutes.POST("/remove-survey", mw.RequirePayload(), removeSurveyFromStudyHandl)
			}
		}
	}

	survey := rg.Group("/survey") // TODO example endpoints
	{
		survey.POST("/submit", surveySubmitHandl)
	}
}

func studySystemCreateStudyHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(api.TokenInfos)

	var req api.NewStudyRequest
	if err := gjpb.JsonToPb(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = &token

	resp, err := clients.studyService.CreateNewStudy(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	gjpb.SendPBAsJSON(c, http.StatusOK, resp)
}

func saveSurveyToStudyHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(api.TokenInfos)

	var req api.AddSurveyReq
	if err := gjpb.JsonToPb(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = &token

	resp, err := clients.studyService.SaveSurveyToStudy(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	gjpb.SendPBAsJSON(c, http.StatusOK, resp)
}

func removeSurveyFromStudyHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(api.TokenInfos)

	var req api.SurveyReferenceRequest
	if err := gjpb.JsonToPb(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = &token

	resp, err := clients.studyService.RemoveSurveyFromStudy(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	gjpb.SendPBAsJSON(c, http.StatusOK, resp)
}

func surveySubmitHandl(c *gin.Context) {
	var req api.SubmitResponseReq
	if err := gjpb.JsonToPb(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: token handling
	req.Token = &api.TokenInfos{
		InstanceId: "default",
		Id:         "testuserid",
	}
	resp, err := clients.studyService.SubmitResponse(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	gjpb.SendPBAsJSON(c, http.StatusOK, resp)
}

func surveyUpdateHandl(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func surveyGetHandl(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func surveyGetAllHandl(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}
