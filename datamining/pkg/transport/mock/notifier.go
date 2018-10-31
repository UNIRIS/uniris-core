package mock

import (
	"log"

	"github.com/uniris/uniris-core/datamining/pkg/mining/master"
)

//NewNotifier creates a new notifier
func NewNotifier() master.Notifier {
	return notifier{}
}

type notifier struct{}

func (n notifier) NotifyTransactionStatus(tx string, status master.TransactionStatus) error {
	log.Printf("Transaction %s with status %s", tx, status.String())
	return nil
}
