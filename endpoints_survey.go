package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/influenzanet/api-gateway/api"
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

// InitSurveyEndpoints creates all API routes on the supplied RouterGroup
func InitSurveyEndpoints(rg *gin.RouterGroup) {
	survey := rg.Group("/survey") // TODO example endpoints
	{
		survey.POST("create", surveyCreateHandl)
		survey.POST("/submit", surveySubmitHandl)
		survey.POST("/update", surveyUpdateHandl)
		survey.POST("/get", surveyGetHandl)
		survey.POST("/get-all", surveyGetAllHandl)
	}
}

func surveyCreateHandl(c *gin.Context) {
	var req api.CreateSurveyReq
	if err := gjpb.JsonToPb(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: token handling
	req.Token = &api.TokenInfos{
		InstanceId: "default",
		Id:         "testuserid",
	}
	resp, err := clients.studyService.CreateSurvey(context.Background(), &req)
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
