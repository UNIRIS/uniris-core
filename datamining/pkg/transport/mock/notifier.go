package mock

import (
	"log"

	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

type notifier struct {
}

//NewNotifier creates a new notifier
func NewNotifier() mining.Notifier {
	return &notifier{}
}

func (n *notifier) NotifyTransactionStatus(tx string, status mining.TransactionStatus) error {
	log.Printf("Transaction %s with status %s", tx, status.String())
	return nil
}
