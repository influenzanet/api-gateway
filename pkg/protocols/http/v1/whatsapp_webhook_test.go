package v1

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/influenzanet/api-gateway/pkg/models"
	messageAPI "github.com/influenzanet/messaging-service/pkg/api/messaging_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Meta calls /v1/webhooks/whatsapp without any token of ours: the verify token guards the
// subscription check, the signature guards every notification, and only a signed body reaches
// messaging, byte for byte.

const (
	testAppSecret   = "test-app-secret"
	testVerifyToken = "test-verify-token"
)

type webhookMessagingFake struct {
	messageAPI.MessagingServiceApiClient
	payloads [][]byte
	deadline bool
	err      error
}

func (f *webhookMessagingFake) ProcessWhatsAppWebhook(ctx context.Context, in *messageAPI.ProcessWhatsAppWebhookReq, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	f.payloads = append(f.payloads, in.Payload)
	_, f.deadline = ctx.Deadline()
	if f.err != nil {
		return nil, f.err
	}
	return &emptypb.Empty{}, nil
}

func webhookRouter(fake *webhookMessagingFake) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := NewHTTPHandler(&models.APIClients{MessagingService: fake}, models.UseEndpoints{}, nil)
	h.AddWhatsAppWebhookAPI(router.Group("/v1"), models.WhatsAppWebhookConfig{AppSecret: testAppSecret, VerifyToken: testVerifyToken})
	return router
}

func sign(body []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func postNotification(router *gin.Engine, body []byte, signature string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/v1/webhooks/whatsapp", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if signature != "" {
		req.Header.Set("X-Hub-Signature-256", signature)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

// The body exactly as signed: spacing and key order matter to the signature, and messaging
// gets these very bytes.
var notification = []byte(`{"object":"whatsapp_business_account",  "entry":[{"id":"1","changes":[]}]}`)

func TestWhatsAppWebhookForwardsASignedNotification(t *testing.T) {
	fake := &webhookMessagingFake{}
	rec := postNotification(webhookRouter(fake), notification, sign(notification, testAppSecret))
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d, want 200", rec.Code)
	}
	if len(fake.payloads) != 1 || !bytes.Equal(fake.payloads[0], notification) {
		t.Fatalf("messaging got %q, want the body as received", fake.payloads)
	}
	if !fake.deadline {
		t.Error("the call to messaging has no deadline")
	}
}

func TestWhatsAppWebhookRefusesWhatIsNotSignedByMeta(t *testing.T) {
	tampered := append([]byte{}, notification...)
	tampered[len(tampered)-2] = ' '
	valid := sign(notification, testAppSecret)
	for _, c := range []struct {
		name      string
		body      []byte
		signature string
	}{
		{"no signature", notification, ""},
		{"signed with another secret", notification, sign(notification, "another-secret")},
		{"the body changed after signing", tampered, valid},
		{"without the sha256= prefix", notification, strings.TrimPrefix(valid, "sha256=")},
		{"another algorithm", notification, "sha1=" + strings.TrimPrefix(valid, "sha256=")},
		{"not hex", notification, "sha256=" + strings.Repeat("z", 64)},
		{"empty", notification, "sha256="},
		{"truncated", notification, valid[:len(valid)-2]},
	} {
		t.Run(c.name, func(t *testing.T) {
			fake := &webhookMessagingFake{}
			rec := postNotification(webhookRouter(fake), c.body, c.signature)
			if rec.Code != http.StatusUnauthorized {
				t.Errorf("status %d, want 401", rec.Code)
			}
			if len(fake.payloads) != 0 {
				t.Error("an unsigned body reached messaging")
			}
		})
	}
}

func TestWhatsAppWebhookRefusesAnOversizedBody(t *testing.T) {
	fake := &webhookMessagingFake{}
	body := bytes.Repeat([]byte("a"), whatsAppWebhookMaxBody+1)
	if rec := postNotification(webhookRouter(fake), body, sign(body, testAppSecret)); rec.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("status %d, want 413", rec.Code)
	}
	if len(fake.payloads) != 0 {
		t.Error("an oversized body reached messaging")
	}

	// Meta's largest notifications (3 MB) still go through.
	body = bytes.Repeat([]byte("a"), 3<<20)
	if rec := postNotification(webhookRouter(fake), body, sign(body, testAppSecret)); rec.Code != http.StatusOK {
		t.Errorf("3 MB body: status %d, want 200", rec.Code)
	}
}

