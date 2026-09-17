package v1

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/influenzanet/api-gateway/pkg/models"
	"github.com/influenzanet/go-utils/pkg/api_types"
	umAPI "github.com/influenzanet/user-management-service/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// F-57: the route -> middleware matrix of the participant API is only expressed
// in AddUserManagementParticipantAPI. Nothing else stops a phone route from
// being registered outside the authenticated group, or without
// CheckAccountConfirmed, which is exactly the M-1 class of regression. These
// tests drive the real router and assert, per route, who is refused and who
// reaches the handler.

// Tokens the fake user-management client knows about. Any other value is
// rejected, which is how an unknown or expired token is simulated.
const (
	confirmedAccountToken   = "token-of-a-confirmed-account"
	unconfirmedAccountToken = "token-of-an-unconfirmed-account"
	unknownToken            = "token-nobody-issued"
)

// Error messages the middlewares emit, so a failure says which guard let the
// request through instead of only showing a status code (three of them use 400).
const (
	msgNoToken       = "no Authorization token found"
	msgInvalidToken  = "error during token validation"
	msgNotConfirmed  = "account not confirmed yet"
	contactRoutePath = "/v1/user/contact"
)

// fakeUserManagementClient records the gRPC methods the handlers call, so a
// test can tell a handler that was actually reached from one the middleware
// chain refused before it. ValidateJWT is deliberately not recorded: it is a
// middleware call, not a sign that the handler ran.
type fakeUserManagementClient struct {
	umAPI.UserManagementApiClient
	calls []string
}

func (f *fakeUserManagementClient) record(method string) {
	f.calls = append(f.calls, method)
}

func (f *fakeUserManagementClient) callsTo(method string) int {
	count := 0
	for _, call := range f.calls {
		if call == method {
			count++
		}
	}
	return count
}

