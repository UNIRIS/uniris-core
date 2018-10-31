package checks

import (
	"errors"
)

//ErrInvalidTransaction is returned a transaction is invalid
var ErrInvalidTransaction = errors.New("Invalid transaction")

//TransactionDataHasher define methods to hash transaction data
type TransactionDataHasher interface {
	HashTransactionData(data interface{}) (string, error)
}

type integrityChecker struct {
	h TransactionDataHasher
}

//NewIntegrityChecker creates an intergrity checker
func NewIntegrityChecker(h TransactionDataHasher) Handler {
	return integrityChecker{h}
}

func (c integrityChecker) CheckData(data interface{}, txHash string) error {
	hash, err := c.h.HashTransactionData(data)
	if err != nil {
		return err
	}
	if hash != txHash {
		return ErrInvalidTransaction
	}
	return nil
}
