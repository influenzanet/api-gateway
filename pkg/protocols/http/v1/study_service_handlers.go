package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/coneno/logger"
	"github.com/gin-gonic/gin"
	"github.com/h2non/filetype"
	"github.com/influenzanet/api-gateway/pkg/utils"
	"github.com/influenzanet/go-utils/pkg/api_types"
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

func (h *HttpEndpoints) unpublishSurveyHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.SurveyReferenceRequest
	req.Token = token
	req.StudyKey = c.Param("studyKey")
	req.SurveyKey = c.Param("surveyKey")

	resp, err := h.clients.StudyService.UnpublishSurvey(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) removeSurveyVersionHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.SurveyVersionReferenceRequest
	req.Token = token
	req.StudyKey = c.Param("studyKey")
	req.SurveyKey = c.Param("surveyKey")
	req.VersionId = c.Param("versionID")

	resp, err := h.clients.StudyService.RemoveSurveyVersion(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) deleteParticipantFilesReq(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)
	var req studyAPI.DeleteParticipantFilesReq
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	req.StudyKey = c.Param("studyKey")

	resp, err := h.clients.StudyService.DeleteParticipantFiles(context.Background(), &req)
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

func (h *HttpEndpoints) removeConfidentialDataHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	p := c.DefaultQuery("profiles", "")
	profiles := []string{}
	if len(p) > 0 {
		profiles = strings.Split(p, ",")
	}
	req := &studyAPI.RemoveConfidentialResponsesForProfilesReq{
		Token:       token,
		ForProfiles: profiles,
	}

	resp, err := h.clients.StudyService.RemoveConfidentialResponsesForProfiles(context.Background(), req)
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

func (h *HttpEndpoints) getStudySurveySurveyKeysHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.GetSurveyKeysRequest
	req.Token = token
	req.StudyKey = c.Param("studyKey")
	resp, err := h.clients.StudyService.GetSurveyKeys(context.Background(), &req)
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

func (h *HttpEndpoints) getReportsForParticipant(c *gin.Context) {
	// ?studies=todo1,todo2&profileIds=todo1,todo2&from=time1&until=time2&reportKey=todo3
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.GetReportsForUserReq
	req.Token = token

	studies := c.DefaultQuery("studies", "")
	if len(studies) > 0 {
		req.OnlyForStudies = strings.Split(studies, ",")
	}
	profileIds := c.DefaultQuery("profileIds", "")
	if len(profileIds) > 0 {
		req.OnlyForProfiles = strings.Split(profileIds, ",")
	}
	req.ReportKeyFilter = c.DefaultQuery("reportKey", "")
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

	limit := c.DefaultQuery("limit", "")
	if len(limit) > 0 {
		n, err := strconv.ParseInt(limit, 10, 64)
		if err == nil {
			req.Limit = n
		}
	}

	ignoreReports := c.DefaultQuery("ignoreReports", "")
	if len(ignoreReports) > 0 {
		req.IgnoreReports = strings.Split(ignoreReports, ",")
	}

	resp, err := h.clients.StudyService.GetReportsForUser(context.Background(), &req)
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

func (h *HttpEndpoints) getSurveyVersionInfosHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.SurveyReferenceRequest
	req.Token = token
	req.StudyKey = c.Param("studyKey")
	req.SurveyKey = c.Param("surveyKey")
	resp, err := h.clients.StudyService.GetSurveyVersionInfos(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) getSurveyDefForStudyHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.SurveyVersionReferenceRequest
	req.Token = token
	req.StudyKey = c.Param("studyKey")
	req.SurveyKey = c.Param("surveyKey")
	req.VersionId = c.Param("versionID")
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

func (h *HttpEndpoints) runCustomStudyRulesForSingleParticipantHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.RunRulesForSingleParticipantReq
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	req.StudyKey = c.Param("studyKey")
	resp, err := h.clients.StudyService.RunRulesForSingleParticipant(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) runCustomStudyRulesForPreviousResponsesHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.RunRulesForPreviousResponsesReq
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var filter studyAPI.RunRulesForPreviousResponsesReq_ResponseFilter
	req.Token = token
	req.StudyKey = c.Param("studyKey")
	surveyKeys := c.DefaultQuery("surveyKeys", "")
	if len(surveyKeys) > 0 {
		filter.SurveyKeys = strings.Split(surveyKeys, ",")
	}
	participantIds := c.DefaultQuery("participantIds", "")
	if len(participantIds) > 0 {
		filter.ParticipantIds = strings.Split(participantIds, ",")
	}
	for i, ele := range filter.ParticipantIds {
		filter.ParticipantIds[i] = strings.TrimSpace(ele)
	}
	participantStatus := c.DefaultQuery("participantStatus", "active")
	if len(participantStatus) > 0 {
		filter.ParticipantStatus = strings.Split(participantStatus, ",")
	}
	for i := range filter.ParticipantStatus {
		filter.ParticipantStatus[i] = strings.TrimSpace(filter.ParticipantStatus[i])
	}
	from := c.DefaultQuery("from", "")
	if len(from) > 0 {
		n, err := strconv.ParseInt(from, 10, 64)
		if err == nil {
			filter.From = n
		}
	}
	until := c.DefaultQuery("until", "")
	if len(until) > 0 {
		n, err := strconv.ParseInt(until, 10, 64)
		if err == nil {
			filter.Until = n
		}
	}

	req.Filter = &filter

	resp, err := h.clients.StudyService.RunRulesForPreviousResponses(context.Background(), &req)
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

func (h *HttpEndpoints) getStudyNotificationSubscriptionsHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.GetResearcherNotificationSubscriptionsReq
	req.Token = token
	req.StudyKey = c.Param("studyKey")
	resp, err := h.clients.StudyService.GetResearcherNotificationSubscriptions(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) saveStudyNotificationSubscriptionsHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.UpdateResearcherNotificationSubscriptionsReq
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	req.StudyKey = c.Param("studyKey")
	resp, err := h.clients.StudyService.UpdateResearcherNotificationSubscriptions(context.Background(), &req)
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

	resps := &studyAPI.SurveyResponses{
		Responses: []*studyAPI.SurveyResponse{},
	}
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Error.Printf("unexpected error during stream: %v", err)
			break
		}
		resps.Responses = append(resps.Responses, r)

	}
	h.SendProtoAsJSON(c, http.StatusOK, resps)
}

