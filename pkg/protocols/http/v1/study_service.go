package v1

/*
func (h *HttpEndpoints) studySystemCreateStudyHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api.TokenInfos)

	var req api.NewStudyRequest
	if err := gjpb.JsonToPb(c, &req); err != nil {
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

	gjpb.SendPBAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) saveSurveyToStudyHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api.TokenInfos)

	var req api.AddSurveyReq
	if err := gjpb.JsonToPb(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token

	resp, err := h.clients.StudyService.SaveSurveyToStudy(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	gjpb.SendPBAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) removeSurveyFromStudyHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api.TokenInfos)

	var req api.SurveyReferenceRequest
	if err := gjpb.JsonToPb(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token

	resp, err := h.clients.StudyService.RemoveSurveyFromStudy(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	gjpb.SendPBAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) getAssignedSurveyHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api.TokenInfos)

	var req api.SurveyReferenceRequest
	if err := gjpb.JsonToPb(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	resp, err := h.clients.StudyService.GetAssignedSurvey(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	gjpb.SendPBAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) submitSurveyResponseHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api.TokenInfos)

	var req api.SubmitResponseReq
	if err := gjpb.JsonToPb(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	resp, err := h.clients.StudyService.SubmitResponse(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	gjpb.SendPBAsJSON(c, http.StatusOK, resp)
}
*/
