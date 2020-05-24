package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/influenzanet/api-gateway/pkg/utils"
	messageAPI "github.com/influenzanet/messaging-service/pkg/api/messaging_service"
	umAPI "github.com/influenzanet/user-management-service/pkg/api"
	"google.golang.org/grpc/status"
)

func (h *HttpEndpoints) getEmailTemplatesHandl(c *gin.Context) {
	token := utils.ConvertTokenInfosForMessageAPI(c.MustGet("validatedToken").(*umAPI.TokenInfos))

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
	token := utils.ConvertTokenInfosForMessageAPI(c.MustGet("validatedToken").(*umAPI.TokenInfos))

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
	token := utils.ConvertTokenInfosForMessageAPI(c.MustGet("validatedToken").(*umAPI.TokenInfos))

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
