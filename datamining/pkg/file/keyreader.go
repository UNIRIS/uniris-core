package file

import (
	"encoding/hex"
	"io/ioutil"
	"path"
	"path/filepath"
)

type reader struct {
	keysDir string
}

//NewReader creates a new file reader
func NewReader() (r reader, err error) {
	keysDir, err := filepath.Abs("../../../keys")
	if err != nil {
		return
	}

	return reader{keysDir}, nil
}

func (r reader) SharedRobotPrivateKey() ([]byte, error) {
	biodPvKey, err := ioutil.ReadFile(path.Join(r.keysDir, "sharedRobot.key"))
	if err != nil {
		return nil, err
	}

	b, err := hex.DecodeString(string(biodPvKey))
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (r reader) SharedRobotPublicKey() ([]byte, error) {
	biodPub, err := ioutil.ReadFile(path.Join(r.keysDir, "sharedRobot.pub"))
	if err != nil {
		return nil, err
	}

	b, err := hex.DecodeString(string(biodPub))
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (r reader) SharedBiodPublicKey() ([]byte, error) {
	robotPub, err := ioutil.ReadFile(path.Join(r.keysDir, "sharedBiod.pub"))
	if err != nil {
		return nil, err
	}

	b, err := hex.DecodeString(string(robotPub))
	if err != nil {
		return nil, err
	}

	return b, nil
}
