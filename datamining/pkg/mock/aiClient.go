package mock

import (
	"net"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/transport/rpc/internalrpc"
)

type aiClient struct{}

//NewAIClient create a new mock of the AI client
func NewAIClient() internalrpc.AIClient {
	return aiClient{}
}

func (c aiClient) GetBiometricStoragePool(personHash string) (datamining.Pool, error) {
	return datamining.NewPool(datamining.Peer{IP: net.ParseIP("127.0.0.1")}), nil
}
func (c aiClient) GetKeychainStoragePool(address string) (datamining.Pool, error) {
	return datamining.NewPool(datamining.Peer{IP: net.ParseIP("127.0.0.1")}), nil
}

func (c aiClient) GetMasterPeer(txHash string) (datamining.Peer, error) {
	return datamining.Peer{IP: net.ParseIP("127.0.0.1")}, nil
}

func (c aiClient) GetValidationPool(txHash string) (datamining.Pool, error) {
	return datamining.NewPool(datamining.Peer{IP: net.ParseIP("127.0.0.1")}), nil
}
