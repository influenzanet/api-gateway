package utils

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

func GRPCStatusToHTTP(status codes.Code) int {
	switch status {
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.PermissionDenied:
		return http.StatusUnauthorized
	case codes.Unimplemented:
		return http.StatusNotImplemented
	}
	return http.StatusInternalServerError
}
