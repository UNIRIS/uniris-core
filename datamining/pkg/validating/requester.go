package validating

import (
	"net"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
)

//Peer to reach for data validating
type Peer struct {
	IP        net.IP
	Port      int
	PublicKey string
}

//ValidationRequester defines methods to reach peer to valid wallet data
type ValidationRequester interface {
	RequestWalletValidation(Peer, *datamining.WalletData) (datamining.Validation, error)
	RequestBioValidation(Peer, *datamining.BioData) (datamining.Validation, error)
}
