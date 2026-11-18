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
// ADDR_MESSAGING_SERVICE. The route is on only when all three are set. With none of them it is
// off and nothing is wrong (a deployment without WhatsApp); with only some, missing names the
// others, so that a half-done configuration is reported instead of silently leaving the
// webhook off.
func WhatsAppWebhookConfigFromEnv(getenv func(string) string) (cfg WhatsAppWebhookConfig, on bool, missing []string) {
	values := []struct {
		name  string
		value *string
	}{
		{"WHATSAPP_APP_SECRET", &cfg.AppSecret},
		{"WHATSAPP_WEBHOOK_VERIFY_TOKEN", &cfg.VerifyToken},
		{"ADDR_MESSAGING_SERVICE", &cfg.MessagingServiceAddr},
	}
	for _, v := range values {
		*v.value = strings.TrimSpace(getenv(v.name))
		if *v.value == "" {
			missing = append(missing, v.name)
		}
	}
	switch len(missing) {
	case 0:
		return cfg, true, nil
	case len(values):
		return WhatsAppWebhookConfig{}, false, nil
	}
	return WhatsAppWebhookConfig{}, false, missing
}
