package v1

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/influenzanet/api-gateway/pkg/models"
	"github.com/influenzanet/go-utils/pkg/api_types"
	messageAPI "github.com/influenzanet/messaging-service/pkg/api/messaging_service"
	"google.golang.org/grpc"
)

// fakeMessagingClient captures the requests the handlers forward to the
// messaging-service so a test can check what survived the JSON -> proto step.
type fakeMessagingClient struct {
	messageAPI.MessagingServiceApiClient
	savedTemplate    *messageAPI.SaveEmailTemplateReq
	savedAutoMessage *messageAPI.SaveAutoMessageReq
}

func (f *fakeMessagingClient) SaveEmailTemplate(ctx context.Context, in *messageAPI.SaveEmailTemplateReq, opts ...grpc.CallOption) (*messageAPI.EmailTemplate, error) {
	f.savedTemplate = in
	return in.Template, nil
}

func (f *fakeMessagingClient) GetEmailTemplates(ctx context.Context, in *messageAPI.GetEmailTemplatesReq, opts ...grpc.CallOption) (*messageAPI.EmailTemplates, error) {
	name := "weekly_survey_it"
	return &messageAPI.EmailTemplates{Templates: []*messageAPI.EmailTemplate{{
		Id:                   "tpl-weekly",
		MessageType:          "weekly",
		WhatsappTemplateName: &name,
		WhatsappParams:       map[string]string{"first_name": "firstName"},
	}}}, nil
}

func (f *fakeMessagingClient) SaveAutoMessage(ctx context.Context, in *messageAPI.SaveAutoMessageReq, opts ...grpc.CallOption) (*messageAPI.AutoMessage, error) {
	f.savedAutoMessage = in
	return in.AutoMessage, nil
}

// The WhatsApp template binding as a management client would send it,
// next to the email fields that already existed in the contract.
const emailTemplateWithWhatsAppJSON = `{
	"id": "tpl-weekly",
	"messageType": "weekly",
	"defaultLanguage": "it",
	"translations": [{"lang": "it", "subject": "Questionario settimanale", "templateDef": "SGVsbG8="}],
	"whatsappTemplateName": "weekly_survey_it",
	"whatsappParams": {"first_name": "firstName", "survey_url": "surveyLink"}
}`

func newMessagingTestHandler(fake *fakeMessagingClient) *HttpEndpoints {
	return NewHTTPHandler(&models.APIClients{MessagingService: fake}, models.UseEndpoints{}, nil)
}

func postJSON(t *testing.T, handler gin.HandlerFunc, body string) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("validatedToken", &api_types.TokenInfos{Id: "admin-id", InstanceId: "italy"})
	handler(c)
	return rec
}

func assertWhatsAppBinding(t *testing.T, tpl *messageAPI.EmailTemplate) {
	t.Helper()
	if tpl == nil {
		t.Fatal("template was not forwarded to the messaging-service")
	}
	if got := tpl.GetWhatsappTemplateName(); got != "weekly_survey_it" {
		t.Errorf("whatsappTemplateName forwarded as %q, want %q", got, "weekly_survey_it")
	}
	params := tpl.GetWhatsappParams()
	if len(params) != 2 || params["first_name"] != "firstName" || params["survey_url"] != "surveyLink" {
		t.Errorf("whatsappParams forwarded as %v, want first_name->firstName and survey_url->surveyLink", params)
	}
	if got := tpl.GetMessageType(); got != "weekly" {
		t.Errorf("messageType forwarded as %q, want %q", got, "weekly")
	}
}

