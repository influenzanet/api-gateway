package v1

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/crewjam/saml/samlsp"
	"github.com/gin-gonic/gin"
	"github.com/influenzanet/api-gateway/pkg/models"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type HttpEndpoints struct {
	clients      *models.APIClients
	useEndpoints models.UseEndpoints
	samlConfig   *models.SAMLConfig
	marshaller   protojson.MarshalOptions
	unmarshaller protojson.UnmarshalOptions
}

func NewHTTPHandler(
	clientRef *models.APIClients,
	useEndpoints models.UseEndpoints,
	samlConfig *models.SAMLConfig,
) *HttpEndpoints {
	m := protojson.MarshalOptions{
		EmitUnpopulated: false,
	}
	um := protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}
	return &HttpEndpoints{
		clients:      clientRef,
		useEndpoints: useEndpoints,
		samlConfig:   samlConfig,
		marshaller:   m,
		unmarshaller: um,
	}
}

func (h *HttpEndpoints) SendProtoAsJSON(c *gin.Context, statusCode int, pbMsg proto.Message) {
	// b, err := .MarshalToString(pbMsg)
	jsonObject, err := h.marshaller.Marshal(pbMsg)

	if err != nil {
		fmt.Println("error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "protobuf message couldn't be transform to json"})
	}
	c.Data(statusCode, "application/json; charset=utf-8", jsonObject)
}

func (h *HttpEndpoints) JsonToProto(c *gin.Context, pbObj interface{}) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	err = h.unmarshaller.Unmarshal(body, (pbObj).(proto.Message))
	return err
}

func (h HttpEndpoints) InitSamlSP() (*samlsp.Middleware, error) {
	if h.samlConfig == nil {
		return nil, errors.New("SAML config not available")
	}

	keyPair, err := tls.LoadX509KeyPair(h.samlConfig.SessionCertPath, h.samlConfig.SessionKeyPath)
	if err != nil {
		return nil, fmt.Errorf("Problem loading session certificate or key: %v", err)
	}

	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		return nil, fmt.Errorf("Problem parsing session certificate or key: %v", err)
	}

	idpMetadataURL, err := url.Parse(h.samlConfig.MetaDataURL)
	if err != nil {
		return nil, fmt.Errorf("Can't parse metadata url: %v", err)
	}

	rootURL, err := url.Parse(h.samlConfig.SPRootUrl)
	if err != nil {
		return nil, fmt.Errorf("Can't parse service provider root URL: %v", err)
	}

	metaData, err := samlsp.FetchMetadata(context.TODO(), http.DefaultClient, *idpMetadataURL)
	if err != nil {
		return nil, fmt.Errorf("Error when fetching metadata: %v", err)
	}

	samlSP, err := samlsp.New(samlsp.Options{
		URL:         *rootURL,
		Key:         keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate: keyPair.Leaf,
		IDPMetadata: metaData,
		EntityID:    h.samlConfig.EntityID,
	})
	if err != nil {
		return nil, err
	}

	// acsURL, err := url.Parse("https://dfki-3094.dfki.de/adminapi/v1/saml/acs")
	acsURL, err := url.Parse(h.samlConfig.SPRootUrl + "/v1/saml/acs")
	if err != nil {
		return nil, fmt.Errorf("Can't parse ACS URL: %v", err)
	}
	samlSP.ServiceProvider.AcsURL = *acsURL

	return samlSP, nil
}
