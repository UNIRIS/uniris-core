package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/crypto"
	"google.golang.org/grpc"
)

//NewSharedHandler creates a new HTTP handler for the shared endpoints
func NewSharedHandler(apiGroup *gin.RouterGroup, internalPort int) {

	apiGroup.GET("/sharedkeys", getSharedKeys(internalPort))
}

func getSharedKeys(internalPort int) func(*gin.Context) {
	return func(c *gin.Context) {

		emPublicKey := c.Query("emitter_public_key")

		//Check the emitter public key parameters
		if _, err := crypto.IsPublicKey(emPublicKey); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("emitter_public_key: %s", err.Error())})
			return
		}

		//Call the internal datamining to get the last shared keys
		serverAddr := fmt.Sprintf("localhost:%d", internalPort)
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		defer conn.Close()

		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		cli := api.NewInternalServiceClient(conn)
		res, err := cli.GetLastSharedKeys(context.Background(), &api.LastSharedKeysRequest{
			EmitterPublicKey: emPublicKey,
			Timestamp:        time.Now().Unix(),
		})
		if err != nil {
			c.JSON(parseGrpcError(err))
			return
		}

		//Building the JSON response
		emKeys := make([]map[string]string, 0)
		for _, k := range res.EmitterKeys {
			emKeys = append(emKeys, map[string]string{
				"public_key":            k.PublicKey,
				"encrypted_private_key": k.EncryptedPrivateKey,
			})
		}

		c.JSON(http.StatusOK, map[string]interface{}{
			"shared_node_public_key": res.NodePublicKey,
			"shared_emitter_keys":    emKeys,
		})
	}
}
