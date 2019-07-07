package main

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

func grpcStatusToHTTP(status codes.Code) int {
	switch status {
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	}
	return http.StatusInternalServerError
}
