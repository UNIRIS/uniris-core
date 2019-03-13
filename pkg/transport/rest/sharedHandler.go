package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/shared"
)

//GetSharedKeysHandler defines an HTTP handler to retrieve the shared keys
func GetSharedKeysHandler(techReader shared.TechDatabaseReader) func(*gin.Context) {
	return func(c *gin.Context) {

		emPubKey, httpErr := extractEmitterPublicKey(c)
		if httpErr != nil {
			c.JSON(httpErr.code, httpErr)
			return
		}

		auth, err := shared.IsEmitterKeyAuthorized(emPubKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, httpError{
				Error:     err.Error(),
				Status:    http.StatusText(http.StatusInternalServerError),
				Timestamp: time.Now().Unix(),
			})
			return
		}
		if !auth {
			c.JSON(http.StatusUnauthorized, httpError{
				Error:     "emitter not authorized",
				Status:    http.StatusText(http.StatusUnauthorized),
				Timestamp: time.Now().Unix(),
			})
			return
		}

		nodeLastKeys, err := techReader.NodeLastKeys()
		if err != nil {
			c.JSON(http.StatusInternalServerError, httpError{
				Error:     err.Error(),
				Status:    http.StatusText(http.StatusInternalServerError),
				Timestamp: time.Now().Unix(),
			})
			return
		}

		sharedEmKeys, err := techReader.EmitterKeys()
		if err != nil {
			c.JSON(http.StatusInternalServerError, httpError{
				Error:     err.Error(),
				Status:    http.StatusText(http.StatusInternalServerError),
				Timestamp: time.Now().Unix(),
			})
			return
		}

		res, err := createSharedKeyResponse(sharedEmKeys, nodeLastKeys)
		if err != nil {
			c.JSON(http.StatusInternalServerError, httpError{
				Error:     err.Error(),
				Status:    http.StatusText(http.StatusInternalServerError),
				Timestamp: time.Now().Unix(),
			})
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

func extractEmitterPublicKey(c *gin.Context) (crypto.PublicKey, *httpError) {
	emPubParam := c.Query("emitter_public_key")
	if emPubParam == "" {
		return nil, &httpError{
			Error:     "emitter public key is missing",
			Status:    http.StatusText(http.StatusBadRequest),
			Timestamp: time.Now().Unix(),
			code:      http.StatusBadRequest,
		}
	}
	emPubBytes, err := hex.DecodeString(emPubParam)
	if err != nil {
		return nil, &httpError{
			Error:     fmt.Sprintf("emitter public key is not in hexadecimal"),
			Status:    http.StatusText(http.StatusBadRequest),
			Timestamp: time.Now().Unix(),
			code:      http.StatusBadRequest,
		}
	}

	emPublicKey, err := crypto.ParsePublicKey(emPubBytes)
	if err != nil {
		return nil, &httpError{
			Error:     fmt.Sprintf("emitter public key is not valid: %s", err.Error()),
			Status:    http.StatusText(http.StatusBadRequest),
			Timestamp: time.Now().Unix(),
			code:      http.StatusBadRequest,
		}
	}
	return emPublicKey, nil
}

func createSharedKeyResponse(sharedEmKeys shared.EmitterKeys, nodeLastKeys shared.NodeKeyPair) (sharedKeysResponse, error) {
	emKeys := make([]emitterSharedKeys, 0)
	for _, k := range sharedEmKeys {
		emPubBytes, err := k.PublicKey().Marshal()
		if err != nil {
			return sharedKeysResponse{}, err
		}
		emKeys = append(emKeys, emitterSharedKeys{
			EncryptedPrivateKey: hex.EncodeToString(k.EncryptedPrivateKey()),
			PublicKey:           hex.EncodeToString(emPubBytes),
		})
	}

	nodePubBytes, err := nodeLastKeys.PublicKey().Marshal()
	if err != nil {
		return sharedKeysResponse{}, err
	}

	return sharedKeysResponse{
		NodePublicKey: hex.EncodeToString(nodePubBytes),
		EmitterKeys:   emKeys,
	}, nil
}
