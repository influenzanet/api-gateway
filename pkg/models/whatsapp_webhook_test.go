package models

import (
	"reflect"
	"testing"
)

// The webhook route is on with all three settings, off without any, and a partial configuration
// names what is missing.
// Fails if the route can come up without its secret, verify token or messaging address, or if a
// partial configuration goes unreported.
func TestWhatsAppWebhookConfigFromEnv(t *testing.T) {
	full := map[string]string{
		"WHATSAPP_APP_SECRET":           " app-secret ",
		"WHATSAPP_WEBHOOK_VERIFY_TOKEN": "verify-token",
		"ADDR_MESSAGING_SERVICE":        "messaging:5004",
	}
	read := func(env map[string]string) func(string) string { return func(k string) string { return env[k] } }

	cfg, on, missing := WhatsAppWebhookConfigFromEnv(read(full))
	want := WhatsAppWebhookConfig{AppSecret: "app-secret", VerifyToken: "verify-token", MessagingServiceAddr: "messaging:5004"}
	if !on || missing != nil || cfg != want {
		t.Errorf("full: %+v on=%v missing=%v, want %+v on", cfg, on, missing, want)
	}

	if cfg, on, missing := WhatsAppWebhookConfigFromEnv(read(nil)); on || missing != nil || cfg != (WhatsAppWebhookConfig{}) {
		t.Errorf("empty: %+v on=%v missing=%v, want off with nothing missing", cfg, on, missing)
	}

	for name := range full {
		partial := map[string]string{}
		for k, v := range full {
			if k != name {
				partial[k] = v
			}
		}
		if name == "WHATSAPP_APP_SECRET" {
			partial[name] = "   " // only spaces counts as unset
		}
		cfg, on, missing := WhatsAppWebhookConfigFromEnv(read(partial))
		if on || cfg != (WhatsAppWebhookConfig{}) || !reflect.DeepEqual(missing, []string{name}) {
			t.Errorf("without %s: %+v on=%v missing=%v, want off with %s missing", name, cfg, on, missing, name)
		}
	}
}
