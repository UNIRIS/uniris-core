package rest

import (
	"net/http"

	"github.com/uniris/uniris-core/api/pkg/crypto"
)

//ErrorMessage define an HTTP error
type ErrorMessage struct {
	Message   string `json:"error_message"`
	Signature string `json:"error_signature"`
	Code      int    `json:"error_code"`
}

func createError(handleErrorCode int, handleErr error, robotPvKey string) ErrorMessage {
	sig, err := crypto.HashAndSign(robotPvKey, handleErr.Error())
	if err != nil {
		return ErrorMessage{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	return ErrorMessage{
		Message:   handleErr.Error(),
		Signature: string(sig),
		Code:      handleErrorCode,
	}
}