// Anything but 200 makes Meta send the notification again; that is wanted exactly when a later
// attempt can succeed.
func TestWhatsAppWebhookAnswerTellsMetaWhetherToRetry(t *testing.T) {
	for _, c := range []struct {
		name string
		err  error
		want int
	}{
		{"not a webhook payload: never retried", status.Error(codes.InvalidArgument, "not a webhook payload"), http.StatusOK},
		{"not processed completely", status.Error(codes.Unavailable, "x"), http.StatusServiceUnavailable},
		{"not configured in messaging", status.Error(codes.FailedPrecondition, "x"), http.StatusServiceUnavailable},
		{"messaging too slow", status.Error(codes.DeadlineExceeded, "x"), http.StatusServiceUnavailable},
		{"messaging without the call", status.Error(codes.Unimplemented, "x"), http.StatusServiceUnavailable},
		{"internal error", status.Error(codes.Internal, "x"), http.StatusServiceUnavailable},
	} {
		t.Run(c.name, func(t *testing.T) {
			fake := &webhookMessagingFake{err: c.err}
			if rec := postNotification(webhookRouter(fake), notification, sign(notification, testAppSecret)); rec.Code != c.want {
				t.Errorf("status %d, want %d", rec.Code, c.want)
			}
		})
	}
}

func TestWhatsAppWebhookSubscriptionCheck(t *testing.T) {
	get := func(query string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		webhookRouter(&webhookMessagingFake{}).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/webhooks/whatsapp"+query, nil))
		return rec
	}

	rec := get("?hub.mode=subscribe&hub.verify_token=" + testVerifyToken + "&hub.challenge=1158201444")
	if rec.Code != http.StatusOK || rec.Body.String() != "1158201444" {
		t.Fatalf("status %d body %q, want 200 and the challenge", rec.Code, rec.Body)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/plain") || rec.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Errorf("Content-Type %q, nosniff %q: the echoed challenge must not be read as HTML", ct, rec.Header().Get("X-Content-Type-Options"))
	}

	for _, query := range []string{
		"",
		"?hub.mode=subscribe&hub.challenge=1",
		"?hub.mode=subscribe&hub.verify_token=wrong&hub.challenge=1",
		"?hub.mode=subscribe&hub.verify_token=" + testVerifyToken[:len(testVerifyToken)-1] + "&hub.challenge=1",
		"?hub.mode=subscribe&hub.verify_token=" + testVerifyToken + "x&hub.challenge=1",
		"?hub.mode=unsubscribe&hub.verify_token=" + testVerifyToken + "&hub.challenge=1",
		"?hub.mode=subscribe&hub.verify_token=" + testVerifyToken,
	} {
		if rec := get(query); rec.Code != http.StatusForbidden || rec.Body.Len() != 0 {
			t.Errorf("%q: status %d body %q, want 403 and nothing echoed", query, rec.Code, rec.Body)
		}
	}
}

// Without an app secret nothing is ever valid, not even a body signed with the empty key.
func TestValidWhatsAppSignatureRefusesAnUnsetSecret(t *testing.T) {
	if validWhatsAppSignature(sign(notification, ""), notification, "") {
		t.Error("a signature made with the empty key passed without an app secret")
	}
}

func TestSameSecretRefusesAnUnsetSecret(t *testing.T) {
	if sameSecret("", "") {
		t.Error("an empty token matched an unset verify token")
	}
	if !sameSecret("abc", "abc") || sameSecret("abc", "abd") {
		t.Error("sameSecret does not compare")
	}
}
