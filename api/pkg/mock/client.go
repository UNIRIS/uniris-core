package mock

import (
	"github.com/uniris/uniris-core/api/pkg/adding"
	"github.com/uniris/uniris-core/api/pkg/listing"
)

//Client a wrapped of the client methods
type Client interface {
	listing.RobotClient
	adding.RobotClient
}

type client struct {
}

//NewClient creates a mock of the blockchain client
func NewClient() Client {
	return client{}
}

func (c client) GetAccount(listing.AccountRequest) (listing.AccountResult, error) {
	return listing.AccountResult{
		EncryptedAESKey:     "encrypted_aes_key",
		EncryptedAddrPerson: "addr_wallet_person",
		EncryptedWallet:     "encrypted_wallet",
	}, nil
}

func (c client) AddAccount(adding.EnrollmentRequest) (adding.EnrollmentResult, error) {
	return adding.EnrollmentResult{
		Hash:             "hash of the udpated wallet",
		SignatureRequest: "signature of the response",
	}, nil
}
