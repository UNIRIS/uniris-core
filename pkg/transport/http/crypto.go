package http

import uniris "github.com/uniris/uniris-core/pkg"

type TransactionDecrypter interface {
	DecryptTransaction(cipher string) (uniris.Transaction, error)
}
