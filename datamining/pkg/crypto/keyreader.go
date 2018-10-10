package crypto

import (
	"encoding/hex"
	"io/ioutil"
	"path"
	"path/filepath"

	robot "github.com/uniris/uniris-core/datamining/pkg"
)

type reader struct {
	keysDir string
}

//NewReader creates a new file reader
func NewReader() (robot.KeyReader, error) {
	keysDir, err := filepath.Abs("../../../keys")
	if err != nil {
		return nil, err
	}

	return reader{keysDir}, nil
}

func (r reader) SharedRobotPrivateKey() (robot.PrivateKey, error) {
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

func (r reader) SharedRobotPublicKey() (robot.PublicKey, error) {
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

func (r reader) SharedBiodPublicKey() (robot.PublicKey, error) {
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
