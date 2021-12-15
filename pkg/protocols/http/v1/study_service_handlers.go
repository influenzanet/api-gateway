package v1

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/influenzanet/api-gateway/pkg/utils"
	"github.com/influenzanet/go-utils/pkg/api_types"
	"github.com/influenzanet/study-service/pkg/api"
	studyAPI "github.com/influenzanet/study-service/pkg/api"
	"google.golang.org/grpc/status"
)

func (h *HttpEndpoints) enterStudyHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.EnterStudyRequest
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.StudyKey = c.Param("studyKey")
	req.Token = token
	resp, err := h.clients.StudyService.EnterStudy(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) studySystemCreateStudyHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.NewStudyRequest
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token

	resp, err := h.clients.StudyService.CreateNewStudy(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) saveSurveyToStudyHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.AddSurveyReq
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	req.StudyKey = c.Param("studyKey")

	resp, err := h.clients.StudyService.SaveSurveyToStudy(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) removeSurveyFromStudyHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.SurveyReferenceRequest
	req.Token = token
	req.StudyKey = c.Param("studyKey")
	req.SurveyKey = c.Param("surveyKey")

	resp, err := h.clients.StudyService.RemoveSurveyFromStudy(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) registerTempParticipant(c *gin.Context) {
	var req studyAPI.RegisterTempParticipantReq

	req.StudyKey = c.Query("study")
	req.InstanceId = c.Query("instance")
	resp, err := h.clients.StudyService.RegisterTemporaryParticipant(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) getTempParticipantSurveys(c *gin.Context) {
	var req studyAPI.GetAssignedSurveysForTemporaryParticipantReq
	req.StudyKey = c.Query("study")
	req.InstanceId = c.Query("instance")
	req.TemporaryParticipantId = c.Query("pid")
	resp, err := h.clients.StudyService.GetAssignedSurveysForTemporaryParticipant(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) getSurveyDefForTempParticipantHandl(c *gin.Context) {
	// ?instance=todo&study=todo&pid=todo&survey=todo
	var req studyAPI.SurveyReferenceRequest
	req.SurveyKey = c.Query("survey")
	req.StudyKey = c.Query("study")
	req.InstanceId = c.Query("instance")
	req.TemporaryParticipantId = c.Query("pid")
	resp, err := h.clients.StudyService.GetAssignedSurvey(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) submitSurveyResponseForTempParticipantHandl(c *gin.Context) {
	var req studyAPI.SubmitResponseReq
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.clients.StudyService.SubmitResponse(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) convertTempToActiveParticipant(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.ConvertTempParticipantReq
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	req.StudyKey = c.Param("studyKey")

	resp, err := h.clients.StudyService.ConvertTemporaryToParticipant(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) getAssignedSurveyHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.SurveyReferenceRequest
	req.Token = token
	req.ProfileId = c.Query("pid")
	req.StudyKey = c.Param("studyKey")
	req.SurveyKey = c.Param("surveyKey")
	resp, err := h.clients.StudyService.GetAssignedSurvey(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) submitSurveyResponseHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.SubmitResponseReq
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	req.StudyKey = c.Param("studyKey")
	resp, err := h.clients.StudyService.SubmitResponse(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) getAllAssignedSurveysHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	resp, err := h.clients.StudyService.GetAssignedSurveys(context.Background(), token)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) leaveStudyHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.LeaveStudyMsg
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	req.StudyKey = c.Param("studyKey")
	resp, err := h.clients.StudyService.LeaveStudy(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) deleteParticipantDataHandl(c *gin.Context) {
	// token := c.MustGet("validatedToken").(*api_types.TokenInfos)
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
	/*
		var req studyAPI.StudyStatusReq
		if err := h.JsonToProto(c, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		req.Token = token
		req.StudyKey = c.Param("studyKey")
		resp, err := h.clients.StudyService.SaveStudyStatus(context.Background(), &req)
		if err != nil {
			st := status.Convert(err)
			c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		h.SendProtoAsJSON(c, http.StatusOK, resp)
	*/
}

func (h *HttpEndpoints) getStudiesForUserHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.GetStudiesForUserReq
	req.Token = token
	resp, err := h.clients.StudyService.GetStudiesForUser(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) getAllActiveStudiesHandl(c *gin.Context) {
	req := c.MustGet("validatedToken").(*api_types.TokenInfos)

	resp, err := h.clients.StudyService.GetActiveStudies(context.Background(), req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) getStudySurveyInfosHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.StudyReferenceReq
	req.Token = token
	req.StudyKey = c.Param("studyKey")
	resp, err := h.clients.StudyService.GetStudySurveyInfos(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) getAllStudiesHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	resp, err := h.clients.StudyService.GetAllStudies(context.Background(), token)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) getStudyHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.StudyReferenceReq
	req.StudyKey = c.Param("studyKey")
	req.Token = token
	resp, err := h.clients.StudyService.GetStudy(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) getSurveyDefForStudyHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.SurveyReferenceRequest
	req.Token = token
	req.StudyKey = c.Param("studyKey")
	req.SurveyKey = c.Param("surveyKey")
	resp, err := h.clients.StudyService.GetSurveyDefForStudy(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) studySaveMemberHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.StudyMemberReq
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	req.StudyKey = c.Param("studyKey")
	resp, err := h.clients.StudyService.SaveStudyMember(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) studyRemoveMemberHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.StudyMemberReq
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.StudyKey = c.Param("studyKey")
	req.Token = token
	resp, err := h.clients.StudyService.RemoveStudyMember(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) saveStudyRulesHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.StudyRulesReq
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	req.StudyKey = c.Param("studyKey")
	resp, err := h.clients.StudyService.SaveStudyRules(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) runCustomStudyRulesHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.StudyRulesReq
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	req.StudyKey = c.Param("studyKey")
	resp, err := h.clients.StudyService.RunRules(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) saveStudyStatusHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.StudyStatusReq
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	req.StudyKey = c.Param("studyKey")
	resp, err := h.clients.StudyService.SaveStudyStatus(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) saveStudyPropsHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.StudyPropsReq
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	req.StudyKey = c.Param("studyKey")
	resp, err := h.clients.StudyService.SaveStudyProps(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) deleteStudyHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.StudyReferenceReq
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	req.StudyKey = c.Param("studyKey")
	resp, err := h.clients.StudyService.DeleteStudy(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) getSurveyResponseStatisticsHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.SurveyResponseQuery
	studyKey := c.Param("studyKey")
	req.StudyKey = studyKey
	from := c.DefaultQuery("from", "")
	if len(from) > 0 {
		n, err := strconv.ParseInt(from, 10, 64)
		if err == nil {
			req.From = n
		}
	}
	until := c.DefaultQuery("until", "")
	if len(until) > 0 {
		n, err := strconv.ParseInt(until, 10, 64)
		if err == nil {
			req.Until = n
		}
	}

	req.Token = token
	resp, err := h.clients.StudyService.GetStudyResponseStatistics(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) getSurveyResponsesHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.SurveyResponseQuery
	studyKey := c.Param("studyKey")
	req.StudyKey = studyKey
	req.SurveyKey = c.DefaultQuery("surveyKey", "")
	from := c.DefaultQuery("from", "")
	if len(from) > 0 {
		n, err := strconv.ParseInt(from, 10, 64)
		if err == nil {
			req.From = n
		}
	}
	until := c.DefaultQuery("until", "")
	if len(until) > 0 {
		n, err := strconv.ParseInt(until, 10, 64)
		if err == nil {
			req.Until = n
		}
	}

	req.Token = token
	stream, err := h.clients.StudyService.StreamStudyResponses(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	resps := &api.SurveyResponses{
		Responses: []*api.SurveyResponse{},
	}
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("getSurveyResponsesHandl(_) = _, %v", err)
			break
		}
		resps.Responses = append(resps.Responses, r)

	}
	h.SendProtoAsJSON(c, http.StatusOK, resps)
}

func (h *HttpEndpoints) getResponseWideFormatCSV(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)
	var req studyAPI.ResponseExportQuery
	studyKey := c.Param("studyKey")
	req.StudyKey = studyKey
	surveyKey := c.Param("surveyKey")
	req.SurveyKey = surveyKey

	from := c.DefaultQuery("from", "")
	if len(from) > 0 {
		n, err := strconv.ParseInt(from, 10, 64)
		if err == nil {
			req.From = n
		}
	}
	until := c.DefaultQuery("until", "")
	if len(until) > 0 {
		n, err := strconv.ParseInt(until, 10, 64)
		if err == nil {
			req.Until = n
		}
	}
	req.IncludeMeta = &studyAPI.ResponseExportQuery_IncludeMeta{
		Position:       c.DefaultQuery("withPositions", "false") == "true",
		ItemVersion:    c.DefaultQuery("withItemVersions", "false") == "true",
		InitTimes:      c.DefaultQuery("withInitTimes", "false") == "true",
		DisplayedTimes: c.DefaultQuery("withDisplayTimes", "false") == "true",
		ResponsedTimes: c.DefaultQuery("withResponseTimes", "false") == "true",
	}
	req.Separator = c.DefaultQuery("sep", "-")
	req.ShortQuestionKeys = c.DefaultQuery("shortKeys", "true") == "true"
	req.Token = token

	stream, err := h.clients.StudyService.GetResponsesWideFormatCSV(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	content := []byte{}
	for {
		chnk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			st := status.Convert(err)
			c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		content = append(content, chnk.Chunk...)
	}

	reader := bytes.NewReader(content)
	contentLength := int64(len(content))
	contentType := "text/csv"

	extraHeaders := map[string]string{
		"Content-Disposition": `attachment; filename=` + fmt.Sprintf("%s_%s.csv", studyKey, surveyKey),
	}

	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
}

func (h *HttpEndpoints) getResponseLongFormatCSV(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)
	var req studyAPI.ResponseExportQuery
	studyKey := c.Param("studyKey")
	req.StudyKey = studyKey
	surveyKey := c.Param("surveyKey")
	req.SurveyKey = surveyKey

	from := c.DefaultQuery("from", "")
	if len(from) > 0 {
		n, err := strconv.ParseInt(from, 10, 64)
		if err == nil {
			req.From = n
		}
	}
	until := c.DefaultQuery("until", "")
	if len(until) > 0 {
		n, err := strconv.ParseInt(until, 10, 64)
		if err == nil {
			req.Until = n
		}
	}
	req.IncludeMeta = &studyAPI.ResponseExportQuery_IncludeMeta{
		Position:       c.DefaultQuery("withPositions", "false") == "true",
		ItemVersion:    c.DefaultQuery("withItemVersions", "false") == "true",
		InitTimes:      c.DefaultQuery("withInitTimes", "false") == "true",
		DisplayedTimes: c.DefaultQuery("withDisplayTimes", "false") == "true",
		ResponsedTimes: c.DefaultQuery("withResponseTimes", "false") == "true",
	}
	req.Separator = c.DefaultQuery("sep", "-")
	req.ShortQuestionKeys = c.DefaultQuery("shortKeys", "true") == "true"
	req.Token = token

	stream, err := h.clients.StudyService.GetResponsesLongFormatCSV(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	content := []byte{}
	for {
		chnk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			st := status.Convert(err)
			c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		content = append(content, chnk.Chunk...)
	}

	reader := bytes.NewReader(content)
	contentLength := int64(len(content))
	contentType := "text/csv"

	extraHeaders := map[string]string{
		"Content-Disposition": `attachment; filename=` + fmt.Sprintf("%s_%s.csv", studyKey, surveyKey),
	}

	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
}

func (h *HttpEndpoints) getResponseFlatJSON(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)
	var req studyAPI.ResponseExportQuery
	studyKey := c.Param("studyKey")
	req.StudyKey = studyKey
	surveyKey := c.Param("surveyKey")
	req.SurveyKey = surveyKey

	from := c.DefaultQuery("from", "")
	if len(from) > 0 {
		n, err := strconv.ParseInt(from, 10, 64)
		if err == nil {
			req.From = n
		}
	}
	until := c.DefaultQuery("until", "")
	if len(until) > 0 {
		n, err := strconv.ParseInt(until, 10, 64)
		if err == nil {
			req.Until = n
		}
	}
	req.IncludeMeta = &studyAPI.ResponseExportQuery_IncludeMeta{
		Position:       c.DefaultQuery("withPositions", "false") == "true",
		ItemVersion:    c.DefaultQuery("withItemVersions", "false") == "true",
		InitTimes:      c.DefaultQuery("withInitTimes", "false") == "true",
		DisplayedTimes: c.DefaultQuery("withDisplayTimes", "false") == "true",
		ResponsedTimes: c.DefaultQuery("withResponseTimes", "false") == "true",
	}
	req.Separator = c.DefaultQuery("sep", "-")
	req.ShortQuestionKeys = c.DefaultQuery("shortKeys", "true") == "true"
	req.Token = token

	stream, err := h.clients.StudyService.GetResponsesFlatJSON(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	content := []byte{}
	for {
		chnk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			st := status.Convert(err)
			c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		content = append(content, chnk.Chunk...)
	}

	reader := bytes.NewReader(content)
	contentLength := int64(len(content))
	contentType := "application/json"

	extraHeaders := map[string]string{
		"Content-Disposition": `attachment; filename=` + fmt.Sprintf("%s_%s.json", studyKey, surveyKey),
	}

	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
}

func (h *HttpEndpoints) getSurveyInfoPreview(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)
	var req studyAPI.SurveyInfoExportQuery
	studyKey := c.Param("studyKey")
	req.StudyKey = studyKey
	surveyKey := c.Param("surveyKey")
	req.SurveyKey = surveyKey
	req.PreviewLanguage = c.DefaultQuery("lang", "en")
	req.ShortQuestionKeys = c.DefaultQuery("shortKeys", "true") == "true"
	req.Token = token

	resp, err := h.clients.StudyService.GetSurveyInfoPreview(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) getSurveyInfoPreviewCSV(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)
	var req studyAPI.SurveyInfoExportQuery
	studyKey := c.Param("studyKey")
	req.StudyKey = studyKey
	surveyKey := c.Param("surveyKey")
	req.SurveyKey = surveyKey
	req.PreviewLanguage = c.DefaultQuery("lang", "en")
	req.ShortQuestionKeys = c.DefaultQuery("shortKeys", "true") == "true"
	req.Token = token

	stream, err := h.clients.StudyService.GetSurveyInfoPreviewCSV(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	content := []byte{}
	for {
		chnk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			st := status.Convert(err)
			c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
		content = append(content, chnk.Chunk...)
	}

	reader := bytes.NewReader(content)
	contentLength := int64(len(content))
	contentType := "text/csv"

	extraHeaders := map[string]string{
		"Content-Disposition": `attachment; filename=` + fmt.Sprintf("survey_%s_%s.csv", studyKey, surveyKey),
	}

	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
}
