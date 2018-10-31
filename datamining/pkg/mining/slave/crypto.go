package slave

import (
	"time"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/mining/slave/checks"
)

//Signer defines methods to handling signatures
type Signer interface {
	SignValidation(v Validation, pvKey string) (string, error)
	checks.Signer
}

//Validation represents a validation before its signature
type Validation struct {
	Status    datamining.ValidationStatus `json:"status"`
	Timestamp time.Time                   `json:"timestamp"`
	PublicKey string                      `json:"pubk"`
}
