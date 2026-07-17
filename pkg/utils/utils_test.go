package utils

import (
	"net/http"
	"testing"

	"google.golang.org/grpc/codes"
)

func TestGRPCStatusToHTTP(t *testing.T) {
	testCases := []struct {
		code     codes.Code
		expected int
	}{
		{codes.Unauthenticated, http.StatusUnauthorized},
		{codes.InvalidArgument, http.StatusBadRequest},
		{codes.Unavailable, http.StatusServiceUnavailable},
		{codes.PermissionDenied, http.StatusUnauthorized},
		{codes.Unimplemented, http.StatusNotImplemented},
		{codes.ResourceExhausted, http.StatusTooManyRequests},
		{codes.Internal, http.StatusInternalServerError},
		{codes.Unknown, http.StatusInternalServerError},
	}

	for _, tc := range testCases {
		t.Run(tc.code.String(), func(t *testing.T) {
			if got := GRPCStatusToHTTP(tc.code); got != tc.expected {
				t.Errorf("wrong status for %s: %d instead of %d", tc.code, got, tc.expected)
			}
		})
	}
}
