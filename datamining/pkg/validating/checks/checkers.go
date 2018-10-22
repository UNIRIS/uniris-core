package checks

import (
	"github.com/uniris/uniris-core/datamining/pkg"
)

//BioChecker defines methods to validate bio data
type BioChecker interface {
	CheckBioWallet(*datamining.BioData) error
}

//DataChecker defines methods to validate wallet data
type DataChecker interface {
	CheckDataWallet(*datamining.WalletData) error
}
