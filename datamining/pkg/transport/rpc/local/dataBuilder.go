package localrpc

import (
	"encoding/json"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	robot "github.com/uniris/uniris-core/datamining/pkg"
	crypto "github.com/uniris/uniris-core/datamining/pkg/crypto"
	wallet "github.com/uniris/uniris-core/datamining/pkg/wallet"
	formater "github.com/uniris/uniris-core/datamining/pkg/walletformating"
)

//DataBuilder defines methods to transform API entities for the domain layer
type DataBuilder struct{}

//ToGetWalletAck constitue details for the GetWalletAck rpc command
func (f DataBuilder) ToGetWalletAck(w wallet.Wallet, bw wallet.BioWallet) *api.GetWalletAck {
	return &api.GetWalletAck{
		Wd: &api.WalletDetails{
			EncWallet:     w.CWallet(),
			EncAeskey:     bw.CipherAesKey(),
			EncWalletaddr: bw.CipherAddrBio(),
		},
	}
}

//FromGetWalletRequest constitue details for the FromGetWalletRequest rpc command
func (f DataBuilder) FromGetWalletRequest(p *api.GetWalletRequest) (biohash robot.BioHash, err error) {
	r, err := crypto.NewReader()
	if err != nil {
		return
	}
	e, err := crypto.Newencrypter()
	if err != nil {
		return
	}
	privk, err := r.SharedRobotPrivateKey()
	if err != nil {
		return
	}
	biohash, err = e.Decrypt(privk, p.EncHashPerson)
	if err != nil {
		return
	}
	return
}

//FromSetWalletRequest constitue details for the SetWalletRequest rpc command
func (f DataBuilder) FromSetWalletRequest(p *api.SetWalletRequest) (fw formater.FormatedWallet, fbw formater.FormatedBioWallet, err error) {
	r, err := crypto.NewReader()
	if err != nil {
		return
	}
	e, err := crypto.Newencrypter()
	if err != nil {
		return
	}
	privk, err := r.SharedRobotPrivateKey()
	if err != nil {
		return
	}
	bioData, err := e.Decrypt(privk, p.EncryptedBioData)
	if err != nil {
		return
	}
	walletData, err := e.Decrypt(privk, p.EncryptedWalletData)
	if err != nil {
		return
	}

	err = json.Unmarshal(p.SignatureBioData.SigBiod, &fbw)
	if err != nil {
		return
	}
	err = json.Unmarshal(p.SignatureBioData.SigPerson, &fbw)
	if err != nil {
		return
	}
	err = json.Unmarshal(bioData, &fbw)
	if err != nil {
		return
	}
	err = json.Unmarshal(p.SignatureWalletData.SigBiod, &fw)
	if err != nil {
		return
	}
	err = json.Unmarshal(p.SignatureWalletData.SigPerson, &fw)
	if err != nil {
		return
	}
	err = json.Unmarshal(walletData, &fw)
	if err != nil {
		return
	}
	return
}
