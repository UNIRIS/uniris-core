package localrpc

import (
	"golang.org/x/net/context"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	crypto "github.com/uniris/uniris-core/datamining/pkg/crypto"
	wallet "github.com/uniris/uniris-core/datamining/pkg/wallet"
)

type localsrvHandler struct {
	db wallet.Database
}

//GetWallet implements the protobuf GetWallet request handler
func (lh localsrvHandler) GetWallet(ctx context.Context, req *api.GetWalletRequest) (*api.GetWalletAck, error) {
	builder := DataBuilder{}
	r, _ := crypto.NewReader()
	e, _ := crypto.Newencrypter()
	privk, _ := r.SharedRobotPrivateKey()

	biohash, err := builder.FromGetWalletRequest(req)
	bw, err := lh.db.GetBioWallet(biohash)
	if err != nil {
		return &api.GetWalletAck{}, nil
	}

	clearaddr, _ := e.Decrypt(privk, bw.CipherAddrRobot())
	w, err := lh.db.GetWallet(clearaddr)
	if err != nil {
		return &api.GetWalletAck{}, nil
	}

	return builder.ToGetWalletAck(w, bw), nil
}

//SetWallet implements the protobuf SetWallet request handler
func (lh localsrvHandler) SetWallet(ctx context.Context, req *api.SetWalletRequest) (*api.SetWalletAck, error) {

	return &api.SetWalletAck{}, nil
}

//NewServerHandler create a new GRPC server handler
func NewServerHandler(d wallet.Database) api.WalletServer {
	return localsrvHandler{
		db: d,
	}
}
