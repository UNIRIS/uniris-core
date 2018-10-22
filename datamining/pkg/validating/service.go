package validating

import (
	"time"

	"github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/validating/checks"
)

//Service is the interface that provide methods for wallets validation
type Service interface {
	EndorseWalletAsMaster(*datamining.WalletData) (*datamining.Endorsement, string, error)
	EndorseBioWalletAsMaster(*datamining.BioData) (*datamining.Endorsement, error)
	AskWalletValidations([]Peer, *datamining.WalletData) ([]datamining.Validation, error)
	AskBioWalletValidations([]Peer, *datamining.BioData) ([]datamining.Validation, error)
}

type service struct {
	bioChecks      []checks.BioChecker
	dataChecks     []checks.DataChecker
	validRequester ValidationRequester
}

//NewService creates a approving service
func NewService(sig checks.Signer, validRequester ValidationRequester) Service {
	bioChecks := make([]checks.BioChecker, 0)
	dataChecks := make([]checks.DataChecker, 0)

	bioChecks = append(bioChecks, checks.NewSignatureChecker(sig))
	dataChecks = append(dataChecks, checks.NewSignatureChecker(sig))

	return &service{bioChecks, dataChecks, validRequester}
}

func (s service) EndorseWalletAsMaster(w *datamining.WalletData) (*datamining.Endorsement, string, error) {
	for _, c := range s.dataChecks {
		err := c.CheckDataWallet(w)
		if err != nil {
			return nil, "", err
		}
	}

	//TODO: defines POW values
	lastTxRvk := make([]string, 0)
	powRobotKey := ""
	powValidation := datamining.Validation{}
	masterValid := datamining.NewMasterValidation(
		lastTxRvk, powRobotKey, powValidation,
	)

	//TODO: retreive the old transaction hash
	oldHash := "old transaction hash"

	//TODO: create transaction hash
	txHash := "hash"

	endor := datamining.NewEndorsement(time.Now(), txHash, masterValid, nil)
	return endor, oldHash, nil
}

func (s service) EndorseBioWalletAsMaster(bw *datamining.BioData) (*datamining.Endorsement, error) {
	for _, c := range s.bioChecks {
		err := c.CheckBioWallet(bw)
		if err != nil {
			return nil, err
		}
	}

	//TODO: defines POW values
	lastTxRvk := make([]string, 0)
	powRobotKey := ""
	powValidation := datamining.Validation{}
	masterValid := datamining.NewMasterValidation(
		lastTxRvk, powRobotKey, powValidation,
	)

	//TODO: create transaction hash
	txHash := "hash"

	endor := datamining.NewEndorsement(time.Now(), txHash, masterValid, nil)
	return endor, nil
}

func (s service) AskWalletValidations(peers []Peer, w *datamining.WalletData) ([]datamining.Validation, error) {
	valids := make([]datamining.Validation, 0)

	for _, p := range peers {
		v, err := s.validRequester.RequestWalletValidation(p, w)
		if err != nil {
			return nil, err
		}
		valids = append(valids, v)
	}

	return valids, nil
}

func (s service) AskBioWalletValidations(peers []Peer, bd *datamining.BioData) ([]datamining.Validation, error) {
	valids := make([]datamining.Validation, 0)

	for _, p := range peers {
		v, err := s.validRequester.RequestBioValidation(p, bd)
		if err != nil {
			return nil, err
		}
		valids = append(valids, v)
	}

	return valids, nil
}
