package file

import (
	"encoding/hex"
	"io/ioutil"
	"path"
	"path/filepath"
)

//Reader create a keys reader struct
type Reader struct {
	keysDir string
}

//NewReader creates a new file reader
func NewReader() (r Reader, err error) {
	keysDir, err := filepath.Abs("../../../keys")
	if err != nil {
		return
	}

	return Reader{keysDir}, nil
}

//SharedRobotPrivateKey return the shared privatekey between robot
func (r Reader) SharedRobotPrivateKey() ([]byte, error) {
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

//SharedRobotPublicKey returns the shared publickey between robot
func (r Reader) SharedRobotPublicKey() ([]byte, error) {
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

//SharedBiodPublicKey returns the shared publickey between Biometric device
func (r Reader) SharedBiodPublicKey() ([]byte, error) {
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
