package interpreter

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"regexp"
	"time"

	"github.com/uniris/uniris-core/pkg/crypto"

	"github.com/uniris/uniris-core/pkg/chain"
)

var stdFunctions = map[string]func(*Scope, map[string]interface{}) (interface{}, error){
	"now":                           nowFunc,
	"hash256":                       sha256HashFunc,
	"matchRegex":                    matchRegexFunc,
	"contains":                      containsFunc,
	"currentAddress":                currentAddressFunc,
	"currentPublicKey":              currentPublicKeyFunc,
	"currentCode":                   currentCodeFunc,
	"currentTriggers":               currentTriggersFunc,
	"currentResponseConditions":     currentResponseConditionsFunc,
	"currentInheritConditions":      currentInheritConditionsFunc,
	"currentActions":                currentActionsFunc,
	"currentKeys":                   currentKeysFunc,
	"currentContent":                currentContentFunc,
	"chainLength":                   chainLengthFunc,
	"countResponses":                countResponsesFunc,
	"responseRetries":               responseRetriesFunc,
	"responsePublicKey":             responsePublicKeyFunc,
	"responseContent":               responseContentFunc,
	"derivateKey":                   derivateKeysFunc,
	"decrypt":                       decryptFunc,
	"publicKey":                     publicKeyFunc,
	"privateKey":                    privateKeyFunc,
	"newUCOLedger":                  newUCOLedgerFunc,
	"newContract":                   newContractFunc,
	"checkMultisig":                 checkMultisigFunc,
	"incomingPostPaidFeeConditions": incomingPostPaidFeeConditionsFunc,
	"incomingContent":               incomingContentFunc,
	"incomingUCOLedger":             incomingnewUCOLedgerFunc,
}

func nowFunc(sc *Scope, args map[string]interface{}) (interface{}, error) {
	return float64(time.Now().Unix()), nil
}

