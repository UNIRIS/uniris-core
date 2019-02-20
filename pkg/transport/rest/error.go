package rest

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func parseGrpcError(err error) (code int, message string) {
	statusErr, _ := status.FromError(err)
	switch statusErr.Code() {
	case codes.InvalidArgument:
		return http.StatusBadRequest, statusErr.Message()
	case codes.NotFound:
		return http.StatusNotFound, statusErr.Message()
	}
	return http.StatusInternalServerError, statusErr.Message()
}

type httpError struct {
	Error     string `json:"error"`
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
	code      int
}
