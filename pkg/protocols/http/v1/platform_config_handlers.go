package v1

import (
	"context"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/influenzanet/api-gateway/pkg/utils"
	umAPI "github.com/influenzanet/user-management-service/pkg/api"
	"google.golang.org/grpc/status"
)

// platformConfigInstancePattern is the shape of an instance id: it names databases, so letters,
// digits, '-' and '_' only. Anything else is refused before reaching user-management.
var platformConfigInstancePattern = regexp.MustCompile(`^[A-Za-z0-9_-]{1,64}$`)

// platformConfigCacheControl lets browsers and proxies reuse the answer briefly: a change of
// the WhatsApp configuration reaches the sites within a minute.
const platformConfigCacheControl = "public, max-age=60"

// platformConfigTimeout bounds the call to user-management: the route is public, so a slow
// user-management must not hold requests open for free.
const platformConfigTimeout = 5 * time.Second

// getPlatformConfigHandl tells the participant site, before anyone logs in, whether WhatsApp is
// on for an instance (?instance=<id>). The JSON key is the one the web core reads.
func (h *HttpEndpoints) getPlatformConfigHandl(c *gin.Context) {
	instanceID := c.Query("instance")
	if !platformConfigInstancePattern.MatchString(instanceID) {
		c.Header("Cache-Control", "no-store")
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing or invalid instance"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), platformConfigTimeout)
	defer cancel()
	resp, err := h.clients.UserManagement.GetPlatformConfig(ctx, &umAPI.PlatformConfigReq{InstanceId: instanceID})
	if err != nil {
		st := status.Convert(err)
		c.Header("Cache-Control", "no-store")
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	c.Header("Cache-Control", platformConfigCacheControl)
	c.JSON(http.StatusOK, gin.H{"whatsAppEnabled": resp.GetWhatsappEnabled()})
}
