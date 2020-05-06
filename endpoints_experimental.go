package main

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

// InitUserEndpoints creates all API routes on the supplied RouterGroup
func InitExperimentalEndpoints(rg *gin.RouterGroup) {
	covidAppRoutes := rg.Group("/covidapp")
	{
		covidAppRoutes.POST("/register-app-user", covidAppRegisterHandl)
		covidAppRoutes.POST("/login-app-user", covidAppLoginHandl)

		covidAppStudyRoutes := covidAppRoutes.Group("/study")
		covidAppStudyRoutes.Use(mw.ExtractToken())
		covidAppStudyRoutes.Use(mw.ValidateToken(clients.AuthService))
		{
			covidAppStudyRoutes.POST("/enter", mw.RequirePayload(), covidAppEnterStudyHandl)
			covidAppStudyRoutes.GET("/fetch-assigned-surveys", covidAppGetAssignedSurveysHandl)
			covidAppStudyRoutes.GET("/fetch-survey-infos", covidAppGetSurveyInfosHandl)
			covidAppStudyRoutes.POST("/status-report", mw.RequirePayload(), covidAppSubmitStatusReportHandl)
		}
	}
}

type appCredentials struct {
	AppToken   string `json:"appToken"`
	InstanceID string `json:"instanceID"`
	UserID     string `json:"userID"`
	Password   string `json:"password"`
}

func hasInstance(instances []string, instanceID string) bool {
	for _, i := range instances {
		if i == instanceID {
			return true
		}
	}
	return false
}

func covidAppRegisterHandl(c *gin.Context) {
	var req appCredentials
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate App token before continue with request
	resp, err := clients.AuthService.ValidateAppToken(context.Background(), &api.AppTokenRequest{Token: req.AppToken})
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	if !hasInstance(resp.Instances, req.InstanceID) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "app token cannot be used for given instanceID"})
	}

	// Register user
	signUpReq := api.SignupWithEmailMsg{
		InstanceId: req.InstanceID,
		Email:      req.UserID + "@app-id.no-reply",
		Password:   req.Password,
	}
	token, err := clients.AuthService.SignupWithEmail(context.Background(), &signUpReq)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	gjpb.SendPBAsJSON(c, http.StatusOK, token)
}

func covidAppLoginHandl(c *gin.Context) {
	var req appCredentials
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate App token before continue with request
	resp, err := clients.AuthService.ValidateAppToken(context.Background(), &api.AppTokenRequest{Token: req.AppToken})
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	if !hasInstance(resp.Instances, req.InstanceID) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "app token cannot be used for given instanceID"})
	}

	loginReq := api.LoginWithEmailMsg{
		InstanceId: req.InstanceID,
		Email:      req.UserID + "@app-id.no-reply",
		Password:   req.Password,
	}
	token, err := clients.AuthService.LoginWithEmail(context.Background(), &loginReq)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	gjpb.SendPBAsJSON(c, http.StatusOK, token)
}

func covidAppEnterStudyHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(api.TokenInfos)

	var req api.EnterStudyRequest
	if err := gjpb.JsonToPb(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = &token

	// Enter study
	resp, err := clients.StudyService.EnterStudy(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	gjpb.SendPBAsJSON(c, http.StatusOK, resp)
}

func covidAppGetAssignedSurveysHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(api.TokenInfos)

	resp, err := clients.StudyService.GetAssignedSurveys(context.Background(), &token)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	gjpb.SendPBAsJSON(c, http.StatusOK, resp)
}

func covidAppGetSurveyInfosHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(api.TokenInfos)

	req := api.StudyReferenceReq{
		Token:    &token,
		StudyKey: "covid-19",
	}
	resp, err := clients.StudyService.GetStudySurveyInfos(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	gjpb.SendPBAsJSON(c, http.StatusOK, resp)
}

func covidAppSubmitStatusReportHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(api.TokenInfos)

	var req api.StatusReportRequest
	if err := gjpb.JsonToPb(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = &token

	resp, err := clients.StudyService.SubmitStatusReport(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	gjpb.SendPBAsJSON(c, http.StatusOK, resp)
}