func sha256HashFunc(sc *Scope, args map[string]interface{}) (interface{}, error) {
	h := sha256.New()

	data, exist := args["data"]
	if !exist {
		return nil, errors.New("data argument is missing")
	}

	switch data.(type) {
	case string:
		h.Write([]byte(data.(string)))
	case float64:
		var buf [8]byte
		binary.BigEndian.PutUint64(buf[:], math.Float64bits(data.(float64)))
		b := buf[:]
		h.Write(b)
	default:
		b, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		h.Write(b)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func matchRegexFunc(sc *Scope, args map[string]interface{}) (interface{}, error) {

	vPattern, exist := args["pattern"]
	if !exist {
		return nil, errors.New("pattern argument is missing")
	}
	pattern, ok := vPattern.(string)
	if !ok {
		return nil, errors.New("pattern argument must be a string")
	}

	vData, exist := args["data"]
	if !exist {
		return nil, errors.New("data argument is missing")
	}
	data, ok := vData.(string)
	if !ok {
		return nil, errors.New("data argument must be a string")
	}

	ok, err := regexp.MatchString(fmt.Sprintf("^%s$", pattern), data)
	if err != nil {
		return nil, err
	}
	return ok, nil
}

func containsFunc(sc *Scope, args map[string]interface{}) (interface{}, error) {
	vSource, exist := args["from"]
	if !exist {
		return nil, errors.New("in argument is missing")
	}

	in, ok := vSource.([]interface{})
	if !ok {
		return nil, errors.New("in argument must be an array")
	}

	value, exist := args["value"]
	if !exist {
		return nil, errors.New("value argument is missing")
	}

	for _, v := range in {
		if v == value {
			return true, nil
		}
	}

	return false, nil
}

func currentAddressFunc(sc *Scope, args map[string]interface{}) (interface{}, error) {
	return sc.contract.tx.Address(), nil
}

func currentPublicKeyFunc(sc *Scope, args map[string]interface{}) (interface{}, error) {
	return sc.contract.tx.PublicKey(), nil
}

func currentCodeFunc(sc *Scope, args map[string]interface{}) (interface{}, error) {
	code := sc.contract.tx.Data()["smartcontract"]
	return code, nil
}

func currentTriggersFunc(sc *Scope, args map[string]interface{}) (interface{}, error) {
	return sc.contract.Triggers, nil
}

func currentResponseConditionsFunc(sc *Scope, args map[string]interface{}) (interface{}, error) {
	return sc.contract.Conditions.Response, nil
}

func currentInheritConditionsFunc(sc *Scope, args map[string]interface{}) (interface{}, error) {
	return sc.contract.Conditions.Inherit, nil
}

func currentActionsFunc(sc *Scope, args map[string]interface{}) (interface{}, error) {
	return sc.contract.actions, nil
}

func currentKeysFunc(sc *Scope, args map[string]interface{}) (interface{}, error) {
	k, exist := sc.contract.tx.Data()["keys"]
	if !exist {
		return nil, errors.New("keys is missing in the transaction data")
	}
	var keys map[string]interface{}
	if err := json.Unmarshal(k, &keys); err != nil {
		return nil, err
	}

	vSelect, exist := args["select"]
	if exist {
		s, ok := vSelect.(string)
		if !ok {
			return nil, errors.New("from argument must be a string")
		}

		val, exist := keys[s]
		if !exist {
			return nil, errors.New("key does not exist")
		}
		return val, nil
	}

	return keys, nil
}

func currentContentFunc(sc *Scope, args map[string]interface{}) (interface{}, error) {
	content, exist := sc.contract.tx.Data()["content"]
	if !exist {
		return nil, errors.New("content is missing in the transaction data")
	}

	return content, nil
}

func chainLengthFunc(sc *Scope, args map[string]interface{}) (interface{}, error) {
	//TODO: implements DB query
	return 0, nil
}

func countResponsesFunc(sc *Scope, args map[string]interface{}) (interface{}, error) {

	// vContReg, exist := args["contentRegexp"]
	// if exist {
	// 	strCtReg, ok := vContReg.(string)
	// 	if !ok {
	// 		return nil, errors.New("contentRegexp argument must be of type string")
	// 	}
	// }

	//TODO: implements DB query
	return 0, nil
}

func responseRetriesFunc(sc *Scope, args map[string]interface{}) (interface{}, error) {
	//TODO: implements DB query
	return nil, nil
}

func derivateKeysFunc(sc *Scope, args map[string]interface{}) (interface{}, error) {
	//TODO: derivate key

	return keypair{}, nil
}

func decryptFunc(sc *Scope, args map[string]interface{}) (interface{}, error) {
	vData, exist := args["data"]
	if !exist {
		return nil, errors.New("data argument is missing")
	}
	data, ok := vData.(crypto.Cipher)
	if !ok {
		return nil, errors.New("data argument must be as []byte format")
	}

	clear, err := sc.sharedNodePvKey.Decrypt(data)
	if err != nil {
		return nil, err
	}
	return clear, nil
}

func publicKeyFunc(sc *Scope, args map[string]interface{}) (interface{}, error) {

	from, exist := args["from"]
	if !exist {
		return nil, errors.New("from argument is missing")
	}

	keypair, ok := from.(keypair)
	if !ok {
		return nil, errors.New("from argument must be as type keypair")
	}

	return keypair.public, nil
}

func privateKeyFunc(sc *Scope, args map[string]interface{}) (interface{}, error) {
	from, exist := args["from"]
	if !exist {
		return nil, errors.New("from argument is missing")
	}

	keypair, ok := from.(keypair)
	if !ok {
		return nil, errors.New("from argument must be as type keypair")
	}

	return keypair.private, nil
}

func responsePublicKeyFunc(sc *Scope, args map[string]interface{}) (interface{}, error) {
	return sc.response.PublicKey(), nil
}

func responseContentFunc(sc *Scope, args map[string]interface{}) (interface{}, error) {
	content, exist := sc.response.Data()["content"]
	if !exist {
		return nil, errors.New("content is missing in the transaction data")
	}
	return content, nil
}

func newUCOLedgerFunc(sc *Scope, args map[string]interface{}) (interface{}, error) {

	var fee float64
	if v, exist := args["fee"]; exist {
		fee = v.(float64)
	}

	var restTo float64
	if v, exist := args["restTo"]; exist {
		restTo = v.(float64)
	}

	return ucoLedger{
		fee:    fee,
		restTo: restTo,
	}, nil
}

func newContractFunc(sc *Scope, args map[string]interface{}) (interface{}, error) {

	vAddr, exist := args["address"]
	if !exist {
		return nil, errors.New("address argument is missing")
	}
	addr, ok := vAddr.(crypto.VersionnedHash)
	if !ok {
		return nil, errors.New("address argument must be as []byte format")
	}
	vTime, exist := args["timestamp"]
	if !exist {
		return nil, errors.New("address argument is missing")
	}
	timestamp, ok := vTime.(float64)
	if !ok {
		return nil, errors.New("timestamp argument must be as float64 format")
	}

	vPublicKey, exist := args["publicKey"]
	if !exist {
		return nil, errors.New("address argument is missing")
	}
	pubKb, ok := vPublicKey.(crypto.VersionnedKey)
	if !ok {
		return nil, errors.New("publicKey argument must be as []byte format")
	}
	pubK, err := crypto.ParsePublicKey(pubKb)
	if err != nil {
		return nil, err
	}

	vUcoLedger, exist := args["ucoLedger"]
	if !exist {
		return nil, errors.New("ucoLedger argument is missing")
	}

	ucoLedger, ok := vUcoLedger.(ucoLedger)
	if !ok {
		return nil, errors.New("ucoLedger argument must be as ucoLedger type")
	}

	var sContract []byte
	vSc, exist := args["smartcontract"]
	if exist {
		scStr, ok := vSc.(string)
		if !ok {
			return nil, errors.New("code argument must be as string type")
		}
		sContract = []byte(scStr)
	}

	var content []byte
	vContent, exist := args["content"]
	if exist {
		contentStr, ok := vContent.(string)
		if !ok {
			return nil, errors.New("content argument must be as string type")
		}
		content = []byte(contentStr)
	}

	var keys map[string][]byte
	vKeys, exist := args["keys"]
	if exist {
		keys, ok = vKeys.(map[string][]byte)
		if !ok {
			return nil, errors.New("content argument must be as string type")
		}
	}

	vPrevKey, exist := args["previousKey"]
	if !exist {
		return nil, errors.New("previousKey argument is missing")
	}
	pvKeyBytes, ok := vPrevKey.(crypto.VersionnedKey)
	if !ok {
		return nil, errors.New("previousKey argument must be as []byte format")
	}
	pvKey, err := crypto.ParsePrivateKey(pvKeyBytes)
	if err != nil {
		return nil, err
	}

	c := contract{
		addr:          addr,
		timestamp:     time.Unix(int64(timestamp), 0),
		publicKey:     pubK,
		ucoLedger:     ucoLedger,
		smartContract: sContract,
		content:       content,
		keys:          keys,
	}

	b, err := c.marshalBeforeSignature()
	if err != nil {
		return nil, err
	}
	sig, err := pvKey.Sign(b)
	if err != nil {
		return nil, err
	}

	c.sig = sig
	b, err = c.marshalBeforeOriginSignature()
	if err != nil {
		return nil, err
	}
	originSig, err := sc.originPvKey.Sign(b)
	if err != nil {
		return nil, err
	}
	c.originSig = originSig

	//!IMPORTANT
	//TODO: send the transaction

	return c, nil

}

func checkMultisigFunc(sc *Scope, args map[string]interface{}) (interface{}, error) {

	// vPubKey, exist := args["publicKey"]
	// if !exist {
	// 	return nil, errors.New("publicKey argument is missing")
	// }

	// var publicKey crypto.PublicKey

	// switch vPubKey.(type) {
	// case string:
	// 	b, err := hex.DecodeString(vPubKey.(string))
	// 	if err != nil {
	// 		return nil, errors.New("publicKey argument is not in hexadecimal")
	// 	}
	// 	publicKey, err = crypto.ParsePublicKey(b)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	// //TODO: check the multisignature section with the given public key

	return true, nil
}

func incomingPostPaidFeeConditionsFunc(sc *Scope, args map[string]interface{}) (interface{}, error) {
	if (sc.incoming.Conditions.PostPaidFee != literalExpression{value: ""}) {
		return sc.incoming.Conditions.PostPaidFee, nil
	}
	return nil, nil
}

func incomingContentFunc(sc *Scope, args map[string]interface{}) (interface{}, error) {
	content, exist := sc.incoming.tx.Data()["content"]
	if !exist {
		return nil, errors.New("content is missing in the transaction data")
	}
	return content, nil
}

func incomingnewUCOLedgerFunc(sc *Scope, args map[string]interface{}) (interface{}, error) {
	ledgerBytes, exist := sc.incoming.tx.Data()["ledger"]
	if !exist {
		return nil, errors.New("ledger is missing in the transaction data")
	}

	var ledger map[string]interface{}
	if err := json.Unmarshal(ledgerBytes, &ledger); err != nil {
		return nil, err
	}

	ucoLedgerBytes, exist := ledger["ucoLedger"]
	if !exist {
		return nil, errors.New("ucoLedger is missing from the ledger in the 0transaction data")
	}

	var ul ucoLedger
	if err := json.Unmarshal(ucoLedgerBytes.([]byte), &ul); err != nil {
		return nil, err
	}

	return ul, nil
}

type keypair struct {
	private crypto.PrivateKey
	public  crypto.PublicKey
}

type contract struct {
	addr          crypto.VersionnedHash
	ucoLedger     ucoLedger
	smartContract []byte
	content       []byte
	keys          map[string][]byte
	timestamp     time.Time
	publicKey     crypto.PublicKey
	sig           crypto.VersionnedHash
	originSig     crypto.VersionnedHash
}

func (c contract) marshalBeforeSignature() ([]byte, error) {
	pubKey, err := c.publicKey.Marshal()
	if err != nil {
		return nil, err
	}

	return json.Marshal(map[string]interface{}{
		"addr": c.addr,
		"type": chain.ContractTransactionType,
		"data": map[string]interface{}{
			"ledger": map[string]interface{}{
				"uco": map[string]interface{}{
					"fee":       c.ucoLedger.fee,
					"transfers": c.ucoLedger.transfers,
					"restTo":    c.ucoLedger.restTo,
				},
			},
			"smartContract": c.smartContract,
			"content":       c.content,
			"keys":          c.keys,
		},
		"timestamp": c.timestamp.Unix(),
		"publicKey": pubKey,
	})
}

func (c contract) marshalBeforeOriginSignature() ([]byte, error) {
	pubKey, err := c.publicKey.Marshal()
	if err != nil {
		return nil, err
	}

	return json.Marshal(map[string]interface{}{
		"addr": c.addr,
		"type": chain.ContractTransactionType,
		"data": map[string]interface{}{
			"ledger": map[string]interface{}{
				"uco": map[string]interface{}{
					"fee":       c.ucoLedger.fee,
					"transfers": c.ucoLedger.transfers,
					"restTo":    c.ucoLedger.restTo,
				},
			},
			"smartContract": c.smartContract,
			"content":       c.content,
			"keys":          c.keys,
		},
		"timestamp": c.timestamp.Unix(),
		"publicKey": pubKey,
		"sig":       c.sig,
	})
}

type ucoLedger struct {
	fee       float64
	transfers []transfer
	restTo    float64
}

type transfer struct {
	to     crypto.VersionnedHash
	amount float64
}
