package mock

import (
	"net"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	addAdding "github.com/uniris/uniris-core/datamining/pkg/account/adding"
	"github.com/uniris/uniris-core/datamining/pkg/transport/rpc"
)

type AIClient interface {
	rpc.AIClient
	addAdding.AIClient
}

type aiClient struct{}

//NewAIClient create a new mock of the AI client
func NewAIClient() AIClient {
	return aiClient{}
}

func (c aiClient) GetStoragePool(personHash string) (datamining.Pool, error) {
	return datamining.NewPool(datamining.Peer{IP: net.ParseIP("127.0.0.1")}), nil
}

func (c aiClient) GetMasterPeer(txHash string) (datamining.Peer, error) {
	return datamining.Peer{IP: net.ParseIP("127.0.0.1")}, nil
}

func (c aiClient) GetValidationPool(txHash string) (datamining.Pool, error) {
	return datamining.NewPool(datamining.Peer{IP: net.ParseIP("127.0.0.1")}), nil
}

func (c aiClient) CheckStorageAuthorization(txHash string) error {
	return nil
}

func (c aiClient) GetMininumValidations(txHash string) (int, error) {
	return 1, nil
}