func (f *fakeUserManagementClient) ValidateJWT(ctx context.Context, in *umAPI.JWTRequest, opts ...grpc.CallOption) (*api_types.TokenInfos, error) {
	switch in.GetToken() {
	case confirmedAccountToken:
		return &api_types.TokenInfos{Id: "user-confirmed", InstanceId: "italy", AccountConfirmed: true}, nil
	case unconfirmedAccountToken:
		return &api_types.TokenInfos{Id: "user-unconfirmed", InstanceId: "italy", AccountConfirmed: false}, nil
	default:
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
}

func (f *fakeUserManagementClient) AddPhoneNumber(ctx context.Context, in *umAPI.PhoneMsg, opts ...grpc.CallOption) (*umAPI.User, error) {
	f.record("AddPhoneNumber")
	return &umAPI.User{}, nil
}

func (f *fakeUserManagementClient) EditPhoneNumber(ctx context.Context, in *umAPI.PhoneMsg, opts ...grpc.CallOption) (*umAPI.User, error) {
	f.record("EditPhoneNumber")
	return &umAPI.User{}, nil
}

func (f *fakeUserManagementClient) DeletePhoneNumber(ctx context.Context, in *api_types.TokenInfos, opts ...grpc.CallOption) (*umAPI.User, error) {
	f.record("DeletePhoneNumber")
	return &umAPI.User{}, nil
}

func (f *fakeUserManagementClient) VerifyWhatsAppCode(ctx context.Context, in *umAPI.VerifyWhatsAppCodeReq, opts ...grpc.CallOption) (*umAPI.User, error) {
	f.record("VerifyWhatsAppCode")
	return &umAPI.User{}, nil
}

// resendWhatsAppCode reads the current phone number first, so GetUser is the
// call that proves the handler was reached.
func (f *fakeUserManagementClient) GetUser(ctx context.Context, in *umAPI.UserReference, opts ...grpc.CallOption) (*umAPI.User, error) {
	f.record("GetUser")
	return &umAPI.User{}, nil
}

func (f *fakeUserManagementClient) ResendContactVerification(ctx context.Context, in *umAPI.ResendContactVerificationReq, opts ...grpc.CallOption) (*umAPI.ServiceStatus, error) {
	f.record("ResendContactVerification")
	return &umAPI.ServiceStatus{}, nil
}

func (f *fakeUserManagementClient) AddEmail(ctx context.Context, in *umAPI.ContactInfoMsg, opts ...grpc.CallOption) (*umAPI.User, error) {
	f.record("AddEmail")
	return &umAPI.User{}, nil
}

func (f *fakeUserManagementClient) RemoveEmail(ctx context.Context, in *umAPI.ContactInfoMsg, opts ...grpc.CallOption) (*umAPI.User, error) {
	f.record("RemoveEmail")
	return &umAPI.User{}, nil
}

func (f *fakeUserManagementClient) UpdateContactPreferences(ctx context.Context, in *umAPI.ContactPreferencesMsg, opts ...grpc.CallOption) (*umAPI.User, error) {
	f.record("UpdateContactPreferences")
	return &umAPI.User{}, nil
}

func (f *fakeUserManagementClient) VerifyContact(ctx context.Context, in *umAPI.TempToken, opts ...grpc.CallOption) (*umAPI.User, error) {
	f.record("VerifyContact")
	return &umAPI.User{}, nil
}

func (f *fakeUserManagementClient) UseUnsubscribeToken(ctx context.Context, in *umAPI.TempToken, opts ...grpc.CallOption) (*umAPI.ServiceStatus, error) {
	f.record("UseUnsubscribeToken")
	return &umAPI.ServiceStatus{}, nil
}

func (f *fakeUserManagementClient) SignupWithEmail(ctx context.Context, in *umAPI.SignupWithEmailMsg, opts ...grpc.CallOption) (*umAPI.TokenResponse, error) {
	f.record("SignupWithEmail")
	return &umAPI.TokenResponse{}, nil
}

// routeExpectation is one row of the matrix: the protection a route must have,
// and the gRPC call that proves its handler ran.
type routeExpectation struct {
	method string
	path   string
	// body is sent unchanged in every case, refusals included, so that
	// RequirePayload can never be mistaken for the guard under test.
	body string
	// grpcCall is the first user-management call the handler makes.
	grpcCall string
	// requiresToken: the route sits in the group carrying ExtractToken +
	// ValidateToken.
	requiresToken bool
	// requiresConfirmedAccount: the route also carries CheckAccountConfirmed.
	requiresConfirmedAccount bool
}

func (r routeExpectation) String() string {
	return r.method + " " + r.path
}

// routeMiddlewareMatrix mirrors AddUserManagementParticipantAPI. Every route
// under /v1/user/contact must appear here (TestContactRoutesAreAllInTheMatrix
// enforces it in both directions), plus the signup and unsubscribe routes as
// unauthenticated controls.
var routeMiddlewareMatrix = []routeExpectation{
	// The WhatsApp / phone routes: authenticated and confirmed-account only.
	{http.MethodPost, "/v1/user/contact/add-phone", `{"newPhone":"+393331234567"}`, "AddPhoneNumber", true, true},
	{http.MethodPost, "/v1/user/contact/change-phone", `{"newPhone":"+393337654321"}`, "EditPhoneNumber", true, true},
	{http.MethodDelete, "/v1/user/contact/delete-phone", "", "DeletePhoneNumber", true, true},
	{http.MethodPost, "/v1/user/contact/verify-whatsapp-code", `{"code":"123456"}`, "VerifyWhatsAppCode", true, true},
	{http.MethodPost, "/v1/user/contact/resend-whatsapp-code", "", "GetUser", true, true},
	// Email contact routes, same group.
	{http.MethodPost, "/v1/user/contact/add-email", `{"contactInfo":{}}`, "AddEmail", true, true},
	{http.MethodPost, "/v1/user/contact/remove-email", `{"contactInfo":{}}`, "RemoveEmail", true, false},
	// Channel preferences: authenticated, and deliberately usable before the
	// account is confirmed. Pinned as-is so a change is a conscious one.
	{http.MethodPost, "/v1/user/contact-preferences", `{"contactPreferences":{}}`, "UpdateContactPreferences", true, false},
	// Public by design: the link in the verification message carries its own
	// one-time token.
	{http.MethodPost, "/v1/user/contact-verification", `{"token":"one-time-token"}`, "VerifyContact", false, false},
	// Controls outside the /v1/user/contact prefix.
	{http.MethodGet, "/v1/user/unsubscribe-newsletter", "", "UseUnsubscribeToken", false, false},
	{http.MethodPost, "/v1/auth/signup-with-email", `{"email":"participant@example.it","password":"a-password","instanceId":"italy","preferredLanguage":"it"}`, "SignupWithEmail", false, false},
}

// newParticipantAPIRouter builds the router exactly as cmd/participant-api does,
// so the middleware chain under test is the registered one. Recovery is added
// on purpose: a route carrying CheckAccountConfirmed outside the ValidateToken
// group panics on c.MustGet, and the test must report that as a failure rather
// than crash the binary.
func newParticipantAPIRouter(fake *fakeUserManagementClient) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(gin.Recovery())
	h := NewHTTPHandler(
		&models.APIClients{UserManagement: fake},
		models.UseEndpoints{SignupWithEmail: true},
		nil,
	)
	h.AddUserManagementParticipantAPI(router.Group("/v1"))
	return router
}

