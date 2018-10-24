package mock

import (
	"log"

	"github.com/uniris/uniris-core/datamining/pkg/leading"
)

//NewNotifier creates a new notifier
func NewNotifier() leading.Notifier {
	return notifier{}
}

type notifier struct{}

func (n notifier) NotifyTransactionStatus(tx string, status leading.TransactionStatus) error {
	log.Printf("Transaction %s with status %s", tx, status.String())
	return nil
}
