package v1

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	umAPI "github.com/influenzanet/user-management-service/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// The weekly unsubscribe link: the token of the query goes to user-management with the weekly
// subscription, and its refusal reaches the page with a status the page can tell apart.
func TestUnsubscribeWeekly(t *testing.T) {
	call := func(fake *fakeUserManagementClient, query string) *httptest.ResponseRecorder {
		router := newParticipantAPIRouter(fake)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/user/unsubscribe-weekly"+query, nil))
		return rec
	}

	t.Run("the token and the weekly subscription are passed", func(t *testing.T) {
		fake := &fakeUserManagementClient{}
		rec := call(fake, "?token=the-token")
		if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "unsubscribed") {
			t.Errorf("got %d %s, want 200 with the service status", rec.Code, rec.Body.String())
		}
		if len(fake.unsubscribeRequests) != 1 {
			t.Fatalf("expected one request, got %d", len(fake.unsubscribeRequests))
		}
		req := fake.unsubscribeRequests[0]
		if req.Token != "the-token" || req.Subscription != umAPI.Subscription_SUBSCRIPTION_WEEKLY {
			t.Errorf("got token %q subscription %v, want \"the-token\" weekly", req.Token, req.Subscription)
		}
	})

	t.Run("no token: user-management decides", func(t *testing.T) {
		fake := &fakeUserManagementClient{unsubscribeErr: status.Error(codes.InvalidArgument, "missing argument")}
		rec := call(fake, "")
		if len(fake.unsubscribeRequests) != 1 || fake.unsubscribeRequests[0].Token != "" {
			t.Errorf("expected one request with no token, got %v", fake.unsubscribeRequests)
		}
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got %d, want 400", rec.Code)
		}
	})

	for _, c := range []struct {
		code codes.Code
		want int
	}{
		{codes.InvalidArgument, http.StatusBadRequest},
		{codes.Internal, http.StatusInternalServerError},
	} {
		t.Run("refused with "+c.code.String(), func(t *testing.T) {
			fake := &fakeUserManagementClient{unsubscribeErr: status.Error(c.code, "token expired")}
			rec := call(fake, "?token=old")
			if rec.Code != c.want || !strings.Contains(rec.Body.String(), "token expired") {
				t.Errorf("got %d %s, want %d with the message", rec.Code, rec.Body.String(), c.want)
			}
		})
	}
}