// callRoute replays one route on a freshly built router with a fresh fake.
// An empty token means no Authorization header at all.
func callRoute(route routeExpectation, token string) (*httptest.ResponseRecorder, *fakeUserManagementClient) {
	fake := &fakeUserManagementClient{}
	router := newParticipantAPIRouter(fake)

	var body *bytes.Buffer
	if route.body == "" {
		body = bytes.NewBuffer(nil)
	} else {
		body = bytes.NewBufferString(route.body)
	}
	req := httptest.NewRequest(route.method, route.path, body)
	if route.body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec, fake
}

func assertRefusedBeforeHandler(t *testing.T, route routeExpectation, rec *httptest.ResponseRecorder, fake *fakeUserManagementClient, wantStatus int, wantMessage string) {
	t.Helper()
	if calls := fake.callsTo(route.grpcCall); calls != 0 {
		t.Errorf("%s: handler was reached (%s called %d times) although the request had to be refused", route, route.grpcCall, calls)
	}
	if rec.Code != wantStatus {
		t.Errorf("%s: got status %d, want %d; body: %s", route, rec.Code, wantStatus, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), wantMessage) {
		t.Errorf("%s: refusal body %s does not mention %q, so a different guard answered", route, rec.Body.String(), wantMessage)
	}
}

func assertHandlerReached(t *testing.T, route routeExpectation, rec *httptest.ResponseRecorder, fake *fakeUserManagementClient) {
	t.Helper()
	if rec.Code != http.StatusOK {
		t.Errorf("%s: got status %d, want 200; body: %s", route, rec.Code, rec.Body.String())
	}
	if calls := fake.callsTo(route.grpcCall); calls != 1 {
		t.Errorf("%s: %s called %d times, want exactly 1 (the handler must be reached)", route, route.grpcCall, calls)
	}
}

// TestParticipantRouteMiddlewareMatrix is the regression guard for M-1: every
// phone / WhatsApp route must refuse an anonymous caller and an unconfirmed
// account before the handler runs, and must serve a confirmed one.
func TestParticipantRouteMiddlewareMatrix(t *testing.T) {
	// Keep the signup route's recaptcha middleware out of the picture.
	t.Setenv("USE_RECAPTCHA", "false")
	t.Setenv("RECAPTCHA_SECRET", "")

	for _, route := range routeMiddlewareMatrix {
		route := route
		t.Run(route.String(), func(t *testing.T) {
			// (a) no credentials at all.
			rec, fake := callRoute(route, "")
			if route.requiresToken {
				// ExtractToken answers 400, not 401, when the header is
				// missing entirely; only an unusable token reaches
				// ValidateToken. Both are pinned, see the next case.
				assertRefusedBeforeHandler(t, route, rec, fake, http.StatusBadRequest, msgNoToken)
			} else {
				assertHandlerReached(t, route, rec, fake)
			}

			// (a bis) a token no one issued: refused by ValidateToken with 401.
			if route.requiresToken {
				rec, fake = callRoute(route, unknownToken)
				assertRefusedBeforeHandler(t, route, rec, fake, http.StatusUnauthorized, msgInvalidToken)
			}

			// (b) authenticated, account not confirmed yet.
			if route.requiresToken {
				rec, fake = callRoute(route, unconfirmedAccountToken)
				if route.requiresConfirmedAccount {
					assertRefusedBeforeHandler(t, route, rec, fake, http.StatusBadRequest, msgNotConfirmed)
				} else {
					assertHandlerReached(t, route, rec, fake)
				}
			}

			// (c) authenticated and confirmed: the handler must run.
			if route.requiresToken {
				rec, fake = callRoute(route, confirmedAccountToken)
				assertHandlerReached(t, route, rec, fake)
			}
		})
	}
}

// TestContactRoutesAreAllInTheMatrix keeps the table honest in both directions:
// a contact route registered without an expectation fails here instead of
// shipping unprotected, and a typo in the table cannot silently disable a row.
func TestContactRoutesAreAllInTheMatrix(t *testing.T) {
	router := newParticipantAPIRouter(&fakeUserManagementClient{})

	expected := map[string]bool{}
	for _, route := range routeMiddlewareMatrix {
		if strings.HasPrefix(route.path, contactRoutePath) {
			expected[route.String()] = false
		}
	}

	for _, registered := range router.Routes() {
		if !strings.HasPrefix(registered.Path, contactRoutePath) {
			continue
		}
		key := registered.Method + " " + registered.Path
		if _, ok := expected[key]; !ok {
			t.Errorf("route %s is registered under %s but has no expectation in routeMiddlewareMatrix: add one stating whether it needs authentication and a confirmed account", key, contactRoutePath)
			continue
		}
		expected[key] = true
	}

	for key, covered := range expected {
		if !covered {
			t.Errorf("routeMiddlewareMatrix expects %s but no such route is registered: the path is stale or misspelled", key)
		}
	}
}
