package mock

import (
	"github.com/uniris/uniris-core/api/pkg/listing"
)

type client struct {
	listing.RobotClient
}

func NewClient() client {
	return client{}
}

func (c client) GetAccount(listing.AccountRequest) (listing.AccountResult, error) {
	return listing.AccountResult{
		EncryptedAESKey:  "encrypted_aes_key",
		AddrWalletPerson: "addr_wallet_person",
		EncryptedWallet:  "encrypted_wallet",
	}, nil
}
