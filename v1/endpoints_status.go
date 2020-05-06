package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/influenzanet/api-gateway/utils"
	gjpb "github.com/phev8/gin-protobuf-json-converter"
	"google.golang.org/grpc/status"
)

// InitUserEndpoints creates all API routes on the supplied RouterGroup
func initStatusEndpoints(rg *gin.RouterGroup) {
	userRoutes := rg.Group("/status")
	userRoutes.GET("/user-management", statusUserMangementServiceHandl)
	userRoutes.GET("/study-service", statusStudyServiceHandl)
	userRoutes.GET("/authentication-service", statusAuthenticationServiceHandl)
}

func statusUserMangementServiceHandl(c *gin.Context) {
	resp, err := clients.UserManagement.Status(context.Background(), &empty.Empty{})
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	gjpb.SendPBAsJSON(c, http.StatusOK, resp)
}

func statusStudyServiceHandl(c *gin.Context) {
	resp, err := clients.StudyService.Status(context.Background(), &empty.Empty{})
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	gjpb.SendPBAsJSON(c, http.StatusOK, resp)
}

func statusAuthenticationServiceHandl(c *gin.Context) {
	resp, err := clients.AuthService.Status(context.Background(), &empty.Empty{})
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	gjpb.SendPBAsJSON(c, http.StatusOK, resp)
}
