package system

import (
	"fmt"
	"net"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

//CheckAutodiscoveryProcess check that the autodiscovery daemon is running
func CheckAutodiscoveryProcess(p discovery.Peer) error {
	selfport := p.Port()
	_, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", selfport))
	if err != nil {
		return err
	}
	return nil
}

//CheckDataProcess check that the autodiscovery daemon is running
func CheckDataProcess() error {
	// TBD
	return nil
}

//CheckMiningProcess check that the autodiscovery daemon is running
func CheckMiningProcess() error {
	//TBD
	return nil
}

//CheckAIProcess check that the autodiscovery daemon is running
func CheckAIProcess() error {
	//TBD
	return nil
}

//CheckScyllaProcess check that the autodiscovery daemon is running
func CheckScyllaProcess() error {
	//TBD
	return nil
}

//CheckRedisProcess check that the autodiscovery daemon is running
func CheckRedisProcess() error {
	//TBD
	return nil
}

//CheckRabitmqProcess check that the autodiscovery daemon is running
func CheckRabitmqProcess() error {
	//TBD
	return nil
}
