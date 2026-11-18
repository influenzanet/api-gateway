package models

import (
	"reflect"
	"testing"
)

// The two webhook secrets decide whether the route is on; the messaging address is required only
// once they are set, since deployments pass it to the participant API for other reasons.
// Fails if the route can come up without its secret, verify token or messaging address, if a
// partial configuration goes unreported, or if the messaging address alone is reported as one.
func TestWhatsAppWebhookConfigFromEnv(t *testing.T) {
	read := func(env map[string]string) func(string) string { return func(k string) string { return env[k] } }
	const (
		secret = "WHATSAPP_APP_SECRET"
		token  = "WHATSAPP_WEBHOOK_VERIFY_TOKEN"
		addr   = "ADDR_MESSAGING_SERVICE"
	)
	for _, c := range []struct {
		name        string
		env         map[string]string
		wantOn      bool
		wantMissing []string
	}{
		{"all three, spaces trimmed", map[string]string{secret: " app-secret ", token: "verify-token", addr: "messaging:5004"}, true, nil},
		{"nothing", nil, false, nil},
		{"only the messaging address: off, nothing to report", map[string]string{addr: "messaging:5004"}, false, nil},
		{"secrets set to spaces count as unset", map[string]string{secret: " ", token: "\t", addr: "messaging:5004"}, false, nil},
		{"no verify token", map[string]string{secret: "s", addr: "messaging:5004"}, false, []string{token}},
		{"no app secret", map[string]string{token: "t", addr: "messaging:5004"}, false, []string{secret}},
		{"no messaging address", map[string]string{secret: "s", token: "t"}, false, []string{addr}},
		{"only the token", map[string]string{token: "t"}, false, []string{secret, addr}},
	} {
		t.Run(c.name, func(t *testing.T) {
			cfg, on, missing := WhatsAppWebhookConfigFromEnv(read(c.env))
			if on != c.wantOn || !reflect.DeepEqual(missing, c.wantMissing) {
				t.Fatalf("on=%v missing=%v, want on=%v missing=%v", on, missing, c.wantOn, c.wantMissing)
			}
			want := WhatsAppWebhookConfig{}
			if c.wantOn {
				want = WhatsAppWebhookConfig{AppSecret: "app-secret", VerifyToken: "verify-token", MessagingServiceAddr: "messaging:5004"}
			}
			if cfg != want {
				t.Errorf("config %+v, want %+v", cfg, want)
			}
		})
	}
}