// F-01: a management POST carrying the WhatsApp binding must be accepted and
// forwarded, instead of being rejected as an unknown JSON field (HTTP 400).
func TestSaveEmailTemplateAcceptsWhatsAppBinding(t *testing.T) {
	fake := &fakeMessagingClient{}
	h := newMessagingTestHandler(fake)

	body := `{"template": ` + emailTemplateWithWhatsAppJSON + `}`
	rec := postJSON(t, h.saveEmailTemplateHandl, body)

	if rec.Code != http.StatusOK {
		t.Fatalf("POST /messaging/email-templates returned %d, want 200; body: %s", rec.Code, rec.Body.String())
	}
	if fake.savedTemplate == nil {
		t.Fatal("SaveEmailTemplate was not called")
	}
	if fake.savedTemplate.GetToken().GetId() != "admin-id" {
		t.Errorf("validated token not attached to the request")
	}
	assertWhatsAppBinding(t, fake.savedTemplate.GetTemplate())
}

// F-01: same guarantee for auto messages, which embed the template the
// scheduler later uses to decide whether a WhatsApp send is possible.
func TestSaveAutoMessageAcceptsWhatsAppBinding(t *testing.T) {
	fake := &fakeMessagingClient{}
	h := newMessagingTestHandler(fake)

	body := `{"autoMessage": {"id": "am-1", "type": "all-users", "label": "weekly", "period": 604800, "template": ` +
		emailTemplateWithWhatsAppJSON + `}}`
	rec := postJSON(t, h.saveAutoMessageHandl, body)

	if rec.Code != http.StatusOK {
		t.Fatalf("POST /messaging/auto-messages returned %d, want 200; body: %s", rec.Code, rec.Body.String())
	}
	if fake.savedAutoMessage == nil {
		t.Fatal("SaveAutoMessage was not called")
	}
	assertWhatsAppBinding(t, fake.savedAutoMessage.GetAutoMessage().GetTemplate())
}

// A template without the binding keeps working exactly as before the bump.
func TestSaveEmailTemplateWithoutWhatsAppBindingStillAccepted(t *testing.T) {
	fake := &fakeMessagingClient{}
	h := newMessagingTestHandler(fake)

	body := `{"template": {"id": "tpl-weekly", "messageType": "weekly", "defaultLanguage": "it"}}`
	rec := postJSON(t, h.saveEmailTemplateHandl, body)

	if rec.Code != http.StatusOK {
		t.Fatalf("POST /messaging/email-templates returned %d, want 200; body: %s", rec.Code, rec.Body.String())
	}
	tpl := fake.savedTemplate.GetTemplate()
	if tpl.WhatsappTemplateName != nil {
		t.Errorf("whatsappTemplateName should stay unset when omitted, got %q", tpl.GetWhatsappTemplateName())
	}
	if len(tpl.GetWhatsappParams()) != 0 {
		t.Errorf("whatsappParams should stay empty when omitted, got %v", tpl.GetWhatsappParams())
	}
}

// Unknown fields are still rejected: the fix is the contract, not a looser parser.
func TestSaveEmailTemplateStillRejectsUnknownFields(t *testing.T) {
	fake := &fakeMessagingClient{}
	h := newMessagingTestHandler(fake)

	body := `{"template": {"id": "tpl-weekly", "messageType": "weekly", "notAField": "x"}}`
	rec := postJSON(t, h.saveEmailTemplateHandl, body)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("POST with an unknown field returned %d, want 400; body: %s", rec.Code, rec.Body.String())
	}
	if fake.savedTemplate != nil {
		t.Error("SaveEmailTemplate must not be called when the payload is rejected")
	}
}

// F-01, scenario step 2: a stored binding must be visible on GET, otherwise a
// client that reads and writes back the template drops it without noticing.
func TestGetEmailTemplatesExposesWhatsAppBinding(t *testing.T) {
	h := newMessagingTestHandler(&fakeMessagingClient{})
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Set("validatedToken", &api_types.TokenInfos{Id: "admin-id", InstanceId: "italy"})

	h.getEmailTemplatesHandl(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /messaging/email-templates returned %d, want 200; body: %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	for _, want := range []string{`"whatsappTemplateName":"weekly_survey_it"`, `"whatsappParams":{"first_name":"firstName"}`} {
		if !strings.Contains(body, want) {
			t.Errorf("GET response lacks %s; body: %s", want, body)
		}
	}
}
