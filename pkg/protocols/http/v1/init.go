package v1

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/influenzanet/api-gateway/pkg/models"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type HttpEndpoints struct {
	clients      *models.APIClients
	useEndpoints models.UseEndpoints
	reCaptcha    *models.RecaptchaProps
	marshaller   protojson.MarshalOptions
	unmarshaller protojson.UnmarshalOptions
}

func NewHTTPHandler(
	clientRef *models.APIClients,
	useEndpoints models.UseEndpoints,
	reCaptcha *models.RecaptchaProps,
) *HttpEndpoints {
	m := protojson.MarshalOptions{
		EmitUnpopulated: false,
	}
	um := protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}
	if reCaptcha == nil {
		reCaptcha = &models.RecaptchaProps{}
	}
	return &HttpEndpoints{
		clients:      clientRef,
		useEndpoints: useEndpoints,
		reCaptcha:    reCaptcha,
		marshaller:   m,
		unmarshaller: um,
	}
}

func (h *HttpEndpoints) SendProtoAsJSON(c *gin.Context, statusCode int, pbMsg proto.Message) {
	// b, err := .MarshalToString(pbMsg)
	b, err := h.marshaller.Marshal(pbMsg)

	if err != nil {
		fmt.Println("error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "protobuf message couldn't be transform to json"})
	}
	c.Data(statusCode, "application/json; charset=utf-8", b)
}

func (h *HttpEndpoints) JsonToProto(c *gin.Context, pbObj interface{}) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	err = h.unmarshaller.Unmarshal(body, (pbObj).(proto.Message))
	return err
}
