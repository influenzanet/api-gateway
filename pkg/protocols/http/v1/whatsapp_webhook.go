package v1

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/coneno/logger"
	"github.com/gin-gonic/gin"
	"github.com/influenzanet/api-gateway/pkg/models"
	messageAPI "github.com/influenzanet/messaging-service/pkg/api/messaging_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// WhatsAppWebhookPath is the route Meta calls, under /v1.
const WhatsAppWebhookPath = "/webhooks/whatsapp"

// whatsAppWebhookMaxBody bounds the body read from Meta. Meta's notifications reach 3 MB; the
// bound stays below the 4 MiB a gRPC message may carry by default.
const whatsAppWebhookMaxBody = 3<<20 + 512<<10

// whatsAppWebhookTimeout bounds the forwarding to messaging. Meta waits for the answer and sends
// the notification again on anything but 200, so a slow messaging means a later retry, not a loss.
const whatsAppWebhookTimeout = 10 * time.Second

const whatsAppSignatureHeader = "X-Hub-Signature-256"

// AddWhatsAppWebhookAPI adds the route Meta calls: GET for the subscription check when the
// webhook is set up in the Meta app, POST for the notifications. The route is public; what makes
// it safe is the verify token on GET and the signature on POST.
func (h *HttpEndpoints) AddWhatsAppWebhookAPI(rg *gin.RouterGroup, cfg models.WhatsAppWebhookConfig) {
	rg.GET(WhatsAppWebhookPath, whatsAppWebhookSubscription(cfg.VerifyToken))
	rg.POST(WhatsAppWebhookPath, h.whatsAppWebhookNotification(cfg.AppSecret))
}

// whatsAppWebhookSubscription answers Meta's subscription check: the challenge comes back, as
// plain text, only when the mode is subscribe and the verify token is the configured one.
func whatsAppWebhookSubscription(verifyToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		challenge := c.Query("hub.challenge")
		if c.Query("hub.mode") != "subscribe" || challenge == "" || !sameSecret(c.Query("hub.verify_token"), verifyToken) {
			logger.Warning.Printf("WhatsApp webhook: subscription check refused")
			c.Status(http.StatusForbidden)
			return
		}
		c.Header("X-Content-Type-Options", "nosniff")
		c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(challenge))
	}
}

// whatsAppWebhookNotification checks the signature of a notification on the body as received,
// before anything reads it, and forwards the body to messaging. The answer tells Meta whether to
// send it again: 200 when it was processed, or when it never can be; 503 when messaging could not
// process it completely, or could not be reached.
func (h *HttpEndpoints) whatsAppWebhookNotification(appSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(http.MaxBytesReader(c.Writer, c.Request.Body, whatsAppWebhookMaxBody))
		if err != nil {
			var tooLarge *http.MaxBytesError
			if errors.As(err, &tooLarge) {
				logger.Warning.Printf("WhatsApp webhook: body over %d bytes refused", whatsAppWebhookMaxBody)
				c.Status(http.StatusRequestEntityTooLarge)
				return
			}
			c.Status(http.StatusBadRequest)
			return
		}
		if !validWhatsAppSignature(c.GetHeader(whatsAppSignatureHeader), body, appSecret) {
			logger.Warning.Printf("WhatsApp webhook: notification with a missing or wrong signature refused")
			c.Status(http.StatusUnauthorized)
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), whatsAppWebhookTimeout)
		defer cancel()
		_, err = h.clients.MessagingService.ProcessWhatsAppWebhook(ctx, &messageAPI.ProcessWhatsAppWebhookReq{Payload: body})
		switch code := status.Code(err); code {
		case codes.OK:
			c.Status(http.StatusOK)
		case codes.InvalidArgument:
			// Not a webhook payload: sending it again cannot change that.
			logger.Warning.Printf("WhatsApp webhook: signed notification that is not a webhook payload, dropped")
			c.Status(http.StatusOK)
		default:
			logger.Error.Printf("WhatsApp webhook: notification not processed (code %s); Meta will send it again", code)
			c.Status(http.StatusServiceUnavailable)
		}
	}
}

// validWhatsAppSignature checks the X-Hub-Signature-256 header: "sha256=" and the hex HMAC-SHA256
// of the body with the app secret. The comparison takes the same time wherever the values differ.
func validWhatsAppSignature(header string, body []byte, appSecret string) bool {
	const prefix = "sha256="
	if appSecret == "" || !strings.HasPrefix(header, prefix) {
		return false
	}
	got, err := hex.DecodeString(header[len(prefix):])
	if err != nil {
		return false
	}
	mac := hmac.New(sha256.New, []byte(appSecret))
	mac.Write(body)
	return hmac.Equal(got, mac.Sum(nil))
}

// sameSecret compares two secrets in a time that depends on neither: both are hashed first, so
// not even their lengths show.
func sameSecret(got, want string) bool {
	if want == "" {
		return false
	}
	g, w := sha256.Sum256([]byte(got)), sha256.Sum256([]byte(want))
	return subtle.ConstantTimeCompare(g[:], w[:]) == 1
}
