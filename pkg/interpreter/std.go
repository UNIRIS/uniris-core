package interpreter

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math"
	"time"
)

var stdFunctions = map[string]callable{
	"timestamp": timestampFunc{},
	"hash256":   sha256HashFunc{},
}

type callable interface {
	call(*Scope, ...interface{}) (interface{}, error)
}

type timestampFunc struct{}

func (f timestampFunc) call(sc *Scope, args ...interface{}) (interface{}, error) {
	return float64(time.Now().Unix()), nil
}

type sha256HashFunc struct{}

func (f sha256HashFunc) call(sc *Scope, args ...interface{}) (interface{}, error) {
	h := sha256.New()
	for _, v := range args {
		switch v.(type) {
		case string:
			h.Write([]byte(v.(string)))
		case float64:
			var buf [8]byte
			binary.BigEndian.PutUint64(buf[:], math.Float64bits(v.(float64)))
			b := buf[:]
			h.Write(b)
		default:
			return nil, errors.New("unsupported type to hash")
		}

	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
