package checks

import (
	"github.com/uniris/uniris-core/datamining/pkg"
)

//BioDataChecker defines methods to validate bio wallet
type BioDataChecker interface {
	CheckBioData(*datamining.BioData) error
}

//WalletDataChecker defines methods to validate wallet data
type WalletDataChecker interface {
	CheckWalletData(*datamining.WalletData) error
}
