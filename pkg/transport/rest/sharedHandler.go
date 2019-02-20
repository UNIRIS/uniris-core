package rest

import (
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

		emPublicKey := c.Query("emitter_public_key")

		if _, err := crypto.IsPublicKey(emPublicKey); err != nil {
			c.JSON(http.StatusBadRequest, httpError{
				Error:     fmt.Sprintf("emitter_public_key: %s", err.Error()),
				Status:    http.StatusText(http.StatusBadRequest),
				Timestamp: time.Now().Unix(),
			})
			return
		}

		auth, err := shared.IsEmitterKeyAuthorized(emPublicKey)
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

		emKeys := make([]emitterSharedKeys, 0)
		for _, k := range sharedEmKeys {
			emKeys = append(emKeys, emitterSharedKeys{
				EncryptedPrivateKey: k.EncryptedPrivateKey(),
				PublicKey:           k.PublicKey(),
			})
		}
		c.JSON(http.StatusOK, sharedKeysResponse{
			NodePublicKey: nodeLastKeys.PublicKey(),
			EmitterKeys:   emKeys,
		})
	}
}
