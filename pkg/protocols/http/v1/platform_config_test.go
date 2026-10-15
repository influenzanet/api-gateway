package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/influenzanet/api-gateway/pkg/models"
	umAPI "github.com/influenzanet/user-management-service/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GET /v1/platform-config is public: the participant site calls it before anyone logs in to
// decide whether to show the phone and WhatsApp controls, and keeps them hidden on any error.

type platformConfigFake struct {
	umAPI.UserManagementApiClient
	requests []*umAPI.PlatformConfigReq
	deadline bool
	resp     *umAPI.PlatformConfig
	err      error
}

func (f *platformConfigFake) GetPlatformConfig(ctx context.Context, in *umAPI.PlatformConfigReq, opts ...grpc.CallOption) (*umAPI.PlatformConfig, error) {
	f.requests = append(f.requests, in)
	_, f.deadline = ctx.Deadline()
	return f.resp, f.err
}

func getPlatformConfig(fake *platformConfigFake, query string) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := NewHTTPHandler(&models.APIClients{UserManagement: fake}, models.UseEndpoints{}, nil)
	h.AddUserManagementParticipantAPI(router.Group("/v1"))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/platform-config"+query, nil))
	return rec
}

func TestPlatformConfigAnswersTheBackendDecision(t *testing.T) {
	for _, enabled := range []bool{true, false} {
		fake := &platformConfigFake{resp: &umAPI.PlatformConfig{WhatsappEnabled: enabled}}
		rec := getPlatformConfig(fake, "?instance=italy")

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200: %s", rec.Code, rec.Body)
		}
		// The web core reads exactly this key; false must be written out, not omitted.
		var body map[string]interface{}
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("invalid JSON %q: %v", rec.Body, err)
		}
		if len(body) != 1 || body["whatsAppEnabled"] != enabled {
			t.Errorf("body = %s, want {\"whatsAppEnabled\":%v}", rec.Body, enabled)
		}
		if got := rec.Header().Get("Cache-Control"); got != "public, max-age=60" {
			t.Errorf("Cache-Control = %q, want a short public cache", got)
		}
		if len(fake.requests) != 1 || fake.requests[0].InstanceId != "italy" {
			t.Errorf("user-management called with %v, want one call for instance italy", fake.requests)
		}
	}
}

func TestPlatformConfigRefusesAMissingOrMalformedInstance(t *testing.T) {
	for _, query := range []string{
		"",
		"?instance=",
		"?instance=a%2Fb",
		"?instance=it%20aly",
		"?instance=%3Cscript%3E",
		"?instance=" + strings.Repeat("a", 65),
	} {
		fake := &platformConfigFake{resp: &umAPI.PlatformConfig{WhatsappEnabled: true}}
		rec := getPlatformConfig(fake, query)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("%q: status = %d, want 400", query, rec.Code)
		}
		if len(fake.requests) != 0 {
			t.Errorf("%q: user-management was called for an invalid instance", query)
		}
		if got := rec.Header().Get("Cache-Control"); got != "no-store" {
			t.Errorf("%q: Cache-Control = %q, want no-store on an error", query, got)
		}
	}
}

func TestPlatformConfigAcceptsTheInstanceIdAlphabet(t *testing.T) {
	fake := &platformConfigFake{resp: &umAPI.PlatformConfig{WhatsappEnabled: true}}
	rec := getPlatformConfig(fake, "?instance=Infl_web-2026&ignored=1")
	if rec.Code != http.StatusOK || len(fake.requests) != 1 || fake.requests[0].InstanceId != "Infl_web-2026" {
		t.Errorf("status = %d, requests = %v", rec.Code, fake.requests)
	}
}

func TestPlatformConfigPassesUserManagementErrorsOn(t *testing.T) {
	for _, tc := range []struct {
		code codes.Code
		want int
	}{
		{codes.InvalidArgument, http.StatusBadRequest},
		{codes.Unavailable, http.StatusServiceUnavailable},
	} {
		fake := &platformConfigFake{err: status.Error(tc.code, "invalid arguments")}
		rec := getPlatformConfig(fake, "?instance=unknown")
		if rec.Code != tc.want {
			t.Errorf("%v: status = %d, want %d", tc.code, rec.Code, tc.want)
		}
		if strings.Contains(rec.Body.String(), "whatsAppEnabled") {
			t.Errorf("%v: an error answer carries a WhatsApp flag: %s", tc.code, rec.Body)
		}
		if got := rec.Header().Get("Cache-Control"); got != "no-store" {
			t.Errorf("%v: Cache-Control = %q, want no-store on an error", tc.code, got)
		}
	}
}

func TestPlatformConfigBoundsTheUserManagementCall(t *testing.T) {
	fake := &platformConfigFake{resp: &umAPI.PlatformConfig{}}
	getPlatformConfig(fake, "?instance=italy")
	if !fake.deadline {
		t.Error("the public route calls user-management without a deadline")
	}
}
