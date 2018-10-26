package mining

import (
	"errors"

	"github.com/uniris/uniris-core/datamining/pkg"
)

//ErrInvalidHash is returned when a invalid data
var ErrInvalidHash = errors.New("Invalid hash")

//PreviousDataHasher defines methods to hash previous data
type PreviousDataHasher interface {
	HashWallet(*datamining.Wallet) (string, error)
}

//PreviousDataChecker defines methods to check previous data
type PreviousDataChecker interface {
	CheckPreviousWallet(w *datamining.Wallet, txHash string) error
}

type previousChecker struct {
	h PreviousDataHasher
}

//NewPreviousDataChecker creates a new previous data checker
func NewPreviousDataChecker(h PreviousDataHasher) PreviousDataChecker {
	return previousChecker{h}
}

func (pc previousChecker) CheckPreviousWallet(w *datamining.Wallet, txHash string) error {
	hash, err := pc.h.HashWallet(w)
	if err != nil {
		return err
	}

	if hash != txHash {
		return ErrInvalidHash
	}

	return nil
}