func (h *HttpEndpoints) getConfidentialResponsesHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.ConfidentialResponsesQuery
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	req.StudyKey = c.Param("studyKey")

	resp, err := h.clients.StudyService.GetConfidentialResponses(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) getParticipantStatesForStudy(c *gin.Context) {
	// ?status=active(opt)
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.ParticipantStateQuery
	studyKey := c.Param("studyKey")
	req.StudyKey = studyKey
	req.Status = c.DefaultQuery("status", "")

	req.Token = token
	stream, err := h.clients.StudyService.StreamParticipantStates(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	resp := &studyAPI.ParticipantStates{
		ParticipantStates: []*studyAPI.ParticipantState{},
	}
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Error.Printf("unexpected error during stream: %v", err)
			break
		}
		resp.ParticipantStates = append(resp.ParticipantStates, r)

	}

	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) getReportsForStudy(c *gin.Context) {
	// ?reportKey=todo&from=time1&until=time2&participant=todo
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.ReportHistoryQuery
	studyKey := c.Param("studyKey")
	req.StudyKey = studyKey
	req.ReportKey = c.DefaultQuery("reportKey", "")
	req.ParticipantId = c.DefaultQuery("participant", "")

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
	stream, err := h.clients.StudyService.StreamReportHistory(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	resps := &studyAPI.ReportHistory{
		Reports: []*studyAPI.Report{},
	}
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Error.Printf("unexpected error during stream: %v", err)
			break
		}
		resps.Reports = append(resps.Reports, r)

	}
	h.SendProtoAsJSON(c, http.StatusOK, resps)
}

func (h *HttpEndpoints) getFileInfosForStudy(c *gin.Context) {
	// ?fileType=todo&from=time1&until=time2&participant=todo
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.FileInfoQuery
	studyKey := c.Param("studyKey")
	req.StudyKey = studyKey
	req.FileType = c.DefaultQuery("fileType", "")
	req.ParticipantId = c.DefaultQuery("participant", "")

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
	stream, err := h.clients.StudyService.StreamParticipantFileInfos(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	resps := &studyAPI.FileInfos{
		FileInfos: []*studyAPI.FileInfo{},
	}
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Error.Printf("unexpected error during stream: %v", err)
			break
		}
		resps.FileInfos = append(resps.FileInfos, r)

	}
	h.SendProtoAsJSON(c, http.StatusOK, resps)
}

