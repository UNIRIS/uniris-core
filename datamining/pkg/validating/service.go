package validating

import (
	"time"

	"github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/validating/checkers"
)

//Validation represents a validation before its signature
type Validation struct {
	Status    datamining.ValidationStatus `json:"status"`
	Timestamp time.Time                   `json:"timestamp"`
	PublicKey string                      `json:"pubk"`
}

//Signer defines methods to handle signatures
type Signer interface {
	SignValidation(v Validation, pvKey string) (string, error)
	checkers.Signer
}

//Service is the interface that provide methods for wallets validation
type Service interface {
	ValidateWalletData(*datamining.WalletData) (datamining.Validation, error)
	ValidateBioData(*datamining.BioData) (datamining.Validation, error)
}

type service struct {
	bioChecks  []checkers.BioDataChecker
	dataChecks []checkers.WalletDataChecker
	robotKey   string
	robotPvKey string
	sig        Signer
}

//NewService creates a approving service
func NewService(sig Signer, robotKey, robotPvKey string) Service {
	bioChecks := make([]checkers.BioDataChecker, 0)
	dataChecks := make([]checkers.WalletDataChecker, 0)

	bioChecks = append(bioChecks, checkers.NewSignatureChecker(sig))
	dataChecks = append(dataChecks, checkers.NewSignatureChecker(sig))

	return service{
		bioChecks,
		dataChecks,
		robotKey,
		robotPvKey,
		sig,
	}
}

func (s service) ValidateWalletData(w *datamining.WalletData) (valid datamining.Validation, err error) {
	for _, c := range s.dataChecks {
		err = c.CheckWalletData(w)
		if err != nil {
			return
		}
	}
	v := Validation{
		PublicKey: s.robotKey,
		Status:    datamining.ValidationOK,
		Timestamp: time.Now(),
	}
	signature, err := s.sig.SignValidation(v, s.robotPvKey)
	if err != nil {
		return
	}
	valid = datamining.NewValidation(v.Status, v.Timestamp, v.PublicKey, signature)
	return
}

func (s service) ValidateBioData(bw *datamining.BioData) (valid datamining.Validation, err error) {
	for _, c := range s.bioChecks {
		err = c.CheckBioData(bw)
		if err != nil {
			return
		}
	}
	v := Validation{
		PublicKey: s.robotKey,
		Status:    datamining.ValidationOK,
		Timestamp: time.Now(),
	}
	signature, err := s.sig.SignValidation(v, s.robotPvKey)
	if err != nil {
		return
	}
	valid = datamining.NewValidation(v.Status, v.Timestamp, v.PublicKey, signature)
	return
}
