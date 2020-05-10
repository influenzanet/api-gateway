package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/influenzanet/api-gateway/pkg/utils"
	"google.golang.org/grpc/status"
)

// AddServiceStatusAPI creates all API routes on the supplied RouterGroup
func (h *HttpEndpoints) AddServiceStatusAPI(rg *gin.RouterGroup) {
	userRoutes := rg.Group("/status")
	userRoutes.GET("/user-management", h.statusUserMangementServiceHandl)
	userRoutes.GET("/study-service", h.statusStudyServiceHandl)
}

func (h *HttpEndpoints) statusUserMangementServiceHandl(c *gin.Context) {
	resp, err := h.clients.UserManagement.Status(context.Background(), &empty.Empty{})
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) statusStudyServiceHandl(c *gin.Context) {
	resp, err := h.clients.StudyService.Status(context.Background(), &empty.Empty{})
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}
