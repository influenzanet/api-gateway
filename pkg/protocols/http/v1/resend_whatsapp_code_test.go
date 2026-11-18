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

// postResendWhatsAppCode sends the resend request of a confirmed account through the real
// router, with fake as user-management.
func postResendWhatsAppCode(t *testing.T, fake *fakeUserManagementClient) *httptest.ResponseRecorder {
	t.Helper()
	router := newParticipantAPIRouter(fake)
	req := httptest.NewRequest(http.MethodPost, "/v1/user/contact/resend-whatsapp-code", nil)
	req.Header.Set("Authorization", "Bearer "+confirmedAccountToken)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

// An account without a phone number has nothing to resend a code to: the gateway answers 400
// and never asks user-management to send to an empty address.
// Fails if the request reaches ResendContactVerification without an address, or the refusal
// is not a 400 that says why.
func TestResendWhatsAppCodeWithoutAPhoneNumber(t *testing.T) {
	for _, tt := range []struct {
		name string
		user *umAPI.User
	}{
		{name: "no contact infos", user: &umAPI.User{}},
		{name: "an e-mail only", user: &umAPI.User{ContactInfos: []*umAPI.ContactInfo{
			{Type: "email", Address: &umAPI.ContactInfo_Email{Email: "participant@example.it"}},
		}}},
		{name: "a phone entry with an empty number", user: &umAPI.User{ContactInfos: []*umAPI.ContactInfo{
			{Type: "phone", Address: &umAPI.ContactInfo_Phone{Phone: ""}},
		}}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			fake := &fakeUserManagementClient{user: tt.user}

			rec := postResendWhatsAppCode(t, fake)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("got status %d, want 400; body: %s", rec.Code, rec.Body.String())
			}
			if !strings.Contains(rec.Body.String(), "no phone number on this account") {
				t.Errorf("the refusal does not say why: %s", rec.Body.String())
			}
			if calls := fake.callsTo("ResendContactVerification"); calls != 0 {
				t.Errorf("user-management was asked to resend %d times, to %q", calls, fake.resendAddresses)
			}
		})
	}
}

// With a phone number, the resend goes to user-management with that number, whatever other
// contact infos come first; and a failing user lookup is reported as it was, without a resend.
// Fails if the no-phone check also refuses an account that has a number, or hides the lookup
// failure behind the no-phone answer.
func TestResendWhatsAppCodeWithAPhoneNumber(t *testing.T) {
	fake := &fakeUserManagementClient{user: &umAPI.User{ContactInfos: []*umAPI.ContactInfo{
		{Type: "email", Address: &umAPI.ContactInfo_Email{Email: "participant@example.it"}},
		{Type: "phone", Address: &umAPI.ContactInfo_Phone{Phone: "+393339876543"}},
	}}}

	rec := postResendWhatsAppCode(t, fake)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200; body: %s", rec.Code, rec.Body.String())
	}
	if len(fake.resendAddresses) != 1 || fake.resendAddresses[0] != "+393339876543" {
		t.Errorf("expected one resend to +393339876543, got %q", fake.resendAddresses)
	}

	failing := &fakeUserManagementClient{getUserErr: status.Error(codes.Unavailable, "user-management down")}
	rec = postResendWhatsAppCode(t, failing)
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("a failed user lookup: got status %d, want 503; body: %s", rec.Code, rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), "no phone number") {
		t.Errorf("a failed user lookup was reported as a missing phone: %s", rec.Body.String())
	}
	if calls := failing.callsTo("ResendContactVerification"); calls != 0 {
		t.Errorf("user-management was asked to resend after the lookup failed")
	}
}
