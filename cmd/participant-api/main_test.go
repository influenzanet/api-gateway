package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/influenzanet/api-gateway/pkg/models"
)

var routerConf = models.Config{AllowOrigins: []string{"https://participant.example.org"}}

var routerWebhook = models.WhatsAppWebhookConfig{AppSecret: "s", VerifyToken: "token-value", MessagingServiceAddr: "m:1"}

const subscriptionCheck = "/v1/webhooks/whatsapp?hub.mode=subscribe&hub.verify_token=token-value&hub.challenge=7"

// The webhook route exists only when it is configured.
// Fails if the route is registered without its configuration, or missing with it.
func TestWebhookRouteOnlyWhenConfigured(t *testing.T) {
	for _, on := range []bool{false, true} {
		router := newRouter(&models.APIClients{}, routerConf, routerWebhook, on, &bytes.Buffer{})
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, subscriptionCheck, nil))
		want := http.StatusNotFound
		if on {
			want = http.StatusOK
		}
		if rec.Code != want {
			t.Errorf("webhook on=%v: status %d, want %d", on, rec.Code, want)
		}
	}
}

// The request log keeps the path and leaves the query string out: Meta's subscription check
// carries the verify token there.
// Fails if the log carries the query string, or if requests are no longer logged at all.
func TestRequestLogLeavesQueryStringsOut(t *testing.T) {
	var log bytes.Buffer
	router := newRouter(&models.APIClients{}, routerConf, routerWebhook, true, &log)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, subscriptionCheck, nil))

	if !strings.Contains(log.String(), "/v1/webhooks/whatsapp") {
		t.Fatalf("the request was not logged: %q", log.String())
	}
	if strings.Contains(log.String(), "token-value") || strings.Contains(log.String(), "hub.") {
		t.Errorf("the log carries the query string: %q", log.String())
	}
}
