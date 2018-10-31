package slave

import (
	"time"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/mining/slave/checks"
)

//Service define for the slave mining process
type Service interface {
	Validate(data interface{}, txType datamining.TransactionType) (valid datamining.Validation, err error)
}

type service struct {
	sig        Signer
	robotKey   string
	robotPvKey string
	checks     map[datamining.TransactionType][]checks.Handler
}

//NewService creates slave mining service
func NewService(sig Signer, robotKey, robotPvKey string) Service {
	checks := map[datamining.TransactionType][]checks.Handler{
		datamining.CreateKeychainTransaction: []checks.Handler{
			checks.NewSignatureChecker(sig),
		},
		datamining.CreateBioTransaction: []checks.Handler{
			checks.NewSignatureChecker(sig),
		},
	}
	return service{sig, robotKey, robotPvKey, checks}
}

func (s service) Validate(data interface{}, txType datamining.TransactionType) (valid datamining.Validation, err error) {
	for _, c := range s.checks[txType] {
		err = c.CheckData(data)
		if err != nil {
			if c.IsCatchedError(err) {
				return s.buildValidation(datamining.ValidationKO)
			}
			return
		}
	}
	return s.buildValidation(datamining.ValidationOK)
}

func (s service) buildValidation(status datamining.ValidationStatus) (valid datamining.Validation, err error) {
	v := Validation{
		PublicKey: s.robotKey,
		Status:    status,
		Timestamp: time.Now(),
	}
	signature, err := s.sig.SignValidation(v, s.robotPvKey)
	if err != nil {
		return
	}
	return datamining.NewValidation(v.Status, v.Timestamp, v.PublicKey, signature), nil
}
