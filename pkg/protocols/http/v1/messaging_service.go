package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/influenzanet/api-gateway/pkg/utils"
	"github.com/influenzanet/go-utils/pkg/api_types"
	messageAPI "github.com/influenzanet/messaging-service/pkg/api/messaging_service"
	"google.golang.org/grpc/status"
)

func (h *HttpEndpoints) getEmailTemplatesHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req messageAPI.GetEmailTemplatesReq
	req.Token = token
	resp, err := h.clients.MessagingService.GetEmailTemplates(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) saveEmailTemplateHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req messageAPI.SaveEmailTemplateReq
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	resp, err := h.clients.MessagingService.SaveEmailTemplate(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) deleteEmailTemplateHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req messageAPI.DeleteEmailTemplateReq
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	resp, err := h.clients.MessagingService.DeleteEmailTemplate(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) getAutoMessagesHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req messageAPI.GetAutoMessagesReq
	req.Token = token
	resp, err := h.clients.MessagingService.GetAutoMessages(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) saveAutoMessageHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req messageAPI.SaveAutoMessageReq
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	resp, err := h.clients.MessagingService.SaveAutoMessage(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) deleteAutoMessageHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req messageAPI.DeleteAutoMessageReq
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	resp, err := h.clients.MessagingService.DeleteAutoMessage(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) sendMessageToAllUsersHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req messageAPI.SendMessageToAllUsersReq
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	resp, err := h.clients.MessagingService.SendMessageToAllUsers(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) sendMessageToStudyParticipantsHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req messageAPI.SendMessageToStudyParticipantsReq
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	resp, err := h.clients.MessagingService.SendMessageToStudyParticipants(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}
