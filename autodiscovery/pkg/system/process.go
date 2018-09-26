package system

import (
	"fmt"
	"net"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

//CheckAutodiscoveryProcess check that the autodiscovery daemon is running
func CheckAutodiscoveryProcess(p discovery.Peer) (bool, error) {
	selfport := p.Port()
	_, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", selfport))
	if err != nil {
		return false, err
	}
	return true, nil
}

//CheckDataProcess check that the autodiscovery daemon is running
func CheckDataProcess() (bool, error) {
	// TBD
	return true, nil
}

//CheckMiningProcess check that the autodiscovery daemon is running
func CheckMiningProcess() (bool, error) {
	//TBD
	return true, nil
}

//CheckAIProcess check that the autodiscovery daemon is running
func CheckAIProcess() (bool, error) {
	//TBD
	return true, nil
}

//CheckScyllaProcess check that the autodiscovery daemon is running
func CheckScyllaProcess() (bool, error) {
	//TBD
	return true, nil
}

//CheckRedisProcess check that the autodiscovery daemon is running
func CheckRedisProcess() (bool, error) {
	//TBD
	return true, nil
}

//CheckRabitmqProcess check that the autodiscovery daemon is running
func CheckRabitmqProcess() (bool, error) {
	//TBD
	return true, nil
}
