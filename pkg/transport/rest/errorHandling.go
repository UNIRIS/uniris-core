package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func parseGrpcError(err error) (int, interface{}) {
	statusErr, _ := status.FromError(err)
	switch statusErr.Code() {
	case codes.Internal:
		return http.StatusInternalServerError, gin.H{"error": statusErr.Message()}
	case codes.InvalidArgument:
		return http.StatusBadRequest, gin.H{"error": statusErr.Message()}
	case codes.NotFound:
		return http.StatusNotFound, gin.H{"error": statusErr.Message()}
	}
	return http.StatusInternalServerError, gin.H{"error": statusErr.Message()}
}
