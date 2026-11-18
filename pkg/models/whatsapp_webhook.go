package models

import (
	"strings"
)

// WhatsAppWebhookConfig is the configuration of the WhatsApp webhook route of the participant API.
type WhatsAppWebhookConfig struct {
	// AppSecret is the Meta app secret, the key of the signature Meta puts on every notification.
	AppSecret string
	// VerifyToken is the value Meta sends back when the webhook is subscribed in the Meta app.
	VerifyToken string
	// MessagingServiceAddr is where the notifications are forwarded.
	MessagingServiceAddr string
}

// WhatsAppWebhookConfigFromEnv reads WHATSAPP_APP_SECRET, WHATSAPP_WEBHOOK_VERIFY_TOKEN and
// ADDR_MESSAGING_SERVICE. The two secrets decide: with neither, the route is off and nothing is
// wrong (a deployment without the webhook; the messaging address alone is common, deployments
// set it for other reasons). With both, and the messaging address, the route is on. Anything in
// between is a half-done configuration: the route stays off and missing names what is lacking,
// so that it is reported instead of silently leaving the webhook off.
func WhatsAppWebhookConfigFromEnv(getenv func(string) string) (cfg WhatsAppWebhookConfig, on bool, missing []string) {
	cfg = WhatsAppWebhookConfig{
		AppSecret:            strings.TrimSpace(getenv("WHATSAPP_APP_SECRET")),
		VerifyToken:          strings.TrimSpace(getenv("WHATSAPP_WEBHOOK_VERIFY_TOKEN")),
		MessagingServiceAddr: strings.TrimSpace(getenv("ADDR_MESSAGING_SERVICE")),
	}
	if cfg.AppSecret == "" && cfg.VerifyToken == "" {
		return WhatsAppWebhookConfig{}, false, nil
	}
	for _, v := range []struct{ name, value string }{
		{"WHATSAPP_APP_SECRET", cfg.AppSecret},
		{"WHATSAPP_WEBHOOK_VERIFY_TOKEN", cfg.VerifyToken},
		{"ADDR_MESSAGING_SERVICE", cfg.MessagingServiceAddr},
	} {
		if v.value == "" {
			missing = append(missing, v.name)
		}
	}
	if len(missing) > 0 {
		return WhatsAppWebhookConfig{}, false, missing
	}
	return cfg, true, nil
}