func (h *HttpEndpoints) getParticipantFile(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)
	var req studyAPI.GetParticipantFileReq
	studyKey := c.Param("studyKey")
	req.StudyKey = studyKey
	req.FileId = c.DefaultQuery("id", "")
	req.Token = token

	stream, err := h.clients.StudyService.GetParticipantFile(context.Background(), &req)
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

	kind, _ := filetype.Match(content)
	reader := bytes.NewReader(content)
	contentLength := int64(len(content))
	contentType := kind.MIME.Value

	extraHeaders := map[string]string{
		"Content-Disposition": `attachment; filename=` + fmt.Sprintf("%s.%s", req.FileId, kind.Extension),
	}

	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
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
	pageSize := c.DefaultQuery("pageSize", "")
	if len(pageSize) > 0 {
		n, err := strconv.ParseInt(pageSize, 10, 32)
		if err == nil {
			req.PageSize = int32(n)
		}
	}
	page := c.DefaultQuery("page", "1")
	if len(page) > 0 {
		n, err := strconv.ParseInt(page, 10, 32)
		if err == nil {
			req.Page = int32(n)
		}
	}
	req.IncludeMeta = &studyAPI.ResponseExportQuery_IncludeMeta{
		Position:       c.DefaultQuery("withPositions", "false") == "true",
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
	pageSize := c.DefaultQuery("pageSize", "")
	if len(pageSize) > 0 {
		n, err := strconv.ParseInt(pageSize, 10, 32)
		if err == nil {
			req.PageSize = int32(n)
		}
	}
	page := c.DefaultQuery("page", "1")
	if len(page) > 0 {
		n, err := strconv.ParseInt(page, 10, 32)
		if err == nil {
			req.Page = int32(n)
		}
	}
	req.IncludeMeta = &studyAPI.ResponseExportQuery_IncludeMeta{
		Position:       c.DefaultQuery("withPositions", "false") == "true",
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
	pageSize := c.DefaultQuery("pageSize", "50")
	if len(pageSize) > 0 {
		n, err := strconv.ParseInt(pageSize, 10, 32)
		if err == nil {
			req.PageSize = int32(n)
		} else {
			req.PageSize = 50
		}
	}
	page := c.DefaultQuery("page", "1")
	if len(page) > 0 {
		n, err := strconv.ParseInt(page, 10, 32)
		if err == nil {
			req.Page = int32(n)
		} else {
			req.Page = 1
		}
	}
	req.IncludeMeta = &studyAPI.ResponseExportQuery_IncludeMeta{
		Position:       c.DefaultQuery("withPositions", "false") == "true",
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

func (h *HttpEndpoints) getResponsesFlatJSONWithPagination(c *gin.Context) {
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
	pageSize := c.DefaultQuery("pageSize", "50")
	if len(pageSize) > 0 {
		n, err := strconv.ParseInt(pageSize, 10, 32)
		if err == nil {
			req.PageSize = int32(n)
		} else {
			req.PageSize = 50
		}
	}
	page := c.DefaultQuery("page", "1")
	if len(page) > 0 {
		n, err := strconv.ParseInt(page, 10, 32)
		if err == nil {
			req.Page = int32(n)
		} else {
			req.Page = 1
		}
	}
	req.IncludeMeta = &studyAPI.ResponseExportQuery_IncludeMeta{
		Position:       c.DefaultQuery("withPositions", "false") == "true",
		InitTimes:      c.DefaultQuery("withInitTimes", "false") == "true",
		DisplayedTimes: c.DefaultQuery("withDisplayTimes", "false") == "true",
		ResponsedTimes: c.DefaultQuery("withResponseTimes", "false") == "true",
	}
	req.Separator = c.DefaultQuery("sep", "-")
	req.ShortQuestionKeys = c.DefaultQuery("shortKeys", "true") == "true"
	req.Token = token

	stream, err := h.clients.StudyService.GetResponsesFlatJSONWithPagination(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	chnk, err := stream.Recv()
	info := chnk.GetInfo()
	if info == nil {
		logger.Error.Println("missing argument")
		//return
	}

	content := []byte{}
	b, err := json.Marshal(info)
	content = b

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
		content = append(content, chnk.GetChunk()...)
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

const chunkSize = 64 * 1024 // 64 KiB

func (h *HttpEndpoints) uploadParticipantFileReq(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)
	studyKey := c.Param("studyKey")

	file, err := c.FormFile("file")
	// file, header, err := c.Request.FormFile("file")
	if err != nil {
		logger.Error.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Debug.Printf("Uploading file '%s' with size '%d'", file.Filename, file.Size)

	// Open file
	f, err := file.Open()
	if err != nil {
		logger.Error.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer f.Close()

	// Get bytes
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	content := buf.Bytes()

	stream, err := h.clients.StudyService.UploadParticipantFile(context.Background())
	if err != nil {
		logger.Error.Printf("cannot upload image: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	kind, _ := filetype.Match(content)

	// Send infos
	req := &studyAPI.UploadParticipantFileReq{
		Data: &studyAPI.UploadParticipantFileReq_Info_{
			Info: &studyAPI.UploadParticipantFileReq_Info{
				Token:                token,
				StudyKey:             studyKey,
				VisibleToParticipant: true,
				FileType: &studyAPI.FileType{
					Type:    kind.MIME.Type,
					Subtype: kind.MIME.Subtype,
					Value:   kind.MIME.Value,
				},
				Participant: &studyAPI.UploadParticipantFileReq_Info_ProfileId{
					ProfileId: token.ProfilId,
				},
			},
		},
	}

	err = stream.Send(req)
	if err != nil {
		st := status.Convert(err)
		logger.Error.Println(st.Message())
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return

	}

	for currentByte := 0; currentByte < len(content); currentByte += chunkSize {
		var currentChnk []byte
		if currentByte+chunkSize > len(content) {
			currentChnk = content[currentByte:]
		} else {
			currentChnk = content[currentByte : currentByte+chunkSize]
		}
		req = &studyAPI.UploadParticipantFileReq{
			Data: &studyAPI.UploadParticipantFileReq_Chunk{
				Chunk: currentChnk,
			},
		}

		if err := stream.Send(req); err != nil {
			if err == io.EOF {
				break
			}
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		st := status.Convert(err)
		logger.Error.Println(st.Message())
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	h.SendProtoAsJSON(c, http.StatusOK, reply)
}

func (h *HttpEndpoints) getParticipantStateByIDHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.ParticipantStateByIDQuery
	studyKey := c.Param("studyKey")
	req.StudyKey = studyKey
	req.ParticipantId = c.DefaultQuery("participantID", "")

	req.Token = token
	state, err := h.clients.StudyService.GetParticipantStateByID(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	h.SendProtoAsJSON(c, http.StatusOK, state)
}

func (h *HttpEndpoints) getParticipantStatesWithPaginationHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.GetPStatesWithPaginationQuery
	req.StudyKey = c.Param("studyKey")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		logger.Error.Println("Could not read page parameter")
		page = 1
	}
	req.Page = int32(page)

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "50"))
	if err != nil {
		logger.Error.Println("Could not read page size parameter")
		pageSize = 50
	}
	req.PageSize = int32(pageSize)

	queryParam := c.DefaultQuery("query", "")
	if len(queryParam) > 0 {
		decodedQP, err := url.QueryUnescape(queryParam)
		if err != nil {
			logger.Error.Printf("Failed to decode query parameter: %v", err)
			decodedQP = ""
		}
		req.Query = decodedQP
	}

	sortParam := c.DefaultQuery("sortBy", "")
	if len(sortParam) > 0 {
		decodedSP, err := url.QueryUnescape(sortParam)
		if err != nil {
			logger.Error.Printf("Failed to decode sort parameter: %v", err)
			decodedSP = ""
		}
		err = json.Unmarshal([]byte(decodedSP), &req.SortBy)
		if err != nil {
			logger.Error.Printf("Failed to parse sort parameter: %v", err)
		}
	}

	req.Token = token
	state, err := h.clients.StudyService.GetParticipantStatesWithPagination(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	h.SendProtoAsJSON(c, http.StatusOK, state)
}

func (h *HttpEndpoints) getCurrentStudyRulesHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.StudyReferenceReq
	studyKey := c.Param("studyKey")
	req.StudyKey = studyKey
	req.Token = token

	studyRules, err := h.clients.StudyService.GetCurrentStudyRules(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	h.SendProtoAsJSON(c, http.StatusOK, studyRules)
}

func (h *HttpEndpoints) getStudyRulesHistoryHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.StudyRulesHistoryReq
	req.StudyKey = c.Param("studyKey")
	req.Token = token
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		logger.Error.Println("Could not read page parameter")
		page = 1
	}
	req.Page = int32(page)
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "50"))
	if err != nil {
		logger.Error.Println("Could not read page size parameter")
		req.PageSize = 50
	}
	req.PageSize = int32(pageSize)
	descending, err := strconv.ParseBool(c.DefaultQuery("descending", "FALSE"))
	if err != nil {
		logger.Warning.Println("Could not read order parameter")
		descending = false
	}
	req.Descending = descending

	since, err := strconv.Atoi(c.DefaultQuery("since", "0"))
	if err != nil {
		logger.Error.Println("Could not read start upload date parameter")
		req.Since = 0
	}
	if since < 0 {
		logger.Warning.Println("Invalid start date parameter")
		req.Since = 0
	} else {
		req.Since = int64(since)
	}
	until, err := strconv.Atoi(c.DefaultQuery("until", "0"))
	if err != nil {
		logger.Error.Println("Could not read end upload date parameter")
		req.Until = 0
	}
	if until < 0 {
		logger.Warning.Println("Invalid end date parameter")
		req.Until = 0
	} else {
		req.Until = int64(until)
	}

	studyRules, err := h.clients.StudyService.GetStudyRulesHistory(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	h.SendProtoAsJSON(c, http.StatusOK, studyRules)
}

func (h *HttpEndpoints) removeStudyRulesVersionHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req studyAPI.StudyRulesVersionReferenceReq
	req.Token = token
	req.Id = c.Param("ID")

	resp, err := h.clients.StudyService.RemoveStudyRulesVersion(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	h.SendProtoAsJSON(c, http.StatusOK, resp)
}
